FROM alpine:3.6
RUN apk update && apk add ca-certificates sqlite
CMD /server.linux -listen :8080
COPY server.linux /
COPY res/ /res/
