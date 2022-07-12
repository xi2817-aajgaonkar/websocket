FROM golang:1.17.1-alpine3.14 AS builder
WORKDIR /build
COPY . .

#RUN apk --no-cache add build-base

RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
RUN chown 1001:1001 /
USER 1001
WORKDIR /app
COPY --from=builder /build/app .
CMD ["./app"]