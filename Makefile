VERSION = base/version.go

branch = $(shell git rev-parse --abbrev-ref HEAD)
revision = $(shell git rev-parse HEAD)
dirty = $(shell test -n "`git diff --shortstat 2> /dev/null | tail -n1`" && echo "*")
base = github.com/zero-os/0-fs
ldflags = '-w -s -X $(base).Branch=$(branch) -X $(base).Revision=$(revision) -X $(base).Dirty=$(dirty)'

build:
	cd cmd && go build -ldflags $(ldflags) -o ../g8ufs
capnp:
	capnp compile -I${GOPATH}/src/zombiezen.com/go/capnproto2/std -ogo:cap.np model.capnp
