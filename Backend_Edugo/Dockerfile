FROM golang:1.23.0-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./main.go
RUN chmod +x main
RUN mkdir -p public/images public/pdfs && chmod -R 777 public
RUN mkdir -p temp/images temp/pdfs && chmod -R 777 temp
EXPOSE 8080
CMD ["./main"]