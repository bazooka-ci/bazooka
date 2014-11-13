Bazooka-server is a component of the Bazooka project.

It is an autonomous Docker container used to expose the Bazooka API.

## API

### GET /project

Returns all the registered projects

#### Request:

    GET /project

#### Response:

    [
      {
        "scm_type":"git",
        "scm_uri":"git@bitbucket.org:bywan/bazooka-lang-example.git",
        "name":"Bazooka sample project",
        "id":"544bb11cc4c1b42765000001"
      },
      {
        "scm_type":"git",
        "scm_uri":"https://github.com/resthub/resthub-spring-stack.git",
        "name":"RESThub Spring stack",
        "id":"545167f7c4c1b423aa000001"
      },
    ]

### POST /project

Registers a new project

#### Request:

    POST /project

Body:

    {
      "scm_type":"git",
      "scm_uri":"https://github.com/resthub/resthub-spring-stack.git",
      "name":"RESThub Spring stack"
    }

#### Response:

    {
      "scm_type":"git",
      "scm_uri":"https://github.com/resthub/resthub-spring-stack.git",
      "name":"RESThub Spring stack",
      "id":"545167f7c4c1b423aa000001"
    }

### GET /project/{id}

Returns project details

#### Request:

    GET /project/{id}

#### Response:

    {
      "scm_type":"git",
      "scm_uri":"https://github.com/resthub/resthub-spring-stack.git",
      "name":"RESThub Spring stack",
      "id":"545167f7c4c1b423aa000001"
    }

### POST /project/{id}/job

Starts a new job

#### Request

    POST /project/{id}/job

Body:

    {
      "reference": "master"
    }

#### Response

TODO

## Contract

### Input environment variables

- BZK_SCM_KEYFILE: Private key file on the host to be used for SCM fetch
- BZK_HOME: Home of bazooka on the host
- BZK_DOCKERSOCK: Path of the Docker socket on the host (usually /var/run/docker.sock)

### Input folder (/bazooka)

Home of bazooka on the host. Builds are stored on this folder.

## Run the container

The container has to be linked to a Mongodb container. The port 3000 is exposing.

    docker run -d --name bzk_mongodb dockerfile/mongodb
    docker run \
        -v $BZK_DOCKERSOCK:/var/run/docker.sock \
        -v $BZK_HOME:/bazooka \
        -e BZK_SCM_KEYFILE=$BZK_SCM_KEYFILE \
        -e BZK_HOME=$BZK_HOME \
        -e BZK_DOCKERSOCK=$BZK_DOCKERSOCK \
        --link bzk_mongodb:mongo \
        -p 3000:3000 \
        bazooka/server

## Build your own version of bazooka-server

The Docker image bazooka/orchestration is available on DockerHub.

However, if you want to build your own version:

    docker build -t bazooka/server .
