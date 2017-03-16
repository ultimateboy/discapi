This project serves as a demo of [Rest-Layer](https://github.com/rs/rest-layer) by [Olivier Poitrey](https://github.com/rs), but also as a demonstration of the [Deis](https://deis.com) Golang workflow.

Almost all Deis projects begin as lines of Golang. The code is built using a standardized docker build container. The resultant binary is then loaded into an application docker container and referenced via a helm chart which can be deployed to Kubernetes.

This is all accomplished by using a few standardized `make` targets. The final chart can be deployed to Kubernetes using `helm`.

## API

This app creates a REST API which exposes endpoints for managing a list of Disc Golf companies and their associated discs. It also allows for user registration and authentication and allows users to manage a list of discs in their personal bag.

Rest-Layer has support for multiple backends, but this example uses MongoDB. The helm chart comes with a `requirements.yaml` file which will deploy the official MongoDB chart to your Kubernetes cluster.

## Testing

There are some basic tests using the [Frisby](github.com/verdverm/frisby) framework, for obvious reasons. Testing uses the in-memory backend.

```bash
$ make test
```

## Building

This is only necessary if you make changes to the code. Otherwise, you can use the pre-built docker containers already referenced in the helm chart.

```bash
$ make bootstrap  # this installs the project's dependencies using glide
$ make build  # this builds the binary followed by building the application container
$ make push  # this pushes the image to a docker registry
```

## Requirements

The following are required to build and deploy.

  - [docker](https://www.docker.com/)
  - [gcloud](https://cloud.google.com/sdk/downloads)
  - [kubectl](https://kubernetes.io/docs/user-guide/prereqs/)
  - [helm](https://github.com/kubernetes/helm/releases)
  - [httpie](https://httpie.org/)

## Create a Kubernetes cluster

This can be on any cloud provider, but for this example, we use Google Cloud.

```bash
$ gcloud auth login
$ gcloud alpha projects create some-unique-project-id
$ gcloud config set project some-unique-project-id
$ gcloud container clusters create discapi --zone=us-central1-b
$ helm init
```

## Deployment

```bash
$ helm repo add charts https://kubernetes-charts.storage.googleapis.com/
$ cd chart/discapi
$ helm dependency update
$ helm upgrade --install discapi ./
```

## Usage

It might take a couple of minutes for the service to be given an external load balaced IP.

```bash
$ export DISCAPI=$(kubectl get svc discapi -o 'go-template={{range .status.loadBalancer.ingress}}{{.ip}}{{end}}')
```

### Create new companies
`http post $DISCAPI/companies name="discraft" website="http://discraft.com/discgolf.html"`

`http post $DISCAPI/companies name="innova" website="http://www.innovadiscs.com/"`

### Create discs
`http post $DISCAPI/companies/innova/discs name="Roc3" speed:=5 glide:=4 turn:=0 fade:=3`

### Get all innova discs
`http get $DISCAPI/companies/innova/discs`

### Get high speed innova discs:
`http get $DISCAPI/companies/innova/discs filter=={\"speed\":{\"\$gt\":6}}`

# License

This project is licensed under the MIT License. See the `LICENSE` file for more information.