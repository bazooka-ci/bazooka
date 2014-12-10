# Development

* Clone this project in `$GOPATH/src/github.com/haklop/`
* Run `dev-setup.sh` script to initialize go projects
* Install npm (via [nvm](https://github.com/creationix/nvm))
* Install gulp via npm `npm install -g gulp`
* Run `make devimages` to build docker images
* Run `make run` to start bazooka

# Build the documentation

## Install [mkdocs](http://www.mkdocs.org/)

```
pip install mkdocs
```

## Serve the docs

```
mkdocs serve --dev-addr=0.0.0.0:8081
```
