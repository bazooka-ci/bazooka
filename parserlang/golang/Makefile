default: docker

docker:
	$(PREFIX) docker build -t bazooka/parser-golang .

devimage: gox
	$(PREFIX) docker build -t bazooka/parser-golang dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
