#!/bin/bash

#export BZK_SCM_KEYFILE=/home/ebellemon/.ssh/id_tf
#export BZK_HOME=/home/ebellemon/Documents/perso/bazooka-home
#export BZK_DOCKERSOCK=/var/run/docker.sock

if [ "$(uname)" != "Darwin" ]; then
  d=sudo
fi

$d docker run -d --name bzk_mongodb dockerfile/mongodb

$d docker run -d \
  -v $BZK_DOCKERSOCK:/var/run/docker.sock \
  -v $BZK_HOME:/bazooka \
  -e BZK_SCM_KEYFILE=$BZK_SCM_KEYFILE \
  -e BZK_HOME=$BZK_HOME \
  -e BZK_DOCKERSOCK=$BZK_DOCKERSOCK \
  --link bzk_mongodb:mongo \
  -p 3000:3000 \
  bazooka/server
