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
	Workers   int
	Storage   storage.Storage
	Blocks    []meta.BlockInfo
	BlockSize uint64
}

// OutputBlock is the result of a Dowloader worker
type OutputBlock struct {
	Raw   []byte
	Index int
}

// DownloadBlock downloads a data block identified by block
func (d *Downloader) DownloadBlock(block meta.BlockInfo) ([]byte, error) {
	log.Debugf("downloading block %x", block.Key)
	body, err := d.Storage.Get(block.Key)
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
		info := d.Blocks[index]
		raw, err := d.DownloadBlock(info)
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
	if len(d.Blocks) == 0 {
		return fmt.Errorf("no blocks provided")
	}

	if d.BlockSize == 0 {
		return fmt.Errorf("block size is not set")
	}

	workers := int(math.Min(float64(d.Workers), float64(len(d.Blocks))))
	if workers == 0 {
		workers = int(math.Min(float64(DefaultDownloadWorkers), float64(len(d.Blocks))))
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

	log.Debugf("downloading %d blocks", len(d.Blocks))

	//feed the workers
	group.Go(func() error {
		defer close(feed)
		for index := range d.Blocks {
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
		log.Debugf("writing block %d/%d of %s", count, len(d.Blocks), output.Name())
		if _, err := output.Seek(int64(result.Index)*int64(d.BlockSize), io.SeekStart); err != nil {
			return err
		}

		if _, err := output.Write(result.Raw); err != nil {
			return err
		}

		count++
	}

	return group.Wait()
}
