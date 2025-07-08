FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

ENV PROJECT_ID="lv-playground-appdev"
ENV driver=postgres 
ENV instance_connection_name="lv-playground-appdev:asia-southeast2:dbinstance-testing"
ENV db_user="cloud-sql-user@lv-playground-appdev.iam"
ENV db_name="userdb"
ENV private=""

RUN go mod tidy

RUN go build -o binary

EXPOSE 8080

ENTRYPOINT ["/app/binary"]