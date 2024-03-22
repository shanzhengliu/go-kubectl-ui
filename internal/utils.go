package internal

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func Kubeconfig() string {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return *kubeconfig
}

func KubeconfigList(configPath string) map[string]*api.Context {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	return config.Contexts
}

func MapToString(data map[string]string) string {
	var result string
	for key, value := range data {
		result += key + ":" + value + "\n"
	}
	return result
}

func GetCacheFileNameByCtxMap(ctxMap map[string]interface{}, kubeContext string) string {

	oidcMap := ctxMap["oidcMap-"+kubeContext].(map[string][]string)
	oidcIssuerUrl := oidcMap["oidc-issuer-url"][0]
	oidcClientId := oidcMap["oidc-client-id"][0]
	oidcExtraScopes := oidcMap["oidc-extra-scope"]
	oidcClientSecret := ""
	if oidcMap["oidc-client-secret"] != nil {
		oidcClientSecret = oidcMap["oidc-client-secret"][0]
	}
	key := Key{
		IssuerURL:    oidcIssuerUrl,
		ClientID:     oidcClientId,
		ExtraScopes:  oidcExtraScopes,
		ClientSecret: oidcClientSecret,
	}
	filename, _ := ComputeFilename(key)
	return filename

}
