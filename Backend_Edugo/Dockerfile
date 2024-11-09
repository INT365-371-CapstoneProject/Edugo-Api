FROM golang:1.23.0-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./main.go
RUN chmod +x main
RUN mkdir -p public/images public/pdfs && chmod -R 777 public/images public/pdfs
EXPOSE 8080
CMD ["./main"]