package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMap struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func ConfigMapList(clientset *kubernetes.Clientset, namespace string) []ConfigMap {
	configMapListClient := clientset.CoreV1().ConfigMaps(namespace)
	configMap, error := configMapListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	configMapList := []ConfigMap{}
	for _, item := range configMap.Items {
		configMap := &ConfigMap{
			Name:      item.Name,
			Namespace: item.Namespace,
		}
		configMapList = append(configMapList, *configMap)
	}
	return configMapList
}

func ConfigMapDetail(clientset *kubernetes.Clientset, namespace string, name string) map[string]string {
	configMapListClient := clientset.CoreV1().ConfigMaps(namespace)
	configMap, error := configMapListClient.Get(context.TODO(), name, apiv1.GetOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	return configMap.Data
}

func ConfigMapListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ConfigMapList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func ConfigMapDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConfigMapDetail(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("configmap")))
}
