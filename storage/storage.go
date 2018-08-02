package storage

import (
	"io"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("storage")
)

type Storage interface {
	Get(key string) (io.ReadCloser, error)
}
