# Creating Flists

There are two ways to create a flist:
- [Have Zero-OS Hub create the flist](#have-zero-os-ub-create-the-flist)
- [Creating a flists manually using JumpScale](#creating-a-flists-manually-using-jumpscale)

## Have Zero-OS Hub create the flist

This is the easiest and only supported way to create a flist.

In order to have Zero-OS Hub create your flist you simply need to provide it (upload) a tar file containing all the files. In return Zero-OS Hub will do two things:
- Store all files in its ARDB storage backend
- Create the flist file and make it available for download

See the [Create a Flist and Start a Container](https://github.com/zero-os/home/blob/master/docs/tutorials/Create_a_Flist_and_Start_a_Container.md) tutorial for an example.

For more information about [Zero-OS Hub](https://hub.gig.tech) see the [0-hub](https://github.com/zero-os/0-hub) repository.


## Creating a flists manually using JumpScale

This option is only documented for your information, revealing how  Zero-OS Hub implements the first option, documented above.

This option requires JumpScale. The easiest way to meet this requirement is using a Docker container with JumpScale preinstalled, as documented in [Create a JS9 Docker Container](https://github.com/Jumpscale/ays9/blob/master/docs/gettingstarted/js9.md).

From within the JS9 Docker container it actually takes 2 steps to create the flist:
- [Creation of the flist DB and upload the files](#create-db)
- [Package the flist DB](#packqage-db)


<a id="create-db"></a>
### Creation of the flist DB

Open a connection to a RocksDB database:
```python
kvs = j.servers.kvs.getRocksDBStore(name='flist', namespace=None, dbpath="/tmp/flist-example.db")
```

Create a flist object and pass the reference of the RocksDB to it:
```python
f = j.tools.flist.getFlist(rootpath='/', kvs=kvs)
```

Add the path you want to include in your flist, you can do multiple calls to `f.add`:
```python
f.add('/opt')
```

Upload your list to an ARDB storage backend:
```python
f.upload("<remote-ardb-server>", 16379)
```

<a id="package-db"></a>
### Package your flist DB

```shell
cd /tmp/flist-example.db
tar -cf ../flist-example.db.tar *
cd .. && gzip flist-example.db.tar
```

The result is `flist-example.db.tar.gz` which is actually a flist for creating the container.
