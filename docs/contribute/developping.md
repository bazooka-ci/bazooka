## Building Bazooka for the first time

Some tools are required to develop on Bazooka

* Go ([Installation instructions](https://golang.org/doc/install))
* Docker ([Installation instructions](https://docs.docker.com/installation/))
* npm (Install via [nvm](https://github.com/creationix/nvm))
* gulp (Install via npm `npm install -g gulp`)

* Clone this project in `$GOPATH/src/github.com/bazooka-ci/`

```bash
git clone git@github.com:bazooka-ci/bazooka.git $GOPATH/src/github.com/bazooka-ci/bazooka
```

* Run `make setup devimages` to build everything necessary to run bazooka
* Run `bzk run --restart` to restart bazooka using the images you just built (make sure `$GOPATH/bin` is in your PATH)

## Already installed Bazooka once ?

* Run `make devimages` to build docker images for go projects
* You made any changes to `runner/` or `scm/` ? Run `make scm` or `make runner` to updates those images as well
* Run `bzk run` to start bazooka (only needed is bazooka was not already running or if you changed some code in the server)
