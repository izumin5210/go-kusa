FROM alpine:3.4

MAINTAINER izumin5210 <masayuki@izumin.info>

ENV WORKDIR /app
RUN mkdir $WORKDIR
WORKDIR $WORKDIR

RUN apk add --update --virtual build-dependencies \
        tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && apk del build-dependencies \
    && apk add \
        ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

ENV SSL_CERT_FILE /etc/ssl/certs/ca-certificates.crt

ARG run_at="00\t22\t*\t*\t*"

RUN echo -e "$run_at\tcd $(pwd); ./kusa" >> /var/spool/cron/crontabs/root

COPY kusa .

CMD ["crond", "-l", "2", "-f"]
