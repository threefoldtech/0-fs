package rofs

import (
	"context"
	"fmt"
	"github.com/golang/snappy"
	"github.com/xxtea/xxtea-go/xxtea"
	"github.com/zero-os/0-fs/meta"
	"github.com/zero-os/0-fs/storage"
	"io/ioutil"
	"math"
	"os"
	"sync"
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

func (d *Downloader) worker(ctx context.Context, feed <-chan *DownloadBlock, out chan<- *OutputBlock) {
	defer func() {
		if err := recover(); err != nil {
			log.Warningf("cancel block download: %s", err)
		}
	}()

	for {
		select {
		case blk := <-feed:
			if blk == nil {
				return
			}
			raw, err := d.DownloadBlock(blk.BlockInfo)
			result := &OutputBlock{
				Index: blk.Index,
				Raw:   raw,
			}
			if err != nil {
				log.Errorf("downloading block %d error: %s", blk.Index+1, err)
				result.Err = err
			}
			out <- result
		case <-ctx.Done():
			return
		}
	}
}
func (d *Downloader) writer(output *os.File, expecting int, results <-chan *OutputBlock) error {
	for i := 0; i < expecting; i++ {
		result := <-results
		log.Debugf("writing result of block %d", result.Index+1)
		if result.Err != nil {
			return result.Err
		}

		if _, err := output.Seek(int64(result.Index)*int64(d.BlockSize), os.SEEK_SET); err != nil {
			return err
		}

		if _, err := output.Write(result.Raw); err != nil {
			return err
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

	ctx, cancel := context.WithCancel(context.Background())

	feed := make(chan *DownloadBlock)
	results := make(chan *OutputBlock)

	defer close(feed)
	defer close(results)
	defer cancel()

	//start workers.
	for i := 1; i <= workers; i++ {
		go d.worker(ctx, feed, results)
	}

	//consume all outputs.
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err = d.writer(output, len(d.Blocks), results); err != nil {
			cancel()
		}

		wg.Done()
	}()

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

	wg.Wait()
	return err
}
