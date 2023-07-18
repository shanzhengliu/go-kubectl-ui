package internal

import (
	"context"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Ingress struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Rules     []Rule `json:"rules"`
}

type Rule struct {
	HOST string `json:"host"`
}

func IngressList(clientset *kubernetes.Clientset, namespace string) []Ingress {
	ingressListClient := clientset.NetworkingV1().Ingresses(namespace)
	ingress, error := ingressListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		panic(error.Error())
	}
	ingressList := []Ingress{}

	for _, item := range ingress.Items {
		rules := item.Spec.Rules
		returnRules := []Rule{}
		for _, rule := range rules {
			tempRule := &Rule{
				HOST: rule.Host,
			}
			returnRules = append(returnRules, *tempRule)
		}

		ingress := &Ingress{
			Name:      item.Name,
			Namespace: item.Namespace,
			Rules:     returnRules,
		}
		ingressList = append(ingressList, *ingress)
	}
	return ingressList
}
