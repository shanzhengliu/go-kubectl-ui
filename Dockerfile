##frontend build
FROM node:20-slim AS NodeBuild
WORKDIR /app
COPY ./new-ui/ /app
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build
##node build end


FROM golang:1.20.6-alpine3.18 AS BuildStage

WORKDIR /app
COPY . .
COPY --from=NodeBuild /app/dist/ /app/frontend-build
RUN apk --no-cache add upx
RUN go mod download
RUN go build -o /app/main .
RUN upx /app/main

FROM alpine:latest as ENVStage

RUN apk --no-cache add curl

RUN OS="$(uname | tr '[:upper:]' '[:lower:]')" && \
    ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/' -e 's/armv[0-9]*/arm/')" && \
    curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/${ARCH}/kubectl" && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl && \
    KUBELOGIN="kubelogin_${OS}_${ARCH}" && \
    echo "Downloading kubelogin ${KUBELOGIN}" && \
    curl -fsSLO "https://github.com/int128/kubelogin/releases/download/v1.28.0/${KUBELOGIN}.zip" && \
    unzip "${KUBELOGIN}.zip" && \
    mv "./kubelogin" "/usr/local/bin/kubectl-oidc_login"  && \ 
    curl -LO  "https://get.helm.sh/helm-v3.14.3-${OS}-${ARCH}.tar.gz" && \
    tar -zxvf "helm-v3.14.3-${OS}-${ARCH}.tar.gz" && \
    mv "${OS}-${ARCH}/helm" /usr/local/bin/helm

FROM alpine:latest

WORKDIR /app

COPY --from=BuildStage app/ app/

COPY config /root/kube/.config

COPY --from=ENVStage /usr/local/bin/kubectl  /usr/local/bin/kubectl

COPY --from=ENVStage /usr/local/bin/helm  /usr/local/bin/helm

COPY --from=ENVStage /usr/local/bin/kubectl-oidc_login  /usr/local/bin/kubectl-oidc_login

EXPOSE 8080

EXPOSE 8000

ENTRYPOINT [ "app/main" ]