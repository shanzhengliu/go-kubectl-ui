package main_

import (
	"context"
	"embed"
	_ "embed"
	"encoding/base64"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"modules/internal"
	"net/http"
	"os/exec"

	"os"

	"github.com/creack/pty"
	"github.com/gorilla/mux"
	"github.com/olahol/melody"
	"github.com/rs/cors"
)

//go:embed frontend-build
var frontend embed.FS

var ctxMap map[string]interface{} = make(map[string]interface{})

func ContextAdd(ctx context.Context) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			contextReq := r.WithContext(ctx)
			f(w, contextReq)
		}
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func main_1() {
	var config string
	var namespace string
	var port string
	var path string
	var websitePassword string

	flag.StringVar(&config, "config", "minikube", "config context: eg: minikube")
	flag.StringVar(&namespace, "namespace", "default", "namespace: eg: namespate")
	flag.StringVar(&port, "port", "8080", "port: eg: 8080")
	flag.StringVar(&path, "path", "NONE", "path: eg: /root/.kube/config")
	flag.StringVar(&websitePassword, "websitePassword", "", "password: eg: 123456")
	flag.Parse()
	if os.Getenv("KUBE_CONFIG") != "" {
		config = os.Getenv("KUBE_CONFIG")
	}
	if os.Getenv("KUBE_NAMESPACE") != "" {
		namespace = os.Getenv("KUBE_NAMESPACE")
	}

	if os.Getenv("KUBE_PORT") != "" {
		port = os.Getenv("KUBE_PORT")
	}
	if os.Getenv("KUBE_CONFIG_PATH") != "" {
		path = os.Getenv("KUBE_CONFIG_PATH")
	}
	if os.Getenv("KUBE_WEBSITE_PASSWORD") != "" {
		websitePassword = os.Getenv("KUBE_WEBSITE_PASSWORD")
	}
	fmt.Println("config: " + config)
	fmt.Println("namespace: " + namespace)
	fmt.Println("port: " + port)
	fmt.Println("path: " + path)
	cmd := exec.Command("kubectl", "config", "use-context", config)
	err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	ctx := context.WithValue(context.Background(), "map", ctxMap)
	ctxMap["environment"] = config
	ctxMap["static"] = frontend
	ctxMap["namespace"] = namespace
	ctxMap["websitePassword"] = websitePassword
	if path == "NONE" {
		path = internal.Kubeconfig()
	}
	internal.RouteInit(ctx, path)
	router := mux.NewRouter()

	subFs, err := fs.Sub(frontend, "frontend-build")
	if err != nil {
		log.Fatal(err)
	}

	c := exec.Command("sh")
	f, err := pty.Start(c)
	if err != nil {
		log.Printf("start pty failed: %v\n", err)
		return
	}
	m := melody.New()
	go func() {
		for {
			buff := make([]byte, 4096)
			read, err := f.Read(buff)
			if err != nil {
				return
			}
			encodedMsg := base64.StdEncoding.EncodeToString(buff[:read])
			m.Broadcast([]byte(encodedMsg))

		}
	}()

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		decodedMsg, err := base64.StdEncoding.DecodeString(string(msg))
		if err != nil {
			log.Printf("decode message failed: %v\n", err)
			return
		}
		f.Write(decodedMsg)
	})
	router.PathPrefix("/assets/").Handler(http.FileServer(http.FS(subFs)))
	router.PathPrefix("/xterm/").Handler(http.FileServer(http.FS(subFs)))
	router.PathPrefix("/js/").Handler(http.FileServer(http.FS(subFs)))
	router.HandleFunc("/", Chain(internal.HomeHandler, ContextAdd(ctx)))
	router.HandleFunc("/auth", Chain(internal.AuthHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/login", Chain(internal.LoginHandler, ContextAdd(ctx)))
	router.HandleFunc("/deployment", Chain(internal.DeploymentHandler, ContextAdd(ctx)))
	router.HandleFunc("/configmap", Chain(internal.ConfigMapListHandler, ContextAdd(ctx)))
	router.HandleFunc("/ingress", Chain(internal.IngressListHandler, ContextAdd(ctx)))
	router.HandleFunc("/pod", Chain(internal.PodListHandler, ContextAdd(ctx)))
	router.HandleFunc("/service", Chain(internal.ServiceListHandler, ContextAdd(ctx)))
	router.HandleFunc("/resource", Chain(internal.ResourceUseageHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/user-info", Chain(internal.UserInfoHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/configmap-detail", Chain(internal.ConfigMapDetailHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/context-change", Chain(internal.ContextChangeHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/context-list", Chain(internal.ContextListHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/current-context", Chain(internal.CurrentContextHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/podLogs", Chain(internal.PodLogHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/podYaml", Chain(internal.PodtoYamlHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/deploymentYaml", Chain(internal.DeploymentYamlHandler, ContextAdd(ctx)))
	router.HandleFunc("/webshell", Chain(internal.WebShellHandler, ContextAdd(ctx)))
	router.HandleFunc("/localshell", Chain(internal.LocalShellHandler, ContextAdd(ctx)))
	router.HandleFunc("/ws/webshell", Chain(internal.ServeWsTerminalHandler, ContextAdd(ctx)))
	router.HandleFunc("/ws/localshell", Chain(func(w http.ResponseWriter, r *http.Request) { m.HandleRequest(w, r) }, ContextAdd(ctx)))
	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
	})
	corHandler := cor.Handler(router)
	fmt.Println("listening: " + port + " port")
	fmt.Println("link: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, corHandler))

}
