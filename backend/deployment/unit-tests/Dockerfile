# WORKBENCH ###################################################
FROM golang:1.14.1

WORKDIR /go/src/gitlab.com/davedamoon/dinghy

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD go test ./... && echo "done"
