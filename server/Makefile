default: docker

docker:
	docker build -t bazooka/server .


devimage: gox
	docker build -t bazooka/server dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
