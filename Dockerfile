FROM scratch
LABEL maintainer="latschesar.atanassov@gmx.at"
ADD ossrv ossrv
EXPOSE 8080
ENTRYPOINT ["/ossrv"]