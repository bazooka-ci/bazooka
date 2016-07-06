default: docker

docker: go
	docker build -t bazooka/orchestration -f Dockerfile .

go:
	GOOS=linux GOARCH=amd64 go build -o=main
