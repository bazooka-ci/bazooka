# Development

## Starting from Scratch ?

We have some prerequisites for a development setup

* Install npm (via [nvm](https://github.com/creationix/nvm))
* Install gulp via npm `npm install -g gulp`

* Clone this project in `$GOPATH/src/github.com/haklop/`
* Run `make setup scm runner devimages` to build everything necessary to run bazooka
* Run `make run` to start bazooka

## You just updated some go code and want to try it out

* Run `make devimages` to build docker images for go projects
* Run `make run` to start bazooka (only need is bazooka was not already running or if you changed some code in the server)

# Build the documentation

## Install [mkdocs](http://www.mkdocs.org/)

```
pip install mkdocs
```

## Serve the docs

```
mkdocs serve --dev-addr=0.0.0.0:8081
```
