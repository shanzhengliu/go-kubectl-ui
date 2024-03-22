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
5: Webshell for docker image you are running which contains "kubectl, Kube oidc-login, helm"
6: Webshell connection for container in the pod
7: OKTA adapter for UI, if you have set up okta in your config, you should set the callback url as http://localhost:8000/ in your okta admin account


## Install
run `go build -o web-kubectl main.go` and the binary file will be built.

## Run
run `./web-kubectl --namespace {your namespace} --config {your config use-context} --port {your port} --path {config path}`  
and access the `http://localhost:8080/`

## Paramter List
| Parameter | Default Value           | Optional | Description                                  |
|-----------|-------------------------|----------|----------------------------------------------|
| namespace | default                 | True     | user default namespace. eg: default          |
| config   |                          | True     | user default context. is current context in the config file     |
| port      | 8080                    | True     | application running port eg:8080.  if you are use docker image, please set it the same as the port you are foward  |
| path      | {homeDir}/.kube/config | True     | use kube config path. eg: /root/.kube/config  |
| kubeDefaultPath | /root/.kube      | True     | this is the path for  `.kube` folder          |  

## Docker  

Docker image can also work by running   
`docker build . -t go-kubectl-ui`  
but before running, please confirm you have mount all of the resource you need like kube/config


## Okta support
Now this project has create a some adapter code to generate the token, which is the same as kubelogin. it will automatilly load the config and support the feature. therefore, follow the instracution from kubectl login and set up the okta config. pkce support only now. 

## Image Run


