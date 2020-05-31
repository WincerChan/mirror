FROM golang:alpine

#RUN apk update
#RUN apk add --no-cache git

WORKDIR /mirror

ADD . /mirror

RUN CGO_ENABLED=0 go build -o /mirror/run

EXPOSE 3000

ENTRYPOINT ["/mirror/run"]