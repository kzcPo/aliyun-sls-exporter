

tidy:
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aliyun-sls-exporter main.go

clear:
	@rm -rf aliyun-sls-exporter