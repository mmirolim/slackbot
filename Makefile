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

pack: info build
# build docker image from default Dockerfile and tag it
	docker build -t slack-bot:$(VER) .

build: info
# set binary name and build version into it
	godep go build -o $(BIN) -ldflags "-X main.BuildVersion=$(FORMAT)"
info:
	go version
# rm test files binary and out files
clean:
	rm -f *.test *.out bot
