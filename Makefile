WAGON = CGO_ENABLED=0 go run ./cmd/wagon

DEBUG = 0
ifeq ($(DEBUG),1)
	WAGON := $(WAGON) --log-level=debug
endif

export CGO_ENABLED = 0
export BUILDKIT_HOST =

wagon.debug:
	WAGON_GRAPH=1 $(WAGON) do go build linux/arm64

wagon.ship:
	$(WAGON) do go ship pushx

wagon.help:
	$(WAGON) do help

wagon.archive:
	$(WAGON) do --output=.wagon/build go archive

install:
	CGO_ENABLED=0 go install ./cmd/wagon

test:
	go test ./pkg/...

gen:
	go run ./internal/cmd/tool gen ./cmd/wagon

lint:
	goimports -w -l ./pkg
	goimports -w -l ./cmd

update:
	go get -u ./pkg/...

up:
	cd ./.wagon/engine && nerdctl compose up