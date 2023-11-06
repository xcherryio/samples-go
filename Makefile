# get rid of default behaviors, they're just noise
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

default: help

BINS += xdb-samples
xdb-samples:
	@echo "compiling xdb-samples with OS: $(GOOS), ARCH: $(GOARCH)"
	@go build -o $@ cmd/server/main.go

.PHONY: bins

bins: $(BINS)