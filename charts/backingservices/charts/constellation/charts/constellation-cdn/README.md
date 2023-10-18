#Constellation Pega UI static content CDN

The mechanism for delivery of the Pega generated UI static content is a Pega-provided Content Delivery Network (CDN). This has many advantages:

* A single url is the only information required.
* World-wide regional caching, giving the same fast performance everywhere.
* No customer install or setup is necessary.

Most clients should configure the Pega Infinity DSS ConstellationPegaStaticURL to https://release.constellation.pega.io and do not need to deploy this chart.

For client cloud customers who have situations with no internet access, the content that we publish to the CDN is available as a Docker image that can be installed as a docker or K8s service. **A local install brings additional operational and maintenance costs to the customer. Also Constellation hotfix updates are not available for local installs. Local installs should be avoided where possible.**


##Overview

The CDN content is available packaged into an nginx webserver in a docker image. A seperate image is released for each Infinity major.minor.patch. This image can be installed and run as a containerised webserver on the customers network. The url for the containerised service is an exact replacement for the Pega CDN. The approach is similar to the App-Static install, and the knowledge and skill prerequisites are the same:

Familiarity with Docker (images, containers, start, stop, background execution, logs, ports, repos) is a prerequisite to trying to install.
For simple Docker installs, a good https certificate that matches the domain the image is deployed on is required.
For K8s or Helm, experience with K8s, cluster admin and network configuration are prerequisites to trying to install
For on-prem installation, understanding of, and access to the corporate network is required.
Download of the docker image from Pega downloads is a prerequisite to deployment.

##Details

* The service should be installed **once** in the customer's network.  Do not install in every namespace or for every application.
* It must be on a URL accessible from all browsers used for Constellation UI.
* It should be on the same domain as Infinity to avoid CORS errors.
* The https certificate used for the service must match the domain the service is on.
* This is read-only content. The image only supports http GET. There is no POST, PUT, DELETE, OPTIONS support.
* To simplify certificate management, TLS certificates must be applied to the load-balancer in front of the webserver.
* The image can be run as a docker container, or K8s pod. This is a very simple read-only service. There is no need for fail-over etc.

Most likely there will be multiple Infinity versions in use across the customer organisation. For each version of Infinity running Constellation UI, a matching Constellation CDN service must be installed. Our examples show how to setup the network routing for multiple versions.

### Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enable the CDN Service deployment as a backing service. Set this parameter to `true` to deploy CDNs.                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `name`                        | Specify the name of your CDN. Your deployment creates resources prefixed with this string.                                                                                                                                                                                                                                                                                                                                                         |
| `pegaStaticPort`                            | Defines the port used by the Service.                                                                                                                                                                                                                                                                                                           |
| `pegaStaticTargetPort`                            | Defines the port used by the Pod and Container.                                                                                                                                                                                                                                                                                                           |
| `ingress`                               |  Allows optional configuration of a domain name, ingressClass, and annotations.  An ingress will be provisioned if a domain name is supplied.  Due to the diversity of network configurations, ingress vendors, and TLS requirements it may be necessary to define your ingress separately from this chart.
| `deployments`                         | Defines the list of CDN images to be deployed.  One deployment is needed for each patch release of Pega supported in the organization. |