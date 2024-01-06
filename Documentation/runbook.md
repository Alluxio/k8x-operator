# Alluxio K8S Operator Runbook

## Background



## How to deploy Alluxio K8S Operator with [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus)

### Install Prometheus Operator to K8S Cluster

Make sure add correct namespace in the jsonnet file and generate new manifest file.

```
  prometheus+:: {
    namespaces: ["default", "kube-system", "monitoring", "alluxio-operator"],
  },
```

##### These errors are expected:
```shell
Error from server (NotFound): namespaces "alluxio-operator" not found
Error from server (NotFound): namespaces "alluxio-operator" not found
```
We need to use next step to deploy the `Alluxio Operator`.

### Deploy Alluxio Operator using Helm
`helm install operator -f operator-config.yaml deploy/charts/alluxio-operator`

### Update Prometheus Operator

### Use method in Dev Handbook to verify everything is running,


### Uninstall Operator:
`helm delete operator `