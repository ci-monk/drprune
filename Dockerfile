FROM golang:1.17.8-alpine3.15 as builder
RUN apk add --no-cache alpine-sdk=1.0-r0
WORKDIR /build
COPY [ ".", "." ]
RUN make build

FROM alpine:3.16.1 as release
RUN apk --no-cache add \
      ca-certificates=20191127-r7 \
      bash=5.1.8-r0 \
      bash-completion=2.11-r4 && \
    echo "source <(drprune completion bash)" >> ~/.bashrc
COPY --from=builder [ "/build/bin/drprune", "/usr/local/bin/drprune" ]
RUN chmod +x /usr/local/bin/drprune

ENTRYPOINT [ "drprune" ]
CMD ["--help"]
