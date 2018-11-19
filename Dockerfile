FROM scratch
LABEL maintainer="latschesar.atanassov@gmx.at"

ADD ./certs/server.crt server.crt
ADD ./certs/server.key server.key

ADD ossrv ossrv
EXPOSE 8080
ENTRYPOINT ["/ossrv"]