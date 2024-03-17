package internal

import (
	"context"
	"fmt"
	"strings"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ResourceQuotaData struct {
	Name        string                    `json:"name"`
	Namespace   string                    `json:"namespace"`
	ResourceMap map[string]ResourceStatus `json:"resourceMap"`
}

type ResourceStatus struct {
	Used int64  `json:"used"`
	Hard int64  `json:"hard"`
	Free int64  `json:"free"`
	Unit string `json:"unit"`
}

func ResourceList(clientset *kubernetes.Clientset, namespace string) []ResourceQuotaData {
	resourceQuota := clientset.CoreV1().ResourceQuotas(namespace)
	resourceQuotaList, err := resourceQuota.List(context.TODO(), apiv1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	resourceQuotaListData := make([]ResourceQuotaData, 0, len(resourceQuotaList.Items))

	for _, item := range resourceQuotaList.Items {
		resourceMap := make(map[string]ResourceStatus)

		for resource, hard := range item.Status.Hard {
			used := item.Status.Used[resource]
			hardValue := hard.Value()
			usedValue := used.Value()
			unit := ""

			if strings.Contains(resource.String(), "memory") {
				hardValue /= 1024 * 1024
				usedValue /= 1024 * 1024
				unit = "MB"
			}
			if strings.Contains(resource.String(), "cpu") {
				unit = "Core"
			}

			resourceMap[resource.String()] = ResourceStatus{
				Used: usedValue,
				Hard: hardValue,
				Free: hardValue - usedValue,
				Unit: unit,
			}
		}

		resourceQuotaListData = append(resourceQuotaListData, ResourceQuotaData{
			Name:        item.Name,
			Namespace:   item.Namespace,
			ResourceMap: resourceMap,
		})
	}

	return resourceQuotaListData
}
