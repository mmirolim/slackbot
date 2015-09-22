# simple makefile to test and build
#set build version
BV=printf '%s-commit#%s' `date -u '+%Y-%m-%d_%I:%M:%S%p'` `git rev-parse --short HEAD`

pack: info build
	docker build -t slack-bot:0.1 .

info:
	go version
build:
	godep go build -ldflags "-X main.BuildVersion="$(BV) -o slackbot
