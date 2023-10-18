# Constellation UI setup

Please refer to 
 > https://documents.constellation.pega.io/static/88/introduction.html

for instructions related to the pegastatic content delivery setup. Once that is complete please proceed with the instructions below for the constellation appstatic service setup.

## ConstellationUI helm chart

The ConstellationUI helm chart is used to deploy an instance of the constellation appstatic service into a Kubernetes environment. The service is capable of serving multiple client environments. We encourage having a single deployment for all your systems. The following readme provides a detailed description of the parameter configurations and their default values as applicable. 

#### Prerequisites
#### AWS Cloud 
1. The ConstellationUI helm charts for AWS provider use an application load balancer to expose the service. Before installing the constellationui helm chart please follow the prerequisites section of the following AWS documentation and make sure your cluster is configured accordingly. 

2. If you are using a custom domain please make sure you have a valid certificate imported or created in ACM.

> https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html

#### Downloading Docker images for your deployment

Clients with appropriate licenses can request access to several required images from the Pega-managed Docker image repository. With your access key, you can log in to the image repository and download these Docker images to install the constellation appstatic service.

To download the constellation appstatic service image, specify the version and use the following command

```bash
$ sudo docker pull pega-docker.downloads.pega.com/constellation-appstatic-service/docker-image:xxxxxxx

Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/constellation-appstatic-service/docker-image:xxxxxxx
```

For details about downloading and then pushing Docker images to your repository for your deployment, see [Using Pega-provided Docker images](https://docs.pega.com/client-managed-cloud/87/pega-provided-docker-images).


#### Constellationui runtime configuration

The values.yaml file provides configuration options to define the values for the deployment resources of the constellation appstatic service.

#### Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enables the constellation appstatic service. Set to true to enable constellation appstatic service in the kubernetes environment.                                                                                                                                                                                                               |
| `cloudProvider`                        | Specify the cloud provider details. Accepted values are aws.                                                                                                                                                                                                                                                                                          |
| `awsCertificateArn`                        | Specify the arn for the AWS ACM certificate.                                                                                                                                                                                                                                                                                          |
| `domainName`                        | Specify your custom domain.                                                                                                                                                                                                                                                                                          |
| `ingressAnnotations`                        | Specify additional annotations to add to the ingress.                                                                                                                                                                                                                                                                                          |
| `customerAssetVolumeClaimName`                        | Specify the volume claim name to be used for storing customer assets.                                                                                                                                                                                                                                                                                          |
| `imagePullSecretNames`                        | Specify a list of existing ImagePullSecrets to be added to the Deployment.                                                                                                                                                                                                                                                                                          |
| `docker`.`registry`.`url`                        | Specify the image registry url.                                                                                                                                                                                                                                                                                          |
| `docker`.`registry`.`username`                        | Specify the username for the docker registry.                                                                                                                                                                                                                                                                                          |
| `docker`.`registry`.`password`                        | Specify the password for the docker registry.                                                                                                                                                                                                                                                                                          |
| `docker`.`constellation`.`image`                        | Specify the image version.                                                                                                                                                                                                                                                                                          |

Example:

```yaml
enabled: true
deployment:
  name: "constellation"
# Cloud provider details. Accepted values are : aws
cloudProvider: aws
# For aws cloud provider enter your acm certificate ARN here.
awsCertificateArn : arn:aws:acm:us-west-2:xxxxx:certificate/xxxxxxx
domainName: YOUR_CUSTOM_DOMAIN_HERE
# Docker repos and tag for image
docker:
    # If using a custom Docker registry, supply the credentials here to pull Docker images.
    registry:
      url: YOUR_REGISTRY_URL_HERE
      username: YOUR_REGISTRY_USERNAME_HERE
      password: YOUR_REGISTRY_PASSWORD_HERE
    # Docker image information for the Pega docker image, containing the application server.
    constellation:
      image: pega-docker.downloads.pega.com/constellation-appstatic-service/docker-image:xxxxxxx
logLevel: info
urlPath: /c11n
replicas: 1
livenessProbe:
  initialDelaySeconds: 5
  timeoutSeconds: 5
  periodSeconds: 30
  successThreshold: 1
  failureThreshold: 3
readinessProbe:
  initialDelaySeconds: 5
  timeoutSeconds: 5
  periodSeconds: 30
  successThreshold: 1
  failureThreshold: 3
```

#### Constellationui helm charts installation steps

1. Create a namespace into which you want to deploy the constellation appstatic service
    > kubectl create namespace <<namespace_name_here>>

2. Install the constellation helm charts using the following command 
    > helm install path-to-pega-helm-charts/charts/backingservices/charts/constellation --values path-to-pega-helm-charts/backingservices/charts/constellation/values.yaml -n <<namespace>> --generate-name

3. Once the charts are deployed an application load balancer should be created automatically. You can then configure the records of your custom domain to route to the loadbalancer accordingly. For AWS a simple Route53 record in the correct hosted zone of your custom domain should do it. 

