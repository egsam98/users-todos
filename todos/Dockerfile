FROM golang:alpine as builder

RUN apk update && \
    apk add build-base && \
    apk add --no-cache git

WORKDIR /github.com/egsam98/users-todos/todos

ENV GO111MODULE=on

COPY go.mod ../
COPY go.sum ../
RUN go mod download

COPY pkg ../pkg
COPY todos .

RUN CGO_ENABLED=0 go build -o bin/todos *.go

FROM scratch
COPY --from=builder /github.com/egsam98/users-todos/todos/bin/todos .
ENTRYPOINT ["./todos"]