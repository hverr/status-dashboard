FROM golang:1.5

MAINTAINER Henri Verroken

RUN apt-get update && apt-get install -y \
    git \
    npm \
    nodejs \
    nodejs-legacy \
    zip

RUN npm install -g gulp bower

EXPOSE 8050
