FROM golang:1.17-alpine AS build_env

WORKDIR $GOPATH/src/github.com/ah8ad3/file_server

COPY . .

RUN apk add --no-cache bash git openssh

RUN go mod download && go mod vendor
RUN sh build/build.sh

FROM alpine
RUN apk update && apk add ca-certificates bash && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build_env /go/src/github.com/ah8ad3/file_server/dist/file_server .
