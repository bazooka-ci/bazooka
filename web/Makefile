default: devimage

devimage: buildapp
	cp conf/* build/
	docker build --rm -t bazooka/web build/

image: devimage

buildapp:
	npm install
	gulp build
