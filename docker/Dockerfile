FROM alpine:3.6

CMD /server.linux -listen :8080
COPY res/ /res/
COPY server.linux /
RUN apk update && apk add ca-certificates sqlite
