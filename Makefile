build: test
	go build -o bin/ssh-aliases

test:
	go test -v -cover ./...

release: build
	bash ./package.sh

