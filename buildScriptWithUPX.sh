go build -ldflags="-s -w" -o web-kubectl main.go

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o web-kubectl.exe main.go