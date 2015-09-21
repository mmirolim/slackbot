FROM fedora:latest

ADD bot ./app

CMD ["./app"]