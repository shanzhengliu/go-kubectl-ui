## Introduction  
Due to the complex set up process for dashboard and permission setting by the cluster owner, This project is a really simple project for view the k8s resource status. With this project, the new k8s learner can easily understand same main resource running in the cluster.

2 key concepts should know before using this application.  

1. `Context`: it means which cluster your are using. you can found the name in the file /{userHomeDir}/.kube/config 
2. `Namespace`: it refers to the space for you to deploy the application. For example, there are a lot of teams in an organization, they will deploy their application under their namespace to confirm they won't affect the environment used by other team.  

More details: https://kubernetes.io/docs/home/

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
| config   | minikube                | True     | user default context. eg: minikube           |
| port      | 8080                    | True     | application running port eg:8080.            |
| path      | {homeDir}/.kube/config | True     | use kube config path. eg: /root/.kube/config |

## Docker  
Docker image can also work by running   
`docker build . -t go-kubectl-ui`  
but before running, please confirm you have mount all of the resource you need like kube/config