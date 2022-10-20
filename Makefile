.PHONY: build debug test

# go env
GOPROXY := "https://goproxy.cn,direct"
GOOS    := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH  := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
CC      := musl-gcc

GOENV := GO111MODULE=on
GOENV += GOPROXY=$(GOPROXY)
GOENV += CC=$(CC)
GOENV += GOOS=$(GOOS) GOARCH=$(GOARCH)

# go
GO := go

# output
OUTPUT := bin/pigeon

# build flags
LDFLAGS := -s -w
LDFLAGS += -extldflags "-static -fpic"

BUILD_FLAGS := -a
BUILD_FLAGS += -trimpath
BUILD_FLAGS += -ldflags '$(LDFLAGS)'
BUILD_FLAGS += $(EXTRA_FLAGS)

# debug flags
GCFLAGS := "all=-N -l"

DEBUG_FLAGS := -gcflags=$(GCFLAGS)

# test flags
TEST_FLAGS := -v
TEST_FLAGS += -p 3

# packages
PACKAGES := $(PWD)/cmd/pigeon/main.go

build:
	$(GOENV) $(GO) build -o $(OUTPUT) $(BUILD_FLAGS) $(PACKAGES)

debug:
	$(GOENV) $(GO) build -o $(OUTPUT) $(DEBUG_FLAGS) $(PACKAGES)

test:
	$(GO) test $(TEST_FLAGS) ./...
