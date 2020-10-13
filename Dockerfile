FROM quay.io/giantswarm/alpine:3.12

RUN apk add --update --no-cache ca-certificates \
    && update-ca-certificates \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk \
    && apk add glibc-2.28-r0.apk \
    && rm -rf /var/cache/apk/*

ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ADD ./build/bin/azure-imds-agent-app-linux-amd64 /usr/local/bin/azure-imds-agent-app

EXPOSE 8000
ENTRYPOINT ["/usr/local/bin/azure-imds-agent-app", "daemon"]
