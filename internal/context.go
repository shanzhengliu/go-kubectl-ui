package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func ContextChange(ctx context.Context, contextName string, namespace string) {
	ctxMap := ctx.Value("map").(map[string]interface{})
	path := ctxMap["configPath"].(string)
	restconfig, err := buildConfigFromFlags(contextName, path)
	if err != nil {
		fmt.Println(err)
	}
	clientset, _ := kubernetes.NewForConfig(restconfig)
	restclient, _ := rest.RESTClientFor(restconfig)
	ctxMap["clientSet"] = clientset
	ctxMap["environment"] = contextName
	ctxMap["namespace"] = namespace
	ctxMap["restConfig"] = restconfig
	ctxMap["restClient"] = restclient
}

type KubeContext struct {
	Context   string `json:"context"`
	Namespace string `json:"namespace"`
}

func CurrentContext(ctx context.Context) KubeContext {
	ctxMap := ctx.Value("map").(map[string]interface{})
	KubeContext := KubeContext{
		Context:   ctxMap["environment"].(string),
		Namespace: ctxMap["namespace"].(string),
	}
	return KubeContext
}

func ContextList(ctx context.Context) []string {
	ctxMap := ctx.Value("map").(map[string]interface{})
	return ctxMap["contextList"].([]string)
}

func ContextListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ContextList(r.Context()))
}

func GetCurrentContextFromKubeCmd() string {
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, _ := clientcmd.LoadFromFile(kubeconfig)
	return config.CurrentContext
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
