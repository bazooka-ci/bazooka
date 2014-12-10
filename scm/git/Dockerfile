FROM debian:jessie

RUN apt-get update
RUN apt-get install -y git

ADD run.sh /run.sh
RUN echo "    IdentityFile /bazooka-key" >> /etc/ssh/ssh_config
RUN echo "    StrictHostKeyChecking no" >> /etc/ssh/ssh_config
RUN echo "    LogLevel DEBUG" >> /etc/ssh/ssh_config

VOLUME ["/bazooka"]

CMD ["/run.sh"]
