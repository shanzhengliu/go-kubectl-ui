package internal

import (
	"context"
	"fmt"

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
			_, exists := oidcMap["oidc-login"]
			ctxMap["isOIDC-"+currentContext] = exists
			return exists
		}
		return false
	} else {
		return ctxMap["isOIDC-"+currentContext].(bool)
	}
}
