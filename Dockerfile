# Futfut application

FROM golang
MAINTAINER kampkoder@gmail.com

# Document that the service listens on port 8080.
EXPOSE 8080

# Assuming workspace (GOPATH) configured at /go.
ADD . /go/src/github.com/clausthrane/futfut/
ADD config.json /go/bin/

# get dependencies and install. TODO use "godep".
RUN go get github.com/gorilla/mux
RUN go get github.com/spf13/viper
RUN go install github.com/clausthrane/futfut

ENTRYPOINT /go/bin/futfut
