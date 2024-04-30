run: build
	@./bin/posturecheck

build:
	@GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/posturecheck .

