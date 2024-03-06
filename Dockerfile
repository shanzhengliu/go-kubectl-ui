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

RUN apk add git

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl 

COPY  krew-install.sh /root/kube/krew-install.sh

RUN command chmod +x /root/kube/krew-install.sh && \
    /root/kube/krew-install.sh

ENV PATH="${PATH}:/root/.krew/bin"

RUN kubectl krew install oidc-login

EXPOSE 8080

EXPOSE 8000

ENTRYPOINT [ "app/main" ]