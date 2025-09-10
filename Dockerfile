FROM golang:1.25
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd/main.go .
COPY ./internal ./internal/.
COPY ./migrations ./migrations/.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main
EXPOSE 8080
CMD ["/app/main"]