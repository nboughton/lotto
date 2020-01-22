FROM golang:latest AS go-build

RUN apt update && apt upgrade -y && apt install -y libsqlite3-dev ca-certificates

WORKDIR /go

RUN go get github.com/nboughton/lotto

RUN go build -o site.app github.com/nboughton/lotto 

FROM node:latest AS node-build

#RUN apt update && apt upgrade -y && apt install -y git

WORKDIR /build

RUN git clone https://github.com/nboughton/lotto .

WORKDIR /build/frontend

RUN npm install -g @quasar/cli && npm install && quasar b

FROM debian:latest

LABEL maintainer="Nick Boughton"

WORKDIR /www/

COPY --from=go-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=go-build /go/site.app /go/src/github.com/nboughton/lotto/results.db ./

COPY --from=node-build /build/frontend/dist/spa public/

EXPOSE 8000

CMD ["/www/site.app"]