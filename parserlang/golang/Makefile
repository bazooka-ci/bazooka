default: docker

docker:
	docker build -t bazooka/parser-golang .

devimage: gox
	docker build -t bazooka/parser-golang dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
