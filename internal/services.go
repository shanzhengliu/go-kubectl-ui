package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
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

func ServiceList(ctxMap map[string]interface{}, clientSet *kubernetes.Clientset, namespace string) []Service {
	servicesClient := clientSet.CoreV1().Services(namespace)
	currentContext := ctxMap["environment"].(string)
	service, err := servicesClient.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		ErrorHandlerFunction(http.StatusInternalServerError, nil, "Failed to get services: "+err.Error())
		return []Service{}
	}

	var serviceList []Service

	for _, item := range service.Items {
		selectorString := MapToString(item.Spec.Selector)
		localPort := ""
		servicePort := ""
		isForward := false
		for key := range ctxMap {
			if strings.Contains(key, fmt.Sprintf("service#forward#%s#%s#%s", currentContext, namespace, item.Name)) {
				fmt.Println(key)
				localPort = strings.Split(key, "#")[5]
				servicePort = strings.Split(key, "#")[6]
				isForward = true
				continue
			}
		}
		service := &Service{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Type:        string(item.Spec.Type),
			Selector:    selectorString,
			IsForward:   isForward,
			LocalPort:   localPort,
			ServicePort: servicePort,
		}
		serviceList = append(serviceList, *service)
	}
	return serviceList
}

func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := ServiceList(ctxMap, clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func ServiceForwardHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	currentContext := ctxMap["environment"].(string)

	namespace := r.URL.Query().Get("namespace")
	serviceName := r.URL.Query().Get("service")
	servicePort := r.URL.Query().Get("servicePort")
	localPort := r.URL.Query().Get("localPort")

	isPortAvailable := IsPortAvailable(localPort)
	if !isPortAvailable {
		ErrorHandlerFunction(http.StatusBadRequest, w, "Port is already in use")
		return
	}

	key := fmt.Sprintf("service#forward#%s#%s#%s#%s#%s", currentContext, namespace, serviceName, localPort, servicePort)
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
	servicePort := r.URL.Query().Get("servicePort")
	if namespace == "" || serviceName == "" || localPort == "" || servicePort == "" {
		http.Error(w, "Namespace, serviceName, localPort, servicePort query parameters are required", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("service#forward#%s#%s#%s#%s#%s", currentContext, namespace, serviceName, localPort, servicePort)

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
