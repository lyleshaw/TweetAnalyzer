ARG GO_VERSION=1.18

FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /builder
WORKDIR /builder

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./server ./service/index.go

FROM --platform=linux/amd64 alpine:latest

WORKDIR /
COPY --from=builder /builder/server .

EXPOSE 8085

CMD ["./server"]