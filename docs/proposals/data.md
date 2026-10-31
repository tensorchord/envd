# Data function support

## Summary

To provide mount data support for envd


## Goals

Design a *unified, declarative* interface and underlying architecture to provide dataset in the development environment in a *scalable way*


Non-goals: 
- Support Git-like version control for data

## Common Scenarios

### Possible sources
- local files
- Object storage (AWS S3)
- NFS-like system (AWS EFS, AWS FSx for OpenZFS)
- Block storage (Ceph)
- HDFS
- Lustre
- API endpoint (http path)
- SQL results
- Other distributed fs (alluxio, juicefs)
- Python SDK

### Possible form
- Images
- Text
- Embedding binarys
- CSV

### Access Pattern

The access pattern of most dataset is write once, read multiple times, and concurrently. Therefore 

### Possible versions/tags
- Version by number, V1, V2, V3
- Version by scale, sample dataset vs full dataset
- Version by time, query range of user activity (7d, 30d) from feature store

We can have a new standard on how to version the data like semver

## Proposal

Each version of dataset is immutable. By assuming the data is immutable, we can cache the data and make replication easily, to increase the read throughput in multiple ways.


### Usage

User need to create the dataset beforehand. Than declare mounting in the build.envd file.

```
envd data add -f mnist.yaml
```

User can create multiple dataset with the same name, but need to be different versions

mnist.yaml
```yaml=
ApiVersion: V1alpha
name: mnist
version: "0.0.1-sample"
sources:
    - type: local # First source will be considered major source, others are the replication of this one
      path: ~/.torch/mnist
    - type: s3
      path: xxx
validation:
    checksum:
        - name: MD5
          value: xxxx
```

build.envd
```
def data():
    return [d.mount("mnist", target="./data")] # User can specify mount multiple datasets
```