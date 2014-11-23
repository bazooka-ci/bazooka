#!/bin/bash

#export BZK_SCM_KEYFILE=/home/ebellemon/.ssh/id_tf
#export BZK_HOME=/home/ebellemon/Documents/perso/bazooka-home
#export BZK_DOCKERSOCK=/var/run/docker.sock

if [ "$(uname)" != "Darwin" ]; then
  d=sudo
fi

$d docker run -d --name bzk_mongodb dockerfile/mongodb  2> /dev/null

rc=$?
# if container exists
if [[ $rc != 0 ]] ; then
	docker start bzk_mongodb 2> /dev/null
fi


$d docker inspect bzk_server &> /dev/null

rc=$?
if [[ $rc == 0 ]] ; then
    $d docker stop bzk_server &> /dev/null
    $d docker rm bzk_server &> /dev/null
fi

$d docker run -d \
  -v $BZK_DOCKERSOCK:/var/run/docker.sock \
  -v $BZK_HOME:/bazooka \
  -e BZK_SCM_KEYFILE=$BZK_SCM_KEYFILE \
  -e BZK_HOME=$BZK_HOME \
  -e BZK_DOCKERSOCK=$BZK_DOCKERSOCK \
  --link bzk_mongodb:mongo \
  -p 3000:3000 \
  --name="bzk_server" \
  bazooka/server
