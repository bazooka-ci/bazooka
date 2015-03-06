# Services

Bazooka allows you to use external services (databases, messaging systems...) within your build environment.

Services are simply Docker containers that are linked with the container in which your build runs. This allows us to provide as many services as there are containers on the Docker Hub. It is also possible to use containers from your local Docker registry

## Internals

To explain how services work, we will start with an example

```yaml
services:
  - mongo
```

Bazooka will parse the configuration file, register all the services declared, start the associated Docker containers, and link them at runtime with your build container. If this was done using the Docker CLI, it would give:

```bash
# Start all services
docker run --name services-mongo -d mongo

# Start the container which is actually building your application
docker run buildcontainer --link services-mongo
```

When linking containers together, several things happen.

* Environment variables are injected into the build container to enable programmatic discovery of the linked container (the bazooka service in our case)
*  A host entry for the linked container (the bazooka service) is add to the /etc/hosts file of the build container

To learn more about container linking and what happens, we recommend to read the [Docker documentation about container linking](https://docs.docker.com/userguide/dockerlinks/#container-linking)
