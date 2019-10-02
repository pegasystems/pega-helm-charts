# Pega deployment on Kubernetes

This project provides Helm charts and basic examples for deploying Pega on Kubernetes. This project **does not include** the required database installation image which you [may obtain from the Pega Community](https://community.pega.com/knowledgebase/products/platform/deploy).  Deploying Pega on Kubernetes requires Pega Infinity 8.2 or newer.

[![Build Status](https://travis-ci.org/pegasystems/pega-helm-charts.svg?branch=master)](https://travis-ci.org/pegasystems/pega-helm-charts)
[![GitHub release](https://img.shields.io/github/release/pegasystems/pega-helm-charts.svg)](https://github.com/pegasystems/pega-helm-charts/releases)

# Supported Kubernetes environments

Pegasystems has validated deployments on the following Kubernetes IaaS and PaaS environments.

* Open-source Kubernetes (and [MiniKube for personal deployments](docs/RUNBOOK_MINIKUBE.md))
* Microsoft Azure Kubernetes Service (AKS)
* Amazon Elastic Kubernetes Service (EKS)
* Google Kubernetes Engine (GKE)
* Red Hat OpenShift
* Pivotal Container Service (PKS)

# Getting started

This project assumes you have an installation of Kubernetes available and have Helm installed locally.  The following commands will verify your installation.  The exact output may be slightly different, but they should return without error.  
```console
$ kubectl get nodes
NAME                              STATUS    ROLES     AGE       VERSION
ip-xxx-xxx-xxx-xxx.ec2.internal   Ready     <none>    2d        v1.11.5
ip-yyy-yyy-yyy-yyy.ec2.internal   Ready     <none>    2d        v1.11.5
ip-zzz-zzz-zzz-zzz.ec2.internal   Ready     <none>    2d        v1.11.5

$ helm version
Client: &version.Version{SemVer:"v2.12.2", GitCommit:"7d2b0c73d734f6586ed222a567c5d103fed435be", GitTreeState:"clean"}
Server: &version.Version{SemVer:"v2.12.2", GitCommit:"7d2b0c73d734f6586ed222a567c5d103fed435be", GitTreeState:"clean"}
```

Start by performing a clone (or download) of the latest Charts.

```bash
git clone https://github.com/pegasystems/pega-helm-charts.git
```

## Update dependencies

The Pega charts depends on other charts supplied by third parties.  These are called out in the requirements yaml file for the [pega](charts/pega/requirements.yaml) and [addons](charts/addons/requirements.yaml) charts.  Individual dependencies may or may not be deployed based on the configuration of your values.yaml files.  When you first setup your helm chart, you will need to update your dependencies to pull down these additional charts from their repositories.  For convenience, the required commands are part of the [Makefile](Makefile) and can run with the following command.

```bash
make dependencies
```

For more information about Helm dependencies, see the [Helm documentation](https://helm.sh/docs/helm/#helm-dependency).

## Configure and install using the charts

There are two charts available in this repository - *addons* and *pega*. 

The addons chart installs a collection of supporting services and tools required for a Pega deployment.  The services you will need to deploy will depend on your cloud environment - for example you may need a load balancer on Minikube, but not for EKS. These supporting services are deployed once per Kubernetes environment, regardless of how many Pega Infinity instances are deployed.

[Instructions to configure the Pega addons](charts/addons/README.md)

To install the addons chart, run the following helm command after configuring your values.yaml file.

```bash
helm install . -n pegaaddons --namespace pegaaddons --values /home/user/my-overridden-values.yaml
```

After installing the addons, you can deploy Pega. Before installing using the chart, it is a good idea to review the detailed [deployment guide](https://community.pega.com/knowledgebase/articles/deploying-pega-platform-using-kubernetes) to understand how Pega deploys as a distributed system. Running a Helm installation using the pega chart installs a Pega Infinity instance into a specified namespace.  

[Instructions to configure the Pega chart](charts/pega/README.md)

To install the pega chart, run the following helm command after configuring your values.yaml file.

```bash
helm install . -n mypega --namespace myproject --values /home/user/my-overridden-values.yaml
```

To delete this chart, enter:

```bash
helm delete mypega --purge
```

Navigate to the project directory and open the values.yaml file.  This is the configuration file that tells Helm what and how to deploy.  For additional documentation covering the different deployment options, see the Pega Community article on [Deploying the Pega Platform by using Kubnernetes](https://community.pega.com/knowledgebase/articles/deploying-pega-platform-using-kubernetes).

# Contributing

This is an open source project and contributions are welcome.  Please see the [contributing guidelines](./CONTRIBUTING.md) to get started.