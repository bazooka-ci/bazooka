default: docker

docker:
	docker build -t bazooka/orchestration .

devimage: gox
	docker build -t bazooka/orchestration dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
