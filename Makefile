ci: lint test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

bootstrap:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

.PHONY: bootstrap

#################################################
# Test and linting
#################################################

test:
	@CGO_ENABLED=0 go test -v ./...

lint:
	bin/golangci-lint run ./...

.PHONY: lint test
