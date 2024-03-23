package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"

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
		fmt.Println(error.Error())
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
		fmt.Println(err.Error())
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
			fmt.Println(err.Error())
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
		fmt.Println(err.Error())
	}
	podYamlString, err := yaml.Marshal(pod)
	if err != nil {
		fmt.Println(err.Error())
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

func PodListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := PodList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func PodLogHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	log := PodLog(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("pod"), r.URL.Query().Get("container"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(log))
}

func PodtoYamlHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	yaml := PodtoYaml(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("pod"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(yaml))
}
