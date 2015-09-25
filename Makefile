# simple makefile to test and build

# namespace
PRJ = xr
# app name
APP = slackbot
# binary name
BIN = $(APP)
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

# get information for tools
# and images
info:
	make -v
	sudo docker version --format 'Client: {{ .Client.Version}} Server: {{ .Server.Version }}'
	godep version
	go version
	git describe --tags
	echo "namespace:"$(PRJ) "appname:"$(APP) "binary-name:"$(BIN) "version:"$(FORMAT)

# rm test files binary and out files
clean: docker-clean
	rm -f *.test *.out $(BIN)

# to reduce space usage by docker
# remove not running containers
docker-clean:
	docker rm $(docker ps -a -q)
# docker images also require space
# remove old untagged images
	docker rmi $(docker images -f "dangling=true" -q)

# cheking code style, try to stick to google code review style
lint:
	golint ./... | grep -v "be unexported"

# run unit tests with coverage
unit-test:
	godep go test -v --cover ./...

# set binary name and build version into it
build:
	CGO_ENABLED=0 godep go build -v -o $(BIN) -ldflags "-X main.BuildVersion=$(FORMAT)"

# build in golang container
# for gitlab ci
build-in-docker:
	sudo docker run --rm -v "$(GOPATH)/bin":/go/bin -v "$(PWD)":/go/src/$(PRJ)/$(APP) -w /go/src/$(PRJ)/$(APP) -e CGO_ENABLED=0 golang godep go build -v -o $(BIN)

# there should not be artifact left
# to be sure that required version binary packed
build-image:
	sudo docker build -t slackbot:$(DOCTAG) .

# force tag local image to tutum remote
# and pushes to tutum registry
deploy-tutum:
	sudo docker tag -f slackbot:$(DOCTAG) tutum.co/mmirolim/slackbot:$(DOCTAG)
	sudo docker push tutum.co/mmirolim/slackbot:$(DOCTAG)

