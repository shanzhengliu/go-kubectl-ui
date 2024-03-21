package internal

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"k8s.io/client-go/tools/clientcmd"
)

func GetUserIsOIDC(ctx context.Context) bool {
	ctxMap := ctx.Value("map").(map[string]interface{})
	currentContext := ctxMap["environment"].(string)
	if ctxMap["isOIDC-"+currentContext] == nil {

		configPath := ctxMap["configPath"].(string)

		config, err := clientcmd.LoadFromFile(configPath)
		if err != nil {
			fmt.Printf(err.Error())
		}
		user := config.Contexts[currentContext].AuthInfo
		userInfo, exists := config.AuthInfos[user]
		if !exists {
			fmt.Printf(user)
		}
		if userInfo.Exec != nil {
			ConfigUserExecArgsMap(userInfo.Exec.Args, ctxMap, "oidcMap")
			oidcMap := ctxMap["oidcMap"].(map[string][]string)
			conf := &oauth2.Config{
				ClientID:    oidcMap["oidc-client-id"][0],
				RedirectURL: "http://localhost:8000",
				Scopes:      oidcMap["oidc-extra-scope"],
				Endpoint: oauth2.Endpoint{
					AuthURL:  oidcMap["oidc-issuer-url"][0] + "/v1/authorize",
					TokenURL: oidcMap["oidc-issuer-url"][0] + "/v1/token",
				},
			}
			ctxMap["oidcConfig"] = conf
			_, exists := oidcMap["oidc-login"]

			ctxMap["isOIDC-"+currentContext] = exists
			return exists
		}
		return false
	} else {
		return ctxMap["isOIDC-"+currentContext].(bool)
	}
}
