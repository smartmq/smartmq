OUTDIR ?= "build"

default: build

build: build-darwin-amd64 build-linux build-win

build-darwin-amd64:
	env goos=darwin goarch=amd64 go build -o ${OUTDIR}/darwin-amd64/smq-mqtt ./cmd/smq-mqtt/main.go
	env goos=darwin goarch=amd64 go build -o ${OUTDIR}/darwin-amd64/smq-rest ./cmd/smq-rest/main.go
	env goos=darwin goarch=amd64 go build -o ${OUTDIR}/darwin-amd64/smq-cli-pub ./cmd/smq-cli-pub/main.go
	env goos=darwin goarch=amd64 go build -o ${OUTDIR}/darwin-amd64/smq-cli-sub ./cmd/smq-cli-sub/main.go

build-linux:
	env GOOS=linux go build -o ${OUTDIR}/linux/smq-mqtt ./cmd/smq-mqtt/main.go
	env GOOS=linux go build -o ${OUTDIR}/linux/smq-rest ./cmd/smq-rest/main.go
	env GOOS=linux go build -o ${OUTDIR}/linux/smq-cli-pub ./cmd/smq-cli-pub/main.go
	env GOOS=linux go build -o ${OUTDIR}/linux/smq-cli-sub ./cmd/smq-cli-sub/main.go

build-win:
	env GOOS=windows go build -o ${OUTDIR}/windows/smp-mqtt ./cmd/smq-mqtt/main.go
	env GOOS=windows go build -o ${OUTDIR}/windows/smq-rest ./cmd/smq-rest/main.go
	env GOOS=windows go build -o ${OUTDIR}/windows/smq-cli-pub ./cmd/smq-cli-pub/main.go
	env GOOS=windows go build -o ${OUTDIR}/windows/smq-cli-sub ./cmd/smq-cli-sub/main.go

clean:
	rm -Rf ${OUTDIR}

test:
	go test -v ./...

.PHONY: test clean build-win build-linux build default
