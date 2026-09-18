FROM golang:1.22-alpine AS build

WORKDIR /src

COPY zero-hunger-contracts ./zero-hunger-contracts
COPY zero-hunger-user-service ./zero-hunger-user-service

WORKDIR /src/zero-hunger-user-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/user-service \
    ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H app

USER app

COPY --from=build /out/user-service /usr/local/bin/user-service

EXPOSE 8081 50051

ENTRYPOINT ["user-service"]