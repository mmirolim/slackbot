# simple makefile to test and build
# binary name
BIN = bot
# get current branch
BR = `git name-rev --name-only HEAD`
# set build version from nearest git tag
VER = `git describe --tags --abbrev=0`
# set commit short
COMMIT =`git rev-parse --short HEAD`
# set build time
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
# format version signature
FORMAT = v$(VER)-$(COMMIT)-$(TIMESTM)
# docker tag version
DOCTAG = $(VER)-$(BR)
# @TODO change to ci, testing with local version
deploy-tutum: pack
	sudo docker tag slackbot:$(DOCTAG) tutum.co/mmirolim/slackbot:$(DOCTAG)
	sudo docker push tutum.co/mmirolim/slackbot:$(DOCTAG)

pack: info build
# build docker image from default Dockerfile and tag it
	sudo docker build -t slackbot:$(DOCTAG) .

build: info
# set binary name and build version into it
	CGO_ENABLED=0 godep go build -o $(BIN) -ldflags "-X main.BuildVersion=$(FORMAT)"
info:
	git describe --tags
	go version
# rm test files binary and out files
clean:
	rm -f *.test *.out bot
