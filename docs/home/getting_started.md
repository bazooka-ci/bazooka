# Getting Started

## Bazooka Overview

Bazooka is a highly modular Continuous Integration and Continuous Delivery Server.
Out of the box, Bazooka supports many SCMs (Git, Mercurial, ...) to fetch the source code of your application, and can easily be extended to support other SCMs.
Bazooka also comes with built-in support for many languages (Java, Go, Python, Node...), with the possiblity of supporting others by creating custom docker images, or easier still, by using any container available on [docker Hub](https://hub.docker.com/) as a build environment for your build process.

## Step 1: Register a new project

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

## Step 2: Setup SCM Hook

(Coming Soon)

## Step 3: Add the .bazooka.yml file to your repository

As described in our philosophy, the configuration of your bazooka build is versioned alongside your code.

Create a `.bazooka.yml` file and commit it in your repository.
The format of this file is described in detail in the section [Configure your build](../home/build_configuration)

## Step 4: Trigger your first build

A build can be triggered either manually or by an SCM hook.

To manually trigger a build, use the Bazooka cli:

```
bzk job start NAME master
```

When an SCM hook is set up, a build is automatically triggered whenver you push a new commit to the remote repository.
