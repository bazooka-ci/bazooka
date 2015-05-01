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

* Run `make setup all` to build everything necessary to run bazooka
* Run `bzk service restart` to restart bazooka using the images you just built (make sure `$GOPATH/bin` is in your PATH)

## Already installed Bazooka once ?

* Run `make images` to build docker images for the go projects
* You made any changes to `runner/` or `scm/` ? Run `make scm` or `make runner` to updates those images as well
* Run `bzk service start` to start bazooka if it isn't already running, or `bzk service restart` to restart bazooka using the newly built images
