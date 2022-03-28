FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN apk add build-base

COPY ./ ./

RUN go build -o /pulseid

EXPOSE 8081

CMD [ "/pulseid" ]
