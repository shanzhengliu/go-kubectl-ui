package internal

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Secret struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func SecretList(clientset *kubernetes.Clientset, namespace string) []Secret {
	secrectListClient := clientset.CoreV1().Secrets(namespace)
	secrects, error := secrectListClient.List(context.TODO(), apiv1.ListOptions{})
	if error != nil {
		fmt.Println(error.Error())
	}
	secrectList := []Secret{}
	for _, item := range secrects.Items {
		secrect := &Secret{
			Name:      item.Name,
			Namespace: item.Namespace,
		}
		secrectList = append(secrectList, *secrect)
	}
	return secrectList
}

func SecretDetail(clientset *kubernetes.Clientset, namespace string, name string) map[string]interface{} {
	secretListClient := clientset.CoreV1().Secrets(namespace)
	secret, err := secretListClient.Get(context.TODO(), name, apiv1.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	decodedData := make(map[string]interface{})
	if !strings.Contains(name, ".helm.") {
		for key, value := range secret.Data {
			decodedData[key] = string(value)
		}
	} else {
		for key, value := range secret.Data {

			gzipdata, error := gzipDecompress(value)

			if error == nil {
				var formatData map[string]interface{}
				json.Unmarshal(gzipdata, &formatData)
				// convertedJSON, _ := convertToMapInterface(string(gzipdata))
				decodedData[key] = formatData

			} else {
				decodedData[key] = string(value)
			}
		}
	}

	return decodedData
}

func convertToMapInterface(jsonStr string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{})
	convertToMapInterfaceRecursive(data, result)

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func convertToMapInterfaceRecursive(data interface{}, result map[string]interface{}) {
	switch value := data.(type) {
	case map[string]interface{}:
		for k, v := range value {
			if _, ok := result[k]; !ok {
				result[k] = make(map[string]interface{})
			}
			convertToMapInterfaceRecursive(v, result[k].(map[string]interface{}))
		}
	case []interface{}:
		for i, v := range value {
			if _, ok := result[strconv.Itoa(i)]; !ok {
				result[strconv.Itoa(i)] = make(map[string]interface{})
			}
			convertToMapInterfaceRecursive(v, result[strconv.Itoa(i)].(map[string]interface{}))
		}
	default:

		for k, v := range result {
			if _, ok := v.(map[string]interface{}); ok {
				result[k] = value
			}
		}
	}
}

func gzipDecompress(data []byte) ([]byte, error) {
	compressedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {

		return nil, err
	}
	defer r.Close()

	uncompressedData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return uncompressedData, nil
}

func SecretListHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	result := SecretList(clientset, ctxMap["namespace"].(string))
	ReturnTypeHandler(r.Context(), result, w, r)
}

func SecretDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	clientset := ctxMap["clientSet"].(*kubernetes.Clientset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SecretDetail(clientset, ctxMap["namespace"].(string), r.URL.Query().Get("secret")))
}
