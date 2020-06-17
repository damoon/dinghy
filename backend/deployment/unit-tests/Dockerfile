# WORKBENCH ###################################################
FROM golang:1.14.1

WORKDIR /go/src/gitlab.com/davedamoon/dinghy

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY notify/pkg/pb    /go/src/gitlab.com/davedamoon/dinghy/notify/pkg/pb
COPY backend/cmd backend/cmd
COPY backend/pkg backend/pkg

CMD go test ./backend/... && echo "done"
