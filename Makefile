# simple makefile to test and build
# binary name
BIN = bot
# set build version from nearest git tag
VER =`git describe --tags --abbrev=0`
# set commit short
COMMIT =`git rev-parse --short HEAD`
# set build time
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
# format version signature
FORMAT = v$(VER)-$(COMMIT)-$(TIMESTM)

# @TODO change to ci, testing with local version
deploy: pack
	docker tag slackbot:$(VER) tutum.co/mmirolim/slackbot:$(VER)
	docker push tutum.co/mmirolim/slackbot:$(VER)

pack: info build
# build docker image from default Dockerfile and tag it
	sudo docker build -t slackbot:$(VER) .

build: info
# set binary name and build version into it
	CGO_ENABLED=0 godep go build -o $(BIN) -ldflags "-X main.BuildVersion=$(FORMAT)"
info:
	go version
# rm test files binary and out files
clean:
	rm -f *.test *.out bot
