# Pega deployment on Kubernetes

This project provides Helm charts and basic examples for deploying Pega on Kubernetes. You will also need to download the required [installation kit](https://community.pega.com/knowledgebase/products/platform/deploy) from the Pega Community which includes rules and data to preload into your relational database. Deploying Pega on Kubernetes requires Pega Infinity 8.2 or newer.

[![Build Status](https://travis-ci.org/pegasystems/pega-helm-charts.svg?branch=master)](https://travis-ci.org/pegasystems/pega-helm-charts)
[![GitHub release](https://img.shields.io/github/release/pegasystems/pega-helm-charts.svg)](https://github.com/pegasystems/pega-helm-charts/releases)

# Supported Kubernetes environments

Pegasystems has validated deployments on the following Kubernetes IaaS and PaaS environments.

* Open-source Kubernetes (and [MiniKube for personal deployments](docs/RUNBOOK_MINIKUBE.md))
* Microsoft Azure Kubernetes Service (AKS) - see the [AKS runbook](docs/Deploying-Pega-on-AKS.md)
* Amazon Elastic Kubernetes Service (EKS)
* Google Kubernetes Engine (GKE) - see the [GKE runbook](docs/Deploying-Pega-on-GKE.md)
* Red Hat OpenShift
* Pivotal Container Service (PKS) - see the [PKS runbook](docs/Deploying-Pega-on-PKS.md)

# Getting started

This project assumes you have an installation of Kubernetes available and have Helm installed locally. The following commands will verify your installation. The exact output may be slightly different, but they should return without error.

```bash
$ helm version
version.BuildInfo{Version:"v3.0.0", GitCommit:"e29ce2a54e96cd02ccfce88bee4f58bb6e2a28b6", GitTreeState:"clean", GoVersion:"go1.13.4"}
```

If this command does not successfully return, install Helm 3 for your operating system.  See [Helm Installation](https://helm.sh/docs/intro/install/) for more information.  If you are running Helm 2.x, you will see both a client and server (tiller) portion returned by the version command.  Some of the commands below will also differ slightly for Helm 2.x.

1. Add the Pega repository to your Helm installation.

```bash
$ helm repo add pega https://dl.bintray.com/pegasystems/pega-helm-charts
```

2. Verify the new repository by searching it.

```bash
$ helm search repo pega
NAME       	CHART VERSION	APP VERSION	DESCRIPTION
pega/pega  	1.2.0        	           	Pega installation on kubernetes
pega/addons	1.2.0        	1.0        	A Helm chart for Kubernetes 
```

There are two charts available in this repository - addons and pega.

The addons chart installs a collection of supporting services and tools required for a Pega deployment. The services you will need to deploy will depend on your cloud environment - for example you may need a load balancer on Minikube, but not for EKS. These supporting services are deployed once per Kubernetes environment, regardless of how many Pega Infinity instances are deployed.

3. Download the values file for pega/pega and pega/addons.

```bash
$ helm inspect values pega/pega > pega.yaml
$ helm inspect values pega/addons > addons.yaml
```

4. Edit your values yaml files to specify all required information and customizations for your environment.

* [Instructions to configure the Pega chart](charts/pega/README.md)
* [Instructions to configure the Pega addons](charts/addons/README.md)

5. Create namespaces for your Pega deployment and the addons (if applicable for your environment).

```bash
$ kubectl create namespace mypega
$ kubectl create namespace pegaaddons
```

6. To install the addons chart, run the following helm command after configuring your values.yaml file (if applicable for your environment). 

```bash
$ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
```

7. Now you can deploy Pega using the Helm chart. Before installing using the chart, it is a good idea to review the detailed [deployment guide](https://community.pega.com/knowledgebase/articles/deploying-pega-platform-using-kubernetes) to understand how Pega deploys as a distributed system. Running a Helm installation using the pega chart installs a Pega Infinity instance into a specified namespace.  

```bash
$ helm install mypega pega/pega --namespace mypega --values pega.yaml
```

*If you want to edit the charts and build using your local copy, replace pega/addons or pega/pega with the path to your chart directory.*

8. If you wish to delete your deployment of Pega nodes, enter the following command (this will not delete your database):

```bash
$ helm delete mypega
```
# Downloading Docker images for your deployment

Clients with appropriate licenses can request access to several required images from the Pega-managed Docker image repository. With your access key, you can log in to the image repository and download docker images that are required to install the Pega Platform onto your database. After you pull the images to your local system you must push them into your private Docker registry.

To download your preferred version of the Pega image to your local system, specify the version tag when you enter:

```bash
$ sudo docker pull pega-docker.downloads.pega.com/platform/pega:<version>
Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/pega:<version>
```

For details, see the examples listed in the runbooks:

* [Preparing your local Linux system](docs/prepping-local-system-runbook-linux.md)
* [Preparing your local Windows 10 system](docs/prepping-local-system-runbook-windows.md)

Pegasystems uses a standard naming practice of hostname/product/image:tag. All Pega images are available from the pega-docker.downloads.pega.com host. The :tag represents the version if Pega being deployed, for example :8.3.1 to download Pega 8.3.1. Pega maintains three types of required Docker images for Client-managed Cloud deployments of Pega Platform:

 Name        | Description                                           |
-------------|-------------------------------------------------------|
platform/pega  | Download required. Deploys the Pega Platform with its customized version of the Tomcat application server |
 platform/search | Download required. Deploys the search engine required for the Pega Platform application’s search and reporting capabilities. This Docker image contains Elasticsearch and includes all required plugins |
 platform/installer   | A utility image Pega Platform deployments use to install or upgrade all of the Pega-specific rules and database tables in the “Pega” database you have configured for your deployment.

 Clients push their three Pega-provided images to their registry so it is available to the deployment. Clients then provide their registry URL, credentials, and then reference each image appropriately in the Pega Helm chart. Example usage details for referencing your three images in a repository are included in the appropriate runbook for your type of deployment.

# Contributing

This is an open source project and contributions are welcome.  Please see the [contributing guidelines](./CONTRIBUTING.md) to get started.
