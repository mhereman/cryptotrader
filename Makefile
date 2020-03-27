APP_NAME := cryptotrader
PACKAGE_PATH := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
APP_PATH := $(PACKAGE_PATH)/cmd/$(APP_NAME)

all: clean build

clean:
	go clean $(APP_PATH)

distclean:
	go clean -r -cache $(APP_PATH)

build:
	go build -o $(APP_PATH)/$(APP_NAME) -i $(APP_PATH)

rebuild: clean build

install:
	go install -i $(APP_PATH)

uninstall:
	go clean -i $(APP_PATH)
