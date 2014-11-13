orchestration is a component of the Bazooka project

It is an autonomous Docker container used to orchestrate a bazooka build.

Steps are :

* Fetch source Code from scm
* Parse configuration file and generate Dockerfiles
* Build generated Dockerfiles

Each step consists of only one action... running a Docker container

# Contract

## Input environment variables

* BZK_SCM           : type of System Content Management (eg. git, svn...)
* BZK_SCM_URL       : URL of the source repository
* BZK_SCM_REFERENCE : Reference on the repository (branch name, SHA1...)
* BZK_SCM_KEYFILE   : Private key file on the host to be used for SCM fetch
* BZK_HOME          : Home of bazooka on the host
* BZK_JOB_ID        : Unique ID of the bazooka JOB (Id of the repository)
* BZK_JOB_NUMBER    : Number of the job
* BZK_DOCKERSOCK    : Path of the Docker socket on the host (usually /var/run/docker.sock)

## Input folder (/bazooka)

The source code of the application. Must contains a configuration file, either
`.bazooka.yml` or `.travis.yml`

## Output folder (/bazooka-output)

None

# Run the container

```
docker run \
  -v $BZK_DOCKERSOCK:/var/run/docker.sock \
  -v $BZK_HOME:/bazooka \
  -e BZK_SCM=$BZK_SCM \
  -e BZK_SCM_URL=$BZK_SCM_URL \
  -e BZK_SCM_REFERENCE=$BZK_SCM_REFERENCE \
  -e BZK_SCM_KEYFILE=$BZK_SCM_KEYFILE \
  -e BZK_HOME=$BZK_HOME \
  -e BZK_JOB_ID=$BZK_JOB_ID \
  -e BZK_JOB_NUMBER=$BZK_JOB_NUMBER \
  -e BZK_DOCKERSOCK=$BZK_DOCKERSOCK \
  bazooka/orchestration
```

# Build your own version of bazooka-orchestration

the Docker image bazooka/orchestration is available on DockerHub.

However, if you want to build your own version

```
docker build -t bazooka/orchestration .
```
