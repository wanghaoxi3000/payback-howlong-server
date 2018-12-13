FROM library/golang

# Godep for vendoring
RUN go get github.com/tools/godep

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR=${GOPATH}/src/howlong
ENV APP_DATA_DIR=/var/howlong
ENV DB_SQLITE_PATH=${APP_DATA_DIR}/howlong.db

RUN mkdir -p $APP_DIR && mkdir -p $APP_DATA_DIR

# Set the entrypoint
ENTRYPOINT (cd $APP_DIR && ./src/howlong)
ADD . $APP_DIR

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 godep go build -ldflags '-d -w -s'

EXPOSE 8080
