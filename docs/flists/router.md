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

## Cache
A `router.yaml` can define a `cache` list. Which lists a set of pools (one or more). A cache is always updated with a block once it's retrieved from the lookup. Usually an flist should not define a `cache` list. It's the user of the flist who would probably need to add his `cache` entries so blocks retrieved from remote are cached locally for faster access next time.

`0-fs` provides a simple way to merge to `router.yaml`, a locally defined `router.yaml` with the one published by the flist itself. This way a user can hook in his local infrastructure for caching. For example, a user can has this local `router.yaml` file

```yaml
pools:
  local:
    00:FF : zdb://192.168.1.1:12345

lookup:
  - local

cache:
  - local
```

Then when starting the `0-fs` process, you can pass the path to this yaml file using the `-local-router` flag.
`0-fs` will merge this with the `router.yaml` providing by the flist. Hence you'll end up with a an `router.yaml` that looks like this

```yaml
pools:
  local:
    00:FF: zdb://192.168.1.1:12345
  remote:
    00:FF: ardb://hub.gig.tech:16379

lookup:
  - local
  - remote

cache:
  - local
```

This will configure `0-fs` to do the following
- When retrieving a block, first try `local` pool.
- If not exist try `remote` pool.
- If block is retrieved successfully from a pool that is not listed in cache, update `local` with that block
- Next time the same block is requested, it will be found in local, no call to remote would be needed.
