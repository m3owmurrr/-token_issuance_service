# Сборка приложения
FROM golang:1.23.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o token_issuance_service ./cmd/token_issuance_service/

# Запуск приложения
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/token_issuance_service .
COPY ./config/config.yaml ./config/

ARG JWT_KEY
EXPOSE 1969
ENV JWT_KEY=${JWT_KEY}
CMD [ "./token_issuance_service" ]


# # Сборка приложения
# FROM golang:1.23.0-alpine AS builder
# WORKDIR /app
# COPY . .
# RUN go mod download
# RUN go build -o token_issuance_service ./cmd/token_issuance_service/
# ARG JWT_KEY
# EXPOSE 1969
# ENV JWT_KEY=${JWT_KEY}
# CMD [ "./token_issuance_service" ]
