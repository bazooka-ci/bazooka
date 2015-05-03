# Bazooka Architecture

## Components

Each component of Bazooka is encapsulated in a Docker container, except the CLI.

The command `bzk service start` of the CLI will start Bazooka.

This is simply done by starting three Docker containers, linking them togethers when need, and exposing two ports.

The three containers are:

* **mongodb**: Persistant data in bazooka is stored in Mongo. This container does not expose any port on the host machine.

* **bazooka/server**: server exposes the Bazooka API, which is used by the Web interface, Bazooka CLI, and any other client that may want to do actions in Bazooka
    * Exposes port 3000
    * is linked with mongodb

* **bazooka/web**: web exposes the Web Interface of Bazooka
    * Exposes port 8080
    * is linked with bazooka/server

![Bazooka architecture](./assets/img/bzk_archi.png)

To show you we tell the truth:

```sh
> docker ps
CONTAINER ID        IMAGE                   PORTS                           NAMES
6c3be3d47ee4        bazooka/web:latest      443/tcp, 0.0.0.0:8000->80/tcp   bzk_web
ac3201746425        bazooka/server:latest   0.0.0.0:3000->3000/tcp          bzk_server
d207368ae0f3        mongo:latest            27017/tcp                       bzk_mongodb
```

## Build Pipeline

There are several actions you can do with Bazooka. One of the most important is to start a job.

Starting a job will trigger a build pipeline, which consists of many non-persistant Docker containers, that will be called to do some actions, which will result in a project being built.

Here is a schema of what happens when a Bazooka job is triggered

![Bazooka architecture](./assets/img/bzk_build_pipeline.png)

* **Step 0**: The server will start a Docker container for the orchestration of the job, *bazooka/orchestration*.

* **Step 1**: The "first" real step of a build pipeline is to fetch the source code from the SCM repository (git, mercurial...).

* **Step 2**: When the code source of the project is available, Bazooka will parse the *.bazooka.yml* file at the root of the project to determine your build configuration.
If the `language` attribute is set in your *.bazooka.yml*, the parser will invoke a dedicated language parser.

* **Step2b**: The language parser will parse all attributes specific to the language (for instance the `jdk` attribute of a java project).
The couple parser/language parser will output a set of Dockerfiles, generated according to the *.bazooka.yml* configuration. Each Dockerfile consists of a **variant** of your build (for instance, you can have one variant per jdk version present in your *.bazooka.yml*)

* **Step3**: The orchestration container will build containers for all Dockerfiles generated by the parser

* **Step4**: The orchestration container will run the build containers. Each run will be a variant build. If the `docker run` fails, it means your build failed.