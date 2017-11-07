build: test
	go build -o bin/ssh-aliases -ldflags "-s -w"

test:
	go test -cover ./...

release: build
	bash ./package.sh

