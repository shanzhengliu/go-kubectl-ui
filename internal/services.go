package internal

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Selector  string `json:"selecter"`
}

func ServiceList(clientset *kubernetes.Clientset, namespace string) []Service {
	servicesClient := clientset.CoreV1().Services(namespace)
	service, error := servicesClient.List(context.TODO(), v1.ListOptions{})
	if error != nil {
		panic(error.Error())
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
