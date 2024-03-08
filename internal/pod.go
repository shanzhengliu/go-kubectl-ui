package internal

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type Pod struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	PodImages  []PodImage  `json:"images"`
	Status     v1.PodPhase `json:"status"`
	CreateTime string      `json:"createTime"`
}

type PodImage struct {
	ContainerName   string `json:"containerName"`
	Name            string `json:"name"`
	Id              string `json:"id"`
	ContainerStatus string `json:"containerStatus"`
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
				Name:            container.Image,
				Id:              container.ImageID,
				ContainerName:   container.Name,
				ContainerStatus: calculateRunningState(container.State),
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

func PodLog(clientset *kubernetes.Clientset, namespace string, name string, container string) string {
	podClient := clientset.CoreV1().Pods(namespace)
	stream, err := podClient.GetLogs(name, &v1.PodLogOptions{Container: container}).Stream(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	message := ""
	for {
		buf := make([]byte, 2000)
		numBytes, err := stream.Read(buf)
		if err == io.EOF {
			break
		}
		if numBytes == 0 {
			continue
		}
		if err != nil {
			panic(err.Error())
		}
		message += string(buf[:numBytes])
	}
	defer stream.Close()
	return message
}

func PodtoYaml(clientset *kubernetes.Clientset, namespace string, name string) string {
	podClient := clientset.CoreV1().Pods(namespace)
	pod, err := podClient.Get(context.Background(), name, apiv1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	podYamlString, err := yaml.Marshal(pod)
	if err != nil {
		panic(err.Error())
	}
	return string(podYamlString)
}

func calculateRunningState(state v1.ContainerState) string {
	if state.Running != nil {
		return "Running"
	}
	if state.Terminated != nil {
		return "Terminated"
	}
	if state.Waiting != nil {
		return "Waiting"
	}
	return "Unknown"
}
