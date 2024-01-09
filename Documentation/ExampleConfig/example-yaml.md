### dataset-1
```
name: my-dataset-1
spec:
  path: s3://test-fuse/
  credentials:
    aws.accessKeyId: XYZ
    aws.secretKey: ABC
```
### dataset-2
```
name: my-dataset-2
spec:
  path: s3://test-fuse-2/
```


### alluxio-1
```
name: alluxio-1
spec:
  dataset: my-dataset-1
  image: alluxio/alluxio
  etcd:
    enabled: true
  alluxio-monitor:
    enabled: true
```
### alluxio-2
```
name: alluxio-2
spec:
  dataset: my-dataset-2
  image: alluxio/alluxio
  etcd:
    enabled: true
  alluxio-monitor:
    enabled: true
```