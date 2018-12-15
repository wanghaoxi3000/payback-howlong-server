FROM golang:1.11

# Recompile the standard library without CGO
#ENV CGO_ENABLED=0
#RUN go install -a std

ENV APP_DIR=/go/src/howlong
ENV APP_RUNMODE=prod

RUN mkdir -p ${APP_DIR}
COPY . $APP_DIR

RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR ${APP_DIR}
RUN dep init -v
RUN go build -ldflags '-w -s' -v


FROM alpine:3.8

RUN apk update && apk --no-cache add tzdata ca-certificates wget \
     && cp -rf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
     && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk \
     && apk add glibc-2.28-r0.apk && rm -f glibc-2.28-r0.apk /etc/apk/keys/sgerrand.rsa.pub

ENV APP_DIR=/opt/howlong
ENV APP_DATA_DIR=/var/howlong
ENV APP_DB_SQLITE_PATH=${APP_DATA_DIR}/howlong.db

COPY --from=0 /go/src/howlong/howlong ${APP_DIR}/howlong
RUN mkdir -p ${APP_DATA_DIR}

EXPOSE 8080
WORKDIR ${APP_DIR}
ENTRYPOINT ["./howlong"]
