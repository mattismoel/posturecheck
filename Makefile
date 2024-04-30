run: build
	@./bin/posturecheck

build:
	@go build -o ./bin/posturecheck .
