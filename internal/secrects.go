package internal

import (
	"context"
	"fmt"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Secrect struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func SecrectList(clientset *kubernetes.Clientset, namespace string) []Secrect {
	secrectListClient := clientset.CoreV1().Secrets(namespace)
	secrects, error := secrectListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	secrectList := []Secrect{}
	for _, item := range secrects.Items {
		secrect := &Secrect{
			Name:      item.Name,
			Namespace: item.Namespace,
		}
		secrectList = append(secrectList, *secrect)
	}
	return secrectList
}
