parser is a component of the Bazooka project

It is an autonomous Docker container used to convert a bazooka configuration
file to one or more Dockerfile

# Contract

## Input environment variables

* BZK_HOME       : Home of bazooka on the host
* BZK_PROJECT_ID : Unique ID of the bazooka Project (Id of the repository)
* BZK_JOB_ID     : ID of the job

## Input folder (/bazooka)

The source code of the application. Must contains a configuration file, either
`.bazooka.yml` or `.travis.yml`

## Output folder (/bazooka-output)

A list of compiled Dockerfile. As many as needed builds. At least one, `Dockerfile0`

* Dockerfile0
* Dockerfile1
* Dockerfile2
* ...

# Run the container

```
docker run \
    -v $BZK_HOME/$BZK_PROJECT_ID/$BZK_JOB_ID/source/:/bazooka \
    -v $BZK_HOME/$BZK_PROJECT_ID/$BZK_JOB_ID/output/:/bazooka-output \
    bazooka/parser-golang
```

# Build your own version of bazooka-parser

the Docker image bazooka/parser is available on DockerHub.

However, if you want to build your own version

```
docker build -t bazooka/parser .
```
