FROM golang:alpine

WORKDIR /api

COPY . .  

RUN go build

EXPOSE 5000

CMD ["./devbook_api"]