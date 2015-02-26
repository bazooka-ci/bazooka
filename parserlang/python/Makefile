default: devimage

devimage: gox
	docker build -t bazooka/parser-python -f Dockerfile.dev .

gox:
	gox -osarch="linux/amd64" -output="main"
