default: devimage

devimage: gox
	docker build -t bazooka/parser-java -f Dockerfile.dev .

gox:
	gox -osarch="linux/amd64" -output="main"
