build:
	cd cmd && go build -tags=embed -o ../g8ufs
capnp:
	capnp compile -I${GOPATH}/src/zombiezen.com/go/capnproto2/std -ogo:cap.np model.capnp