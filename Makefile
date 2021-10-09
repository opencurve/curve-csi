.PHONY: all

CSI_IMAGE_NAME=$(if $(ENV_CSI_IMAGE_NAME),$(ENV_CSI_IMAGE_NAME),harbor.cloud.netease.com/qzprod-k8s/k8scsi/curve-csi)
CSI_IMAGE_VERSION=$(if $(ENV_CSI_IMAGE_VERSION),$(ENV_CSI_IMAGE_VERSION),csi-v1.1.0-rc2)

GO_PROJECT=github.com/opencurve/curve-csi
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date '+%Y%m%d.%H%M%S.%Z')
GO_VERSION=$(shell go version|sed 's/ /-/g')
TESTARGS_DEFAULT := "-v"
export TESTARGS ?= $(TESTARGS_DEFAULT)

# go build flags
LDFLAGS ?=
LDFLAGS += -X $(GO_PROJECT)/pkg/util.GitCommit=$(GIT_COMMIT)
# CSI_IMAGE_VERSION will be considered as the driver version
LDFLAGS += -X $(GO_PROJECT)/pkg/util.Version=$(CSI_IMAGE_VERSION)
LDFLAGS += -X $(GO_PROJECT)/pkg/util.BuildTime=$(BUILD_TIME)
LDFLAGS += -X $(GO_PROJECT)/pkg/util.GoVersion=$(GO_VERSION)

all: curve-csi

test:
	go test -tags=unit $(shell go list ./...) $(TESTARGS)

curve-csi:
	if [ ! -d ./vendor ]; then (go mod tidy && go mod vendor); fi
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -ldflags "$(LDFLAGS) -extldflags '-static'"  -o _output/curve-csi ./cmd/

image: curve-csi
	cp _output/curve-csi images/curve-csi
	docker build --network host -t $(CSI_IMAGE_NAME):$(CSI_IMAGE_VERSION) ./images

push-image: image
	docker push $(CSI_IMAGE_NAME):$(CSI_IMAGE_VERSION)

clean:
	go clean -r -x
	rm -f _output/curve-csi
	rm -f images/curve-csi
