## Introduction  
This project is a really simple project for view the k8s resource status.

## Install
run `go build main.go` and the binary file will be built.

## Run
run `./main --namespace {your namespace} --config {your config use-context}`  
and access the `http://localhost:8080/`

## Docker  
Docker image can also work by running   
`docker build . -t go-kubectl-ui`  
but before running, please confirm you have mount all of the resource you need like kube/config