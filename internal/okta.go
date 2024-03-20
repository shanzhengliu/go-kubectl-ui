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

	"golang.org/x/oauth2"
)

var currentState string
var currentNonce string

func ComputeFilename(key Key) (string, error) {
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

type CacheToken struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func test() {
	ctx := context.Background()

	var err error
	currentState, err = NewRand32()
	if err != nil {
		log.Fatalf("无法生成state: %v", err)
	}

	currentNonce, err = NewRand32()

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

	params, nil := NewParam([]string{"S256"})

	url := conf.AuthCodeURL(currentState, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", params.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", params.CodeChallengeMethod), oauth2.SetAuthURLParam("nonce", currentNonce))
	fmt.Printf("请访问此URL以进行认证: %v\n", url)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		receivedState := r.URL.Query().Get("state")
		if receivedState != currentState {

			http.Error(w, "State Incorrect", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")

		token, err := conf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", params.CodeVerifier))
		if err != nil {
			println(err.Error())
			log.Printf("Exchange Failed: %v", err)
			http.Error(w, "Can Get Token", http.StatusInternalServerError)
			return
		}

		key := Key{
			IssuerURL:   "https://resmed.okta.com/oauth2/auslcfr7vn87JUxjc2p7",
			ClientID:    "0oalcfi71iDpdoa0K2p7",
			ExtraScopes: []string{"email", "offline_access", "profile", "openid"},
		}
		filename, _ := ComputeFilename(key)
		idToken := token.Extra("id_token").(string)

		fmt.Fprintf(w, "获取到的Token: %+v\n", idToken)
		fmt.Fprintf(w, "获取到的Filename: %+v\n", filename)
		cacheToken := CacheToken{
			AccessToken: token.AccessToken,
			IdToken:     idToken,
		}
		jsonToken, _ := json.Marshal(cacheToken)

		os.WriteFile(filename, jsonToken, 777)

	})

	log.Println("listening on :8000...")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
