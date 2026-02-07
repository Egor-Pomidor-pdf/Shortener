FROM golang:1.25.6-alpine AS build
WORKDIR /app

# Кэш зависимостей: сначала go.mod/go.sum
COPY go.mod go.sum ./
RUN go mod download

# Потом исходники
COPY . .

# Собираем статический бинарь
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/shortener ./cmd

# ---- runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /bin/shortener /bin/shortener

COPY --from=build /app/internal/migrations/db /app/internal/migrations/db

COPY --from=build /app/internal/migrations/clickhouse /app/internal/migrations/clickhouse



EXPOSE 8089

ENTRYPOINT ["/bin/shortener"]


