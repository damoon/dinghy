FROM golang:1.17.3-buster@sha256:48f6285a02b7682373e97f153238f3a831268ca2e7dd6204bf540b6801ac63bc
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
