FROM golang:1.23.0-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./main.go
RUN chmod +x main
RUN chmod -R 777 public/images
RUN chmod -R 777 public/pdfs
EXPOSE 8080
CMD ["./main"]