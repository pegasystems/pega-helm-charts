# Constellationui helm chart

The Pega Helm chart is used to deploy an instance of constellationui into a Kubernetes environment.  This readme provides a detailed description of possible configurations and their default values as applicable. 

### Constellationui runtime configuration

The values.yaml file provides configuration options to define the values for the deployment resources of the constellationui service.

### Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enables the constellationui service. Set to true to enable constellationui service in the kubernetes environment.                                                                                                                                                                                                               |
| `image`                        | Specify the image repository and the version.                                                                                                                                                                                                                                                                                          |

Example:

```yaml
image: pega-docker.downloads.pega.com/constellationui/service:8.7.3-ga-37
# log level : error, warn, info, debug.  use error for production
logLevel: info
enabled: true
replicas: 2
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