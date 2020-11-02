package rofs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/snappy"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
	"github.com/xxtea/xxtea-go/xxtea"
	"golang.org/x/crypto/blake2b"

	"golang.org/x/sync/errgroup"
)

const (
	// DefaultDownloadWorkers define the default number of workload to use to downloads data blocks
	DefaultDownloadWorkers = 4
	//DefaultBlockSize is the default block size
	DefaultBlockSize = 512 //KB
)

// Downloader allows to get some data blocks using a pool of workers
type Downloader struct {
	workers   int
	storage   storage.Storage
	blocks    []meta.BlockInfo
	blockSize uint64
}

// DownloaderOption interface
type DownloaderOption interface {
	apply(d *Downloader)
}

type workersOpt struct {
	workers uint
}

func (o workersOpt) apply(d *Downloader) {
	d.workers = int(o.workers)
}

// WithWorkers set number of workers
func WithWorkers(nr uint) DownloaderOption {
	return workersOpt{nr}
}

// NewDownloader creates a downloader for this meta from this storage
func NewDownloader(storage storage.Storage, m meta.Meta, opts ...DownloaderOption) *Downloader {
	downloader := &Downloader{
		storage:   storage,
		blockSize: m.Info().FileBlockSize,
		blocks:    m.Blocks(),
	}

	for _, opt := range opts {
		opt.apply(downloader)
	}

	return downloader
}

// OutputBlock is the result of a Dowloader worker
type OutputBlock struct {
	Raw   []byte
	Index int
}

// downloadBlock downloads a data block identified by block
func (d *Downloader) downloadBlock(block meta.BlockInfo) ([]byte, error) {
	log.Debugf("downloading block %x", block.Key)
	body, err := d.storage.Get(block.Key)
	if err != nil {
		return nil, err
	}

	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	data, err = snappy.Decode(nil, xxtea.Decrypt(data, block.Decipher))
	if err != nil {
		return nil, err
	}

	hasher, err := blake2b.New(16, nil)
	if err != nil {
		return nil, err
	}

	if _, err := hasher.Write(data); err != nil {
		return nil, err
	}

	hash := hasher.Sum(nil)
	if !bytes.Equal(hash, block.Decipher) {
		return nil, fmt.Errorf("block key(%x), cypher(%x) hash is wrong hash(%x)", block.Key, block.Decipher, hash)
	}

	return data, nil
}

func (d *Downloader) worker(ctx context.Context, feed <-chan int, out chan<- *OutputBlock) error {
	for index := range feed {
		info := d.blocks[index]
		raw, err := d.downloadBlock(info)
		if err != nil {
			log.Errorf("downloading block %d error: %s", index+1, err)
			return err
		}

		result := &OutputBlock{
			Index: index,
			Raw:   raw,
		}

		select {
		case out <- result:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

//Download download the file into this output file
func (d *Downloader) Download(output *os.File) error {
	if len(d.blocks) == 0 {
		return fmt.Errorf("no blocks provided")
	}

	if d.blockSize == 0 {
		return fmt.Errorf("block size is not set")
	}

	workers := int(math.Min(float64(d.workers), float64(len(d.blocks))))
	if workers == 0 {
		workers = int(math.Min(float64(DefaultDownloadWorkers), float64(len(d.blocks))))
	}
	group, ctx := errgroup.WithContext(context.Background())

	feed := make(chan int)
	results := make(chan *OutputBlock)

	//start workers.
	for i := 1; i <= workers; i++ {
		group.Go(func() error {
			return d.worker(ctx, feed, results)
		})
	}

	log.Debugf("downloading %d blocks", len(d.blocks))

	//feed the workers
	group.Go(func() error {
		defer close(feed)
		for index := range d.blocks {
			select {
			case feed <- index:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	go func() {
		group.Wait()
		close(results)
	}()

	count := 1
	for result := range results {
		log.Debugf("writing block %d/%d of %s", count, len(d.blocks), output.Name())
		if _, err := output.Seek(int64(result.Index)*int64(d.blockSize), io.SeekStart); err != nil {
			return err
		}

		if _, err := output.Write(result.Raw); err != nil {
			return err
		}

		count++
	}

	return group.Wait()
}
