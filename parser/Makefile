default: docker

commons:
	cd ../commons && go install ./...

docker: gox
	docker build -t bazooka/parser -f Dockerfile .

gox: commons
	gox -osarch="linux/amd64" -output="main"
