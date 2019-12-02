FROM golang:1.13.2-alpine3.10 as builder

WORKDIR /project
COPY . /project
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags="-s -w" -o /graphql-coverage ./cmd/graphql-coverage/

FROM alpine:3.10.1
RUN apk update && apk add --no-cache tzdata ca-certificates
COPY --from=builder /graphql-coverage /
