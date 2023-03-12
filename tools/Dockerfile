FROM golang:1.20.2-alpine3.16

RUN apk --update add git bash &&\
    cd / && GO111MODULE=off go get github.com/githubnemo/CompileDaemon

ENV APP_DIR=/opt/service
WORKDIR $APP_DIR

COPY . $APP_DIR

ENTRYPOINT ["/opt/service/tools/run"]
