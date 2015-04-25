# Installation instructions

## Pre-requisites

### Install Docker

Docker is an essential part of Bazooka. To install it, follow [these instructions](https://docs.docker.com/installation/)

## Download

Everything in Bazooka can be done through the CLI (Command Line Interface), including installing bazooka itself

Download Bazooka CLI with the following links and add it to your `PATH`

* [darwin_386](https://bintray.com/artifact/download/bazooka/bazooka/bzk_darwin_386/bzk)
* [darwin_amd64](https://bintray.com/artifact/download/bazooka/bazooka/bzk_darwin_amd64/bzk)
* [linux_386](https://bintray.com/artifact/download/bazooka/bazooka/bzk_linux_386/bzk)
* [linux_amd64 ](https://bintray.com/artifact/download/bazooka/bazooka/bzk_linux_amd64/bzk)
* [linux_arm](https://bintray.com/artifact/download/bazooka/bazooka/bzk_linux_arm/bzk)

## Installation

Installing bazooka is a one-liner

```bash
$ bzk service start
```

You will be prompted for minimal information need for bazooka to run

* Bazooka Home Folder: The path of a directory on your host where bazooka will work. It will contain workspaces of your build, artefacts...
* Docker Socket Path: The path to the docker socket on your local host, usually `/var/run/docker.sock`
* Bazooka Default SCM private key: The path to the private key bazooka will try to use by default when fetching SCM data, for instance with git

## Upgrade

### Upgrade bazooka to the latest version

If you already installed bazooka once and want to upgrade it to the latest version, it's as simple as

```bash
$ bzk service upgrade
```
