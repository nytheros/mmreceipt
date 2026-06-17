PLUGIN_ID ?= readreceipt
BUNDLE_NAME ?= $(PLUGIN_ID).tar.gz
GOOS ?= linux
GOARCH ?= amd64

.PHONY: all dist server webapp clean

all: dist

server:
	mkdir -p server/dist
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -trimpath -o server/dist/plugin-linux-amd64 ./server

webapp:
	cd webapp && npm install && npm run build

dist: server webapp
	rm -rf dist bundle
	mkdir -p bundle/server/dist bundle/webapp/dist dist
	cp plugin.json bundle/plugin.json
	cp -r server/dist/* bundle/server/dist/
	cp -r webapp/dist/* bundle/webapp/dist/
	tar -C bundle -czf dist/$(BUNDLE_NAME) .

clean:
	rm -rf dist bundle server/dist webapp/dist
