.PHONY: all
.DEFAULT_GOAL := all

##### Global variables #####
OS := linux
ARCH := amd64

PV_PROVISIONER := ez-cloud/hostpath-provisioner
VERSION := v1.0.0
##### Public rules #####

all: image

push:
	$(DOCKER) push "$(IMAGE):$(VERSION)"

image:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -gcflags "-N -l" -o docker/provisioner-amd64.bin ./cmd/main.go
	docker build --tag $(PV_PROVISIONER):$(VERSION) --file docker/provisioner/Dockerfile  ./docker/
