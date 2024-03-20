package main

import (
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"modules/internal"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// 全局变量用于存储state，注意：实际应用中应避免使用全局变量来存储重要信息。
var currentState string
var currentNonce string

func computeFilename(key Key) (string, error) {
	s := sha256.New()
	e := gob.NewEncoder(s)
	if err := e.Encode(&key); err != nil {
		return "", fmt.Errorf("could not encode the key: %w", err)
	}
	h := hex.EncodeToString(s.Sum(nil))
	return h, nil
}

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

type cacheToken struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func main() {
	ctx := context.Background()

	// Generate a new state for this auth session
	var err error
	currentState, err = internal.NewState() // 生成state
	if err != nil {
		log.Fatalf("无法生成state: %v", err)
	}

	currentNonce, err = internal.NewNonce() // 生成nonce
	if err != nil {
		log.Fatalf("无法生成nonce: %v", err)
	}

	// codeVerifier, codeChallenge, _ := createCodeVerifierAndChallenge()
	conf := &oauth2.Config{
		ClientID:    "0oalcfi71iDpdoa0K2p7",
		RedirectURL: "http://localhost:8000",
		Scopes:      []string{"email", "offline_access", "profile", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://resmed.okta.com/oauth2/auslcfr7vn87JUxjc2p7/v1/authorize",
			TokenURL: "https://resmed.okta.com/oauth2/auslcfr7vn87JUxjc2p7/v1/token",
		},
	}

	params, nil := internal.NewParam([]string{"S256"})

	url := conf.AuthCodeURL(currentState, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", params.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", params.CodeChallengeMethod), oauth2.SetAuthURLParam("nonce", currentNonce))
	fmt.Printf("请访问此URL以进行认证: %v\n", url)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		receivedState := r.URL.Query().Get("state")
		if receivedState != currentState {
			log.Printf("无效的state: 收到的state=%s, 期望的state=%s", receivedState, currentState)
			http.Error(w, "State不匹配", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")

		token, err := conf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", params.CodeVerifier))
		if err != nil {
			println(err.Error())
			log.Printf("Exchange失败: %v", err)
			http.Error(w, "无法获取token", http.StatusInternalServerError)
			return
		}

		key := Key{
			IssuerURL:   "https://resmed.okta.com/oauth2/auslcfr7vn87JUxjc2p7",
			ClientID:    "0oalcfi71iDpdoa0K2p7",
			ExtraScopes: []string{"email", "offline_access", "profile", "openid"},
		}
		filename, _ := computeFilename(key)
		idToken := token.Extra("id_token").(string)

		fmt.Fprintf(w, "获取到的Token: %+v\n", idToken)
		fmt.Fprintf(w, "获取到的Filename: %+v\n", filename)
		cacheToken := cacheToken{
			AccessToken: token.AccessToken,
			IdToken:     idToken,
		}
		// 写入文件
		jsonToken, _ := json.Marshal(cacheToken)

		os.WriteFile(filename, jsonToken, 777)

	})

	log.Println("启动HTTP服务器,监听8000端口")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
