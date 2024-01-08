# K8s Alluxio Operator Dashboard Wiki

## Background and Motivation

Kubernetes extensively relies on a Command Line Interface (CLI) for interacting with cluster resources. This approach often requires users to navigate complex `kubectl` commands and configuration files to perform CRUD (Create, Read, Update, Delete) operations on Kubernetes resources. Such complexity can pose challenges, especially for those new to Kubernetes or those who prefer a more streamlined, user-friendly interface.

To address this challenge, we've developed a full-stack application that provides a graphical user interface for users to efficiently perform operations on Alluxio Clusters and Datasets. This user-friendly approach significantly simplifies the management of Alluxio on Kubernetes, making it accessible and effective for a broader range of users, regardless of their technical expertise.

## Development Handbook

### Overview
The Alluxio Operator Dashboard is a full-stack project comprising two main components: **the API Server** and **the User Interface**. Once the User Interface is compiled into production-ready files, it can be hosted by the API Server. This setup allows users to directly access the dashboard, providing a seamless and integrated experience.

### API Server on K8s Operator
A new **Kubernetes Deployment** that handles Restful API Request from the User Interface and also hosts the User Interface. Similar to _alluxio-controller_ and _dataset-controller_, it has its own deployment file: `api-server-controller.yaml`. 

#### Restuful WebService
The API Server uses _Restuful WebService_ to communicates with GUI with endpoints: `api/dataset` and `api/alluxio_cluster`. Each endpoint currently supports `GET`, `POST` and `DELETE` HTTP request methods.

When API Server gets a Restful API Request, it uses Kubernetes `controller-runtime` to interact with Kubernetes Cluster. 

We also have converter functions that can simplify the Kubernetes controller-runtime data or transform user input to the defined CRD format for Kubernetes `controller-runtime`.

#### Host User Interface
Embedding production build of User Interface by using `go:embed`. And the endpoint of User Interface is set to `/`. 

#### Run Locally
API-Server has its own `main.go`, you can simply start this application along with the Alluxio Operator deployed in the Kubernetes Cluster. The defualt port has been set to `8080`.

- You can use `curl` to communicate with API Server without the User Interface.


### User Interface
The User Interface is built on **React** with JavaScript, and is using Redux to manange shared state. 

In current version, we are using *continuous polling* to fetch data from the API Server.

#### Run Locally
Using `npm start`. The proxy in `package.json` has been set to `http://localhost:8080` to match the API Server setting.

#### Integrate UI into API Server
The API Server embed the static file under `cmd/api_server/api_server/gui` folder.

- Preferred Method:

  We have `"build": "BUILD_PATH='./gui' react-scripts build && cp -r gui ../cmd/api_server/api_server/"` in `package.json`.

  Now, `npm build` will generate production build of User Interface to the `gui` folder and also copy it to `cmd/api_server/api_server/`


- Legacy Method:

  Copy generated production build file to `cmd/api_server/api_server/gui` folder.


## Deployment
### Deploy Operator with API-Server
#### Please follow [this document](https://tachyonnexus.atlassian.net/wiki/spaces/ENGINEERIN/pages/86147073/K8S+Operator+Wiki#Prerequisites) for a more detailed reference.
- Get operator deployment tarball. Untar and Name it `alluxio-operator`

- Create a private repository on [AWS ECR](https://docs.aws.amazon.com/AmazonECR/latest/userguide/repository-create.html). Name it `alluxio-operator`.

- Push Operator with API-Server Docker Image to ECR.
  * Follow the steps in Appendix to learn how to build local docker image.

- Create a `operator-config.yaml`:
```yaml
image: <LINK_TO_ECR_REPO>
imageTag: <IMG_TAG>
imagePullPolicy: Always
alluxio-csi:
  enabled: false
```

- Execute `helm install operator -f operator-config.yaml alluxio-operator` to install Alluxio Operator.

### Access Dashboard
#### Port Forward Method
- Use  `kubectl get po -A` to find API Server Pod Name: `<api-server-controller-XXX>`
- Run: `kubectl port-forward -n alluxio-operator <api-server-controller Pod Name> <LOCAL PORT>:8080`
- Go to `http://localhost:<LOCAL PORT>/`

<br>

-----------
## Appendix

### Generate Alluxio K8s Operator Docker Image
#### Step 0
Create a dockerhub account, and login in terminal

#### Step 1
Generate Helm Chart files by running `./dev/build/generate.sh` under project root.

#### Step 2
Build dokcer image by running `docker build -t <docker username>/alluxio-operator:<tag> -f dev/build/Dockerfile .` under project root.

* For Apple Silicon Chip: `docker buildx build --platform linux/amd64 -t <docker username>/alluxio-operator:<tag> -f dev/build/Dockerfile .`

#### Step 3
Push image to docker hub : `docker push <docker username>/alluxio-operator:<tag>`.

#### Step 4
Update image url and tage in ```operator-config.yaml```

### Test Operator Docker Image Locally
#### Install Operator:
Install Alluxion Operator via Helm Chart under `deploy/charts/alluxio-operator`
  
  `helm install operator -f operator-config.yaml deploy/charts/alluxio-operator`

#### Uninstall Operator:
  `helm delete operator `