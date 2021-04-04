FROM golang:1.16 AS build-env

ADD . /work
WORKDIR /work

RUN make

FROM debian:buster-slim
COPY --from=build-env /work/bin/gitpanda /app/gitpanda

ENTRYPOINT ["/app/gitpanda"]