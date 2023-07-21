package internal

import (
	"context"

	v1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	PodImages  []PodImage  `json:"images"`
	Status     v1.PodPhase `json:"status"`
	CreateTime string      `json:"createTime"`
}

type PodImage struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func PodList(clientset *kubernetes.Clientset, namespace string) []Pod {
	podListClient := clientset.CoreV1().Pods(namespace)
	pod, error := podListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		panic(error.Error())
	}
	podList := []Pod{}
	for _, item := range pod.Items {
		containers := item.Status.ContainerStatuses
		returnImages := []PodImage{}
		for _, container := range containers {
			tempPodImage := &PodImage{
				Name: container.Image,
				Id:   container.ImageID,
			}
			returnImages = append(returnImages, *tempPodImage)
		}

		tempPod := &Pod{
			Name:       item.Name,
			Namespace:  item.Namespace,
			PodImages:  returnImages,
			Status:     item.Status.Phase,
			CreateTime: item.CreationTimestamp.String(),
		}
		podList = append(podList, *tempPod)
	}
	return podList
}
