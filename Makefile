# simple makefile to test and build
# binary name
BIN = bot
# set build version from nearest git tag
VER = `git describe --tags`
# set build time
TIMESTM = `date -u '+%Y-%m-%d_%I:%M:%S%p'`
# set commit short
COMMIT = `git rev-parse --short HEAD`
# format version signature
FORMAT = "v$(VER)-commit\#$(COMMIT)-$(TIMESTM)"

pack: info build
	docker build -t slack-bot:0.1 .

build: info
# set binary name and build version into it
	godep go build -o $(BIN) -ldflags "-X main.BuildVersion=$(FORMAT)"
info:
	go version
# rm test files binary and out files
clean:
	rm -f *.test *.out bot
