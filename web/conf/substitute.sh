#!/bin/sh

set -eu
echo "injecting env vars"
/bin/sed -i "s/<bzk_server_placeholder>/${SERVER_PORT_3000_TCP_ADDR}:${SERVER_PORT_3000_TCP_PORT}/" /etc/nginx/conf.d/default.conf

nginx -g "daemon off;"
