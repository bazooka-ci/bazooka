FROM busybox:ubuntu-14.04

COPY main /bin/main

COPY template/bazooka_phase.sh /template/bazooka_phase.sh
COPY template/bazooka_run.sh /template/bazooka_run.sh
COPY template/Dockerfile /template/Dockerfile

ENTRYPOINT /bin/main
