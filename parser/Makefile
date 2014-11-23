default: docker

docker:
	docker build -t bazooka/parser .

devimage: gox
	docker build -t bazooka/parser dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"