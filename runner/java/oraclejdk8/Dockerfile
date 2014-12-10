FROM dockerfile/java:oracle-java8

RUN apt-get update
RUN apt-get install -y curl maven ant

ENV GRADLE_VERSION 2.2.1

RUN cd /usr/bin && \
wget https://services.gradle.org/distributions/gradle-$GRADLE_VERSION-all.zip && \
unzip gradle-$GRADLE_VERSION-all.zip && \
ln -s gradle-$GRADLE_VERSION gradle && \
rm gradle-$GRADLE_VERSION-all.zip

# Set Appropriate Environmental Variables
ENV GRADLE_HOME /usr/bin/gradle
ENV PATH $PATH:$GRADLE_HOME/bin
