FROM golang:latest as  builder
RUN apt-get update && apt-get install -y nocache git ca-certificates && update-ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
COPY .env .env ./
RUN go mod download
COPY . .
RUN go build -o ./bin/ci-agent .




FROM golang:latest
RUN apt-get update \
    && apt-get install -y git \
    && apt-get install make



COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/bin /app
EXPOSE 8080
CMD ["./klovercloud-ci-agent"]