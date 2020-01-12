PKG=github.com/manifoldco/ansiwrap

LINTERS=\
	gofmt \
	golint \
	gosimple \
	vet \
	misspell \
	ineffassign \
	deadcode

ci: $(LINTERS) test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
	golang.org/x/lint/golint \
	honnef.co/go/tools/cmd/gosimple \
	github.com/client9/misspell/cmd/misspell \
	github.com/gordonklaus/ineffassign \
	github.com/tsenart/deadcode \
	github.com/alecthomas/gometalinter

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)

.PHONY: bootstrap $(BOOTSTRAP)

#################################################
# Test and linting
#################################################

test:
	@CGO_ENABLED=0 go test -v ./...

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
	 ./... --enable

$(LINTERS):
	$(METALINT) $@

.PHONY: $(LINTERS) test
