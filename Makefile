WAGON=CGO_ENABLED=0 go run ./cmd/wagon --log-level=debug

debug:
	$(WAGON) build --platform=linux/arm64 ./cmd/hello

wagon:
	$(WAGON) build --platform=linux/arm64,linux/amd64 ./cmd/hello

gen:
	go run ./internal/cmd/tool gen ./cmd/ship

vet:
	goimports -w .