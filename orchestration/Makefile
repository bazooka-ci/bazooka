default: docker

docker:
	$(PREFIX) docker build -t bazooka/orchestration .

devimage: gox
	$(PREFIX) docker build -t bazooka/orchestration dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
