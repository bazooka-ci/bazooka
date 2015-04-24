# Docker in docker

In bazooka, everything is done in containers. And containers themselves can orchestrate and start other containers if needed.

To do this, we need to do "Docker in Docker": to be able, within a Docker container, to build images, run other containers, call other "plugin" containers.

One way to do this with Docker is to build your container from a special image, like `jpetazzo/dind`, in *privileged* mode. This technique has some issues obviously.

We choose another way, and we heavily rely on this technique for the entire bazooka Architecture.

To run Docker containers within another Container, we mount as a volume the Docker socket from the host, within each container of Bazooka.

```bash
docker run -v /var/run/docker.socker:/var/run/docker.sock <my_image>
```
![](./assets/img/docker_in_docker.png)

One thing you should consider when building bazooka: Each container can have access to all the containers on the host, even bazooka server, web or the MongoDB database.

We plan this address this issue in the future. Follow [#182](https://github.com/bazooka-ci/bazooka/issues/182) if you want to know more
