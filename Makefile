.PHONY: build build-debug image image-all release release-all proto

TARGETOS_LINUX = linux
TARGETARCH_AMD64 = amd64
TARGETARCH_ARM64 = arm64
TARGETMOD = debug

VERSION_TAG = 1.0.0
MILESTONE_TAG = 04.2022
REGISTRY = rubinus

# auto-generated
COMMIT_ID := $(shell git rev-parse HEAD)
BUILD_TS := $(shell date +'%Y%m%d%H%M%S')
BRANCH_TAG := $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG := $(shell git describe --tags --exact-match 2> /dev/null \
                     || git symbolic-ref -q --short HEAD \
                     || git rev-parse --short HEAD)
VERSION := $(VERSION_TAG)-build-$(BRANCH_TAG)-$(BUILD_TS)
BUILD := $(VERSION_TAG)-$(MILESTONE_TAG)-build-$(BUILD_TS)-$(BRANCH_TAG)-$(COMMIT_ID)
PROJECT_ROOT := $(shell pwd -L)
OUTPUT := $(PROJECT_ROOT)/build
MAIN_SRC_FILE=cmd/main.go

default: build

build:
	@echo "... build normal mod binary ..."
	go build -p 8  -o $(OUTPUT)/origin $(MAIN_SRC_FILE)
	@echo "Compile Done !!!"

build-debug:
	@echo "... build debug mod binary ..."
	go build -p 8 -gcflags="all=-N -l"  -o $(OUTPUT)/origin-debug $(MAIN_SRC_FILE)
	@echo "Compile Done !!!"

image: linux-amd64	# make image只编译并制作amd64的镜像

image-debug: linux-amd64-debug	# make image-debug只编译并制作amd64的镜像

image-all: linux-amd64 linux-arm64

linux-amd64-debug:
	@echo "Build debug image of version $(VERSION)"
	@echo "Build origin debug for $(TARGETOS_LINUX) $(TARGETARCH_AMD64)"
	@GOOS=$(TARGETOS_LINUX) GOARCH=$(TARGETARCH_AMD64) go build -p 8 -gcflags="all=-N -l" -o $(OUTPUT)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_AMD64) $(MAIN_SRC_FILE)
	@echo "Build image $(REGISTRY)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_AMD64)-$(TARGETMOD):$(VERSION)"
	@docker build -f $(PROJECT_ROOT)/Dockerfile.debug --build-arg BUILD=$(BUILD) -t $(REGISTRY)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_AMD64)-$(TARGETMOD):$(VERSION) . >/dev/null
	@echo "Done"

linux-amd64:
	@echo "Build image of version $(VERSION)"
	@echo "Build origin for $(TARGETOS_LINUX) $(TARGETARCH_AMD64)"
	@GOOS=$(TARGETOS_LINUX) GOARCH=$(TARGETARCH_AMD64) CGO_ENABLED=0 go build -p 8 -o $(OUTPUT)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_AMD64) $(MAIN_SRC_FILE)
	@echo "Build image $(REGISTRY)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_AMD64):$(VERSION)"
	@docker build -f $(PROJECT_ROOT)/Dockerfile --build-arg BUILD=$(BUILD) -t $(REGISTRY)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_AMD64):$(VERSION) . >/dev/null
	@echo "Done"

linux-arm64:
	@echo "Build image of version $(VERSION)"
	@echo "Build origin for $(TARGETOS_LINUX) $(TARGETARCH_ARM64)"
	@GOOS=$(TARGETOS_LINUX) GOARCH=$(TARGETARCH_ARM64) CGO_ENABLED=0 go build -p 8 -o $(OUTPUT)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_ARM64) $(MAIN_SRC_FILE)
	@echo "Build image $(REGISTRY)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_ARM64):$(VERSION)"
	@docker build -f $(PROJECT_ROOT)/Dockerfile-arm64 --build-arg BUILD=$(BUILD) -t $(REGISTRY)/origin-$(TARGETOS_LINUX)-$(TARGETARCH_ARM64):$(VERSION) . >/dev/null
	@echo "Done"

proto:
	./compile_proto.sh
