GO := go
GOFMT := gofmt
ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)

.PHONY: all

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

all: build

build: build_windows build_linux build_macos

build_windows:
	mkdir -p build
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/redis-transmission.exe github.com/QiNiuQVMSolutionTeam/Redis-Transmission

build_linux:
	mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/redis-transmission github.com/QiNiuQVMSolutionTeam/Redis-Transmission

build_macos:
	mkdir -p build
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/redis-transmission-mac github.com/QiNiuQVMSolutionTeam/Redis-Transmission

pack: build
	mkdir -p build/release/windows
	mkdir -p build/release/linux
	mkdir -p build/release/macosx
	mv build/redis-transmission.exe build/release/windows/
	mv build/redis-transmission build/release/linux/
	mv build/redis-transmission-mac build/release/macosx/redis-transmission
	cd build/release &&	rm -f redis-transmission-release.tar.gz
	cd build/release &&	tar -czf redis-transmission-release.tar.gz windows linux macosx