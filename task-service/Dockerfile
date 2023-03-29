FROM golang:1.20-alpine as builder
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./app ./app

RUN go build -o tasks ./app

FROM scratch

WORKDIR /srv

COPY --from=builder /build/tasks /srv/tasks

CMD ["/srv/tasks"]