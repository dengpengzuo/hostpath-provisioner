.PHONY: all
.DEFAULT_GOAL := all

##### Global variables #####
OS := linux
ARCH := amd64
##### Public rules #####

all: image

build:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -gcflags "-N -l" -o docker/csi.hostfs.bin ./cmd/main.go

image: build
	docker build --tag ezcloud/hostfs:v1.0.0 --file docker/csi.hostfs/Dockerfile  ./docker/
