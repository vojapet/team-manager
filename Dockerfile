FROM golang
COPY server.go /go/src/team-manager/
WORKDIR /go/src/team-manager/
RUN go get
RUN go build -o server
CMD /go/src/team-manager/server