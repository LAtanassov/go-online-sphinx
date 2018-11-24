FROM scratch
LABEL maintainer="latschesar.atanassov@gmx.at"

ADD ossvc ossvc
EXPOSE 8080
ENTRYPOINT ["/ossvc"]