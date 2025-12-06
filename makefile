.PHONY: build

build:
	rm build/*
	GOOS=windows GOARCH=amd64 go build -o build/wslgogit-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build -o build/wslgogit-arm64.exe main.go