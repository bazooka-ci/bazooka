FROM phusion/baseimage:0.9.9

RUN apt-get update
RUN apt-get install -y git

ADD run.sh /run.sh
RUN echo "    IdentityFile /bazooka-key" >> /etc/ssh/ssh_config
RUN echo "    StrictHostKeyChecking no" >> /etc/ssh/ssh_config

VOLUME ["/bazooka"]

CMD ["/run.sh"]
