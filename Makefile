.PHONY: build plugin

BINARY_NAME=loglinter
PLUGIN_NAME=$(BINARY_NAME).so
CMD_DIR=./cmd/$(BINARY_NAME)
PLUGIN_DIR=./plugin

build:
	go build -o $(BINARY_NAME) $(CMD_DIR)

plugin:
	go build -buildmode=plugin -o $(PLUGIN_NAME) $(PLUGIN_DIR)