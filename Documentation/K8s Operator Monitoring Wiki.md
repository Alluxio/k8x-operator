# K8s Operator Monitoring Wiki

## Overview
To achieve Monitoring, we need to have a Prometheus Operator as well as create custom metrics, alerts and recording rules for Alluxio Operator.

There are two `ServiceMonitor` : `alluxio-controller-monitor` and  `dataset-controller-monitor`

They have their own `Service` which is `alluxio-operator-metrics` and `dataset-operator-metrics`

## Deployment
### Install Prometheus Operator to K8S Cluster
We recommand to use [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus).

- Make sure to add correct namespace in the `.jsonnet` and generate new a manifest.
    ```
    prometheus+:: {
        namespaces: ["default", "kube-system", "monitoring", "alluxio-operator"],
    },
    ```

-  Install Prometheus Operator: 
    ```shell
    kubectl apply --server-side -f manifests/setup -f manifests
    ```

    -  These errors are expected since we have not   yet deply `alluxio-operator`:
        ```shell
        Error from server (NotFound): namespaces "alluxio-operator" not found
        Error from server (NotFound): namespaces "alluxio-operator" not found
        ```

    - Uninstall
        ```shell
        kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup
        ```

### Deploy Alluxio Operator using Helm
```shell
helm install operator -f operator-config.yaml deploy/charts/alluxio-operator
```

### Update Prometheus Operator
```shell
kubectl apply --server-side -f manifests/setup -f manifests
```


### Access prometheus GUI
`kubectl -n monitoring port-forward svc/prometheus-k8s 9090`

### Verify the Deployment
**Check endpoints:**

`kubectl get endpoints -A`

**Check if Prometheus is running:**

`kubectl get prometheus -o yaml`
You should see:
```
.....
items:
- apiVersion: monitoring.coreos.com/v1
  kind: Prometheus
  metadata:
    ......
    labels:
      app.kubernetes.io/managed-by: Helm
      prometheus: k8s
    name: k8s
    namespace: default
    ......
```

**Check if several prometheus operators are running on the cluster:**

`kubectl get pods --all-namespaces | grep 'prom.*operator'`


## Development Handbook
### Customize Metrics

- Declare Metrics in `monitoring/metrics.go`.

    [Metric Types](https://prometheus.io/docs/concepts/metric_types/)

    [Reference](https://github.com/operator-framework/operator-sdk/blob/master/testdata/go/v4/monitoring/memcached-operator/monitoring/metrics.go)

- Update Metrics in Controller `Reconcile` logic.


### Customize Alerts
- Declare Alert in `monitoring/alerts.go`.

    [Alert Rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)

    [Reference](https://github.com/operator-framework/operator-sdk/blob/master/testdata/go/v4/monitoring/memcached-operator/monitoring/alerts.go)

- Update Alert in Controller `Reconcile` logic.


<br>

-----------
## Appendix

### Customize `kube-prometheus`
#### Prepare Environment
```shell
$ mkdir my-kube-prometheus; cd my-kube-prometheus
$ jb init  # Creates the initial/empty `jsonnetfile.json`
# Install the kube-prometheus dependency
$ jb install github.com/prometheus-operator/kube-prometheus/jsonnet/kube-prometheus@main # Creates `vendor/` & `jsonnetfile.lock.json`, and fills in `jsonnetfile.json`

$ wget https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/main/example.jsonnet -O example.jsonnet
$ wget https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/main/build.sh -O build.sh
$ chmod +x build.sh
```
#### Edit the manifests
Edit `example.jsonnet`.

#### Compile the manifests
`./build.sh example.jsonnet`

#### Install
```shell
kubectl apply --server-side -f manifests/setup -f manifests
```
#### Delete
```shell
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup
```
