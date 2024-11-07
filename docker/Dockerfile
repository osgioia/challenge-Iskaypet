FROM golang:1.23-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod tidy

RUN go build -o golangApp .

EXPOSE 8080

CMD ["./golangApp"]
