# Copyright 2019 Iguazio
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
GOOS ?= linux
GOARCH ?= amd64
V3CTL_GIT_COMMIT = $(shell git rev-parse HEAD)
V3CTL_TAG ?= latest
V3CTL_SRC_PATH ?= /v3ctl
V3CTL_BIN_PATH ?= /v3ctl

GO_LINK_FLAGS_INJECT_VERSION := -s -w \
	-X github.com/v3io/version-go.gitCommit=$(V3CTL_GIT_COMMIT) \
	-X github.com/v3io/version-go.label=$(V3CTL_TAG) \
	-X github.com/v3io/version-go.os=$(GOOS) \
	-X github.com/v3io/version-go.arch=$(GOARCH)

V3CTL_BUILD_COMMAND ?= CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="${GO_LINK_FLAGS_INJECT_VERSION}" -o ${V3CTL_BIN_PATH}/v3ctl-$(V3CTL_TAG)-$(GOOS)-$(GOARCH) $(V3CTL_SRC_PATH)/cmd/v3ctl/main.go

.PHONY: lint
lint:
	./hack/lint/install.sh
	./hack/lint/run.sh

.PHONY: fmt
fmt:
	@go fmt $(shell go list ./... | grep -v /vendor/)

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
		gcr.io/iguazio/golang:1.19 \
		make v3ctl-bin
