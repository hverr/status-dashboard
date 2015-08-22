FROM golang:1.4.2-cross

MAINTAINER Henri Verroken

RUN apt-get update && apt-get install -y \
    npm \
    nodejs \
    nodejs-legacy \
    zip

RUN npm install -g gulp
