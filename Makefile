.PHONY: fmt vet clean dev setdev test

all: fmt vet test discordemotes discordemotes.exe

discordemotes.exe: *.go
	GOOS=windows GOARCH=amd64 go build -o discordemotes.exe

discordemotes: *.go
	GOOS=linux GOARCH=386 go build -o discordemotes

clean:
	-rm MovieNight.exe MovieNight

fmt:
	goimports -w .

vet:
	go vet ./...
	GOOS=js GOARCH=wasm go vet ./...

test:
	go test ./...
