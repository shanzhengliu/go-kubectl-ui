package internal

import (
	"context"
	"fmt"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Selector  string `json:"selector"`
}

func ServiceList(clientset *kubernetes.Clientset, namespace string) []Service {
	servicesClient := clientset.CoreV1().Services(namespace)
	service, error := servicesClient.List(context.TODO(), v1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}

	serviceList := []Service{}

	for _, item := range service.Items {
		selectorString := MapToString(item.Spec.Selector)
		service := &Service{
			Name:      item.Name,
			Namespace: item.Namespace,
			Type:      string(item.Spec.Type),
			Selector:  selectorString,
		}
		serviceList = append(serviceList, *service)
	}
	return serviceList
}

func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ServiceList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}
