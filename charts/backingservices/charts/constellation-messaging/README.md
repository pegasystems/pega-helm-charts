# Constellation Messaging Service

The optional Messaging service acts as a middle man between information publishers and information subscribers; accepting simple http publish of information to forward to 1000's of websocket subscribers. The publishers are typically the Infinity case engine or 3rd party integration services, and the subscribers are UI components in browsers running Constellation UI.

Once the service routing (with TLS) is set up, configure the Pega Infinity ConstellationMessageSvcHostPath DSS to your service URL (e.g.  yourhostname.com/c11n-messaging). Do not include a protocol for this setting.

Only a single Messaging Service deployment is necessary to support an entire organization.  Do not install the service in every namespace or for every application or project.

Complete information on the design of the service including architecture, scalability, reliability, operations and troubleshooting is available at [https://documents.constellation.pega.io/messaging/introduction.html](https://documents.constellation.pega.io/messaging/introduction.html).

## Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enable the Messaging Service deployment as a backing service. Set this parameter to `true` to deploy the service.                                                                                                                                                                                                                |
| `provider`                              | Enter your Kubernetes provider. Accepted values are aws, gke or k8s.  |
| `name`                                  | Deprecated, use `deployment.name`. Specify the name of your messaging service. Your deployment creates resources prefixed with this string.                                                                                                                                                                                                                                                                                                                                                          |
| `deployment.name`                        | Specify the name of your messaging service. Your deployment creates resources prefixed with this string.                                                                                                                                                                                                                                                                                                                                                          |
| `imagePullSecretNames`                            | Deprected, use `docker.imagePullSecretNames`. List pre-existing secrets to be used for pulling docker images.                                                                                                                                                                                                                                                                                                           |
| `affinity`                            | Define pod affinity so that it is restricted to run on particular node(s), or to prefer to run on particular nodes.                                                                                                                                                                                                                                       |
| `docker.imagePullSecretNames`                            | List pre-existing secrets to be used for pulling docker images.                                                                                                                                                                                                                                                                                                           |
| `docker.registry.url`                        | Specify the image registry url.                                                                                                                                                                                                                                                                                          |
| `docker.registry.username`                        | Specify the username for the docker registry.                                                                                                                                                                                                                                                                                          |
| `docker.registry.password`                        | Specify the password for the docker registry.                                                                                                                                                                                                                                                                                          |
| `docker.messaging.image`                        | Specify the image version.                                                                                                                                                                                                                                                                                          |
| `docker.messaging.imagePullPolicy`                        | Specify the image pull policy configurations for the image.                                                                                                                                                                                                                                                                                          |
| `pegaMessagingPort`                            | Deprecated, use `service.port`. Defines the port used by the Service.                                                                                                                                                                                                                                                                                                           |
| `service.port`                            | Defines the port used by the Service.                                                                                                                                                                                                                                                                                                           |
| `pegaMessagingTargetPort`                            | Deprecated, use `service.targetPort`. Defines the port used by the Pod and Container.                                                                                                                                                                                                                                                                                                           |
| `service.targetPort`                            | Defines the port used by the Pod and Container.                                                                                                                                                                                                                                                                                                           |
| `service.serviceType`                        | The [type of service](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) you wish to expose.                                                                                                                                                                                                                                               |
| `service.annotations`                        | Optionally add custom annotations for advanced configuration. Specifying a custom set of annotations will result in them being used instead of the default configurations.                                                                                                                                                                                                                                               |
| `ingress.enabled`                        | Set to true in order to deploy an ingress. Due to the diversity of network configurations, ingress vendors, and TLS requirements it may be necessary to define your ingress separately from this chart.                                                                                                                                                                                                                                              |
| `ingress.ingressClassName`                        | Ingress class to be used in place of the annotation.                                                                                                                                                                                                                                               |
| `ingress.tls.enabled`                        | Specify the use of HTTPS for ingress connectivity. If the tls block is omitted, TLS will not be enabled.                                                                                                                                                                                                                                               |
| `ingress.tls.secretName`                        | Specify the Kubernetes secret you created in which you store your SSL certificate for your deployment.                                                                                                                                                                                                                                               |
| `ingress.annotations`                        | Specify additional annotations to add to the ingress.                                                                                                                                                                                                           |
| `ingress.domain`                        | Specify your custom domain.                                                                                                                                                                                                                                                                                         |

```yaml
enabled: true
deployment:
  name: "constellation-messaging"
# Cloud provider details. Accepted values are : aws
provider: aws

# Docker repos and tag for image
docker:
  # If using a custom Docker registry, supply the credentials here to pull Docker images.
  registry:
    url: YOUR_REGISTRY_URL_HERE
    username: YOUR_REGISTRY_USERNAME_HERE
    password: YOUR_REGISTRY_PASSWORD_HERE
  # Docker image information for the Pega docker image, containing the application server.
  constellation:
    image: pega-docker.downloads.pega.com/constellation-messaging/docker-image:5.4.0
    imagePullPolicy: Always

urlPath: /c11n-messaging
replicas: 1

```

### Liveness and readiness probes

Constellation messaging service uses liveness and readiness to determine application health in your deployments. For an overview of these probes, see [Configure Liveness and Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/). Configure a probe for *liveness* to determine if a Pod has entered a broken state; configure it for *readiness* to determine if the application is available to be exposed. If not explicitly configured, default probes are used during the deployment. Set the following parameters as part of a `livenessProbe` or `readinessProbe` configuration.

Notes:
* `timeoutSeconds` cannot be greater than `periodSeconds` in some GCP environments. For details, see [this API library from Google](https://developers.google.com/resources/api-libraries/documentation/compute/v1/csharp/latest/classGoogle_1_1Apis_1_1Compute_1_1v1_1_1Data_1_1HttpHealthCheck.html#a027a3932f0681df5f198613701a83145).

Parameter             | Description    | Default `livenessProbe` | Default `readinessProbe`
---                   | ---            | ---                     | ---
`initialDelaySeconds` | Number of seconds after the container has started before probes are initiated. | `5` | `5`
`timeoutSeconds`      | Number of seconds after which the probe times out. | `5` | `5`
`periodSeconds`       | How often (in seconds) to perform the probe. | `30` | `30`
`successThreshold`    | Minimum consecutive successes for the probe to be considered successful after it determines a failure. | `1` | `1`
`failureThreshold`    | The number consecutive failures for the pod to be terminated by Kubernetes. | `3` | `3`

Example:

```yaml
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
