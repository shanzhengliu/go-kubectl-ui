package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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

func OpenapiHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	port := r.URL.Query().Get("port")
	if path == "" || port == "" {
		http.Error(w, "Path and port query parameters are required", http.StatusBadRequest)
		return
	}
	StartOpenAPIHandler(path, port, r)

	w.Write([]byte("OK"))
}

func StopOpenapiHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	port := r.URL.Query().Get("port")
	if path == "" || port == "" {
		http.Error(w, "Path and port query parameters are required", http.StatusBadRequest)
		return
	}
	StopOpenAPIHandler(path, port, r, w)
}

func GetOpenAPIValueFromRequest(r *http.Request, doc *openapi3.T) (int, string, string) {
	status := 200
	contentType := "application/json"
	example := ""

	if len(doc.Servers) > 0 && r.Header.Get("openapi-server") == "" {
		if strings.Contains(doc.Servers[0].URL, "://") {
			r.Host = strings.Split(doc.Servers[0].URL, "://")[1]
		} else {
			r.Host = doc.Servers[0].URL
		}

	}
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

func IsPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false // Port is likely taken
	}
	ln.Close()  // Close the listener and release the port
	return true // Port is available
}

func StartOpenAPIHandler(path string, port string, r *http.Request) {
	if IsPortAvailable(port) {
		fmt.Printf("Port %s is available.\n", port)

	} else {
		fmt.Printf("Port %s is not available.\n", port)
		return
	}
	ctxMap := r.Context().Value("map").(map[string]interface{})
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		fmt.Printf("Failed to load OpenAPI document: %v", err)
	}
	err = doc.Validate(loader.Context)
	if err != nil {
		fmt.Printf("Failed to Valid OpenAPI document: %v", err)
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		fmt.Printf("Failed to create route: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status, contentType, exampleKey := GetOpenAPIValueFromRequest(r, doc)
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
				response = example.Value
			}
			if key == exampleKey {
				response = example.Value
			}
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
		Handler: nil,
	}
	ctxMap["openapi#server#"+path+"#"+port] = srv
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

}

func StopOpenAPIHandler(path string, port string, r *http.Request, w http.ResponseWriter) {
	ctxMap := r.Context().Value("map").(map[string]interface{})
	srv := ctxMap["openapi#server#"+path+"#"+port].(*http.Server)
	if srv == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No server found"))
	}
	if err := srv.Shutdown(r.Context()); err != nil {
		log.Fatalf("Shutdown(): %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to shutdown server"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server shutdown successfully"))
}

func fileListInDirectory(dir string) ([]string, error) {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, nil
}
