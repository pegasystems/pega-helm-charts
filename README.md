# Pega deployment on Kubernetes

This project provides Helm charts and basic examples for deploying Pega on Kubernetes. You will also need to download the required [Pega-provided Docker images](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/pega-docker-images-manage.html) which include rules and data to preload into your relational database. Deploying Pega on Kubernetes requires Pega Infinity 8.3 or later.

[![Build Status](https://github.com/pegasystems/pega-helm-charts/actions/workflows/github-actions-build.yml/badge.svg)](https://github.com/pegasystems/pega-helm-charts/actions/workflows/github-actions-build.yml)
[![GitHub release](https://img.shields.io/github/release/pegasystems/pega-helm-charts.svg)](https://github.com/pegasystems/pega-helm-charts/releases)

# Supported Kubernetes environments

Pegasystems has validated deployments on the following Kubernetes IaaS and PaaS environments.

* Open-source Kubernetes 
* MiniKube for personal deployments - see the [MiniKube runbook](docs/RUNBOOK_MINIKUBE.md)
* Microsoft Azure Kubernetes Service (AKS) - see the [AKS runbook](docs/Deploying-Pega-on-AKS.md)
* Amazon Elastic Kubernetes Service (EKS) - see the [EKS runbook](docs/Deploying-Pega-on-EKS.md)
* Google Kubernetes Engine (GKE) - see the [GKE runbook](docs/Deploying-Pega-on-GKE.md)
* Red Hat OpenShift Container Platform (Self-managed) - see the [OpenShift runbook](docs/Deploying-Pega-on-openshift.md)
* VMware Tanzu Kubernetes Grid Integrated Edition (TKGI) - see the [TKGI runbook](docs/Deploying-Pega-on-PKS.md) (**Note:** VMware Tanzu Kubernetes Grid Integrated Edition (TKGI) was deprecated for all releases in February 2024. Current deployments on TKGI continue to be supported, but as a best practice, do not use TKGI for new deployments of Pega Platform.)

Pega provides multiple runbooks that are Kubernetes-environment specific for clients looking for minimum guidance for configuring and deploying Pega Platform in their organization's Kubernetes account. Pega also provides runbooks for patching or updating Pega Platform with near-zero-downtime. To access the runbooks, click on the appropriate link in the list above.

Because open-source Kubernetes depends on your specific configuration, Pega does not provide a single open-source Kubernetes runbook. Use the Pega helm charts documentation to configure your deployment, with [Preparing your local Linux system](docs/prepping-local-system-runbook-linux.md), [Preparing your local Windows 10 system](docs/prepping-local-system-runbook-windows.md), and the Getting started section below as your starting points.

The helm charts support running on Kubernetes releases that are in Maintenance and Active Support at the time of the helm chart’s release. To review the date of a Pega helm chart release, see [Releases](https://github.com/pegasystems/pega-helm-charts/releases); to review the latest matrix of Kubernetes release support,  see [Kubernetes EOL policy](https://endoflife.date/kubernetes).

# Getting started

This project assumes you have an installation of Kubernetes available and have Helm installed locally. The following commands will verify your installation. The exact output may be slightly different, but they should return without error.

```bash
$ helm version
version.BuildInfo{Version:"v3.0.0", GitCommit:"e29ce2a54e96cd02ccfce88bee4f58bb6e2a28b6", GitTreeState:"clean", GoVersion:"go1.13.4"}
```

If this command does not successfully return, install Helm 3 for your operating system.  See [Helm Installation](https://helm.sh/docs/intro/install/) for more information.  If you are running Helm 2.x, you will see both a client and server (tiller) portion returned by the version command.  Some of the commands below will also differ slightly for Helm 2.x.

1. Add the Pega repository to your Helm installation.

```bash
$ helm repo add pega https://pegasystems.github.io/pega-helm-charts
```

2. Verify the new repository by searching it.

```bash
$ helm search repo pega
NAME       	            CHART VERSION	APP VERSION	DESCRIPTION
pega/pega  	              1.4.4        	           	Helm chart to configure required installation and deployment configuration settings in your environment for your deployment.
pega/addons	              1.4.4        	1.0        	Helm chart to configure supporting services and tools in your environment for your deployment.
pega/backingservices      1.4.4        	            Helm Chart to provision the latest Search and Reporting Service (SRS) for your Pega Infinity deployment
```

There are three charts available in this repository - addons, backingservices, and pega.

The addons chart installs a collection of supporting services and tools for a Pega deployment. The services you will need to deploy will depend on your cloud environment - for example you may need a load balancer on Minikube, but not for EKS. These supporting services are deployed once per Kubernetes environment, regardless of how many Pega Infinity instances are deployed.

The backingservices chart installs services like 'Search and Reporting Service' (SRS) that you can configure with one or more Pega deployments. You can deploy these backing services in their own namespace; you can isolate the services to a single environment or share them across multiple Pega Infinity environments.

**Example:**
_Single backing service shared across all pega environments:_

backingservice 'Search and Reporting Service' deployed and the service endpoint configured across dev, staging and production pega environments. The service provides isolation of data in a shared setup.

_Multiple backing service deployments:_

You can deploy more than one instance of backing service deployments, in case you want to host a separate deployment of 'Search and Reporting Service' for non-production and production deployments of Pega Infinity. You must configure the appropriate service endpoint using the Pega Infinity deployment values.

3. Download the values file for pega/pega, pega/addons and pega/backingservices.

```bash
$ helm inspect values pega/pega > pega.yaml
$ helm inspect values pega/addons > addons.yaml
$ helm inspect values pega/backingservices > backingservices.yaml
```

4. Edit your values yaml files to specify all required information and customizations for your environment.

* [Instructions to configure the Pega chart](charts/pega/README.md)
* [Instructions to configure the Pega addons](charts/addons/README.md)
* [Instructions to configure the Pega backingservices](charts/backingservices/README.md)

When making customizations for your environment, check the [Pega Platform Support Guide](https://docs.pega.com/bundle/platform/page/platform/deployment/platform-support-guide/platform-support-guide.html) to verify that those changes are supported by your Pega Platform version.

5. Create namespaces for your Pega deployment, backingservices and the addons (if applicable for your environment).

```bash
$ kubectl create namespace mypega
$ kubectl create namespace pegaaddons
$ kubectl create namespace pegabackingservices
```

6. To install the backingservices chart, run the following helm command after configuring your values.yaml file (if applicable for your environment). 

```bash
$ helm install backingservices pega/backingservices --namespace pegabackingservices --values backingservices.yaml
```

7. To install the addons chart, run the following helm command after configuring your values.yaml file (if applicable for your environment). 

```bash
$ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
```

8. With addons and backservices deployed, you are ready to deploy Pega Infinity using the pega chart. Before installing using the chart, it is a good idea to review [Containerized deployments in Kubernetes environments](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/containerized-deployments-kubernetes.html) to understand how Pega deploys as a distributed system. Running a Helm installation using the pega chart installs a Pega Infinity instance into a specified namespace. After you edit the chart with your configuration requirements, run the following command to install the pega chart. 

```bash
$ helm install mypega pega/pega --namespace mypega --values pega.yaml
```

**Important**: To edit the charts and deploy using a local version of the pega/addons, pega/backingservices, or pega/pega charts, you must include the filepath to your local chart directory in your Helm chart reference.

**Tip:** To delete your deployment of Pega nodes, enter the command (this will not delete your database):

```bash
$ helm delete release --namespace mypega
```

# Staying current with a Pega Platform upgrade or patching in zero-downtime

## Upgrades

To upgrade Pega Platform software deployed in a Kubernetes environment with a zero-downtime process, you must do the following:

1. Download the latest Pega software from Pega Digital Software Delivery.
2. Update your repository to use the latest Helm charts and add several parameters to your `pega.yaml` Helm chart.
3. Invoke the upgrade process by using the `helm upgrade release --namespace mypega` command.

For complete details, see the Pega-provided runbook, [Upgrading Pega Platform in your deployment with zero-downtime](/docs/upgrading-pega-deployment-zero-downtime.md).

To upgrade your strategic application, use the latest Upgrade Guide available for your strategic application, which is separate from Pega Platform software. You can locate the appropriate upgrade guide for your installed application from the page, [All Products](https://community.pega.com/knowledgebase/products).

## Patches

To apply a Pega Platform patch with zero downtime to your existing Pega platform software, you must download the latest installer Docker images from Pega Digital Software Delivery and change several options in your Pega Helm chart. For details and helpful resources explaining the Pega Platform patch process, including the Pega Infinity patch policy, see [Pega software maintenance program](https://docs.pega.com/bundle/keeping-current/page/keeping-current/kc/pega-software-maintenance.html). For step-by-step guidance to apply a Pega Platform patch, see the Pega-provided runbook, [Patching Pega Platform in your deployment](/docs/patching-pega-deployment.md).

# Downloading Docker images for your deployment

Clients with appropriate licenses can request access to several required images from the Pega-managed Docker image repository. With your access key, you can log in to the image repository and download these Docker images to install the Pega Platform onto your database. After you pull the images to your local system, you must push them into your private Docker registry.

To download your preferred version of the Pega image to your local system, specify the version tag when by entering:

```bash
$ sudo docker pull pega-docker.downloads.pega.com/platform/pega:<version>

Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/pega:<version>
```

All Docker images for Pega Platform releases that are in Standard Support undergo a nightly rebuild that applies the latest available updates and patches to all third-party components. To take advantage of these updates, you must redeploy your Pega Platform with the latest available images. Pega does not guarantee nightly rebuilds for Pega Platform releases in Extended Support and stops rebuilding images for Pega Platform releases that are out of Extended Support.

For details about downloading and then pushing Docker images to your repository for your deployment, see [Using Pega-provided Docker images](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/pega-docker-images-manage.html).

From Helm chart versions `2.2.0` and above, update your Pega Platform version to the latest patch version.

# Debugging failed upgrades using helm commands

Upgrades using helm charts may fail due to a variety issues, including an invalid configuration, a networking issue, or a platform issue. 
To diagnose the issue or issues, review failure events in the logs and check for a detailed error; after understanding its cause, you can begin troubleshooting the issue or issues.

To help diagnose an issue, you can find the best information by retrieving relevant logs by trying different contexts such as the following 'kubectl log' options.
```kubectl logs <pega-zdt-upgrade-podname> -n <namespace>```

You can use the following option to improve the usefulness of the log output
-f, --follow=false: Specify if the logs should be streamed.

```kubectl logs -f <pod-id> -n <namespace>```

--tail:
Print the last number of lines that you specify in the log file, for example, that last 100 lines in the specified pod. 
By default this option (with no selector specified so tail=-1) displays all lines of the log file

```kubectl logs --tail=100 <pod-id>  -n <namespace>```


If a container is running, you can log in and run the kubectl log command in the container to review the results, or
you can copy the log file in /tmp/foo from a pod to a local temporary directory, /tmp/bar to access the files.
```kubectl cp <some-namespace>/<some-pod>:/tmp/foo /tmp/bar```

# Contributing

This is an open source project and contributions are welcome.  Please see the [contributing guidelines](./CONTRIBUTING.md) to get started.

# Versioning

New versions of this Helm Chart may be released at any time. Versions are defined using [semantic versioning](https://semver.org/):

* Major: Pega introduces new features or functionality that results in breaking backwards compatibility with previous chart versions. Examples:
  * The new values.yaml or maps (config/deploy) cannot be used to deploy previously released Docker image versions.
  * A new, required dependency on a specific Pega Platform version or dependent docker image.
* Minor: Pega adds new functionality and maintains backwards compatibility. Examples:
  * Support for new features while maintaining existing functionality.
  * Support for new, opt-in configurations.
* Patch: Pega fixes bugs and maintains backwards compatibility between minor releases. Examples:
  * Bug fixes or known issue resolutions.
  * Security vulnerability enhancements.

## Helm charts and Docker images compatibility

Both Helm charts and Docker images undergo frequent updates; new Helm chart releases may appear at any time, and the latest patch versions of the Docker images are rebuilt nightly as part of software maintenance (for more information, see [Pega-provided Docker images](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/pega-docker-images-manage.html)). This might result in incompatibility issues during the install and upgrade process. To ensure that Helm charts and Docker images are compatible, do one of the following actions:

* If you update your Helm charts to the latest version to take advantage of new features, update your Docker images to the latest version as well (latest build of the latest patch).
* If you upgrade your Pega Platform to a later version, use the latest Docker image build and latest Helm chart version.
