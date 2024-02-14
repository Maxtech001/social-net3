FROM golang:alpine3.18

RUN mkdir -p /app
WORKDIR /app

COPY backend/ .

RUN apk add --no-cache gcc g++

RUN go build -o main .

EXPOSE 5000

CMD ["/app/main"]