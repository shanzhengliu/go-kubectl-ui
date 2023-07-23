## Introduction  
This project is a really simple project for view the k8s resource status.

## Function  
1: Display Pod, Deployment, Configmap, Service, Ingress  
2: Context && Namespace switch via Web UI  
3: View Configmap Data  
4: View Pod Container Log  

## Install
run `go build -o web-kubectl main.go` and the binary file will be built.

## Run
run `./web-kubectl --namespace {your namespace} --config {your config use-context} --port {your port} --path {config path}`  
and access the `http://localhost:8080/`

## Paramter List
| Parameter | Default Value           | Optional | Description                                  |
|-----------|-------------------------|----------|----------------------------------------------|
| namespace | default                 | True     | user default namespace. eg: default          |
| context   | minikube                | True     | user default context. eg: minikube           |
| port      | 8080                    | True     | application running port eg:8080.            |
| path      | {homeDir}/.kube/config | True     | use kube config path. eg: /root/.kube/config |

## Docker  
Docker image can also work by running   
`docker build . -t go-kubectl-ui`  
but before running, please confirm you have mount all of the resource you need like kube/config