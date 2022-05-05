# BackingServices Helm chart

The Pega Infinity backing service is a feature which you can deploy as an independent service module. For example `Search and Reporting Service` or `SRS` backing service can replace the embedded search feature of Pega Infinity Platform. To use it in your deployment, you provision and deploy it independently as an external service which provides search and reporting capabilities with a Pega Infinity environment.

The backingservices chart supports deployment options for Search and Reporting Service (SRS). You configure this SRS into the `pega` namespace for your Pega Infinity deployment.

## Configuring a backing service with your pega environment

You can provision this SRS into your `pega` environment namespace, with the SRS endpoint configured with the Pega Infinity environment. When you include the SRS into your pega namespace, the service endpoint is included within your Pega Infinity environment network to ensure isolation of your application data.

## Search and Reporting Service support

The Search and Reporting Service provides next generation search and reporting capabilities for Pega Infinity 8.6 and later.

This service replaces the legacy search module from the platform with an independently deployable and scalable service along with the built-in capabilities to support more than one Pega environments with its data isolation features in Pega Infinity 8.6 and later.
The service deployment provisions runtime service pods along with a dependency on a backing technology ElasticSearch service for storage and retrieval of data.

### SRS Version compatibility matrix

| Pega Infinity version | SRS version | ElasticSearch version | Description                                                                                                                                                                                                                                                                                                           |
|-----------------------|-------------|-----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| < 8.6                 | NA          | NA                    | SRS can be used with Pega Infinity 8.6 and later                                                                                                                                                                                                                                                                      |
| \>= 8.6               | \>=1.17.10  | 7.10.2 or 7.16.3      | Pega Infinity 8.6 and later supports using a Pega-provided platform-services/search-n-reporting-service Docker Image. While all SRS Docker image versions starting at 1.12.0 are certified against Elasticsearch versions 7.10.2 and 7.16.3, Pega recommends using the latest available SRS image - 1.17.10 or later. |

### SRS runtime configuration

The values.yaml provides configuration options to define the deployment resources along with option to either provision ElasticSearch cluster automatically for data storage, or you can choose to configure an existing externally managed elasticsearch cluster to use as a datastore with the SRS runtime.

If an externally managed elasticsearch cluster is being used, make sure the service is accessible to the k8s cluster where SRS is deployed.

You may enable the component of [Elasticsearch](https://github.com/helm/charts/tree/master/stable/elasticsearch/values.yaml) in the backingservices by configuring the 'srs.srsStorage' section in values.yaml file to deploy ElasticSearch cluster automatically. For more configuration options available for each of the components, see their Helm Charts.

Note: Pega does **not** actively update the elasticsearch dependency in `requirements.yaml`. To leverage SRS, you must do one of the following:

* To use the internally-provided Elasticsearch service in the SRS cluster, use the default `srs.enabled.true` parameter and set the Elasticsearch version by updating the `elasticsearch.imageTag` parameter in the [values.yaml](./values.yaml) to match the `dependencies.version` parameter in the [requirements.yaml](./requirements.yaml).
* To use an externally-provided Elasticsearch service from the SRS cluster, update the `srs.srsStorage.provisionInternalESCluster` parameter in the [values.yaml](./values.yaml) to `false` and then provide connection details as documented below.

### Deploying SRS with Pega-provided busybox images
To deploy Pega Platform with the SRS backing service, the SRS helm chart requires the use of the busybox image.  For clients who want to pull this image from a registry other than Docker Hub, they must tag and push their image to another registry, and then pull it by specifying `busybox.image` and `busybox.imagePullPolicy`.


### Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enable the Search and Reporting Service deployment as a backing service.                                                                                                                                                                                                                                                                                                                                                               |
| `deploymentName`                        | Specify the name of your SRS cluster. Your deployment creates resources prefixed with this string. This is also the service name for the SRS.                                                                                                                                                                                                                                                                                          |
| `srsRuntime`                            | Use this section to define specific resource configuration options like image, replica count, cpu and memory resource settings in the SRS.                                                                                                                                                                                                                                                                                             |
| `elasticsearch`                         | Define the elasticsearch cluster configurations. The [Elasticsearch](https://github.com/helm/charts/tree/master/stable/elasticsearch/values.yaml) chart defines the values for Elasticsearch provisioning in the SRS cluster. For internally provisioned Elasticsearch the default version is set to `7.10.2`. Set the `elasticsearch.imageTag` parameter in values.yaml to `7.16.3` to use this supported version in the SRS cluster. |
| `srsStorage.provisionInternalESCluster` | <ul><li>Set to `true` to enable this parameter to provision an internally managed and secured Elasticsearch cluster to be used with the SRS cluster; this Requires you to run `$ make es-prerequisite NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT>`.</li><li>Set to `false` to disable this parameter to use your own Elasticsearch service from the SRS cluster.</li></ul>                                                               |
| `busybox`                               | When provisioning an internally managed Elasticsearch cluster, you can customize the location and pull policy of the BusyBox image used during the deployment process by specifying `busybox.image` and `busybox.imagePullPolicy`.                                                                                                                                                                                                     |

Example:

```yaml
srs:
  enabled: true
  deploymentName: "YOUR_SRS_DEPLOYMENT_NAME"

  busybox:
    image: "busybox:1.31.0"
    imagePullPolicy: "IfNotPresent"

  srsRuntime:
    #srs-service values
    replicaCount: 2
    srsImage: "YOUR_SRS_IMAGE:TAG"
    imagePullPolicy: "IfNotPresent"
    resources:
      limits:
        cpu: 1300m
        memory: "2Gi"
      requests:
        cpu: 650m
        memory: "2Gi"
    serviceType: "ClusterIP"
    env:
      #AuthEnabled may be set to true when there is an authentication mechanism in place between SRS and Pega Infinity.
      AuthEnabled: false
      PublicKeyURL:

  srsStorage:
    # To configure authentication for the externally managed Elasticsearch cluster to use Basic Authentication then uncomment
    # add the parameters, srs.srsStorage.basicAuthentication.username and srs.srsStorage.basicAuthentication.password
    #    basicAuthentication:
    #      username: "BASIC_AUTH_USERNAME"
    #      password: "BASIC_AUTH_PASSWORD"
    # To configure authentication for the externally managed Elasticsearch cluster to use an AWS IAM Role then uncomment and
    # add the parameters, srs.srsStorage.awsIAM and srs.srsStorage.awsIAM.region
    # srs.srsStorage.awsIAM's srs.srsStorage.awsIAM.region value
    #    awsIAM:
    #      region: "AWS_ELASTICSEARCH_REGION"

    # set `requireInternetAccess` to true when the elasticsearch domain is outside of the Kubernetes cluster network  and is available over internet
#    requireInternetAccess: true
```
