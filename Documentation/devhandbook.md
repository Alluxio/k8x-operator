# Developemnt Handbook

#### Make Sure you have a k8s cluster

## Local Testing

### Step 0
Generate Helm Chart files by running `./dev/build/generate.sh` under project root.

### Step 1
Install Alluxion Operator via Helm Chart under `deploy/charts/alluxio-operator`
  ```shell
  helm install operator -f operator-config.yaml deploy/charts/alluxio-operator
  ```

### Uninstall Operator:
  ```shell
  helm delete operator 
  ```


## Make Alluxio K8s Operator Docker Image

### Step 0
Create a dockerhub account, and login in terminal

### Step 1
Generate Helm Chart files by running `./dev/build/generate.sh` under project root.

### Step 2
Build dokcer image by running `docker build -t <docker username>/alluxio-operator:<tag> -f dev/build/Dockerfile .` under project root.

* For Apple Silicon Chip: `docker buildx build --platform linux/amd64 -t <docker username>/alluxio-operator:<tag> -f dev/build/Dockerfile .`


* Example:
  ```shell
  docker buildx build --platform linux/amd64 -t kshou433/alluxio-operator:v1.6 -f dev/build/Dockerfile .
  ```

### Step 3
Push image to docker hub : `docker push <docker username>/alluxio-operator:<tag>`.

* Example:
  ```shell
  docker push kshou433/alluxio-operator:v1.6
  ```

### Step 4
Update image url and tage in ```operator-config.yaml```


## Verify the Deployment

### Check endpoints
```kubectl get endpoints -A```

### Check if prometheus is running
```kubectl get prometheus -o yaml```

### Check if several prometheus operators are running on the cluster:
```kubectl get pods --all-namespaces | grep 'prom.*operator'```

### Access prometheus GUI
```kubectl -n monitoring port-forward svc/prometheus-k8s 9090```
