package internal

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	terminal "github.com/maoqide/kubeutil/pkg/terminal"
	wsterminal "github.com/maoqide/kubeutil/pkg/terminal/websocket"
	"golang.org/x/exp/maps"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type RenderResult struct {
	ResultList     interface{} `json:"resultList"`
	ContextList    []string    `json:"contextList"`
	CurrentContext string      `json:"currentContext"`
	Namespace      string      `json:"namespace"`
}

func RouteInit(ctx context.Context, path string) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	ctxMap := ctx.Value("map").(map[string]interface{})
	ctxMap["restConfig"] = config
	ctxMap["configPath"] = path
	ctxMap["clientSet"] = clientset
	ctxMap["contextList"] = maps.Keys(KubeconfigList(path))
	oktaCacheInitFromOS(ctxMap)
	for _, currentCtx := range ctxMap["contextList"].([]string) {
		GenerateKubeUserAuthMap(ctx, currentCtx)
	}

}
func createDirIfNotExist(dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func oktaCacheInitFromOS(ctxMap map[string]interface{}) {
	cachePath := ctxMap["kubeDefaultPath"].(string) + "/cache/oidc-login"
	err := createDirIfNotExist(cachePath)
	files, err := os.ReadDir(cachePath)

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "oidc-login" {
			continue
		}

		cache := LoadCacheToken(ctxMap["kubeDefaultPath"].(string) + "/cache/oidc-login/" + file.Name())

		if cache.AccessToken != "" {
			ctxMap["cacheToken-"+file.Name()] = cache
		}
	}
}

func LoadCacheToken(path string) CacheToken {
	file, err := os.ReadFile(path)
	var cacheToken CacheToken
	if err != nil {
		return cacheToken
	}

	err = json.Unmarshal(file, &cacheToken)
	if err != nil {
		return cacheToken
	}

	return cacheToken

}

func TemplateRender(ctx context.Context, path string, resultList interface{}, w http.ResponseWriter, r *http.Request) {
	ctxMap := ctx.Value("map").(map[string]interface{})
	tplblob := ctxMap["static"].(embed.FS)

	template, err := template.ParseFS(tplblob, "frontend-build/"+path+".html")
	if err != nil {
		fmt.Println(err)
	}
	template.Execute(w, RenderResultInit(ctx, resultList))
}

func RenderResultInit(ctx context.Context, resultList interface{}) *RenderResult {
	ctxMap := ctx.Value("map").(map[string]interface{})
	renderResult := &RenderResult{
		ResultList:     resultList,
		ContextList:    ctxMap["contextList"].([]string),
		CurrentContext: ctxMap["environment"].(string),
		Namespace:      ctxMap["namespace"].(string),
	}
	return renderResult
}

func ReturnTypeHandler(context context.Context, resultList interface{}, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(resultList)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonData)
	return

}

func WebShellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "webshell", "", w, r)
}

func LocalShellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "localshell", "", w, r)
}

func DyPodLogTemplateHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "dypodlog", "", w, r)
}

func ServeWsTerminalHandler(w http.ResponseWriter, r *http.Request) {
	cmd := []string{"sh"}
	ctxMap := r.Context().Value("map").(map[string]interface{})
	namespace := ctxMap["namespace"].(string)
	podName := r.URL.Query().Get("pod")
	containerName := r.URL.Query().Get("container")
	pty, err := wsterminal.NewTerminalSession(w, r, nil)
	if err != nil {
		log.Printf("get pty failed: %v\n", err)
		return
	}
	defer func() {
		log.Println("close session.")
		pty.Close()
	}()
	client := ctxMap["clientSet"].(*kubernetes.Clientset)
	if err != nil {
		log.Printf("get kubernetes client failed: %v\n", err)
		return
	}
	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		log.Printf("get kubernetes client failed: %v\n", err)
		return
	}
	ok, err := terminal.ValidatePod(pod, containerName)
	if !ok {
		msg := fmt.Sprintf("Validate pod error! err: %v", err)
		log.Println(msg)
		pty.Write([]byte(msg))
		pty.Done()
		return
	}
	restConfig := ctxMap["restConfig"].(*rest.Config)

	err = PodExec(client, restConfig, cmd, pty, namespace, podName, containerName)
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		log.Println(msg)
		pty.Write([]byte(msg))
		pty.Done()
	}
}

func PodExec(clientset *kubernetes.Clientset, restconfig *rest.Config, cmd []string, ptyHandler terminal.PtyHandler, namespace string, podName string, containerName string) error {
	defer func() {
		ptyHandler.Done()
	}()

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&apiv1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     !(ptyHandler.Stdin() == nil),
		Stdout:    !(ptyHandler.Stdout() == nil),
		Stderr:    !(ptyHandler.Stderr() == nil),
		TTY:       ptyHandler.Tty(),
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
	if err != nil {
		return err
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler.Stdin(),
		Stdout:            ptyHandler.Stdout(),
		Stderr:            ptyHandler.Stderr(),
		TerminalSizeQueue: ptyHandler,
		Tty:               ptyHandler.Tty(),
	})
	return err
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "index", "", w, r)
}

type ProxyRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    json.RawMessage   `json:"body"`
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	var reqData ProxyRequest

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientReq, err := http.NewRequest(reqData.Method, reqData.URL, bytes.NewReader(reqData.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range reqData.Headers {
		clientReq.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(clientReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range resp.Header {
		w.Header().Set(key, value[0])
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
