# Bazooka, a new Generation Continuous Integration and Continuous Deployment server

Bazooka is an **Open Source Continuous Integration and Continuous Deployment Server** designed to let you install it wherever you want; On your local computer, on dedicated servers within your enterprise network, or on virtualized Cloud instances.

## The Philosophy

We believe your build configuration should reside along your code, and be versionned as well. Tools like [Travis](https://travis-ci.org/) or [CircleCI](https://circleci.com/) leaded the way.

But most of this tools are hosted services, and the "install yourself" alternatives are few and did not match what we think a Continuous Integration Tool like this should be.

This his How Bazooka was created.

## [Development Instructions](docs/home/developping.md)

## [Installation Instructions](docs/home/installation.md)

## Build the documentation

### Install [mkdocs](http://www.mkdocs.org/)

```
pip install mkdocs
```

### Serve the docs

```
mkdocs serve --dev-addr=0.0.0.0:8081
```
