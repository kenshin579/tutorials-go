FROM ubuntu:latest

RUN apt-get update
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install tzdata
RUN apt-get install -y apache2 && apt-get clean

EXPOSE 80
CMD apachectl -D FOREGROUND