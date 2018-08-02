# Introduction
To allow 3rd party to create and host flists, an flist can have a `router.yaml` file which defines a set of locations
where need to looked for the files blocks using the block hash.

The `router.toml` defines a set of pools, each pool defines one or more ranges and the corresponding node that hosts this
hash range.

# `router.yaml`
the `router.yaml` must be found (if provided) under the flist root directly. and has the following syntax
```yaml
pools:
  pool-one:
    00:69: ardb://host:port
    70:FF: zdb://host:port

lookup:
  - pool.name
```
