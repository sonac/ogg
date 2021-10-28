FROM golang:1.16 AS dev

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0 go build -o ogg app/main.go

FROM alpine:3.14 AS production

COPY --from=dev /app/ogg /ogg

CMD ["./ogg"]