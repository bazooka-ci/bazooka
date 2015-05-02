default: docker

docker: gox
	docker build -t bazooka/orchestration -f Dockerfile .

gox:
	gox -osarch="linux/amd64" -output="main"
