# Bazooka

Bazooka is a **Continuous Integration and Continuous Delivery Server** designed to let you install it wherever you want; On your local computer, on dedicated servers within your enterprise network, or on virtualized Cloud instances.

Bazooka is also technology agnostic: out of the box, it supports Go, Java, Python and Node, and can easily [be extended](internals/plugin_architecture.md) to support other languages.

## The Philosophy

We believe the build configuration should reside alongside your code, and be versioned as well. Tools like [Travis](https://travis-ci.org/) or [CircleCI](https://circleci.com/) led the way.

But most of these tools are hosted services, and the "install yourself" alternatives are few and did not match what we think a Continuous Integration Tool like this should be.

This his why Bazooka was created.

## Architecture & Design

Bazooka uses Docker as a runtime and extension mechanism.

In Bazooka, everything is a plugin, and every plugin is a Docker container. You can read more about our [plugin architecture here](internals/plugin_architecture.md)
