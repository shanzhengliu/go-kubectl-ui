package internal

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"k8s.io/client-go/tools/clientcmd"
)

func GetUserIsOIDC(ctx context.Context, currentContext string) bool {
	ctxMap := ctx.Value("map").(map[string]interface{})
	return ctxMap["oidcMap-"+currentContext] != nil
}

func GenerateKubeUserAuthMap(ctx context.Context, context string) {

	ctxMap := ctx.Value("map").(map[string]interface{})
	configPath := ctxMap["configPath"].(string)
	config, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		fmt.Printf(err.Error())
	}
	user := config.Contexts[context].AuthInfo
	userInfo, exists := config.AuthInfos[user]
	if !exists {
		fmt.Printf(user)
	}
	if userInfo.Exec != nil {
		ConfigUserExecArgsMap(userInfo.Exec.Args, ctxMap, "oidcMap-"+context)

		oidcMap := ctxMap["oidcMap-"+context].(map[string][]string)
		clientSecret := ""
		if oidcMap["oidc-client-secret"] != nil {
			clientSecret = oidcMap["oidc-client-secret"][0]
		}

		conf := &oauth2.Config{
			ClientID:     oidcMap["oidc-client-id"][0],
			RedirectURL:  "http://localhost:8000",
			ClientSecret: clientSecret,
			Scopes:       oidcMap["oidc-extra-scope"],
			Endpoint: oauth2.Endpoint{
				AuthURL:  oidcMap["oidc-issuer-url"][0] + "/v1/authorize",
				TokenURL: oidcMap["oidc-issuer-url"][0] + "/v1/token",
			},
		}
		ctxMap["oidcConfig-"+context] = conf
	}

}
