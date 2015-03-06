# Bazooka

Bazooka is a **Continuous Integration and Continuous Deployment Server** designed to let you install it wherever you want; On your local computer, on dedicated servers within your enterprise network, or on virtualized Cloud instances.

## The Philosophy

We believe your build configuration should reside along your code, and be versioned as well. Tools like [Travis](https://travis-ci.org/) or [CircleCI](https://circleci.com/) leaded the way.

But most of this tools are hosted services, and the "install yourself" alternatives are few and did not match what we think a Continuous Integration Tool like this should be.

This his How Bazooka was created.

## Architecture & Design

Bazooka relies heavily on Docker and its ecosystem and has been designed around it, to take advantage of everything Docker is capable, and then focus on the CI/CD part.

In Bazooka, everything is a plugin, and each plugin is a Docker container. You can read more about our [plugin architecture here](docs/internals/plugin.md)

## [Development Instructions](docs/contribute/developping.md)

## [Installation Instructions](docs/home/installation.md)
