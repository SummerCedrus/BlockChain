GOPATH := $(shell pwd)
all:
	GOPATH=$(GOPATH) go install main

