VERSION = base/version.go

GOPATH := $(shell go env GOPATH)
branch = $(shell git rev-parse --abbrev-ref HEAD)
revision = $(shell git rev-parse HEAD)
dirty = $(shell test -n "`git diff --shortstat 2> /dev/null | tail -n1`" && echo "*")
base = github.com/threefoldtech/0-fs
ldflags = '-w -s -X $(base).Branch=$(branch) -X $(base).Revision=$(revision) -X $(base).Dirty=$(dirty)'

getdeps:
	@echo "Installing golint" && go install golang.org/x/lint/golint
	@echo "Installing gocyclo" && go install github.com/fzipp/gocyclo
	@echo "Installing misspell" && go install github.com/client9/misspell/cmd/misspell
	@echo "Installing ineffassign" && go install github.com/gordonklaus/ineffassign
	@echo "Installing statik" && go install github.com/rakyll/statik
	@echo "Installing open api generator" && sudo npm install @openapitools/openapi-generator-cli -g

verifiers: vet fmt lint cyclo spelling staticcheck

vet:
	@echo "Running $@"
	@go vet -atomic -bool -copylocks -nilfunc -printf -rangeloops -unreachable -unsafeptr -unusedresult $(shell go list ./...)

fmt:
	@echo "Running $@"
	@gofmt -d $(shell ls **/*.go)

lint:
	@echo "Running $@"
	golint -set_exit_status $(shell go list ./...)

ineffassign:
	@echo "Running $@"
	ineffassign .

cyclo:
	@echo "Running $@"
	gocyclo -over 100 .


spelling:
	misspell -i monitord -error $(shell ls **/*.go)

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck -- ./...

check: test


build:
	cd cmd && go build -ldflags $(ldflags) -o ../g8ufs

install:
	cd cmd && go install -ldflags $(ldflags)

capnp:
	capnp compile -I${GOPATH}/src/zombiezen.com/go/capnproto2/std -ogo:cap.np model.capnp

test: verifiers
	# we already ran vet separately, so safe to turn it off here
	@CGO_ENABLED=1 go test -v -vet=off ./...

test-race: verifiers
	@echo "Running unit tests with -race flag"
	# we already ran vet separately, so safe to turn it off here
	@CGO_ENABLED=1 go test -v -vet=off -race ./...