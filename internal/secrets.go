package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Secret struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func SecretList(clientset *kubernetes.Clientset, namespace string) []Secret {
	secrectListClient := clientset.CoreV1().Secrets(namespace)
	secrects, error := secrectListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	secrectList := []Secret{}
	for _, item := range secrects.Items {
		secrect := &Secret{
			Name:      item.Name,
			Namespace: item.Namespace,
		}
		secrectList = append(secrectList, *secrect)
	}
	return secrectList
}

func SecretDetail(clientset *kubernetes.Clientset, namespace string, name string) map[string]string {
	secretListClient := clientset.CoreV1().Secrets(namespace)
	secret, err := secretListClient.Get(context.TODO(), name, apiv1.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	decodedData := make(map[string]string)
	for key, value := range secret.Data {
		decodedData[key] = string(value)
	}

	return decodedData
}

func SecretListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := SecretList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func SecretDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SecretDetail(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("secret")))
}
