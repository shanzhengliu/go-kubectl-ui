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
		ErrorHandlerFunction(http.StatusMethodNotAllowed, w, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Max upload size ~10MB
	if err != nil {
		ErrorHandlerFunction(http.StatusBadRequest, w, "Failed to parse form")
		return
	}
	err = UnzipAndSave(r, w, "/tmp/kubectl-go-upload")
	if err != nil {
		ErrorHandlerFunction(http.StatusBadRequest, w, "Failed to unzip and save file")
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
		ErrorHandlerFunction(http.StatusBadRequest, w, "Invalid request body")
		return
	}

	if requestBody.Path == "" || requestBody.Port == "" {
		http.Error(w, "Path and port are required in the request body", http.StatusBadRequest)
		return
	}
	if !IsPortAvailable(requestBody.Port) {
		ErrorHandlerFunction(http.StatusBadRequest, w, "Port is already in use")
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
		ErrorHandlerFunction(http.StatusBadRequest, w, "Invalid request body")
		return
	}

	if requestBody.Path == "" || requestBody.Port == "" {
		ErrorHandlerFunction(http.StatusBadRequest, w, "Path and port are required in the request body")
		return
	}
	StopOpenAPIFunction(requestBody.Path, requestBody.Port, r, w)
}

func StopAllOpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	openAPIList := GetCurrentOpenAPIListPortAndFileName(r)
	for _, openAPIItem := range openAPIList {

		ctxMap := r.Context().Value("map").(map[string]interface{})
		key := "openapi#server#" + "/tmp/kubectl-go-upload" + openAPIItem.Path + "#" + openAPIItem.Port
		srv, ok := ctxMap[key].(*http.Server)
		if !ok {
			ErrorHandlerFunction(http.StatusNotFound, w, "No server found")
			return
		}
		if err := srv.Shutdown(r.Context()); err != nil {
			ErrorHandlerFunction(http.StatusInternalServerError, w, "Failed to shutdown server")
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
		ErrorHandlerFunction(http.StatusInternalServerError, w, "Failed to load OpenAPI document")
		return

	}
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		fmt.Printf("Failed to create route: %v", err)
		ErrorHandlerFunction(http.StatusInternalServerError, w, "Failed to create route")
		return
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status, contentType, exampleKey := GetOpenAPIValueFromRequest(r, doc)

		route, pathParams, err := router.FindRoute(r)
		if err != nil {
			ErrorHandlerFunction(http.StatusNotFound, w, "Route not found")
			return
		}

		ctx := r.Context()
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			ErrorHandlerFunction(http.StatusBadRequest, w, fmt.Sprintf("Request validation failed: %v", err))
			return
		}
		responseContent := MethodResponse(route.Method, route.PathItem).Status(status).Value.Content[contentType]
		var response interface{}

		if responseContent.Example != nil {
			example := responseContent.Example
			response = exampleHandler(example, path, w)
		} else {
			examples := responseContent.Examples
			response = examplesHandler(examples, exampleKey, path, w)
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
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

func handleResponse(response interface{}, path string, w http.ResponseWriter) interface{} {
	if refMap, ok := response.(map[string]interface{}); ok {
		if ref, ok := refMap["$ref"]; ok {
			refPath := filepath.Join(filepath.Dir(path), ref.(string))
			if err := readAndUnmarshalFile(refPath, &response, w); err != nil {
				return response
			}
		}
	}

	if err := replaceRefs(response, filepath.Dir(path)); err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return response
	}

	return response
}

func readAndUnmarshalFile(path string, target interface{}, w http.ResponseWriter) error {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
		return err
	}

	if err := json.Unmarshal(file, target); err != nil {
		fmt.Printf("Failed to unmarshal file: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to unmarshal file: %v", err), http.StatusInternalServerError)
		return err
	}

	return nil
}

func exampleHandler(example interface{}, path string, w http.ResponseWriter) interface{} {
	return handleResponse(example, path, w)
}

func examplesHandler(examples openapi3.Examples, exampleKey string, path string, w http.ResponseWriter) interface{} {
	var response interface{}

	for key, example := range examples {
		if exampleKey == "" || key == exampleKey {
			response = example.Value.Value
			break
		}
	}

	return handleResponse(response, path, w)
}

func StopOpenAPIFunction(path string, port string, r *http.Request, w http.ResponseWriter) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	key := "openapi#server#" + path + "#" + port
	srv, ok := ctxMap[key].(*http.Server)
	if !ok {
		ErrorHandlerFunction(http.StatusNotFound, w, "No server found")
		return
	}
	if err := srv.Shutdown(r.Context()); err != nil {
		ErrorHandlerFunction(http.StatusInternalServerError, w, "Failed to shutdown server")
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
		ErrorHandlerFunction(http.StatusInternalServerError, w, "Failed to marshal file tree")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetFileTree(directoryPath string) FileTree {
	fileTree := make(FileTree)

	files, err := os.ReadDir(directoryPath)
	if err != nil {
		ErrorHandlerFunction(http.StatusInternalServerError, nil, "Failed to read directory")
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

func GetOpenAPIListHandler(w http.ResponseWriter, r *http.Request) {
	openAPIList := GetCurrentOpenAPIListPortAndFileName(r)

	jsonResponse, err := json.Marshal(openAPIList)
	if err != nil {
		ErrorHandlerFunction(http.StatusInternalServerError, w, "Failed to marshal openapi list")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
