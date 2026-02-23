# Search and Reporting Service Helm chart

The Pega `Search and Reporting Service` or `SRS` backing service provides the search and reporting capabilities of Pega Infinity Platform. To use it in your deployment, you provision and deploy it independently as an external service in a Pega Infinity environment.

## Configuring a backing service with your pega environment

You can provision this SRS into your `pega` environment namespace, with the SRS endpoint configured with the Pega Infinity environment. When you include the SRS into your pega namespace, the service endpoint is included within your Pega Infinity environment network to ensure isolation of your application data.

## Search and Reporting Service support

The Search and Reporting Service provides next generation search and reporting capabilities for Pega Infinity 8.6 and later.

This service replaces the legacy search module from the platform with an independently deployable and scalable service along with the built-in capabilities to support more than one Pega environments with its data isolation features in Pega Infinity 8.6 and later.
The service deployment provisions runtime service pods along with a dependency on a backing technology ElasticSearch service for storage and retrieval of data.

### SRS Version compatibility matrix

<table>
    <thead>
        <tr>
            <th>Pega Infinity version</th>
            <th>SRS version</th>
            <th>Docker image</th>
            <th>Authentication</th>
            <th>Certified Elasticsearch/OpenSearch version</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>< 8.6</td>
            <td>NA</td>
            <td>NA</td>
            <td>NA</td>
            <td>NA</td>
            <td>SRS can be used with Pega Infinity 8.6 and later.</td>
        </tr>
        <tr>
            <td rowspan=5> >= 8.6 </td>
            <td rowspan=5>1.42.3</td>
            <td rowspan=3> <b>search-n-reporting-service</b></td>
        </tr>
        <tr>
            <td>Not enabled</td>
            <td>Elasticsearch 7.17.9, 7.17.29</td>
            <td>As a best practice, use Elasticsearch version 7.17.29. <b> Deployments without authentication are not recommended for production environments. </b> </td>
        </tr>
        <tr>
            <td>Enabled</td>
            <td>Elasticsearch 7.17.9, 7.17.29, 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 & 8.19.11</td>
            <td>As a best practice, use Elasticsearch version 8.19.11.</td>
        </tr>
        <tr>
            <td rowspan=1> <b>search-n-reporting-service-os</b></td>
            <td>Enabled</td>
            <td><ul><li>Elasticsearch 7.10 on AWS OpenSearch service</li><li>OpenSearch 1.3 </li><li>OpenSearch 2.15</li><li>OpenSearch 2.19</li></ul></td>
            <td> The <b>search-n-reporting-service-os</b> replaced the earlier <b>search-n-reporting-service-aws</b> image.
The image supports all OpenSearch deployment options (cloud or self-managed).
As a best practice, use OpenSearch 2.19. </td>
        </tr>
    </tbody>
</table>

**Note:** Elasticsearch 7.10.2 and 7.16.3 reached end of vendor support and are no longer supported. Elasticsearch 7.17.9 will reach vendor end of support on Jan 15, 2026.

**Important:** Pega supports only official Elasticsearch and OpenSearch Docker images. Custom images, for example, `bitnami/elasticsearch`, are not supported.

### If your deployment uses the internally-provisioned Elasticsearch: ###
To migrate to Elasticsearch version 7.17.9, 7.17.29, 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 or 8.19.11 from the Elasticsearch version 7.10.2 or 7.16.3, perform the following steps:
1. Update the SRS Docker image version to use v1.31.2. This version has backward compatibility with Elasticsearch versions 7.10.x and 7.16.x, so your SRS will continue to work even before you update your Elasticsearch service.
2. To update Elasticsearch version to 7.17.9 perform the following actions:
   * Update the Elasticsearch `dependencies.version` parameter in the [requirement.yaml](../../requirements.yaml) to 7.17.3.

     Note: This parameter references the Elasticsearch Helm chart version and not the Elasticsearch cluster version.
   * Update the elasticsearch.imageTag in the Backing Services Helm chart to 7.17.9.
3. To update Elasticsearch version to 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 or 8.19.11, perform the following actions:
   * Update the Elasticsearch `dependencies.version` parameter in the [requirement.yaml](../../requirements.yaml) to 8.5.1.

     Note: This parameter references the Elasticsearch Helm chart version and not the Elasticsearch cluster version.
   * Update the elasticsearch.imageTag in the Backing Services Helm chart to 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 or 8.19.11.
4. Restart the SRS pods.

### If your deployment connects to an externally-managed Elasticsearch service: ###
To migrate to Elasticsearch version 7.17.9, 7.17.29, 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 or 8.19.11 from the Elasticsearch version 7.10.2 or 7.16.3, perform the following steps:
1. Update the SRS Docker image version to use v1.31.2. This version has backward compatibility with Elasticsearch versions 7.10.x and 7.16.x, so your SRS will continue to work even before you update your Elasticsearch service.
2. To use Elasticsearch version 7.17.9, upgrade your external Elasticsearch cluster to 7.17.9 according to your organization’s best practices. For more information, see official Elasticsearch version 7.17 documentation.
3. To use Elasticsearch version 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 or 8.19.11, upgrade your external Elasticsearch cluster to 8.10.3, 8.15.1, 8.15.5, 8.18.2, 8.18.3, 8.18.6 or 8.19.11 according to your organization’s best practices. For more information, see official Elasticsearch version 8.x documentation.
4. Restart the SRS pods.

### SRS runtime configuration

The values.yaml provides configuration options to define the deployment resources along with option to either provision ElasticSearch cluster automatically for data storage, or you can choose to configure an existing externally managed elasticsearch cluster to use as a datastore with the SRS runtime.

If an externally managed elasticsearch cluster is being used, make sure the service is accessible to the k8s cluster where SRS is deployed.

You may enable the component of [Elasticsearch](https://github.com/helm/charts/tree/master/stable/elasticsearch/values.yaml) in the backingservices by configuring the 'srs.srsStorage' section in values.yaml file to deploy ElasticSearch cluster automatically. For more configuration options available for each of the components, see their Helm Charts.

Note: Pega does **not** actively update the elasticsearch dependency in `requirements.yaml`. To leverage SRS, you must do one of the following:

* To use the internally-provided Elasticsearch service in the SRS cluster, use the default `srs.enabled.true` parameter and set the Elasticsearch version by updating the `elasticsearch.imageTag` parameter in the [values.yaml](./values.yaml) to match the `dependencies.version` parameter in the [requirements.yaml](../../requirements.yaml). This method streamlines the deployment process for development and testing environments, but it is not suitable for production environments, which require a fully external Elasticsearch cluster. Additionally, even though you deploy SRS and Elasticsearch together, Pega does not license the Elasticsearch cluster deployed using this method and does not maintain it as part of the Pega Platform support. <b>Note: You cannot use OpenSearch for the internally provisioned cluster.</b>
* To use an externally-provided Elasticsearch service with SRS, use the default `srs.enabled.true` parameter, update the `srs.srsStorage.provisionInternalESCluster` parameter in the [values.yaml](./values.yaml) to `false` and then provide connection details as documented below. This is the recommended method and is suitable for production environments.

### Deploying SRS with Pega-provided busybox images
To deploy Pega Platform with the SRS backing service, the SRS helm chart requires the use of the busybox image.  For clients who want to pull this image from a registry other than Docker Hub, they must tag and push their image to another registry, and then pull it by specifying `busybox.image` and `busybox.imagePullPolicy`.

### Configuration settings

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
|-----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`                               | Enable the Search and Reporting Service deployment as a backing service. Set this parameter to `true` to use SRS.                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `deploymentName`                        | Specify the name of your SRS cluster. Your deployment creates resources prefixed with this string. This is also the service name for the SRS.                                                                                                                                                                                                                                                                                                                                                          |
| `srsRuntime`                            | Use this section to define specific resource configuration options like image, replica count, pod affinity, cpu and memory resource settings in the SRS. The default minimum required number of replicas is 2, but as a best practice, deploy 3 replicas to maintain high availability.                                                                                                                                                                                                                                                                                                                                                            |
| `busybox`                               | When provisioning an internally managed Elasticsearch cluster, you can customize the location and pull policy of the Alpine image used during the deployment process by specifying `busybox.image` and `busybox.imagePullPolicy`.                                                                                                                                                                                                                                                                      |
| `elasticsearch`                         | Define the elasticsearch cluster configurations. The [Elasticsearch](https://github.com/helm/charts/tree/master/stable/elasticsearch/values.yaml) chart defines the values for Elasticsearch provisioning in the SRS cluster. For internally provisioned Elasticsearch the default version is set to `7.17.9`. Set the `elasticsearch.imageTag` parameter in values.yaml to `7.16.3` to use this supported version in the SRS cluster.                                                                |
| `k8sProvider`                           | Specify your Kubernetes provider name. Supported values are [`eks`, `aks`, `minikube`, `gke`, `openshift`, `pks`].
| `enableSecureCryptoMode`                | Set to true if you require a highly secured connection that complies with NIST SP 800-53 and NIST SP 800-131. Otherwise, set to false.
| `javaOpts`                              | Use this parameter to configure values for Java options.

### Enabling security between SRS and Elasticsearch
Enabling a secure connection between SRS and your Elasticsearch service depends on the method you chose to deploy the Elasticsearch cluster.
To configure a secure connection between the SRS cluster and internally provisioned Elasticsearch, configure the following parameters.

| Configuration                            | Usage                                                                                                                                                                                                                                            |
|------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `tls`                                    | Set to `true` to enable the SRS service to authenticate to your organization's available Elasticsearch service.                                                                                                                                  |
| `srsStorage.provisionInternalESCluster`  | <b>Note: You cannot use OpenSearch for the internally provisioned cluster.</b><br><br><ol><li>Set the `srsStorage.provisionInternalESCluster` parameter to `true` to provide an internally managed and secured Elasticsearch cluster.</li><li>In the [requirements.yaml](../../requirements.yaml) file, set the `dependencies.version` parameter to the same version you configured for the `elasticsearch.imageTag` version in the Backing Services Helm chart [values.yaml](../../values.yaml) file.</li><li>From the Backing Services Helm chart directory in your environment, run the following command to create your Elasticsearch certificates and pass them to secrets: <p>`$ make es-prerequisite NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT> ELASTICSEARCH_VERSION=<ELASTICSEARCH_VERSION>`</p><p>Where `NAMESPACE` references your deployment namespace of the SRS cluster and `ELASTICSEARCH_VERSION` matches the Elasticsearch version you want to use in [values.yaml](../../values.yaml) and [requirements.yaml](../../requirements.yaml).</p></li></ol> |

To configure a secure connection(tls) between SRS and an external Elasticsearch cluster, configure the following parameters.

| Configuration                           | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
|-----------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `tls`                                   | Set to `true` to enable the SRS service to authenticate to your organization's available Elasticsearch service.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `certificateName`                       | Enter the tls certificate name. Default certificate name will be "elastic-certificates.p12" if not used.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `certificatePassword`                   | Enter the tls certificate password if any. Default value will be empty if not used.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |                                             
| `certsSecret`                | To specify a certificate using a secret, uncomment the certsSecret parameter and provide the secret name containing your certificate and certificate password. Use the full name of the certificate file (together with file extension, for example, “certificate.p12” or“certificate.jks”) as a key name in the secret. Use this key name to configure the “certificateName”parameter.Use a key name “password” to provide the certificate password in the secret. Defaults to "srs-certificates".|
| `authSecret`                | Specify the secret with your Elasticsearch credentials. Use “username” and “password” as keys for your secret.This parameter applies to both basic authentication and TLS-based authentication. Defaults to "srs-elastic-credentials".|
| `esCredentials.username`                | Enter the username for your available Elasticsearch service. This username value must match the values you set in the connection info section of esCredentials. <br><b>Note:</b> This parameter will be deprecated in future releases, so as a best practice, use `authCredentials.username`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| `esCredentials.password`                | Enter the required password for your available Elasticsearch service. This password value must match the values you set in the connection info section of esCredentials. <br><b>Note:</b> This parameter will be deprecated in future releases, so as a best practice, use `authCredentials.password`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      
| `authCredentials.username`                | Enter the username for your available Elasticsearch/OpenSearch service. This username value must match the values you set in the connection info section of authCredentials.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `authCredentials.password`                | Enter the required password for your available Elasticsearch/OpenSearch service. This password value must match the values you set in the connection info section of authCredentials.     |
| `srsStorage.provisionInternalESCluster` | <ol><li>Set the `srsStorage.provisionInternalESCluster` parameter to `false` to disable the internally provisioned Elasticsearch cluster and connect to your available external Elasticsearch service.</li><li>To secure the connection between SRS and your external Elasticsearch service, you must provide the appropriate TLS certificates in an accessible location, for example, /home/certs.</li><li>To pass the required certificates to the cluster using a secrets file, run the following command: <p>`$ make external-es-secrets NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT> ELASTICSEARCH_VERSION=<ELASTICSEARCH_VERSION> PATH_TO_CERTIFICATE=<PATH_TO_CERTS>`</p><p>Where NAMESPACE references your deployment namespace of the SRS cluster, `ELASTICSEARCH_VERSION` matches the Elasticsearch version you want to use, and `PATH_TO_CERTIFICATE` points to the location where you copied the required certificates on your location machine, for example:</p><p>`$ make external-es-secrets NAMESPACE=pegabackingservices ELASTICSEARCH_VERSION=7.10.2 PATH_TO_CERTIFICATE=/home/certs/truststore.jks`</p></li><li>To update the SRS and External Elasticsearch certificates, use the following command: <p>`$ make update-external-es-secrets NAMESPACE=<NAMESPACE_OF EXISTING_DEPLOYMENT> PATH_TO_CERTIFICATE=<PATH_TO_THE_UPDATED_CERTIFICATES>`</p><p>Note: Only .p12 and .jks certificates are supported.</p></li></ol> |
| `domain`                                | Enter the DNS entry associated with your external Elasticsearch service.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |

To configure a mutual Transport Layer Security(mTLS) between SRS and an external Elasticsearch/OpenSearch cluster, configure the following parameters.

| Configuration              | Usage                                                                                                                                                                                                                                                                                                                                                                                                              |
|----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `srsStorage.mtls.enabled`  | Enable this flag if the external Elasticsearch/OpenSearch is provisioned by enforcing the mTLS.                                                                                                                                                                                                                                                                                                                    | 
| `mtls.keystore.file`       | Enter the Elasticsearch/OpenSearch keystore file name that is signed against the CA certificate used for the Elasticsearch/OpenSearch server. When mTLS is enabled, it is mandatory to pass the keystore and truststore in order to establish a mutual handshake between SRS and Elasticsearch/OpenSearch.                                                                                                         | 
| `mtls.keystore.password`   | Enter the keystore password, if any. Pass an empty string('') if no password is set to the keystore.                                                                                                                                                                                                                                                                                                               | 
| `mtls.truststore.file`     | Enter the Elasticsearch/OpenSearch truststore file name that is signed against the CA certificate used for the Elasticsearch/OpenSearch server. When mTLS is enabled, it is mandatory to pass the keystore and truststore in order to establish a mutual handshake between SRS and Elasticsearch/OpenSearch.                                                                                                       | 
| `mtls.truststore.password` | Enter the truststore password, if any. Pass an empty string('') if no password is set to the truststore.                                                                                                                                                                                                                                                                                                           |
| `mtls.certsSecret`         | If you prefer(recommended) to pass the keystore and truststore information through a kubernetes secret, please provide the secret name containing your certificates and their respective passwords. If there is no password for the keystore/truststore, please pass empty string value('') to the passwords of keystore & truststore in the secret with keys being **keystorePassword** & **truststorePassword**. |

**Note:**
- mTLS between SRS and Elasticsearch is not supported in case of internal elasticsearch provisioning mode. It is only supported for an external Elasticsearch/OpenSearch.
- The truststore being used for configuring an external Elasticsearch/OpenSearch server should be signed against the CA certificate used for SRS, if the CA certificate is different.
- The keystore and truststore being passed here under the mtls section is used to establish secure mTLS connection between SRS and Elasticsearch/OpenSearch. You can either use a different keystore, truststore files or reuse the same keystore, truststore used for SRS under `srsRuntime.ssl.keystore.file` & `srsRuntime.ssl.truststore.file`, but make sure to sign the CA certificate used for Elasticsearch/OpenSearch against the SRS truststore.
- Maintaining a different keystore, truststore files is strongly advised to ensure a higher level of security."

To configure mutual Transport Layer Security with PKI Authentication (mTLS with PKI) between SRS and an external Elasticsearch/OpenSearch cluster, configure the following parameters.

| Configuration                                         | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
|-------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `srsStorage.mtlsWithPKIAuthentication.enabled`       | Enable this flag if the external Elasticsearch/OpenSearch is provisioned with PKI-based authentication. This authentication method uses only client certificates (PKI) without requiring username/password credentials.                                                                                                                                                                                                                                                                                                        |
| `mtlsWithPKIAuthentication.keystore.file`            | Enter the Elasticsearch/OpenSearch keystore file name that contains the client certificate signed by a Certificate Authority (CA) that is trusted by the Elasticsearch/OpenSearch server for PKI authentication. When mTLS with PKI is enabled, it is mandatory to pass the keystore and truststore to establish secure communication between SRS and Elasticsearch/OpenSearch using certificate-based authentication only.                                                                                               |
| `mtlsWithPKIAuthentication.truststore.file`          | Enter the Elasticsearch/OpenSearch truststore file name that contains the CA certificates used to verify the Elasticsearch/OpenSearch server's certificate. When mTLS with PKI is enabled, it is mandatory to pass the keystore and truststore to establish secure communication between SRS and Elasticsearch/OpenSearch.                                                                                                                                                                                                   |
| `mtlsWithPKIAuthentication.certsSecret`              | Provide the secret name containing your certificates and their respective passwords. Use the full name of the certificate files (together with file extension, for example, "keystore.p12" or "truststore.jks") as key names in the secret. If there is no password for the keystore/truststore, pass empty string value('') to the passwords in the secret with keys being **keystorePassword** & **truststorePassword**. This is the only supported method for passing certificate credentials for PKI authentication. |

**Note:**

* mTLS with PKI authentication between SRS and Elasticsearch/OpenSearch is not supported for internally-provisioned Elasticsearch clusters. It is only supported for external Elasticsearch/OpenSearch clusters.
* PKI authentication does not require username/password credentials as authentication is performed solely through client certificates.
* Certificate credentials (keystore and truststore passwords) must be provided via Kubernetes secrets only. Plain text password configuration is not supported for PKI authentication.
* Ensure that the client certificate in the keystore is properly signed by a Certificate Authority (CA) that is trusted by the Elasticsearch/OpenSearch server for PKI authentication.
* The keystore must contain a client certificate that is trusted by the Elasticsearch/OpenSearch server, and the truststore must contain the CA certificates used to verify the server's certificate.
* This method provides enhanced security by eliminating the need for shared passwords and relying entirely on certificate-based authentication.

### Enable request authentication/authorization mechanism using identity provider(IdP) between SRS and Pega Infinity
To configure authentication/authorization mechanism using identity provider(IdP) between SRS and Pega Infinity, add the following the settings in your backingservices configuration file and then in pega chart's `values.yml` / pega helm inspected `pega.yaml`.

| Configuration                      | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
|------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `srsRuntime.env.AuthEnabled`       | <ul><li>Set to `false` to disable authentication/authorization between SRS and Pega Infinity.</li><li>Set to `true` to enable authentication/authorization using identity provider(IdP) between SRS and Pega Infinity.</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `srsRuntime.env.OAuthPublicKeyURL` | <ul><li>Provide a valid OAuth Public Key URL from IdP, which can be used in SRS to fetch public keys from IdP and validate the request that is coming into SRS.</li><li> Make sure the authorization server of IdP is having 1) A scope with name `pega.search:full` 2) `scp` & `guid` claims under `pega.search:full` scope with scope name(`pega.search:full`) and `<Customer Deployment Id>` respectively</li><li>Example of JWT token payload having `scp` and `guid` claims <pre>{   "ver": 1,   "jti": "AT.EmHCGDFHE18hC5j3stjbarVonh46twW7tWutB9v8hsw",   "iss": "https://prod-pega.okta.com/oauth2/aus8ahm2k777777",   "aud": "srs-cmc-stg",   "iat": 1678097157,   "exp": 1678100757,   "cid": "0oa8acv8gbBPRCZ7I5d7",   "scp": [     "pega.search:full"   ],   "auth_time": 1000,   "sub": "0oa8acv8gbBPRCZ7I5d7",   "guid": "bf0d4cb6-e09f-1111-ab19-9aac5156b618" } </pre></ul> |

### Enable TLS/HTTPS between SRS and Pega Infinity
There are multiple ways to enable TLS/HTTPS for service which deployed in k8s. It all depends on the k8s setup, environment setup and network policies/restrictions etc. You can work with IT departments to get this done. Some ways to enable TLS/HTTPS are 1. [Using ingress with TLS ](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls). 2. [Using load balancer ](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) etc. Some cloud provided k8s specific ways are 1) AWS: Using combination of [AWS Certificate Manager](https://aws.amazon.com/certificate-manager/), [AWS Load Balancer Controller](https://aws.amazon.com/blogs/opensource/kubernetes-ingress-aws-alb-ingress-controller/), [Kubernetes Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/#aws). 2. Azure: Using combination of [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/), [Kubernetes Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/#azure), [Azure Application Gateway](https://docs.microsoft.com/en-us/azure/application-gateway/). 3. GCP: Using combination of [Google-managed SSL certificates](https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs), [GCP HTTP(S) Load Balancer](https://cloud.google.com/load-balancing/docs/https/), [Kubernetes Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/#gce-gke) etc.

Example:

```yaml
srs:
  # Set srs.enabled=true to enable SRS
  enabled: true

  # specify unique name for the deployment based on org app and/or srs applicable environment name. eg: acme-demo-dev-srs
  deploymentName: "YOUR_SRS_DEPLOYMENT_NAME"

  # Configure the location of the busybox image that is used during the deployment process of
  # the internal Elasticsearch cluster
  busybox:
    image: "alpine:3.16.0"
    imagePullPolicy: "IfNotPresent"

  srsRuntime:
    # Number of pods to provision
    replicaCount: 2

    # docker image of the srs-service, platform-services/search-n-reporting-service:dockerTag
    srsImage: "YOUR_SRS_IMAGE:TAG"

    env:
      # AuthEnabled may be set to true when there is an authentication mechanism in place between SRS and Pega Infinity.
      AuthEnabled: false
      # When `AuthEnabled` is `true`, enter the appropriate public key URL. When `AuthEnabled` is `false`(default), leave this parameter empty.
      OAuthPublicKeyURL: ""
    
    # Set to `true` if Highly secured connection complying NIST SP 800-53 and NIST SP 800-131 is required; otherwise leave set to `false`
    enableSecureCryptoMode: false

    # This is used to configure Java options values.
    javaOpts: ""

  # This section specifies the elasticsearch cluster configuration.
  srsStorage:
    # Setting srsStorage.provisionInternalESCluster to true will provision an internal elasticsearch cluster using the configuration
    # specified in the `elasticsearch` section
    # IF you do not enable SRS and the srs.enabled parameter is set to false, always set srs.srsStorage.provisionInternalESCluster=false
    provisionInternalESCluster: true
    # To use your own Elasticsearch cluster, set srsStorage.provisionInternalESCluster to false and then
    # set the external Elasticsearch cluster URL and port details below when using an externally managed elasticsearch
    # Ensure that the specified endpoint is accessible from the kubernetes cluster pods.
    # domain: ""
    # port: 9200
    # protocol: https
    # The elasticsearch connection supports three authentication methods: basic authentication , AWS IAM role-based authentication and Elasticsearch secure connection(tls).
    # Set srs.srsStorage.tls.enabled: true to enable the use of TLS-based authentication to your Elasticsearch service whether is it running as an internalized or externalized service in your SRS cluster.
    tls:
      enabled: false
    # To specify a certificate used to authenticate an external Elasticsearch service (with tls.enabled: true and srsStorage.provisionInternalESCluster: false), uncomment the following lines to specify the TLS certificate name with password for your Elasticsearch service.
    # Default certificatePassword value will be empty if not used.
    # certificateName: "Certificate_Name"
    # certificatePassword: "password"
    # To specify a certificate using a secret, uncomment the certsSecret parameter and provide the secret name containing your certificate and certificate password.
    # Use the full name of the certificate file (together with file extension, for example, “certificate.p12” or “certificate.jks”) as a key name in the secret. Use this key name to configure the “certificateName” parameter.
    # Use a key name “password” to provide the certificate password in the secret.
    # certsSecret: srs-certificates    
    # Set srs.srsStorage.basicAuthentication.enabled: true to enable the use of basic authentication to your Elasticsearch service whether is it running as an internalized or externalized service in your SRS cluster.
    basicAuthentication:
      enabled: true
    # To configure basic authentication or TLS-based authentication to your externally-managed Elasticsearch/OpenSearch service in your SRS cluster,
    # uncomment and add the parameter details: srs.srsStorage.authCredentials.username and srs.srsStorage.authCredentials.password
    # Auth Credentials added under authCredentials field which supports both Elasticsearch and OpenSearch credentials.
    # authCredentials:
    #   username: "username"
    #   password: "password"
    # for your externally managed Elasticsearch cluster.
    # uncomment and add the parameter details: srs.srsStorage.esCredentials.username and srs.srsStorage.esCredentials.password for your externally managed elasticsearch cluster.
    # esCredentials will be deprecated in future releases, please switch to authCredentials.
    # esCredentials:
    #   username: "username"
    #   password: "password"
    # To use a secret to configure basic authentication or TLS-based authentication between your external Elasticsearch service and SRS,
    # uncomment the authSecret parameter and set it to the secret name. Use "username" and "password" as keys for your secret.    
    # authSecret: srs-elastic-credentials    
    # To configure AWS IAM role-based authentication to your externally-managed Elasticsearch cluster, uncomment
    # and add the parameter details: srs.srsStorage.awsIAM and its associated region, srs.srsStorage.awsIAM.region
    # awsIAM:
    #   region: "AWS_ELASTICSEARCH_REGION"
    # To configure either authentication method, when the elasticsearch domain requires an open internet connection,set the requireInternetAccess parameter to "true".
    requireInternetAccess: false

```

To enable Transport Layer Security(TLS) or Mutual TLS(mTLS) for SRS, configure the following parameters.

| Configuration            | Usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
|--------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `srsRuntime.ssl.enabled` | Set this flag to `true` to enable either TLS or mTLS for SRS. This will enable HTTPS for SRS.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `clientAuthentication`   | To enforce mTLS between Infinity and SRS, enable `ssl.enabled` and set client-authentication to **'need'**. And for a simple TLS(Https) mode, set client-authentication to **'want'**. When mutual TLS(mTLS) is enabled, SRS expects the client(in this case Infinity) to present a valid certificate that it can verify against its trusted certificate authority before completing the TLS handshake.                                                                                                                                                                                                                                           | 
| `keystore.file`          | Enter the keystore file name. When SSL is enabled, it is mandatory to pass the keystore and truststore in order to establish SSL for SRS.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `keystore.password`      | Enter the keystore password, if any. Pass an empty string('') if no password is set to the keystore.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | 
| `keystore.type`          | Type refers to the type/extension of the keystore. Ex : PKCS12.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | 
| `truststore.file`        | Enter the truststore file name. When SSL is enabled, it is mandatory to pass the keystore and truststore in order to establish SSL for SRS.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `truststore.password`    | Enter the truststore password, if any.  Pass an empty string('') if no password is set to the truststore.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         | 
| `truststore.type`        | Type refers to the type/extension of the truststore. Ex : JKS.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | 
| `certsSecret`            | If you prefer(recommended) to pass the keystore and truststore information through a kubernetes secret, please provide the secret name containing your certificates and their respective passwords. If there is no password for the keystore/truststore, please pass empty string value('') to the passwords of keystore & truststore in the secret with keys being **keystorePassword** & **truststorePassword**. Despite the certificates and passwords being passed through secrets, you must still configure "srs.srsRuntime.ssl.keystore.type" and "srs.srsRuntime.ssl.truststore.type" with the type of store file it is (Ex : PKCS12, JKS) |

#### Pre-requisite Step to enable mTLS for SRS
In order to enable mTLS for SRS, you need to perform the following make command.
  ```bash
  make srs-mtls-prerequisite CA_CERT_PATH=<PATH_TO_THE_CA_CERT_OF_SRS> KEYSTORE_PATH=<PATH_TO_THE_KEYSTORE_OF_SRS> KEYSTORE_PASS=<PASSWORD_TO_THE_KEYSTORE_OF_SRS> NAMESPACE=<KUBERNETES_NAMESPACE_WHERE_SRS_IS_DEPLOYED>
  ```
- Since mTLS is enforced, we need to pass the client certificate & it's key along with the CA certificate to be able to check the health of an SRS container.
- We need to run the above make command to make sure that required certificates are created for the readiness probe to check the health of SRS containers.
- Here `PATH_TO_THE_CA_CERT_OF_SRS` refers to the CA certificate path which is used to create the SRS keystore(`srsRuntime.ssl.keystore.file`) and truststore(`srsRuntime.ssl.truststore.file`).
- `PATH_TO_THE_KEYSTORE_OF_SRS` refers to the SRS keystore(`srsRuntime.ssl.keystore.file`) file path. And `PASSWORD_TO_THE_KEYSTORE_OF_SRS` refers to the keystore's password.
- When enforcing mTLS for SRS, please make sure to pass the required keystore & truststore certificates to the infinity charts where ever it is applicable.

### Steps to upgrade SRS (with Internal Elasticsearch) to Kubernetes Cluster Version >=1.25

To support SRS on Kubernetes version >=1.25 you need to use Elasticsearch server version 7.17.9. If you are using an earlier version (7.10.2 or 7.16.3) of Elasticsearch in your deployment, to upgrade to 7.17.9, you need to perform the following steps:
1. Get the latest backingservices Helm chart which supports `k8s version >=1.25`.
2. Update the SRS and Elasticsearch certificates by running the following Make command:
  ```bash
  make update-secrets NAMESPACE=<NAMESPACE_OF EXISTING_DEPLOYMENT> ELASTICSEARCH_VERSION=7.17.9
  ```
3.   To use Elasticsearch version 7.17.9, inspect the values.yaml file from the latest backingservices helm chart and confirm if the imageTag parameter in the values.yaml file is same as in the example below:
  ```yaml
  elasticsearch:
    imageTag: 7.17.9
  ```
4. Upgrade your deployment using the following Helm command:
  ```bash
  helm upgrade backingservices pega/backingservices --version <CHART_VERSION> --namespace <NAMESPACE_OF EXISTING_DEPLOYMENT> --values <YAML_FILE_WITH_ES_IMAGE_TAG_CHANGES>
  ```
5. Verify that the Elasticsearch pods status is Running.
6. Restart the old SRS pods and verify that the status of the new pods is Running.
7. Verify all pods are running and working as expected.
