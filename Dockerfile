FROM golang:1.25-bookworm AS build-env

ADD . /work
WORKDIR /work

RUN make

FROM debian:bookworm-slim
COPY --from=build-env /work/bin/gitpanda /app/gitpanda
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/gitpanda"]
