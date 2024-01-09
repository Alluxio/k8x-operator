
## Install

### Prep Evn
```bash
$ mkdir my-kube-prometheus; cd my-kube-prometheus
$ jb init  # Creates the initial/empty `jsonnetfile.json`
# Install the kube-prometheus dependency
$ jb install github.com/prometheus-operator/kube-prometheus/jsonnet/kube-prometheus@main # Creates `vendor/` & `jsonnetfile.lock.json`, and fills in `jsonnetfile.json`

$ wget https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/main/example.jsonnet -O example.jsonnet
$ wget https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/main/build.sh -O build.sh
$ chmod +x build.sh
```

###  Compile the manifests:
`./build.sh example.jsonnet`



## Deployment
### Install
```shell
kubectl apply --server-side -f manifests/setup -f manifests
```
### Delete
```shell
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup
```

