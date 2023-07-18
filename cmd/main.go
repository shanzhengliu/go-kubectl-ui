package main

import (
	"fmt"
	"log"
	"modules/internal"
	"net/http"
	"os/exec"
)

func main() {

	cmd := exec.Command("kubectl", "config", "use-context", "nonprod")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

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
