FROM golang:1.22-alpine AS builder

ARG PROJECT_VERSION

COPY . /src/
WORKDIR /src
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=${PROJECT_VERSION}'" \
    -o app cmd/main.go
#RUN go test -cover -v ./...

FROM alpine:3.17.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/app .

ENTRYPOINT ["./app"]
