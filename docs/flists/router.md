# Block routing
As you know by now, 0-fs will retrieve the files it needs to access
for the first time from the internet. 0-fs does not know where to retrieve the blocks, and it's actually up to the flist to tell the 0-fs where to retrieve the blocks from.

> To avoid manipulation, you should only use signed flists from trusted sources.

and `flist` can provide the system with a `router.yaml` file that describe to the filesystem where to retrieve certain files. The `router.yaml` file SHOULD exist under the `/` of the flist base and has the exact name `router.yaml` if the file is not provided, the `g8ufs` process will fallback to the backward compatibility flag `storage-url` and use it as a router.

## Syntax
A `router.yaml` must define at least one pool, each pool can has one or more `routing rules`. The `router.yaml` must define a `lookup` order

```yaml
pools:
  pool.name:
    <hash-range>: destination.1
    <hash-range>: destination.2

lookup:
 - pool.name
 - ...
```

### Example
A simple router.yaml that points to the hub.git.tech for all flists
that are hosted by hub.gig.tech

```yaml
pools:
  hub:
    00:FF: ardb://hub.gig.tech:16379

lookup:
  - hub
```

## Hash range syntax
- A hash match can be exact (match exact prefix), for example a valid exact range is `AB` which will match all hashes that is prefixed with `AB`. The exact match can be of any length. A `123` is a valid range
- A hash match can define a range, for example a range can be `00:9F` will match all hashes that has prefixes [00, 9F]. The range can also have any length, for example a `000:FFF` is a valid range, As long as start and end prefixes are of the same length.

> In a single pool, ranges can overlap. In that case, all valid destination will be tried.