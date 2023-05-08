FROM golang:1.20-alpine3.17

RUN apk add curl

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /application
ADD . /application

RUN go mod vendor

EXPOSE 60013

ENTRYPOINT ["air"]