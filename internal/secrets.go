package internal

import (
	"context"
	"fmt"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Secret struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func SecretList(clientset *kubernetes.Clientset, namespace string) []Secret {
	secretListClient := clientset.CoreV1().Secrets(namespace)
	secret, error := secretListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	secretList := []Secret{}
	for _, item := range secret.Items {
		secret := &Secret{
			Name:      item.Name,
			Namespace: item.Namespace,
		}
		secretList = append(secretList, *secret)
	}
	return secretList
}
