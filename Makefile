pack: build
	docker build -t slack-bot:0.1 .

build:
	go build -o bot main.go
