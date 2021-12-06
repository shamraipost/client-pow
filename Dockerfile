FROM golang:1.17.1-alpine3.14

ENV SERVER_HOST "server-pow"
ENV SERVER_PORT 50005
ENV NAME "client"

WORKDIR /usr/src/myapp

CMD ["go", "run", "main.go"]
