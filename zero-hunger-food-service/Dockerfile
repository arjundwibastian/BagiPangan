FROM golang:1.26.2-alpine AS build

WORKDIR /src

COPY zero-hunger-contracts ./zero-hunger-contracts
COPY zero-hunger-food-service ./zero-hunger-food-service

WORKDIR /src/zero-hunger-food-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/food-service \
    ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H app

USER app

COPY --from=build /out/food-service /usr/local/bin/food-service

EXPOSE 8082 50052

ENTRYPOINT ["food-service"]