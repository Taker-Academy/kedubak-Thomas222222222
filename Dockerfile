FROM golang:1.22.2-bullseye

WORKDIR /app

COPY . .
RUN go mod download && go mod verify
RUN go build KeDuBak

EXPOSE 8080

CMD ["./KeDuBak"]
