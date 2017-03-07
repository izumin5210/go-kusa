FROM alpine:3.4

MAINTAINER izumin5210 <masayuki@izumin.info>

ENV ENTRYKIT_VERSION 0.4.0

WORKDIR /
ENV WORKDIR /app
RUN mkdir $WORKDIR
WORKDIR $WORKDIR

RUN apk add --update \
        ca-certificates \
        tzdata \
    && apk add --update --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ \
        entrykit \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

ENV SSL_CERT_FILE /etc/ssl/certs/ca-certificates.crt

COPY kusa .
COPY start.sh .

CMD ["./start.sh"]
