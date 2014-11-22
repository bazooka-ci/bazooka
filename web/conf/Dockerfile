FROM nginx

COPY web/ /usr/share/nginx/html/

COPY bzk-web.conf /etc/nginx/conf.d/default.conf
COPY substitute.sh /
RUN chmod +x /substitute.sh

ENTRYPOINT ["/substitute.sh"]
