# Build stage
# --------------------------------
FROM golang:1.8-alpine AS builder

MAINTAINER izumin5210 <masayuki@izumin.info>

RUN apk --update --virtual build-deps add \
  build-base \
  git

WORKDIR /app

ADD . .
RUN make build

RUN apk del build-deps \
  && rm -rf /var/cache/apk/*

# Runtime stage
# --------------------------------
FROM alpine:3.6

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/kusa /usr/local/bin/kusa

ENTRYPOINT ["/usr/local/bin/kusa"]
