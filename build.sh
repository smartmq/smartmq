#!/usr/bin/env bash

OUTDIR="$(pwd)/build"

rm -Rf $OUTDIR

env GOOS=linux go build -x -a -o $OUTDIR/linux/smq-mqtt ./cmd/smq-mqtt/main.go
env GOOS=linux go build -x -a -o $OUTDIR/linux/smq-rest ./cmd/smq-rest/main.go
env GOOS=linux go build -x -a -o $OUTDIR/linux/smq-cli-pub ./cmd/smq-cli-pub/main.go
env GOOS=linux go build -x -a -o $OUTDIR/linux/smq-cli-sub ./cmd/smq-cli-sub/main.go

env GOOS=windows go build -x -a -o $OUTDIR/windows/smp-mqtt ./cmd/smq-mqtt/main.go
env GOOS=windows go build -x -a -o $OUTDIR/windows/smq-rest ./cmd/smq-rest/main.go
env GOOS=windows go build -x -a -o $OUTDIR/windows/smq-cli-pub ./cmd/smq-cli-pub/main.go
env GOOS=windows go build -x -a -o $OUTDIR/windows/smq-cli-sub ./cmd/smq-cli-sub/main.go
