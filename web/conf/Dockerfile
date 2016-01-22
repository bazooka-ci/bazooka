FROM nginx

COPY web/ /usr/share/nginx/html/

COPY bzk-web.conf /etc/nginx/conf.d/default.conf
COPY wait-for-server.sh /
RUN chmod +x /wait-for-server.sh

ENTRYPOINT ["/wait-for-server.sh"]
