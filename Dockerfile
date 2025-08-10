# Etapa de build
FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o zone4 ./cmd/server

# Etapa final
FROM debian:stable-slim
WORKDIR /app
COPY --from=build /app/zone4 /app/zone4
COPY .env /app/.env

# Expor portas alinhadas
EXPOSE 9090        # TCP (Canal/Jogo)
EXPOSE 9091/udp    # UDP
EXPOSE 6060        # pprof/metrics

CMD ["/app/zone4"]
