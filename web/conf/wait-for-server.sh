#!/bin/sh

set -eu

/bin/sed -i "s/<bzk_server_placeholder>/${BZK_SERVER_HOST}/" /etc/nginx/conf.d/default.conf

exec nginx -g "daemon off;"
