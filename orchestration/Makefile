default: docker

commons:
	cd ../commons && go install ./...

docker: gox
	docker build -t bazooka/orchestration -f Dockerfile .

gox: commons
	gox -osarch="linux/amd64" -output="main"
