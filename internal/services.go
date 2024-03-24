package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"syscall"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Type        string `json:"type"`
	Selector    string `json:"selector"`
	IsForward   bool   `json:"isForward"`
	LocalPort   string `json:"localPort"`
	ServicePort string `json:"servicePort"`
}

var (
	mu sync.Mutex
)

func ServiceList(clientset *kubernetes.Clientset, namespace string) []Service {
	servicesClient := clientset.CoreV1().Services(namespace)
	service, error := servicesClient.List(context.TODO(), v1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}

	serviceList := []Service{}

	for _, item := range service.Items {
		selectorString := MapToString(item.Spec.Selector)
		service := &Service{
			Name:      item.Name,
			Namespace: item.Namespace,
			Type:      string(item.Spec.Type),
			Selector:  selectorString,
		}
		serviceList = append(serviceList, *service)
	}
	return serviceList
}

func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ServiceList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func ServiceForwardHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	currentContext := ctxMap["environment"].(string)

	namespace := r.URL.Query().Get("namespace")
	serviceName := r.URL.Query().Get("service")
	servicePort := r.URL.Query().Get("servicePort")
	localPort := r.URL.Query().Get("localPort")
	key := fmt.Sprintf("service-forward-%s-%s-%s-%s", currentContext, namespace, serviceName, localPort)
	if ctxMap[key] != nil {
		http.Error(w, "Port-forward process already exists for given service", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("kubectl", "port-forward", "--address", "0.0.0.0", fmt.Sprintf("svc/%s", serviceName), fmt.Sprintf("%s:%s", localPort, servicePort), "-n", namespace)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start port-forward process: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	ctxMap[key] = cmd.Process.Pid
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Port forwarding started successfully with key %s and PID %d", key, cmd.Process.Pid)))

}

func StopPortForward(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	currentContext := ctxMap["environment"].(string)
	namespace := r.URL.Query().Get("namespace")
	serviceName := r.URL.Query().Get("service")
	localPort := r.URL.Query().Get("localPort")

	if namespace == "" || serviceName == "" || localPort == "" {
		http.Error(w, "Namespace, serviceName and localPort query parameters are required", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("service-forward-%s-%s-%s-%s", currentContext, namespace, serviceName, localPort)

	mu.Lock()
	pid, exists := ctxMap[key]
	if !exists {
		http.Error(w, "No port-forward process found with given key", http.StatusNotFound)
		mu.Unlock()
		return
	}
	delete(ctxMap, key)
	mu.Unlock()

	if err := syscall.Kill(pid.(int), syscall.SIGTERM); err != nil {
		http.Error(w, "Failed to stop port-forward process: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Successfully stopped port forwarding for key %s", key)))
}
