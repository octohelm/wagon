WAGON = CGO_ENABLED=0 go run ./cmd/wagon

DEBUG = 0
ifeq ($(DEBUG),1)
	WAGON := $(WAGON) --log-level=debug
endif

wagon.ship:
	$(WAGON) do go ship pushx

wagon.archive:
	$(WAGON) do go archive

wagon.debug:
	$(WAGON) do go build linux/arm64

install:
	CGO_ENABLED=0 go install ./cmd/wagon

test:
	go test ./pkg/...

gen:
	go run ./internal/cmd/tool gen ./cmd/wagon

lint:
	goimports -w -l ./pkg
	goimports -w -l ./cmd

up:
	cd ./.wagon/engine && nerdctl compose up