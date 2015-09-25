FROM scratch

ADD ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
ADD slackbot ./app

CMD ["./app"]
