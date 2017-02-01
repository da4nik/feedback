FROM alpine:3.5

RUN apk add --no-cache --update ca-certificates && \
    mkdir -p /etc/ssl/certs/ && update-ca-certificates --fresh \
    update-ca-certificates && \    
    rm -rf /var/cache/apk/* /tmp/*

RUN mkdir -p /app

WORKDIR /app
COPY feedback /app/

EXPOSE 9000

CMD ["./feedback"]
