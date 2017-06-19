# 0-fs

0-fs is the fuse file system of Zero-OS.

0-fs can be mounted only using a relatively small meta data database (currently support RocksDB). On accessing
the file it fetches the required file chunks from a remote store, and cache it locally. The idea of using this file system
is to speed up container creation by just mounting the container root from any image metadata file (we call it a `flist` file) and once
the container starts, it fetches only the required files from the remote store. So no need to clone large images locally.

## Design

The fuse mount point is actually a `unionfs` mount of two layers:
- **RW** (read-write) layer that is just an actual directory on the raw file system of your hard disk
- **RO** (read-only) layer that is the actual fuse mount point. The read-only layer will download the files into a cache when they are opened for reading the first time

By `merging` those 2 layers on top of each other, (read-write on top) the merged mount point will
expose a read-write file system where all file changes, and new files get written to the RW layer,
while reading file operations will be forwarded to the underlaying read-only layer. Once a file is opened
for writing (that is only available on the read-only layer) it will be copied (copy on write) to the
read-write layer and afterwards all read and write operations will be handled directly by the RW layer.

## Building

Make sure you have `librocksdb` v5.2.1 or higher.

```bash
godep restore
make
```

## Mounting the file system

```
$ ./g8ufs -h
Usage of ./g8ufs:
  -backend string
    	Working directory of the filesystem (cache and others) (default "/tmp/backend")
  -debug
    	Print debug messages
  -meta string
    	Path to metadata database (rocksdb)
  -reset
    	Reset filesystem on mount
  -storage-url string
    	Storage url (default "ardb://hub.gig.tech:16379")
```

## More

All documentation is in the [`/docs`](./docs) directory, including a [table of contents](/docs/SUMMARY.md).

In [Getting Started with 0-fs](/docs/gettingstarted/README.md) you find the recommended path to quickly get up and running.
