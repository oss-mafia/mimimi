# Copyright 2019 The OSS Mafia team
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

BINARY := mimimi
BINARY_STATIC := $(BINARY)-static
BUILD_DIR := bin

# override to push to a different registry or tag the image differently
CONTAINER_REGISTRY ?= gcr.io/ignasi-permanent-ffxxqcd8
CONTAINER_TAG ?= v0.1

# Make sure we pick up any local overrides.
-include .makerc

# Enable go modules when building, even if the repo is copied into a GOPATH.
export GO111MODULE=on

##@ Build targets

build: $(BUILD_DIR)/$(BINARY)  ## Build the bot
$(BUILD_DIR)/$(BINARY):
	go build -v -o $@ github.com/oss-mafia/mimimi/cmd

test:  ## Run all unit tests
	go test `go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...`

build-static: $(BUILD_DIR)/$(BINARY_STATIC)  ## Build the statically linked Linux binary
$(BUILD_DIR)/$(BINARY_STATIC):
	CGO_ENABLED=0 GOOS=linux go build \
		-a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo \
		-o $@ github.com/oss-mafia/mimimi/cmd
	chmod +x $@

clean:  ## Clean all binary articats
	rm -f $(BUILD_DIR)/$(BINARY)

##@ Packaging and distribution

docker-build: $(BUILD_DIR)/$(BINARY_STATIC)  ## Build the Docker image
	docker build -t $(CONTAINER_REGISTRY)/$(BINARY):$(CONTAINER_TAG) -f Dockerfile .

docker-push:  ## Push the Docker image to the configured registry
	docker push $(CONTAINER_REGISTRY)/$(BINARY):$(CONTAINER_TAG)

##@ Code quality and integrity

LINTER := bin/golangci-lint
$(LINTER):
	wget -q -O- https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.13

lint: $(LINTER)  ## Run the linters on all projects
	bin/golangci-lint run --config golangci.yml

##@ Others

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
			/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
			/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

.PHONY: build build-static docker-build docker-push docker-run clean help
