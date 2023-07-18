package main

import (
	"fmt"
	"log"
	"modules/internal"
	"net/http"
)

func main() {
	internal.RouteInit()
	router := http.NewServeMux()
	router.HandleFunc("/", internal.DeploymentHandler)
	router.HandleFunc("/deployment", internal.DeploymentHandler)
	router.HandleFunc("/configmap", internal.ConfigMapListHandler)
	router.HandleFunc("/ingress", internal.IngressListHandler)
	router.HandleFunc("/api/configmap-detail", internal.ConfigMapDetailHandler)
	fmt.Println("listening 8080 port")
	log.Fatal(http.ListenAndServe(":8080", router))

}
