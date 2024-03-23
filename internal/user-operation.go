package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func OktaLogout(ctx context.Context, currentContext string) {
	ctxMap := ctx.Value("map").(map[string]interface{})
	filename := GetCacheFileNameByCtxMap(ctxMap, currentContext)
	cacheFilePath := ctxMap["kubeDefaultPath"].(string) + "/cache/oidc-login/" + filename
	os.Remove(cacheFilePath)
	ctxMap["cacheToken-"+filename] = nil
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})

	if ctxMap["oidcMap-"+ctxMap["environment"].(string)] == nil {
		w.WriteHeader(200)
		w.Write([]byte("don't need to login"))
		return
	}
	conf := ctxMap["oidcConfig-"+ctxMap["environment"].(string)].(*oauth2.Config)
	filename := GetCacheFileNameByCtxMap(ctxMap, ctxMap["environment"].(string))
	cacheToken := ctxMap["cacheToken-"+filename]
	if cacheToken == nil {
		w.WriteHeader(401)
		w.Write([]byte("need to login, cacheToken is nil"))
		return
	}
	oidcIssuerUrl := ctxMap["oidcMap-"+ctxMap["environment"].(string)].(map[string][]string)["oidc-issuer-url"][0]
	accessToken := cacheToken.(CacheToken).AccessToken
	userInfo, err := conf.Client(context.Background(), &oauth2.Token{AccessToken: accessToken}).Get(oidcIssuerUrl + "/v1/userinfo")
	if err != nil || userInfo.StatusCode != 200 {
		w.WriteHeader(401)
		w.Write([]byte("need to login"))
		return
	}
	defer userInfo.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(userInfo.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func OIDCLogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	OktaLogout(r.Context(), ctxMap["environment"].(string))
	w.WriteHeader(200)
	w.Write([]byte("logout"))
}

func OIDCLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	isOIDC := GetUserIsOIDC(r.Context(), ctxMap["environment"].(string))

	if isOIDC {
		oidcMap := ctxMap["oidcMap-"+ctxMap["environment"].(string)].(map[string][]string)
		currentState, currentNonce, params := GenerateStateAndNonce()
		ctxMap["state"] = currentState
		ctxMap["nonce"] = currentNonce
		ctxMap["params"] = params
		oidcClientSecret := ""
		if oidcMap["oidc-client-secret"] != nil {
			oidcClientSecret = oidcMap["oidc-client-secret"][0]
		}
		url := OIDCLoginUrlGenerate(r.Context(), oidcMap["oidc-issuer-url"][0], oidcMap["oidc-client-id"][0], oidcClientSecret, "http://localhost:8000", oidcMap["oidc-extra-scope"], params, currentNonce, currentState)
		w.WriteHeader(200)
		//response := map[string]string{"url": url}
		w.Write([]byte(url))

	} else {
		w.WriteHeader(201)
		w.Write([]byte("not oidc"))
	}
}
