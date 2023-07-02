.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/create main/create.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/checkphone main/checkphone.go
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose --aws-profile cloudguruR
