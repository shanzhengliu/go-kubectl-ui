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
8: view serects in the namespace as plaintext (Helm template secret support)

## Docker

Docker image can also work by running   
`docker build . -t go-kubectl-ui`  
but before running, please confirm you have mount all of the resource you need like kube/config

## Docker Run
You can run the image easily with `docker-compose up -d` after changing the specific place holder in the file


## ENV variable List
| Parameter         | Default Value      | Optional | Description                                  |
|-------------------|--------------------|----------|----------------------------------------------|
| KUBE_NAMESPACE    | default            | True     | user default namespace. eg: default          |
| KUBE_CONFIG       |                    | True     | user default context. is current context in the config file     |
| KUBE_PORT         | 8080               | True     | application running port eg:8080.  if you are use docker image, please set it the same as the port you are foward  |
| KUBE_CONFIG_PATH  | /root/.kube/config | True     | use kube config path. eg: /root/.kube/config  |
| KUBE_DEFAULT_PATH | /root/.kube        | True     | this is the path for  `.kube` folder          |  



## Okta support
Now this project has created a some adapter code to generate the token, which is the same as kube oidc-login. it will automatilly load the config and support the feature. therefore, follow the instracution from kubectl login and set up the okta config. pkce support only now.


## Screenshot
Pod Function (shell to connect to the pod container, rolling logs, yaml file, rolling logs are support also)
![screenshot](./screenshot/pod.png)


Deployment Function
![screenshot](./screenshot/deployment.png)

Deployment Function (can forward service in K8S to localhost, local port range could be 7000-7100 if you start application via docker compose)
![screenshot](./screenshot/service.png)

Http Helper (light weight http request tool)
![screenshot](./screenshot/http%20helper.png)

Docker Shell (talk to the container running this web app, you can use "kubectl" , "helm" command. it means you don't need to install kubectl in your local machine)
![screenshot](./screenshot/docker%20shell.png)

Open API helper Allows you to test the API. you can upload `zip` or `yaml` file to the docker container, and select the file you are going to mock, then select the port which you have expose to local, (from 7001-7100). you can test the api define in the open api doc via http tools, like `Postman` `curl` for write the code for the request you need.
![screenshot](./screenshot/openapi%20helper.png)


