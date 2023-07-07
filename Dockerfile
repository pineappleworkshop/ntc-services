FROM golang:1.20.2-alpine3.17 as build-env

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates openssl \
    && update-ca-certificates 2>/dev/null || true

RUN mkdir /services
WORKDIR /services
RUN apk add git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/services

FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /go/bin/services /go/bin/services
ENTRYPOINT ["/go/bin/services"]