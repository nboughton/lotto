FROM debian:latest

LABEL maintainer="Nick Boughton"

#RUN apt update && apt upgrade -y

WORKDIR /www/

ADD site.app results.db ./

ADD frontend/dist/spa/ public/

EXPOSE 8000

CMD ["/www/site.app"]