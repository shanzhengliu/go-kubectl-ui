FROM golang:1.20.6-alpine3.18 AS BuildStage

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main .

FROM alpine:latest

WORKDIR /app

COPY --from=BuildStage app/ app/

COPY config /root/kube/.config

RUN apk --no-cache add curl

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl 

EXPOSE 8080

ENTRYPOINT [ "app/main" ]