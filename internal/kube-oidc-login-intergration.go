package internal

import (
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

type Key struct {
	IssuerURL      string
	ClientID       string
	ClientSecret   string
	Username       string
	ExtraScopes    []string
	CACertFilename string
	CACertData     string
	SkipTLSVerify  bool
}

type CacheToken struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

type OidcLoginSuccess struct {
	AccessToken string                 `json:"access_token"`
	UserInfo    map[string]interface{} `json:"user_info"`
}

func GenerateStateAndNonce() (string, string, Params) {
	currentState, _ := NewRand32()
	currentNonce, _ := NewRand32()
	params, _ := NewParam([]string{"S256"})
	return currentState, currentNonce, params
}

func ComputeFilename(key Key) (string, error) {
	s := sha256.New()
	e := gob.NewEncoder(s)
	if err := e.Encode(&key); err != nil {
		return "", fmt.Errorf("could not encode the key: %w", err)
	}
	h := hex.EncodeToString(s.Sum(nil))
	return h, nil
}

func ConfigUserExecArgsMap(args []string, ctxMap map[string]interface{}, mapKey string) {
	oidcMap := make(map[string][]string)
	for _, arg := range args {
		key, value := "", ""

		if strings.Contains(arg, "=") {
			key = strings.Replace(strings.Split(arg, "=")[0], "--", "", -1)
			value = strings.Split(arg, "=")[1]
			if _, exist := oidcMap[key]; !exist {
				oidcMap[key] = []string{value}
			} else {
				oidcMap[key] = append(oidcMap[key], value)
			}
		} else {
			key, value = arg, ""
			oidcMap[key] = []string{value}
		}

	}
	ctxMap[mapKey] = oidcMap
}

func OIDCLoginUrlGenerate(context context.Context, oidcIssuerUrl string, clientId string, ClientSecret string, redirectUrl string, scopes []string, params Params, currentNonce string, currentState string) string {
	ctxMap := context.Value("map").(map[string]interface{})
	currentContext := ctxMap["environment"].(string)
	conf := ctxMap["oidcConfig-"+currentContext].(*oauth2.Config)
	url := ""
	if ClientSecret != "" {
		url = conf.AuthCodeURL(currentState, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", params.CodeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", params.CodeChallengeMethod), oauth2.SetAuthURLParam("nonce", currentNonce), oauth2.SetAuthURLParam("client_secret", ClientSecret))
	} else {
		url = conf.AuthCodeURL(currentState, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", params.CodeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", params.CodeChallengeMethod), oauth2.SetAuthURLParam("nonce", currentNonce))
	}
	return url

}

func OktaCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	receivedState := r.URL.Query().Get("state")
	currentContext := ctxMap["environment"].(string)
	params := ctxMap["params"].(Params)
	conf := ctxMap["oidcConfig-"+currentContext].(*oauth2.Config)

	if receivedState != ctxMap["state"] {

		http.Error(w, "State Incorrect", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", params.CodeVerifier))
	if err != nil {
		println(err.Error())
		log.Printf("Exchange Failed: %v", err)
		http.Error(w, "Can Get Token", http.StatusInternalServerError)
		return
	}

	filename := GetCacheFileNameByCtxMap(ctxMap, currentContext)
	idToken := token.Extra("id_token").(string)
	kubeDefaultPath := ctxMap["kubeDefaultPath"].(string)
	writePath := kubeDefaultPath + "/cache/oidc-login/" + filename
	accessToken := token.AccessToken
	cacheToken := CacheToken{
		IdToken:     idToken,
		AccessToken: accessToken,
	}

	ctxMap["cacheToken-"+filename] = cacheToken
	jsonToken, _ := json.Marshal(cacheToken)
	os.WriteFile(writePath, jsonToken, 0777)
	w.Write([]byte("<html><script>(function(){ window.location.href=\"http://localhost:" + ctxMap["applicationPort"].(string) + "\"; })()</script></html>"))
}
