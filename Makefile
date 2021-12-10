CSI_IMAGE_NAME?=harbor.cloud.netease.com/qzprod-k8s/k8scsi/curve-csi

# VERSION is the git tag
VERSION?=$(shell git describe --tags --match "v*")

GO_PROJECT=github.com/opencurve/curve-csi
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date '+%Y%m%d.%H%M%S.%Z')

# go build flags
LDFLAGS ?=
LDFLAGS += -X $(GO_PROJECT)/pkg/util.Version=$(VERSION)
LDFLAGS += -X $(GO_PROJECT)/pkg/util.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X $(GO_PROJECT)/pkg/util.BuildTime=$(BUILD_TIME)

# test args
TESTARGS_DEFAULT := "-v"
export TESTARGS ?= $(TESTARGS_DEFAULT)

.PHONY: all
all: build

.PHONY: test
test:
	go test -tags=unit $(shell go list ./...) $(TESTARGS)

.PHONY: build
build:
	if [ ! -d ./vendor ]; then (go mod tidy && go mod vendor); fi
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -ldflags "$(LDFLAGS) -extldflags '-static'"  -o _output/curve-csi ./cmd/curve-csi.go

.PHONY: release-image
release-image:
	docker build --network host -f ./build/curve-csi/Dockerfile \
		--build-arg VERSION="$(VERSION)" \
		-t $(CSI_IMAGE_NAME):$(VERSION) .

.PHONY: push-image
push-image:
	docker push $(CSI_IMAGE_NAME):$(VERSION)

.PHONY: clean
clean:
	go clean -r -x
	rm -f _output/curve-csi
	rm -f images/curve-csi
