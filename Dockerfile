FROM golang:1.17.3-alpine3.14 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

RUN go build -o /backend cmd/server/main.go

FROM alpine:3.14
COPY --from=build /backend /backend
EXPOSE 8080

CMD [ "/backend" ]