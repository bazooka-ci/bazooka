## Building Bazooka for the first time

Some tools are required to develop on Bazooka

* Go ([Installation instructions](https://golang.org/doc/install))
* Docker ([Installation instructions](https://docs.docker.com/installation/))
* npm (Install via [nvm](https://github.com/creationix/nvm))
* gulp (Install via npm `npm install -g gulp`)

* Clone this project in `$GOPATH/src/github.com/haklop/`

```bash
git clone git@github.com:haklop/bazooka.git $GOPATH/src/github.com/haklop/
```

* Run `make setup scm runner devimages` to build everything necessary to run bazooka
* Run `make run` to start bazooka

## Already installed Bazooka once ?

* Run `make devimages` to build docker images for go projects
* You made any changes to `runner/` or `scm/` ? Run `make scm runner` to updates those images as well
* Run `make run` to start bazooka (only need is bazooka was not already running or if you changed some code in the server)
