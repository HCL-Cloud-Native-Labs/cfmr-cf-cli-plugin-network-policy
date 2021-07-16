.PHONY: build

PLUGIN_NAME := $(shell basename "$(PWD)")
PLUGIN_BIN_PATH := ${GOBIN}

build = GOOS=$(1) GOARCH=$(2) go build -o ${PLUGIN_BIN_PATH}/$(PLUGIN_NAME)-$(1)-$(2)$(3) -ldflags="-X main.pluginVersion=v1.0.0"

build: clean release

release: 
	@echo "Compiling for every OS and Platform"
	$(call build,linux,amd64)
	$(call build,linux,386)
	$(call build,windows,amd64,.exe)
	$(call build,windows,386,.exe)
	$(call build,darwin,amd64)	
	
clean:
	rm -rf ${PLUGIN_BIN_PATH}/
	
lint:
	go vet ./...

test:
	go test ./...
