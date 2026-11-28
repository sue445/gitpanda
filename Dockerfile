FROM golang:1.27-trixie AS build-env

ADD . /work
WORKDIR /work

RUN make

FROM gcr.io/distroless/base-debian13
COPY --from=build-env /work/bin/gitpanda /app/gitpanda
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/gitpanda"]
