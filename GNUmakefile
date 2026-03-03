default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run unit tests
.PHONY: test
test:
	go test ./... -v $(TESTARGS) -timeout 30m

# Build the provider
.PHONY: build
build:
	go build -o terraform-provider-moneat

# Install the provider locally for development
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/moneat-io/moneat/0.0.1/$$(go env GOOS)_$$(go env GOARCH)
	mv terraform-provider-moneat ~/.terraform.d/plugins/registry.terraform.io/moneat-io/moneat/0.0.1/$$(go env GOOS)_$$(go env GOARCH)/

# Generate documentation
.PHONY: docs
docs:
	go generate ./...

# Lint
.PHONY: lint
lint:
	golangci-lint run ./...

# Format
.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: tidy
tidy:
	go mod tidy
