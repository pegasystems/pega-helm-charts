# Constellation messaging service helm chart

The Pega Helm chart is used to deploy an instance of constellation messaging service into a Kubernetes environment.  This readme provides a detailed description of possible configurations and their default values as applicable. 

### Constellation messaging service runtime configuration

The values.yaml file provides configuration options to define the values for the deployment resources of the constellation messaging service.

### Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enables the constellation messaging service. Set to true to enable constellation messaging service in the kubernetes environment.                                                                                                                                                                                                               |
| `image`                        | Specify the image repository and the version.                                                                                                                                                                                                                                                                                          |
| `global.imageCredentials.registry`                        | Specify the docker image registry.                                                                                                                                                                                                                                                                                          |
| `global.imageCredentials.username`                        | Specify the username for the registry.                                                                                                                                                                                                                                                                                          |
| `global.imageCredentials.password`                        | Specify the password for the image registry.                                                                                                                                                                                                                                                                                          |
Example:

```yaml
global:
  imageCredentials:
    registry: "YOUR_DOCKER_REGISTRY"
    username: "YOUR_DOCKER_REGISTRY_USERNAME"
    password: "YOUR_DOCKER_REGISTRY_PASSWORD"

enabled: true
deploymentName: c11n-messaging

c11nMessagingRuntime:
  replicaCount: 2
  c11nMessagingImage: pega-docker.downloads.pega.com/constellation-messaging/docker-image:0.0.5-20220225103826317
  imagePullPolicy: Always
  resources:
    limits:
      cpu: 1300m
      memory: "2Gi"
    requests:
      cpu: 650m
      memory: "2Gi"
  serviceType: "LoadBalancer"
```