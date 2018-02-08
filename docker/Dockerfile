FROM alpine:3.6

CMD /server.linux
COPY res/ /res/
COPY server.linux /
RUN apk update && apk add ca-certificates
