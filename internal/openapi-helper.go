package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Max upload size ~10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = UnzipAndSave(r, w, "/tmp/kubectl-go-upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully!"))
}

func StartOpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Path string `json:"path"`
		Port string `json:"port"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Path == "" || requestBody.Port == "" {
		http.Error(w, "Path and port are required in the request body", http.StatusBadRequest)
		return
	}
	if !IsPortAvailable(requestBody.Port) {
		http.Error(w, "Port is already in use", http.StatusBadRequest)
		return
	}
	StartOpenAPIFunction(requestBody.Path, requestBody.Port, r, w)

	w.Write([]byte("OK"))
}

func StopOpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Path string `json:"path"`
		Port string `json:"port"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Path == "" || requestBody.Port == "" {
		http.Error(w, "Path and port are required in the request body", http.StatusBadRequest)
		return
	}
	StopOpenAPIFunction(requestBody.Path, requestBody.Port, r, w)
}

func StopAllOpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	openAPIList := GetCurrentOpenAPIListPortAndFileName(r)
	for _, openAPIItem := range openAPIList {

		ctxMap := r.Context().Value("map").(map[string]interface{})
		key := "openapi#server#" + openAPIItem.Path + "#" + openAPIItem.Port
		srv, ok := ctxMap[key].(*http.Server)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No server found"))
			return
		}
		if err := srv.Shutdown(r.Context()); err != nil {
			log.Fatalf("Shutdown(): %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to shutdown server"))
			return
		}
		delete(ctxMap, key)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Server shutdown successfully"))
}

func GetOpenAPIValueFromRequest(r *http.Request, doc *openapi3.T) (int, string, string) {
	status := 200
	contentType := "application/json"
	example := ""
	r.Host = ""
	if r.Header.Get("openapi-status-code") != "" {
		status, _ = strconv.Atoi(r.Header.Get("openapi-status-code"))

	}
	if r.Header.Get("openapi-content-type") != "" {
		contentType = r.Header.Get("openapi-content-type")
	}
	if r.Header.Get("openapi-example") != "" {
		example = r.Header.Get("openapi-example")
	}
	return status, contentType, example
}

func MethodResponse(method string, pathItem *openapi3.PathItem) *openapi3.Responses {
	switch method {
	case "GET":
		return pathItem.Get.Responses
	case "POST":
		return pathItem.Post.Responses
	case "PUT":
		return pathItem.Put.Responses
	case "DELETE":
		return pathItem.Delete.Responses
	default:
		return pathItem.Options.Responses
	}
}

func StartOpenAPIFunction(path string, port string, r *http.Request, w http.ResponseWriter) {

	ctxMap := r.Context().Value("map").(map[string]interface{})
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(path)
	if len(doc.Servers) > 0 {
		doc.Servers[0].URL = ""

	}
	if err != nil {
		fmt.Printf("Failed to load OpenAPI document: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to load OpenAPI document"))
		return

	}
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		fmt.Printf("Failed to create route: %v", err)
	}
	mux := NewEnhancedMux(path)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status, contentType, exampleKey := GetOpenAPIValueFromRequest(r, doc)
		fmt.Println("status", status, "contentType", contentType, "exampleKey", exampleKey)

		route, pathParams, err := router.FindRoute(r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Not found: %v", err)
			return
		}

		ctx := r.Context()
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid request: %v", err)
			return
		}
		examples := MethodResponse(route.Method, route.PathItem).Status(status).Value.Content[contentType].Examples
		var response interface{}
		for key, example := range examples {
			if response == nil && exampleKey == "" {
				response = example.Value.Value
			}
			if key == exampleKey {
				response = example.Value.Value
			}
		}
		dirPath := filepath.Dir(path)
		if refMap, ok := response.(map[string]interface{}); ok {

			if ref, ok := refMap["$ref"]; ok {
				refPath := filepath.Join(dirPath, ref.(string))
				//read file from the path to json
				file, err := os.ReadFile(refPath)
				if err != nil {
					fmt.Println("Failed to read file: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to read file: %v", err)
					return
				}
				err = json.Unmarshal(file, &response)
				if err != nil {
					fmt.Println("Failed to unmarshal response: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to unmarshal response: %v", err)
					return

				}
			} else {
				fmt.Println("'$ref' not found")
			}
		} else {
			fmt.Println("response not the expect type")
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to marshal response: %v", err)
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	ctxMap["openapi#server#"+path+"#"+port] = srv
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

}

func StopOpenAPIFunction(path string, port string, r *http.Request, w http.ResponseWriter) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	key := "openapi#server#" + path + "#" + port
	srv, ok := ctxMap[key].(*http.Server)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No server found"))
		return
	}
	if err := srv.Shutdown(r.Context()); err != nil {
		log.Fatalf("Shutdown(): %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to shutdown server"))
		return
	}
	delete(ctxMap, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server shutdown successfully"))
}

type FileTree map[string]interface{}

func GetFileTreeHandler(w http.ResponseWriter, r *http.Request) {
	directoryPath := "/tmp/kubectl-go-upload/"
	fileTree := GetFileTree(directoryPath)

	jsonResponse, err := json.Marshal(fileTree)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetFileTree(directoryPath string) FileTree {
	fileTree := make(FileTree)

	files, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Println("Error reading directory:", err)
		return fileTree
	}

	for _, file := range files {
		filePath := filepath.Join(directoryPath, file.Name())

		if file.IsDir() {
			fileTree[file.Name()] = GetFileTree(filePath)
		} else {
			if strings.Split(file.Name(), ".")[1] == "yaml" ||
				strings.Split(file.Name(), ".")[1] == "yml" ||
				strings.Split(file.Name(), ".")[1] == "json" {
				fileTree[file.Name()] = true
			}
		}
	}

	return fileTree
}

type OpenAPIListenItem struct {
	Path string `json:"path"`
	Port string `json:"port"`
}

func GetCurrentOpenAPIListPortAndFileName(r *http.Request) []OpenAPIListenItem {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	var openAPIList []OpenAPIListenItem
	for key := range ctxMap {
		if strings.Contains(key, "openapi#server#") {
			openAPIList = append(openAPIList, OpenAPIListenItem{
				Path: strings.Replace(strings.Split(key, "#")[2], "/tmp/kubectl-go-upload", "", -1),
				Port: strings.Split(key, "#")[3],
			})
		}
	}
	return openAPIList
}

type GlobalVarsKey struct{}

// 全局变量结构
type GlobalVars struct {
	// 在此处添加你想要的任何全局变量
}

func GetOpenAPIListHandler(w http.ResponseWriter, r *http.Request) {
	openAPIList := GetCurrentOpenAPIListPortAndFileName(r)

	jsonResponse, err := json.Marshal(openAPIList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
