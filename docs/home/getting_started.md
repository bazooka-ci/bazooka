# Getting Started

## Bazooka Overview

Bazooka is a Continuous Integration and Continuous Deployment Server that is highly modular. You can choose between many SCM to fetch the source code of your application, or create your own plugin for your own SCM. You can use many of our built-in build environments (Java, Go, Node...), create one of your own, or use any container available on [docker Hub](https://hub.docker.com/) as a build environment for your build process

## Register a new project

Once Bazooka is up and running, you can register your project in Bazooka

### Register a new project with CLI
```
bzk project create NAME SCM_TYPE SCM_URI [SCM_KEY]
```

For instance, if you want to build bazooka on bazooka itself:
```
bzk project create bazooka git \
  git@github.com:bazooka-ci/bazooka.git ~/.ssh/id_github
```

The private SCM Key is optional. If none is provided, Bazooka will try to use the default SCM Key provided during [Bazooka installation](../home/installation)

### Register a new project with the Web Interface

(Coming Soon)
