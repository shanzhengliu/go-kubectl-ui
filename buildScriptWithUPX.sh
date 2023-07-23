go build -ldflags="-s -w" -o main main.go

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o main.exe main.go