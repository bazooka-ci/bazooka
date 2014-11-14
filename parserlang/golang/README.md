parser-golang is a component of the Bazooka project

# Contract

## Input environment variables

None

## Input folder (/bazooka)

The source code of the application. Must contains a configuration file, either
`.bazooka.yml` or `.travis.yml`

## Output folder (/bazooka-output)

A list of compiled bazooka configuration files. As many as needed builds. At least one, `.bazooka.0.yml`

* .bazooka.0.yml
* .bazooka.1.yml
* .bazooka.2.yml
* ...

# Run the container

```
docker run -v $BZK_HOME/$BZK_JOB_ID/$BZK_JOB_NUMBER/source/:/bazooka -v $BZK_HOME/$BZK_JOB_ID/$BZK_JOB_NUMBER/output:/bazooka-output bazooka/parser-golang
```

# Build your own version of bazooka-parser-golang

the Docker image bazooka/parser-golang is available on DockerHub.

However, if you want to build your own version

```
docker build -t bazooka/parser-golang .
```
