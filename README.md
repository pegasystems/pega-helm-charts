# Pega Deployment on Kubernetes

This project provides Helm charts and basic examples for deploying Pega on Kubernetes. This project **does not include** the required database installation image which you [may obtain from the Pega Community](https://community.pega.com/knowledgebase/products/platform/deploy).

[![Build Status](https://travis-ci.org/pegasystems/pega-helm-charts.svg?branch=master)](https://travis-ci.org/pegasystems/pega-helm-charts)
![GitHub release](https://img.shields.io/github/release/pegasystems/pega-helm-charts.svg)

## Getting Started

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

`` git clone https://github.com/pegasystems/pega-helm-charts.git ``

Navigate to the project directory and open the values.yaml file.  This is the configuration file that tells Helm what and how to deploy.  For additional documentation covering the different deployment options, see the Pega Community article on [Deploying the Pega Platform by using Kubnernetes](https://community.pega.com/knowledgebase/articles/deploying-pega-platform-using-kubernetes).
