.PHONY: all install build env clean

APP_NAME=venus-forest
OUT_ROOT=_build
BUILD_TIME=`date +%FT%T%z`
#APP_HOME=$(pwd)

PROJ_MAIN=./
PROJ_APP="github.com/venusforest2013/config/application"
PROJ_REV=`git rev-parse --short HEAD`

all: install

install: build
	cp -r conf $(OUT_ROOT)
	./overwrite.sh

build: env
	export GO111MODULE=on
	go build -v -x -tags=jsoniter -ldflags "-X $(PROJ_APP).Revision=$(PROJ_REV) -X $(PROJ_APP).BuildTime=$(BUILD_TIME)" \
		-mod=vendor -o $(OUT_ROOT)/bin/$(APP_NAME) $(PROJ_MAIN)

env:
	mkdir -p $(OUT_ROOT)/{bin,conf}
	# swag init -d ./web/controller -g query.go
	mkdir -p $(OUT_ROOT)/$(APP_NAME)
	cp go.mod $(OUT_ROOT)/$(APP_NAME)

clean:
	rm -rf ./_build

test:
	echo "forest"
	cp -r conf $(OUT_ROOT)

run:
	$(OUT_ROOT)/bin/$(APP_NAME) -c _build/conf/settings.conf
