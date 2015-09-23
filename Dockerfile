FROM scratch

ADD ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
ADD bot ./app

CMD ["./app"]