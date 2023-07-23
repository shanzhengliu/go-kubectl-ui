package main

import (
	"context"
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"modules/internal"
	"net/http"
	"os/exec"
)

//go:embed static
var static embed.FS

var ctxMap map[string]interface{} = make(map[string]interface{})

func loadTeamplate(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "static", static)
	return ctx
}

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

func main() {
	var config string
	var namespace string
	var port string
	var path string

	flag.StringVar(&config, "config", "minikube", "config context: eg: minikube")
	flag.StringVar(&namespace, "namespace", "default", "namespace: eg: namespate")
	flag.StringVar(&port, "port", "8080", "port: eg: 8080")
	flag.StringVar(&path, "path", "NONE", "path: eg: /root/.kube/config")
	flag.Parse()
	cmd := exec.Command("kubectl", "config", "use-context", config)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), "map", ctxMap)
	ctxMap["environment"] = config
	ctxMap["static"] = static
	ctxMap["namespace"] = namespace
	if path == "NONE" {
		path = internal.Kubeconfig()
	}
	internal.RouteInit(ctx, path)
	router := http.NewServeMux()
	router.HandleFunc("/", Chain(internal.DeploymentHandler, ContextAdd(ctx)))
	router.HandleFunc("/deployment", Chain(internal.DeploymentHandler, ContextAdd(ctx)))
	router.HandleFunc("/configmap", Chain(internal.ConfigMapListHandler, ContextAdd(ctx)))
	router.HandleFunc("/ingress", Chain(internal.IngressListHandler, ContextAdd(ctx)))
	router.HandleFunc("/pod", Chain(internal.PodListHandler, ContextAdd(ctx)))
	router.HandleFunc("/service", Chain(internal.ServiceListHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/configmap-detail", Chain(internal.ConfigMapDetailHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/context-change", Chain(internal.ContextChangeHandler, ContextAdd(ctx)))
	fmt.Println("listening: " + port + " port")
	log.Fatal(http.ListenAndServe(":"+port, router))

}
