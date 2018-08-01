package rofs_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"
	"github.com/xxtea/xxtea-go/xxtea"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/rofs"
)

const (
	LedisAddr = "127.0.0.1:6789"
	ChunkSize = 512
	Decipher  = "test-decipher-key"
)

type TestStorage struct {
	data map[string][]byte
	hash []byte
}

func (t *TestStorage) Get(key string) (io.ReadCloser, error) {
	if data, ok := t.data[key]; ok {
		return ioutil.NopCloser(bytes.NewBuffer(data)), nil
	}
	return nil, fmt.Errorf("not found")
}

func MakeStorage(chunks int) (*TestStorage, []meta.BlockInfo) {
	s := TestStorage{
		data: make(map[string][]byte),
	}
	var blocks []meta.BlockInfo

	hash := md5.New()

	for i := 0; i < chunks; i++ {
		buf := make([]byte, ChunkSize)
		rand.Read(buf)
		hash.Write(buf)

		key := fmt.Sprintf("block-%d", i)

		block := meta.BlockInfo{
			Key:      []byte(key),
			Decipher: []byte(Decipher),
		}

		blocks = append(blocks, block)

		s.data[key] = xxtea.Encrypt(snappy.Encode(nil, buf), block.Decipher)
	}

	s.hash = hash.Sum(nil)
	return &s, blocks
}

func TestDownloadSuccess(t *testing.T) {
	//initialize test data
	storage, blocks := MakeStorage(20)

	downloader := rofs.Downloader{
		Storage:   storage,
		Blocks:    blocks,
		BlockSize: ChunkSize,
	}

	out, err := ioutil.TempFile("", "dt-")
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	defer func() {
		out.Close()
		os.RemoveAll(out.Name())
	}()

	err = downloader.Download(out)
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	hash := md5.New()
	out.Seek(0, 0) // rewind file

	io.Copy(hash, out)

	if ok := assert.Equal(t, storage.hash, hash.Sum(nil)); !ok {
		t.Error("wrong hash")
	}
}

func TestDownloadFailure(t *testing.T) {
	//initialize test data
	storage, blocks := MakeStorage(20)

	downloader := rofs.Downloader{
		Storage:   storage,
		Blocks:    blocks,
		BlockSize: ChunkSize,
	}

	//drop some blocks
	delete(storage.data, "block-1")
	delete(storage.data, "block-19")

	out, err := ioutil.TempFile("", "dt-")
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	defer func() {
		out.Close()
		os.RemoveAll(out.Name())
	}()

	err = downloader.Download(out)
	if ok := assert.Error(t, err); !ok {
		t.Fatal()
	}
}

func TestDownloadSingle(t *testing.T) {
	//initialize test data
	storage, blocks := MakeStorage(20)

	downloader := rofs.Downloader{
		Storage:   storage,
		Blocks:    blocks,
		BlockSize: ChunkSize,
		Workers:   1,
	}

	out, err := ioutil.TempFile("", "dt-")
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	defer func() {
		out.Close()
		os.RemoveAll(out.Name())
	}()

	err = downloader.Download(out)
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	hash := md5.New()
	out.Seek(0, 0) // rewind file

	io.Copy(hash, out)

	if ok := assert.Equal(t, storage.hash, hash.Sum(nil)); !ok {
		t.Error("wrong hash")
	}
}
