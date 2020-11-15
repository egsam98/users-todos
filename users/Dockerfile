FROM golang:alpine as builder

RUN apk update && \
    apk add build-base && \
    apk add --no-cache git

WORKDIR /github.com/egsam98/users-todos/users

ENV GO111MODULE=on

COPY go.mod ../
COPY go.sum ../
RUN go mod download

COPY pkg ../pkg
COPY users .

RUN CGO_ENABLED=0 go build -o bin/users *.go

FROM scratch
COPY --from=builder /github.com/egsam98/users-todos/users/bin/users .
ENTRYPOINT ["./users"]