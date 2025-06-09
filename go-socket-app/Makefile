# Makefile for go-socket-app

APP_NAME=go-socket-app
SRC_DIR=./src
CLIENT_DIR=$(SRC_DIR)/client
SERVER_DIR=$(SRC_DIR)/server

.PHONY: all build-client build-server clean

all: build-client build-server

build-client:
	cd $(CLIENT_DIR) && go build -o ../../bin/client client.go

build-server:
	cd $(SERVER_DIR) && go build -o ../../bin/server server.go

clean:
	rm -rf ./bin/*
