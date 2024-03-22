package internal

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	terminal "github.com/maoqide/kubeutil/pkg/terminal"
	wsterminal "github.com/maoqide/kubeutil/pkg/terminal/websocket"
	"golang.org/x/exp/maps"
	"golang.org/x/oauth2"
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
	// ctxMap["restClient"] = restClient
	ctxMap["contextList"] = maps.Keys(KubeconfigList(path))
	oktaCacheInitFromOS(ctxMap)
	for _, currentCtx := range ctxMap["contextList"].([]string) {
		GenerateKubeUserAuthMap(ctx, currentCtx)
	}

}

func oktaCacheInitFromOS(ctxMap map[string]interface{}) {
	files, err := os.ReadDir(ctxMap["kubeDefaultPath"].(string) + "/cache/oidc-login")
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

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	TemplateRender(r.Context(), "auth", "", w, r)
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

func DeploymentHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := DeploymentList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func ConfigMapListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ConfigMapList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func IngressListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := IngressList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func PodListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := PodList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ServiceList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func ConfigMapDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConfigMapDetail(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("configmap")))
}

func ContextChangeHandler(w http.ResponseWriter, r *http.Request) {
	ContextChange(r.Context(), r.URL.Query().Get("context"), r.URL.Query().Get("namespace"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func CurrentContextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CurrentContext(r.Context()))
}

func ContextListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ContextList(r.Context()))
}

func PodLogHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	log := PodLog(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("pod"), r.URL.Query().Get("container"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(log))
}

func PodtoYamlHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	yaml := PodtoYaml(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("pod"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(yaml))
}

func DeploymentYamlHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	yaml := DeploymentToYaml(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("deployment"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(yaml))
}

func WebShellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "webshell", "", w, r)
}

func LocalShellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "localshell", "", w, r)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	websitePassword := ctxMap["websitePassword"].(string)
	if websitePassword == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
		return
	}
	password := r.URL.Query().Get("password")
	if websitePassword == password {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	} else {
		w.WriteHeader(401)
		w.Write([]byte("false"))
	}
}

func ResourceUseageHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ResourceList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	TemplateRender(r.Context(), "index", "", w, r)
}

func OIDCLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	isOIDC := GetUserIsOIDC(r.Context(), ctxMap["environment"].(string))

	if isOIDC {
		oidcMap := ctxMap["oidcMap-"+ctxMap["environment"].(string)].(map[string][]string)
		currentState, currentNonce, params := GenerateStateAndNonce()
		ctxMap["state"] = currentState
		ctxMap["nonce"] = currentNonce
		ctxMap["params"] = params
		oidcClientSecret := ""
		if oidcMap["oidc-client-secret"] != nil {
			oidcClientSecret = oidcMap["oidc-client-secret"][0]
		}
		url := OIDCLoginUrlGenerate(r.Context(), oidcMap["oidc-issuer-url"][0], oidcMap["oidc-client-id"][0], oidcClientSecret, "http://localhost:8000", oidcMap["oidc-extra-scope"], params, currentNonce, currentState)
		w.WriteHeader(200)
		//response := map[string]string{"url": url}
		w.Write([]byte(url))

	} else {
		w.WriteHeader(201)
		w.Write([]byte("not oidc"))
	}

}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})

	if ctxMap["oidcMap-"+ctxMap["environment"].(string)] == nil {
		w.WriteHeader(200)
		w.Write([]byte("don't need to login"))
		return
	}
	oidcMap := ctxMap["oidcMap-"+ctxMap["environment"].(string)].(map[string][]string)

	oidcIssuerUrl := oidcMap["oidc-issuer-url"][0]
	oidcClientId := oidcMap["oidc-client-id"][0]
	oidcClientSecret := ""
	if oidcMap["oidc-client-serect"] != nil {
		oidcClientSecret = oidcMap["oidc-client-serect"][0]
	}
	oidcExtraScopes := oidcMap["oidc-extra-scope"]
	conf := ctxMap["oidcConfig-"+ctxMap["environment"].(string)].(*oauth2.Config)

	key := Key{
		IssuerURL:    oidcIssuerUrl,
		ClientID:     oidcClientId,
		ExtraScopes:  oidcExtraScopes,
		ClientSecret: oidcClientSecret,
	}
	filename, _ := ComputeFilename(key)
	cacheToken := ctxMap["cacheToken-"+filename]
	if cacheToken == nil {
		w.WriteHeader(401)
		w.Write([]byte("need to login, cacheToken is nil"))
		return
	}
	accessToken := cacheToken.(CacheToken).AccessToken
	userInfo, err := conf.Client(context.Background(), &oauth2.Token{AccessToken: accessToken}).Get(oidcIssuerUrl + "/v1/userinfo")
	if err != nil || userInfo.StatusCode != 200 {
		w.WriteHeader(401)
		w.Write([]byte("need to login"))
		return
	}
	defer userInfo.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(userInfo.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
