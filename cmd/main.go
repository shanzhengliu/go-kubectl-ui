package main

import (
	"log"
	"modules/internal"
	"net/http"
)

func main() {

	// config, err := clientcmd.BuildConfigFromFlags("", internal.Kubeconfig())

	// if err != nil {
	// 	panic(err)
	// }
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(internal.DeploymentList(clientset, "shared-helios"))

	internal.RouteInit()
	router := http.NewServeMux()
	router.HandleFunc("/", internal.DeploymentHandler)
	router.HandleFunc("/deployment", internal.DeploymentHandler)
	router.HandleFunc("/configmap", internal.ConfigMapListHandler)
	router.HandleFunc("/ingress", internal.IngressListHandler)
	router.HandleFunc("/api/configmap-detail", internal.ConfigMapDetailHandler)
	log.Fatal(http.ListenAndServe(":8080", router))

}
