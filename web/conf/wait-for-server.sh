#!/bin/sh

set -eu

until grep $BZK_SERVER_HOST /etc/hosts
do
    echo "waiting for server to come up"
    sleep 1
done

/bin/sed -i "s/<bzk_server_placeholder>/${BZK_SERVER_HOST}/" /etc/nginx/conf.d/default.conf

exec nginx -g "daemon off;"
