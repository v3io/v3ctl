GOOS ?= linux
GOARCH ?= amd64
V3CTL_GIT_COMMIT = $(shell git rev-parse HEAD)
V3CTL_TAG ?= latest
V3CTL_SRC_PATH = /v3ctl
V3CTL_BIN_PATH = /v3ctl

GO_LINK_FLAGS_INJECT_VERSION := -s -w \
	-X github.com/v3io/version-go.gitCommit=$(V3CTL_GIT_COMMIT) \
	-X github.com/v3io/version-go.label=$(V3CTL_TAG) \
	-X github.com/v3io/version-go.os=$(GOOS) \
	-X github.com/v3io/version-go.arch=$(GOARCH)

V3CTL_BUILD_COMMAND ?= CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="${GO_LINK_FLAGS_INJECT_VERSION}" -o ${V3CTL_BIN_PATH}/v3ctl-$(V3CTL_TAG)-$(GOOS)-$(GOARCH) $(V3CTL_SRC_PATH)/cmd/v3ctl/main.go

.PHONY: lint
lint:
	docker run --rm \
		--volume ${shell pwd}:/go/src/github.com/v3io/v3ctl \
		golang:1.14 \
		bash /go/src/github.com/v3io/v3ctl/hack/lint.sh

	@echo Done.

.PHONY: get-dependencies
get-dependencies:
	go get ./...

.PHONY: v3ctl-bin
v3ctl-bin:
	$(V3CTL_BUILD_COMMAND)

.PHONY: v3ctl
v3ctl:
	docker run \
		--volume $(shell pwd):${V3CTL_SRC_PATH} \
		--volume $(shell pwd):$(V3CTL_BIN_PATH) \
		--workdir ${V3CTL_SRC_PATH} \
		--env GOOS=$(GOOS) \
		--env GOARCH=$(GOARCH) \
		--env V3CTL_TAG=$(V3CTL_TAG) \
		golang:1.14 \
		make v3ctl-bin
