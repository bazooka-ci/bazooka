# Installation instructions

## Pre-requisites

### Install Docker

Docker is an essential part of Bazooka. To install it, follow [these instructions](https://docs.docker.com/installation/)

## Download

Everything in Bazooka can be done through the CLI (Command Line Interface), including installing bazooka itself

Download Bazooka CLI with the following links, rename it simply to `bzk`add it to your `PATH`

* [darwin_386](https://bintray.com/artifact/download/bazooka/bazooka/bzk_darwin_386)
* [darwin_amd64](https://bintray.com/artifact/download/bazooka/bazooka/bzk_darwin_amd64)
* [linux_386](https://bintray.com/artifact/download/bazooka/bazooka/bzk_linux_386)
* [linux_amd64 ](https://bintray.com/artifact/download/bazooka/bazooka/bzk_linux_amd64)
* [linux_arm](https://bintray.com/artifact/download/bazooka/bazooka/bzk_linux_arm)


## Installation

Install bazooka is a oneliner

```bash
$ bzk run
```

## Upgrade

### Upgrade bazooka to the latest version

If you already installed bazooka once and want to upgrade it to the latest version, it's as simple as

```bash
$ bzk run --update
```
