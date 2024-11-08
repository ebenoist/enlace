.EXPORT_ALL_VARIABLES:

GOFLAGS=-mod=vendor
GOPROXY="off"

all: vet test install
test:
	@go test ./... -v -race

install:
	@go install -mod=vendor .

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

setup: install

update: install

run: install
	source .env && enlace
