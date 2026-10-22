# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/app .

FROM alpine:3.21
# tzdata: configs sets time.Local to Asia/Jakarta; ca-certificates: outbound TLS.
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
USER app
COPY --from=build /bin/app /usr/local/bin/app
EXPOSE 8080
ENTRYPOINT ["app"]
CMD ["server"]
