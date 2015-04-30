# Bazooka Plugin Architecture

As described in the previous page, each step of the build will start a dedicated container, which has one task to do. For instance the SCM container has the duty to fetch the source code of the project.

Each Bazooka component is encapsulated in a Docker container, it is then really easy to swap between one container to another. For instance, if our source code is hosted in a git repository, we will use the *bazooka/scm-git* container for the SCM step. If it is hosted on a mercurial repository, we will use *bazooka/scm-hg*.

*bazooka/scm-git* and *bazooka/scm-hg* are described as **plugins** in Bazooka

Each step of a Bazooka build is pluggable. Potentially, even the *orchestration*, *server* or *web* containers could be pluggable.

## Definition of a Bazooka plugin

**A plugin is a Docker container!** nothing more...

There several types of plugins for each customizable step of Bazooka:

* SCM Plugin
* Language Parser Plugin

To ensure that any plugin, provided by bazooka, or not, will work with Bazooka, each one defines a **contract**.

Since a plugin is a Docker container, there are only a few interfaces we can use to define a contract with a container.

* Volumes
* Environment variables
* Exit code

Volumes can be of two types; Input volumes where Bazooka will put some files, and Ouput volumes, where Bazooka will expect the plugin to write some files

Each type a plugin container is run, Bazooka will mount volumes, pass environment variables to the container, and expect some result in output volumes.

## SCM Plugin

The contract of a SCM plugin is:

![SCM Plugin Contrat](./assets/img/scm_plugin.png)

* **Input Volumes**
    * **/bazooka-key** *(optional)* - The private key defined at the project level

* **Environnement variables**
    * **BZK_SCM_URL** - The URL of the SCM repository where is hosted the source code of the project
    * **BZK_SCM_REFERENCE** - The reference to checkout (Can be a branch, tag, commit ID...)

* **Output Volumes**
    * **/bazooka/** - The directory where the source code of the project should be
    * **/meta/** - Metadata needed by Bazooka to provide information to the user. Bazooka expects to find a `scm` file here, containing the following information in yaml format.

<!--- json displays like plain text. `no-highlight` doesn't work and if no language is set, it will autodetect applescript ??!!-->
```json
origin: <the remote url of the repository>
reference: <the reference that has been checkout>
commit_id: <id of the commit>
author:
  name: <author name>
  email: <author email>
committer:
  name: <committer name>
  email: <committer email>
date: <commit date>
message: |
<commit message>
```

## Language parser plugin

(TODO)


## Other pluggable endpoints in Bazooka

* Runners
* Services
