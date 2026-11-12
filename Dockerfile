FROM golang:1.26.2-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY pkg ./pkg
RUN CGO_ENABLED=0 go build -o /scheduler .

FROM alpine:3.23

WORKDIR /app

COPY --from=build /scheduler ./scheduler
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

RUN mkdir -p /data

CMD ["./scheduler"]
