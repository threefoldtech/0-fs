# G8ufs
G8os fuse filesystem.
g8ufs can be mounted only using a relatively small meta data database (currently support rocksdb). On accessing
the file it fetches the required file chunks from a remote store, and cache it locally. The idea of using this filesystem
is to speed up `containers` creation by just mounting the container root from any image meta (we call it `flist`) and once
the container starts, it fetches only the required files from the remote store. So no need to clone large images locally.
 
# Design
The fuse mount point is actually a `unionfs` mount of two layers
- `RW` (read-write) layer that is just an actual directory on the raw filesystem of your hard disk
- `RO` (read-only) layer that is the actual fuse mount point. The readonly layer will download the files into a cache when
 they are opened for reading the first time
 
By `merging` those 2 layers on top of each other, (read-write on top) the merged mount point will
expose a `read-write` file system where all file edits, and new files will be written on the `RW` layer,
while reading file operations will be forwarded to the underlaying `read-only` layer. Once a file is opened
for writing (that is only available on the `read-only` layer) it will be copied (copy on write) to the
`read-write` layer and afterwards all `read` and `write` operations will be handled directly by the `RW` layer.

# Building
## Requirements
Make sure you have librocksdb >= v5.2.1
```bash
godep restore
make
```

# Mounting the filesystem
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