package internal

import (
	"context"
	"fmt"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ResourceQuotaData struct {
	Name        string                    `json:"name"`
	Namespace   string                    `json:"namespace"`
	ResourceMap map[string]ResourceStatus `json:"resourceMap"`
}

type ResourceStatus struct {
	Used int64 `json:"used"`
	Hard int64 `json:"hard"`
	Free int64 `json:"free"`
}

func ResourceList(clientset *kubernetes.Clientset, namespace string) []ResourceQuotaData {
	resourceQuota := clientset.CoreV1().ResourceQuotas(namespace)
	resourceQuotaList, error := resourceQuota.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	resourceQuotaListData := []ResourceQuotaData{}
	resourceMap := make(map[string]ResourceStatus)
	for _, item := range resourceQuotaList.Items {
		for resource, _ := range item.Status.Hard {
			used := item.Status.Used[resource]
			hard := item.Status.Hard[resource]
			resourceMap[resource.String()] = ResourceStatus{
				Used: used.Value(),
				Hard: hard.Value(),
				Free: hard.Value() - used.Value(),
			}
		}

		resourceQuota := &ResourceQuotaData{
			Name:        item.Name,
			Namespace:   item.Namespace,
			ResourceMap: resourceMap,
		}
		resourceQuotaListData = append(resourceQuotaListData, *resourceQuota)
	}
	return resourceQuotaListData
}
