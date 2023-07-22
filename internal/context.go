package internal

import (
	"context"

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
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(restconfig)
	ctxMap["clientSet"] = clientset
	ctxMap["environment"] = contextName
	ctxMap["namespace"] = namespace
}
