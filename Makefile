.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/create main/create.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/delete main/delete.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/get main/get.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/list main/list.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/update main/update.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose --aws-profile cloudguruR
