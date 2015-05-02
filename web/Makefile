default: docker

docker: buildapp
	cp conf/* build/
	docker build -t bazooka/web build/

buildapp:
	npm install
	gulp build
