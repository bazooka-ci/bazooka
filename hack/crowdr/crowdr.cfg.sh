#!/bin/bash

if [[ $DOCKER_HOST =~ tcp://(.*):.* ]]
then
    DOCKER0=${BASH_REMATCH[1]}
else
    DOCKER0=$(ip -4 a show docker0 | grep -oP "(?<=inet ).*(?=/)")
fi


NET=bzk_net

if docker network inspect ${NET} > /dev/null 2>&1
then
    true
else
    docker network create -d=bridge ${NET} > /dev/null 2>&1
fi

cat <<EOF
global project bzk

#
# DB container
#
db image mongo:3.1
db net ${NET}

#
# Queue container
#
queue image jawher/beanstalkd
queue net ${NET}

#
# server container
#
server image bazooka/server
server net ${NET}
server publish 3000:3000
server publish 3001:3001
server env BZK_DB_URL=bzk_db:27017
server env BZK_QUEUE_URL=bzk_queue:11300

#
# worker container
#
worker image bazooka/worker
worker net ${NET}
worker volume ${BZK_HOME}:/bazooka
worker volume /var/run/docker.sock:/var/run/docker.sock
worker env BZK_DOCKERSOCK=/var/run/docker.sock
worker env BZK_API_URL=http://bzk_server:3000
worker env BZK_QUEUE_URL=bzk_queue:11300
worker env BZK_SYSLOG_URL=tcp://${DOCKER0}:3001
worker env BZK_NETWORK=bzk_net
worker env BZK_SCM_KEYFILE=${BZK_SCM_KEYFILE}
worker env BZK_HOME=${BZK_HOME}

#
# Web container
#
web image bazooka/web
web net ${NET}
web publish 8000:80
web env BZK_SERVER_HOST=bzk_server
EOF