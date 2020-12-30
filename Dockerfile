FROM alpine:latest


RUN apk --update --no-cache upgrade && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

COPY sshare /

WORKDIR /

ENTRYPOINT ["/sshare", "server"]
