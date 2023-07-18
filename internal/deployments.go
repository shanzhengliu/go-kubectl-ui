package internal

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployment struct {
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
	Status     int32       `json:"status"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func DeploymentList(clientset *kubernetes.Clientset, namespace string) []Deployment {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployment, error := deploymentsClient.List(context.TODO(), v1.ListOptions{})
	if error != nil {
		panic(error.Error())
	}
	deploymentList := []Deployment{}
	for _, item := range deployment.Items {
		deployment := &Deployment{
			Name:       item.Name,
			Containers: []Container{},
			Status:     item.Status.AvailableReplicas,
		}
		for _, container := range item.Spec.Template.Spec.Containers {

			container := &Container{
				Name:  container.Name,
				Image: container.Image,
			}
			deployment.Containers = append(deployment.Containers, *container)
		}
		deploymentList = append(deploymentList, *deployment)
	}
	return deploymentList
}
