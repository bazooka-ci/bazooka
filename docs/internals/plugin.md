# Bazooka Plugin Architecture



As described in the previous page, each step of the build will start a dedicated container, which has one task to do. For instance the SCM container has the duty to fetch the source code of the project.

Each Bazooka component is encapsulated in a Docker container, it is then really easy to swap between one container to another. For instance, if our source code is hosted in a git repository, we will use the *bazooka/scm-git* container for the SCM step. If it is hosted on a mercurial repository, we will use *bazooka/scm-hg*.

*bazooka/scm-git* and *bazooka/scm-hg* are described as **plugins** in Bazooka

Each step of a Bazooka build is pluggable. Potentially, even the *orchestration*, *server* or *web* containers could be pluggable.

## Definition of a plugin

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
