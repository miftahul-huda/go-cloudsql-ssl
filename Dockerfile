# ==== Dockerfile ====
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o user-app .

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/user-app .
COPY templates ./templates
COPY config.yaml .
COPY certs ./certs

EXPOSE 8080
CMD ["/app/user-app"]
