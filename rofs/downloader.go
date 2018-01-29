package rofs

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/snappy"
	"github.com/xxtea/xxtea-go/xxtea"
	"github.com/zero-os/0-fs/meta"
	"github.com/zero-os/0-fs/storage"

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

type DownloadBlock struct {
	meta.BlockInfo
	Index int
}

type OutputBlock struct {
	Raw   []byte
	Index int
	Err   error
}

func (d *Downloader) DownloadBlock(block meta.BlockInfo) ([]byte, error) {
	log.Debugf("downloading block %s", string(block.Key))
	body, err := d.Storage.Get(string(block.Key))
	if err != nil {
		return nil, err
	}

	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	data = xxtea.Decrypt(data, block.Decipher)
	return snappy.Decode(nil, data)
}

func (d *Downloader) worker(ctx context.Context, feed <-chan *DownloadBlock, out chan<- *OutputBlock) error {
	for blk := range feed {
		raw, err := d.DownloadBlock(blk.BlockInfo)
		result := &OutputBlock{
			Index: blk.Index,
			Raw:   raw,
		}
		if err != nil {
			log.Errorf("downloading block %d error: %s", blk.Index+1, err)
			result.Err = err
		}
		select {
		case out <- result:
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func (d *Downloader) writer(ctx context.Context, output *os.File, results <-chan *OutputBlock) error {

	for result := range results {
		log.Debugf("writing result of block %d", result.Index+1)
		if result.Err != nil {
			return result.Err
		}

		select {
		default:
			if _, err := output.Seek(int64(result.Index)*int64(d.BlockSize), os.SEEK_SET); err != nil {
				return err
			}

			if _, err := output.Write(result.Raw); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

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
	downloaderGroup, _ := errgroup.WithContext(ctx)

	feed := make(chan *DownloadBlock)
	results := make(chan *OutputBlock)

	//start workers.
	for i := 1; i <= workers; i++ {
		downloaderGroup.Go(func() error {
			return d.worker(ctx, feed, results)
		})
	}

	group.Go(func() error {
		err := downloaderGroup.Wait()
		close(results)
		return err
	})

	//consume all outputs.
	group.Go(func() error {
		return d.writer(ctx, output, results)
	})

	for i, block := range d.Blocks {
		downloadBlock := &DownloadBlock{
			BlockInfo: block,
			Index:     i,
		}
		select {
		case feed <- downloadBlock:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	close(feed)

	return group.Wait()
}
