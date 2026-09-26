FROM golang:1.26.2-alpine AS build

WORKDIR /src

COPY zero-hunger-contracts ./zero-hunger-contracts
COPY zero-hunger-claim-service ./zero-hunger-claim-service

WORKDIR /src/zero-hunger-claim-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/claim-service \
    ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H app

USER app

COPY --from=build /out/claim-service /usr/local/bin/claim-service

EXPOSE 8084

ENTRYPOINT ["claim-service"]