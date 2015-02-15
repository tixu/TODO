FROM golang

ADD . /go/src/todo
RUN go get github.com/gorilla/mux
RUN go get github.com/boltdb/bolt
RUN go install todo
ENTRYPOINT /go/bin/todo

EXPOSE 8080
