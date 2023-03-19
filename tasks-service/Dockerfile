FROM golang:1.20-alpine as builder

ADD . /build
WORKDIR /build

RUN go build -o tasks-api ./app

FROM golang:1.20-alpine

WORKDIR /srv

COPY --from=builder /build/tasks-api /srv/tasks-api

CMD ["/srv/tasks-api"]