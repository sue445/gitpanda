FROM golang:1.18 AS build-env

ADD . /work
WORKDIR /work

RUN make

FROM debian:buster-slim
COPY --from=build-env /work/bin/gitpanda /app/gitpanda
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/gitpanda"]
