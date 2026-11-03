FROM golang:1.24-alpine AS build
ARG SERVICE
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /service ./$SERVICE/cmd
FROM alpine:3.21
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /service /app/service
COPY --from=build /src/api_gateway_service/config/config.yaml /app/api_gateway_service/config/config.yaml
COPY --from=build /src/reader_service/config/config.yaml /app/reader_service/config/config.yaml
COPY --from=build /src/writer_service/config/config.yaml /app/writer_service/config/config.yaml
USER app
ENTRYPOINT ["/app/service"]
