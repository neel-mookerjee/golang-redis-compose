FROM alpine:3.9

RUN apk add --update wget ca-certificates unzip git bash curl openssh jq && \
    apk del --purge wget ca-certificates unzip

CMD ["/bin/bash"] 