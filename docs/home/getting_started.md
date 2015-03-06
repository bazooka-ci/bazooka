# Getting Started

## Bazooka Overview

Bazooka is a Continuous Integration and Continuous Deployment Server that is highly modular. You can choose between many SCM to fetch the source code of your application, or create your own plugin for your own SCM. You can use many of our built-in build environments (Java, Go, Node...), create one of your own, or use any container available on [docker Hub](https://hub.docker.com/) as a build environment for your build process

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

As described in our philosophy, the configuration of your bazooka build is versioned along your code, in a `.bazooka.yml` file. Create this file and commit it in your repository. The content of this file is described in details in the section [Configure your build](../home/build_configuration)

## Step 4: Trigger your first build

Once you SCM hook is set up, push your commit that adds .bazooka.yml to your repository. Bazooka will then start a new build based on your build configuration contained in your `.bazooka.yml` each time new commits are available on your SCM
