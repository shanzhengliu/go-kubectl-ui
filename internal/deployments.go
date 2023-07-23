package internal

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type Deployment struct {
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
	Selector   string      `json:"selector"`
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
		selectString := MapToString(item.Spec.Selector.MatchLabels)
		deployment := &Deployment{
			Name:       item.Name,
			Containers: []Container{},
			Selector:   selectString,
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

func DeploymentToYaml(clientset *kubernetes.Clientset, namespace string, name string) string {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployment, error := deploymentsClient.Get(context.TODO(), name, v1.GetOptions{})
	if error != nil {
		panic(error.Error())
	}
	deploymentYaml, err := yaml.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	return string(deploymentYaml)
}
