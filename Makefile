build:
	go build -o bin/sqs-consumer

macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/sqs-consumer

clean:
	rm -rf bin
