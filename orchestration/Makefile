default: docker

commons:
	cd ../commons && go install ./...

docker:
	docker build -t bazooka/orchestration -f Dockerfile ..

devimage: gox
	docker build -t bazooka/orchestration -f Dockerfile.dev .

gox: commons
	gox -osarch="linux/amd64" -output="main"
