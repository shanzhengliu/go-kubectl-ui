package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"modules/internal"
	"net/http"
	"os/exec"
)

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

	flag.StringVar(&config, "config", "minikube", "config context: eg: minikube")
	flag.StringVar(&namespace, "namespace", "default", "namespace: eg: namespate")
	flag.Parse()
	cmd := exec.Command("kubectl", "config", "use-context", config)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), "namespace", namespace)
	ctx = context.WithValue(ctx, "environment", config)

	internal.RouteInit()
	router := http.NewServeMux()

	router.HandleFunc("/", Chain(internal.DeploymentHandler, ContextAdd(ctx)))
	router.HandleFunc("/deployment", Chain(internal.DeploymentHandler, ContextAdd(ctx)))
	router.HandleFunc("/configmap", Chain(internal.ConfigMapListHandler, ContextAdd(ctx)))
	router.HandleFunc("/ingress", Chain(internal.IngressListHandler, ContextAdd(ctx)))
	router.HandleFunc("/pod", Chain(internal.PodListHandler, ContextAdd(ctx)))
	router.HandleFunc("/api/configmap-detail", Chain(internal.ConfigMapDetailHandler, ContextAdd(ctx)))
	fmt.Println("listening 8080 port")
	log.Fatal(http.ListenAndServe(":8080", router))

}
