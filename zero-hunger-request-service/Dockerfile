FROM golang:1.26.2-alpine AS build

WORKDIR /src

COPY zero-hunger-contracts ./zero-hunger-contracts
COPY zero-hunger-request-service ./zero-hunger-request-service

WORKDIR /src/zero-hunger-request-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/request-service \
    ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H app

USER app

COPY --from=build /out/request-service /usr/local/bin/request-service

EXPOSE 8083 50053

ENTRYPOINT ["request-service"]