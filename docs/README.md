# 0-fs

0-fs is the Zero-OS file system used for containers, which is actually a FUSE file system.

Mounting 0-fs is done by using a flist, which is a relatively small RocksDB database file, containing the metadata of the actual files and directories. On accessing a file Zero-OS fetches the required file chunks from a remote key-value store, and caches it locally. The remote key-value store that is used is configured with the `storage` global Zero-OS parameter, documented in [Main Configuration](https://github.com/zero-os/0-core/blob/master/docs/config/main.md). The default is set to the ARDB storage cluster implemented in [Zero-OS Hub](https://hub.gig.tech).

The idea of using this approach is to speed up container creation by just mounting the container root from the image metadata contained in the flist file and once the container starts, it fetches only the required files from the remote key-value store. So no need to clone large images locally.

The FUSE mount point is actually a UnionFS mount of two layers:
- **RW (read-write) layer**, which is just a directory on the cache disk of the Zero-OS node.
- **RO (read-only) layer**, which is the actual FUSE mount point. The read-only layer will download the files into a cache when they are opened for reading the first time

By merging those 2 layers on top of each other, (read-write on top) the merged mount point will expose a read-write file system where all file edits, and new files will be written on the RW layer, while reading file operations will be forwarded to the underlaying read-only layer. Once a file is opened for writing (that is only available on the read-only layer) it will be copied (copy on write) to the read-write layer and afterwards all read and write operations will be handled directly by the RW layer.


See the [Table of Contents](SUMMARY.md) for more documentation on 0-fs.

In [Getting Started with 0-fs](gettingstarted/README.md) you find the recommended path to quickly get up and running.
