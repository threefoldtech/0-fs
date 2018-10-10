# Getting Started with 0-fs
## Mounting an flist
```bash
> g8ufs --help
Usage of g8ufs:
  -backend string
    	Working directory of the filesystem (cache and others) (default "/tmp/backend")
  -cache backend
    	Optional external (common) cache directory, if not provided a temporary cache location will be created under backend
  -debug
    	Print debug messages
  -local-router string
    	Path to local router.yaml to merge with the router.yaml from the flist. This will allow adding some caching layers
  -meta string
    	Path to metadata database (optional)
  -reset
    	Reset filesystem on mount
  -storage-url string
    	Fallback storage url in case no router.yaml available in flist (default "ardb://hub.gig.tech:16379")
  -version
    	Print version and exit
```

- `backend` is a location on physical disk used as a working directory for g8ufs. Backend has the read/write layer of g8ufs.
- `cache` a optional cache directory where downloaded files are stored for later use. A cache directory will be created under `backend` if no one is provided. A cache directory can be shared between multiple instance of g8ufs.
- `debug` prints useful debug information
- `meta` path to flist, or extraced flist
- `reset` if set, the `backend` directory is cleaned up on start, which will causes the mount point to reset to initial flist state. - `storage-url` URL to a store where file blocks can be reached. Supported services are `zdb`, `ardb`, and `redis`. The storage-url is used __ONLY__ if an flist didn't provide a `router.yaml` file. This option is mainly here for backward compatibility with older flist that does not provide router.yaml file.
- `local-router` An optionaly `router.yaml` file that is layerd on top of the `router.yaml` file provided by the flist. This will allow the user of the filesystem to configure local store replication for faster access. Please check the [router](../flist/router.md) for more details.
- `version` print version number and exit


## Creating an flist
It is recommended to first learn how to create a flist, as documented in [Creating Flists](../flists/creating.md).

Then see the [Create a Flist and Start a Container](https://github.com/zero-os/home/blob/master/docs/tutorials/Create_a_Flist_and_Start_a_Container.md) tutorial for an example.

