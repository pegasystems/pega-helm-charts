Pega Helm Chart
===========

Creates a helm chart for PRPC.

# Quickstart

All below examples assume you're executing from within the `src/main/helm/pega` folder.

To install this chart invoke the following:
```
helm install . -n mypega --namespace myproject --values /home/user/my-overridden-values.yaml
```

To delete this chart invoke the following:
```
helm delete mypega --purge
```
# Usability
Following actions are supported for this helm chart.
```
actions:
  # install - install the pega platform database.
  # deploy - deploy the full pega cluster.
  # install-deploy - installation followed by complete pega cluster deployment.
  # upgrade - Upgrades the pega platform database to the specified higher version.
  # upgrade-deploy - Upgrades the pega platform database followed by rolling updates of existing deployments.
  ```
 **Use  Case 1.** User who wants to perform only the installation has to provide the action [***install***] in values.yaml.
```
actions:
  execute: "install"
```   
 **Use  Case 2.** User who wants to deploy the complete pega cluster has to provide the action [***deploy***] in values.yaml.
```
actions:
  execute: "deploy"
```   
 **Use  Case 3.** User who wants to perform the installation followed by pega cluster deployment has to provide the action [***install-deploy***] in values.yaml.
```
actions:
  execute: "install-deploy"
```   
**Use  Case 4.** User who wants to perform only the upgrade to some higher version has to provide the action [***upgrade***] in values.yaml.
```
actions:
  execute: "upgrade"
```   
**Use  Case 5.** User who wants to perform upgrade followed by rolling updates of the existing deployments has to provide the action [***upgrade-deploy***] in values.yaml.
```
actions:
  execute: "upgrade-deploy"
```   
   
# On Building

The following command will build the chart @ helm/build/chart/pega-<version>.tgz
```
./gradlew -PincludeOnly=infrastructure/distribution/helm build
```
# Provider Notes

### kubernetes

N/A

### openshift

Before launching the provided template yaml, a small configuration change will have to be made to use the public Cassandra, metrics-server and Elastic Search images.
The image expects to be run as the cassandra user (UID: 999), or for elastic search (UID:1000) and the default UID range only allows 1000120000-1000129999

To remedy this:
1. Login as the system admin user: `$ oc login -u system:admin`
2. Edit the "restricted" scc: `$ oc edit scc restricted`
    1. Find the `allowPrivilegedContainer` setting and set its **value** to `true`
    2. Find the `RunAsUser` setting and set its **type** to `RunAsAny`
3. Add the "anyuid" scc as cassandra needs to run as root: `$ oc adm policy add-scc-to-user anyuid system:serviceaccount:<namespace>:default`

### eks

N/A

### gke

N/A

### pks

Pivotal Kubernetes is built atop Google Kubernetes Engine (gke) so anything which works there should work here as well.

### aks

N/A
