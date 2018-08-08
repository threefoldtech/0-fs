package rofs

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/snappy"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
	"github.com/xxtea/xxtea-go/xxtea"

	"golang.org/x/sync/errgroup"
)

const (
	DefaultDownloadWorkers = 4
	DefaultBlockSize       = 512 //KB
)

type Downloader struct {
	Workers   int
	Storage   storage.Storage
	Blocks    []meta.BlockInfo
	BlockSize uint64
}

type OutputBlock struct {
	Raw   []byte
	Index int
}

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

	return snappy.Decode(nil, xxtea.Decrypt(data, block.Decipher))
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

	for result := range results {
		if _, err := output.Seek(int64(result.Index)*int64(d.BlockSize), os.SEEK_SET); err != nil {
			return err
		}

		if _, err := output.Write(result.Raw); err != nil {
			return err
		}
	}

	return group.Wait()
}
