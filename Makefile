.EXPORT_ALL_VARIABLES:

GOFLAGS=-mod=vendor
GOPROXY="off"

GOTEST = ${GOPATH}/bin/gotestsum
GOTEST_FLAGS = --format-icons --format pkgname

all: vet test install

${GOTEST}:
	$(shell go install gotest.tools/gotestsum@latest)

test: ${GOTEST}
	${GOTEST}

install:
	@go install -mod=vendor .

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

setup: install

update: install

run: install
	source .env && enlace
