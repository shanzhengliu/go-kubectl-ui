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

type RenderResult struct {
	ResultList     interface{} `json:"resultList"`
	ContextList    []string    `json:"contextList"`
	CurrentContext string      `json:"currentContext"`
	Namespace      string      `json:"namespace"`
}

func RouteInit(ctx context.Context, path string) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ctxMap := ctx.Value("map").(map[string]interface{})
	ctxMap["configPath"] = path
	ctxMap["clientSet"] = clientset
	ctxMap["contextList"] = maps.Keys(KubeconfigList(path))
}

func TemplateRender(ctx context.Context, path string, resultList interface{}, w http.ResponseWriter, r *http.Request) {
	ctxMap := ctx.Value("map").(map[string]interface{})
	tplblob := ctxMap["static"].(embed.FS)

	template, err := template.ParseFS(tplblob, "static/"+path+".html", "static/tpl/navigator.html", "static/tpl/contextSwitch.html", "static/tpl/style.html")
	if err != nil {
		panic(err)
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

func DeploymentHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := DeploymentList(clientset, ctxMap["namespace"].(string))
	TemplateRender(r.Context(), "deployment", result, w, r)
}

func ConfigMapListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ConfigMapList(clientset, ctxMap["namespace"].(string))
	TemplateRender(r.Context(), "configmap", result, w, r)

}

func IngressListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := IngressList(clientset, ctxMap["namespace"].(string))
	TemplateRender(r.Context(), "ingress", result, w, r)
}

func PodListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := PodList(clientset, ctxMap["namespace"].(string))
	TemplateRender(r.Context(), "pod", result, w, r)
}

func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ServiceList(clientset, ctxMap["namespace"].(string))
	TemplateRender(r.Context(), "service", result, w, r)
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
