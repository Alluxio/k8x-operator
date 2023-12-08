# How to Set up Local Testing

### Step 0
注册一个dockerhub的账号，在terminal里docker login

### Step 1
在project home的路径下运行 
```shell
./dev/build/generate.sh
```

### Step 2
在project home的路径下运行 docker:  
```shell
docker build -t <docker username>/alluxio-operator:<tag> -f dev/build/Dockerfile .
```

* For Apple Silicon Chip: 
  ```shell
  docker buildx build --platform linux/amd64 -t kshou433/alluxio-operator:rev2 -f dev/build/Dockerfile .
  ```


### Step 3
```shell
docker push kshou433/alluxio-operator:rev2
```

### Step 4
在operator的config里面update image和imageTag




