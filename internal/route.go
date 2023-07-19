package internal

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"text/template"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset = nil

func RouteInit() {
	config, err := clientcmd.BuildConfigFromFlags("", Kubeconfig())
	if err != nil {
		panic(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

}

func TemplateRender(ctx context.Context, path string) *template.Template {
	tplblob := ctx.Value("static").(embed.FS)

	template, err := template.ParseFS(tplblob, "static/"+path+".html", "static/navigator.tpl")
	if err != nil {
		panic(err)
	}
	return template
}

func DeploymentHandler(w http.ResponseWriter, r *http.Request) {
	result := DeploymentList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "deployment").Execute(w, result)
}

func ConfigMapListHandler(w http.ResponseWriter, r *http.Request) {
	result := ConfigMapList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "configmap").Execute(w, result)

}

func IngressListHandler(w http.ResponseWriter, r *http.Request) {
	result := IngressList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "ingress").Execute(w, result)
}

func PodListHandler(w http.ResponseWriter, r *http.Request) {
	result := PodList(clientset, r.Context().Value("namespace").(string))
	TemplateRender(r.Context(), "pod").Execute(w, result)
}

func ConfigMapDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConfigMapDetail(clientset, r.Context().Value("namespace").(string), r.URL.Query().Get("configmap")))
}
