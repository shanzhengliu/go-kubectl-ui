package internal

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"text/template"

	"golang.org/x/exp/maps"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset = nil

type RenderResult struct {
	ResultList     interface{} `json:"resultList"`
	ContextList    []string    `json:"contextList"`
	CurrentContext string      `json:"currentContext"`
	Namespace      string      `json:"namespace"`
}

func RouteInit(ctx context.Context) context.Context {
	path := Kubeconfig()
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		panic(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, "clientSet", clientset)
	ctx = context.WithValue(ctx, "contextList", maps.Keys(KubeconfigList(path)))
	return ctx
}

func TemplateRender(ctx context.Context, path string, resultList interface{}, w http.ResponseWriter, r *http.Request) {
	tplblob := ctx.Value("static").(embed.FS)

	template, err := template.ParseFS(tplblob, "static/"+path+".html", "static/navigator.tpl", "static/contextSwitch.tpl")
	if err != nil {
		panic(err)
	}
	template.Execute(w, RenderResultInit(ctx, resultList))
}

func RenderResultInit(ctx context.Context, resultList interface{}) *RenderResult {
	renderResult := &RenderResult{
		ResultList:     resultList,
		ContextList:    ctx.Value("contextList").([]string),
		CurrentContext: ctx.Value("environment").(string),
		Namespace:      ctx.Value("namespace").(string),
	}
	return renderResult
}

func DeploymentHandler(w http.ResponseWriter, r *http.Request) {
	result := DeploymentList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "deployment", result, w, r)
}

func ConfigMapListHandler(w http.ResponseWriter, r *http.Request) {
	result := ConfigMapList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "configmap", result, w, r)

}

func IngressListHandler(w http.ResponseWriter, r *http.Request) {
	result := IngressList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "ingress", result, w, r)
}

func PodListHandler(w http.ResponseWriter, r *http.Request) {
	result := PodList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "pod", result, w, r)
}

func ConfigMapDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConfigMapDetail(clientset, r.Context().Value("namespace").(string), r.URL.Query().Get("configmap")))
}
