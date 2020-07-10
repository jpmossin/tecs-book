
.PHONY: build clean deploy gomodgen

PATH = "./Scratch.vm"

build:
	go build -o bin/runvm runvm.go

run: build
	bin/runvm $(PATH)
