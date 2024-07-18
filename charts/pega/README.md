# Pega Helm chart

The Pega Helm chart is used to deploy an instance of Pega Infinity into a Kubernetes environment.  This readme provides a detailed description of possible configurations and their default values as applicable. You reference the Pega Helm chart to deploy using the parameter settings in the Helm chart using the `helm --set` command to specify a one-time override specific parameter settings that you configured in the Pega Helm chart.

## Supported providers

Enter your Kubernetes provider which will allow the Helm charts to configure to any differences between deployment environments. These values are case-sensitive and must be lowercase.

Value       | Deployment target
---         | ---
k8s         | Open-source Kubernetes
openshift   | Red Hat Openshift
eks         | Amazon Elastic Kubernetes Service (EKS)
gke         | Google Kubernetes Engine (GKE)
pks         | VMware Tanzu Kubernetes Grid Integrated Edition (TKGI), which used to be Pivotal Container Service (PKS) (**Note:** VMware Tanzu Kubernetes Grid Integrated Edition (TKGI) was deprecated for all releases in February 2024. Current deployments on TKGI continue to be supported, but as a best practice, do not use TKGI for new deployments of Pega Platform.)
aks         | Microsoft Azure Kubernetes Service (AKS)

Example for a kubernetes environment:

```yaml
provider: "k8s"
```

## Actions

Use the `action` section in the helm chart to specify a deployment action.  The standard actions is to deploy against an already installed database, but you can also install a Pega system. These values are case-sensitive and must be lowercase.

For additional, required installation parameters, see the [Installer section](#installations).

Value             | Action
---               | ---
deploy            | Start the Pega containers using an existing Pega database installation.
install           | Install Pega Platform into your database without deploying.
install-deploy    | Install Pega Platform into your database and then deploy.
upgrade           | Upgrade or patch Pega Platform in your database without deploying.
upgrade-deploy    | Upgrade or patch Pega Platform in your database and then deploy.
<!--upgrade           | Upgrade the Pega Platform installation in your database.
upgrade-deploy    | Upgrade the Pega Platform installation in your database, and then deploy.
-->
Example:

```yaml
action: "deploy"
```
## Kerberos Configuration

Use the `kerberos` section to configure Kerberos authentication for Decisioning data flows that fetch data from Kafka or HBase streams. For more information on Decisioning data flows that use Kerberos, see [Data Set types](https://docs.pega.com/bundle/platform/page/platform/decision-management/data-set-types.html).

To configure Kerberos authentication, provide the contents of your krb5.conf file n the `krb5.conf` parameter. For more information, see official Kerberos documentation.

For example:
```yaml
global:
  kerberos:
    krb5.conf: |
      ----SAMPLE KRB5.CONF----
```

## JDBC Configuration

Use the `jdbc` section  of the values file to specify how to connect to the Pega database. Pega must be installed to this database before deploying on Kubernetes.  

### URL and Driver Class
These required connection details will point Pega to the correct database and provide the type of driver used to connect. Examples of the correct format to use are provided below. 

Example for Oracle:
```yaml
jdbc:
  url: jdbc:oracle:thin:@//YOUR_DB_HOST:1521/YOUR_DB_NAME
  driverClass: oracle.jdbc.OracleDriver
```

Example for Microsoft SQL Server:

```yaml
jdbc:
  url: jdbc:sqlserver://YOUR_DB_HOST:1433;databaseName=YOUR_DB_NAME;selectMethod=cursor;sendStringParametersAsUnicode=false
  driverClass: com.microsoft.sqlserver.jdbc.SQLServerDriver
```

Example for IBM DB2 for LUW:
```yaml
jdbc:
  url: jdbc:db2://YOUR_DB_HOST:50000/YOUR_DB_NAME:fullyMaterializeLobData=true;fullyMaterializeInputStreams=true;progressiveStreaming=2;useJDBC4ColumnNameAndLabelSemantics=2;
  driverClass: com.ibm.db2.jcc.DB2Driver
```

Example for IBM DB2 for z/OS:
```yaml
jdbc:
  url: jdbc:db2://YOUR_DB_HOST:50000/YOUR_DB_NAME
  driverClass: com.ibm.db2.jcc.DB2Driver
```

Example for PostgreSQL:
```yaml
jdbc:
  url: jdbc:postgresql://YOUR_DB_HOST:5432/YOUR_DB_NAME
  driverClass: org.postgresql.Driver
```

#### (Optional) Support for providing credentials/certificates using External Secrets Operator

To avoid directly entering your confidential content in your Helm charts such as passwords or certificates in plain text, Pega supports Kubernetes secrets to secure credentials and related information.
Use secrets to represent credentials for your database, Docker registry, SSL certificates, externalized kafka service, or any other token or key that you need to pass to a deployed application. Your secrets can be stored in any secrets manager provider. 
Pega supports two methods of passing secrets to your deployments; choose the method that best suits you organization's needs:

• Mount secrets into your Docker containers using the External Secrets Operator([https://external-secrets.io/v0.5.3/](https://external-secrets.io/v0.5.1/)).

To support this option,

1) Create two files following the Kubernetes documentation for External Secrets Operator :
   - An external secret file that specifies what information in your secret to fetch.
   - A secret store to define access how to access the external and placing the required files in your Helm directory.
2) Copy both files into the pega-helm-charts/charts/pega/templates directory of your local Helm repository.
3) Update your local Helm repository to the latest version using the command: 
   - helm repo update pega https://pegasystems.github.io/pega-helm-charts
4) Update the `external_secret_name` parameter in the values.yaml file to refer to the `spec.target.name` defined in the External Secret file you created in step 1. Update the parameter for each section where you want to use the External Secrets Operator.

•  Pass secrets directly to your deployment using your organization's recommend practices. Pega supports the providers listed under the [Provider tab]( https://external-secrets.io/v0.8.1) as long as your implementation meets the documented guidelines for a given provider.

##### Things to note in case of providing keystore, certificates for Enabling encryption of traffic between Ingress/LoadBalancer and Pod
1. Configure the CA certificate and keystore as a base64 encrypted string inside your preferred secret manager (AWS Secret Manager, Azure Key Vault etc). For details, see [this section.](#enabling-encryption-of-traffic-between-ingressloadbalancer-and-pod)
2. Have the keystore password as plaintext.
3. The secret key should be TOMCAT_KEYSTORE_CONTENT, TOMCAT_KEYSTORE_PASSWORD and ca.crt for keystore, keystore password and CA certificate respectively.
4. For alternate configuration the keys should be TOMCAT_CERTIFICATE_FILE, TOMCAT_CERTIFICATE_KEY_FILE and TOMCAT_CERTIFICATE_CHAIN_FILE, ca.crt(in case of traefik addon enabled) for certificate and key files.


### Driver URI

Pega requires a database driver JAR to be provided for connecting to the relational database.  This JAR may either be baked into your image by extending the Pega provided Docker image, or it may be pulled in dynamically when the container is deployed.  If you want to pull in the driver during deployment, you will need to specify a URL to the driver using the `jdbc.driverUri` parameter.  This address must be visible and accessible from the process running inside the container.

Use the `customArtifactory.authentication.basic` section to provide access credentials or use `customArtifactory.authentication.apiKey` to provide an APIKey value and dedicated APIKey header details if you host the driver in a custom artifactory that requires Basic or APIKey Authentication.

If you configured a secret in an external secrets operator for customArtifactory credentials, enter the secret name in `customArtifactory.authentication.external_secret_name` parameter. For details, see [this section.](#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)

If your artifactory domain server certificate is not issued by Certificate Authority, you must provide the server certificate using the `customArtifactory.certificate` parameter. To disable SSL verification, you can set `customArtifactory.enableSSLVerification` to `false` and leave the `CustomArtifactory.certificate` parameter blank.

The Pega Docker images use Java 11, which requires that the JDBC driver that you specify is compatible with Java 11.

### Authentication

The simplest way to provide database authorization is via the `jdbc.username` and `jdbc.password` parameters. These values will create a Kubernetes Secret and at runtime will be obfuscated and stored in a secrets file.

### Connection Properties

You may optionally set your connection properties that will be sent to our JDBC driver when establishing new connections.  The format of the string is `[propertyName=property;]`. Otherwise, refer to the `URL and Driver Class` section above to determine the adequate connection properties.

### Schemas

It is standard practice to have separate schemas for your rules and data.  You may specify them as `rulesSchema` and `dataSchema`.  If desired, you may also optionally set the `customerDataSchema` for your database. The `customerDataSchema` defaults to value of `dataSchema` if not specified. Additional schemas can be defined within Pega.
 
 Example:
 
 ```yaml
jdbc:
  ...
  rulesSchema: "rules"
  dataSchema: "data"
  customerDataSchema: ""
```

### JDBC Connections

JDBC Connections facilitate the communication between Java applications and relational databases. They enable Java applications to run SQL statements, retrieve results, and propagate changes to the database in a standardized manner. A Connection interface represents a session with a specific database. It provides methods to create statements, commit or rollback transactions, and manage database connections.

Pega provides the following three JDBC Connection types:

JDBC Connection type     | Description
---           |---
`PegaRULES`           | The default Pega JDBC Connection.
`PegaRULESLongRW`         | Provides a longer database connection, which results in less frequent timeouts. Starting in Pega Helm Charts v3.23.0, `PegaRULESLongRW` is enabled for all deployment tiers.
`PegaRULESReadOnly`        | Provides a read only database connection.

The JDBC Connection configurations are defined in [charts/pega/config/deploy/context.xml.tmpl](https://github.com/pegasystems/pega-helm-charts/blob/ab0cb220fe3d297a2e8d8be1c278bcdba96bd646/charts/pega/config/deploy/context.xml.tmpl#L4).

## Docker

Specify the location for the Pega Docker image.  This image is available on DockerHub, but can also be mirrored and/or extended with the use of a private registry.  Specify the url of the image with `docker.pega.image`. You may optionally specify an imagePullPolicy with `docker.pega.imagePullPolicy`.

When using a private registry that requires a username and password, specify them using the `docker.registry.username` and `docker.registry.password` parameters.

To avoid specifying Docker registry credentials in values.yaml, create secrets for Docker registry credentials to avoid exposing these credentials. Refer to [Kubernetes secrets](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) to create Docker registry secrets. Specify secret names as a list of comma-separated strings using the `docker.imagePullSecretNames` parameter. Kubernetes checks each registry secret of image pull secrets to pull an image from the repository. If the specified image is available in one of the provided secrets, kubernetes will pull it from that repository. To create Docker registry secrets from external secrets, refer to [this section](#optional-support-for-providing-credentialscertificates-using-external-secrets-operator).

When you download Docker images, it is recommended that you specify the absolute image version and the image name instead of using the `latest` tag; for example: `pegasystems/pega:8.4.4` or `platform-services/search-n-reporting-service:1.12.0`. When you download these images with these details from the Pega repository, you pull the latest available image. If you pull images only specifying `latest`, you may not get the image you wanted.

For this reason, it is also recommended that you specify the `docker.pega.imagePullPolicy: "IfNotPresent"` option in production, since it will ensure that a new generic tagged image will not overwrite the locally cached version.

Example:

 ```yaml
docker:
  registry:
    url: "YOUR_DOCKER_REGISTRY"
    username: "YOUR_DOCKER_REGISTRY_USERNAME"
    password: "YOUR_DOCKER_REGISTRY_PASSWORD"
  imagePullSecretNames: []
  pega:
    image: "pegasystems/pega:8.4.4"
    imagePullPolicy: "Always"
```

## Deploying with busybox and k8s-wait-for utility images from a private registry
To deploy Pega Platform, the Pega helm chart requires the use of the busybox and k8s-wait-for images. For clients who want to pull these images from a registry other than Docker Hub, they must tag and push these images to another registry, and then pull these images by specifying `busybox` and `k8s-wait-for` values as described below.

Example:

 ```yaml
utilityImages:
  busybox:
    image: "busybox:1.31.0"
    imagePullPolicy: "IfNotPresent"
  k8s_wait_for:
    image: "pegasystems/k8s-wait-for"
    imagePullPolicy: "IfNotPresent"
```

## Deployment Name (Optional)

Specify a deployment name that is used to differentiate this deployment in your environment. This name will be prepended to the various Pega tiers and the associated k8s objects in your deployment. Your deployment name should be constrained to lowercase alphanumeric and '-' characters.

This is meant as an alternative to renaming individual deployment tiers (see [Tiers of a Pega deployment](#tiers-of-a-pega-deployment)).

For example:
```yaml
global:
  deployment:
    name: app1-dev
```
will result in:
```
>kubectl get pods -n test
NAME                              READY   STATUS    RESTARTS   AGE
app1-dev-search-0                 1/1     Running   0          24m
app1-dev-batch-86584dcd6b-dsvdd   1/1     Running   0          24m
app1-dev-batch-86584dcd6b-lfwjg   1/1     Running   0          7m31s
app1-dev-stream-0                 1/1     Running   0          24m
app1-dev-stream-1                 1/1     Running   0          18m
app1-dev-web-788cfb8cc4-6c5nz     1/1     Running   0          8m57s
app1-dev-web-788cfb8cc4-gcltx     1/1     Running   0          24m
```
The default value is "pega" if it is unset.

## Tiers of a Pega deployment

Pega supports deployments using a multi-tier architecture model that separates application processing from k8s functions. Isolating processing in its own tier supports unique deployment configurations, including the Pega application prconfig, resource allocations, and scaling characteristics. To avoid misconfiguration, use a single helm chart to deploy all your tiers simultaneously. Pega does not support using separate charts for different tiers in a single deployment. Use the tier section in the helm chart to complete your tier deployment with appropriate nodes dedicated to the logical tasks of the tier.

### Tier examples

Three values.yaml files are provided to showcase real world deployment examples.  These examples can be used as a starting point for customization and are not expected to deployed as-is.

For more information about the architecture for how Pega Platform runs in a Pega cluster, see [How Pega Platform and applications are deployed on Kubernetes](https://community.pega.com/knowledgebase/articles/cloud-choice/how-pega-platform-and-applications-are-deployed-kubernetes).

#### Standard deployment using two tiers

To provision a three tier Pega cluster, use the default example in the helm chart, which is a good starting point for most deployments:

Tier name     | Description
---           |---
web           | Interactive, foreground processing nodes that are exposed to the load balancer. Pega recommends that these node use the node classification “WebUser” `nodetype`.
batch         | Background processing nodes which handle workloads for non-interactive processing. Pega recommends that these node use the node classification “BackgroundProcessing” `nodetype`. These nodes should not be exposed to the load balancer.
stream (Deprecated)        | For Pega Platform '23, the use of the 'Stream' node classification is deprecated; new deployments running version 8.8 and later should not use "Stream" nodes. New deployments connect to a Kafka service that you manage in your organization. For existing deployments using an embedded Kafka deployment which are not exposed to the deployment cluster load balancer, Pega will continue to support the "Stream" node classification nodetype.

#### Small deployment with a single tier

To get started running a personal deployment of Pega on kubernetes, you can handle all processing on a single tier.  This configuration provides the most resource utilization efficiency when the characteristics of a production deployment are not necessary.  The [values-minimal.yaml](./values-minimal.yaml) configuration provides a starting point for this simple model.

Tier Name   | Description
---         | ---
pega        | With embedded Kafka, which is currently deprecated, one tier handles all foreground and background processing using the nodeType classification "WebUser,BackgroundProcessing,search,Stream". For newer Pega Platform deployments using a configuration that connects to a Kafka service managed in your organization, "Stream" nodetype not supported.

#### Large deployment for production isolation of processing

To run a larger scale Pega deployment in production, you can split additional processing out to dedicated tiers.  The [values-large.yaml](./values-large.yaml) configuration provides an example of a multi-tier deployment that Pega recommends as a good starting point for larger deployments.

Tier Name   | Description
---         | ---
web         | Interactive, foreground processing nodes that are exposed to the load balancer. Pega recommends that these node use the node classification “WebUser” `nodetype`.
batch       | Background processing nodes which handle some of the non-interactive processing.  Pega recommends that these node use the node classification   “BackgroundProcessing,Search,Batch” `nodetype`. These nodes should not be exposed to the load balancer.
stream (Deprecated)     | For Pega Platform '23, the use of the 'Stream' node classification is deprecated; new deployments running version 8.8 and later should not use "Stream" nodes. New deployments connect to a Kafka service that you manage in your organization. For existing deployments using an embedded Kafka deployment which are not exposed to the deployment cluster load balancer, Pega will continue to support the "Stream" node classification nodetype.
bix         | Nodes dedicated to BIX processing can be helpful when the BIX workload has unique deployment or scaling characteristics. Pega recommends that these node use the node classification “Bix” `nodetype`. These nodes should not be exposed to the load balancer.

### Name (*Required*)

Use the `tier` section in the helm chart to specify the name of each tier configuration in order to label a tier in your Kubernetes deployment.  This becomes the name of the tier's replica set in Kubernetes. This name must be unique across all Pega deployments to ensure compatibility with logging and monitoring tools.

Example:

```yaml
name: "mycrm-prod-web"
```

### nodeType (*Required*)

Node classification is the process of separating nodes by purpose, predefining their behavior by assigning node types. When you associate a work resource with a specific node type,you optimize work performance in your Pega application. For more information, see
[Node classification](https://community.pega.com/sites/default/files/help_v83/procomhelpmain.htm#engine/node-classification/eng-node-classification-con.htm).

Specify the list of Pega node types for this deployment.  For more information about valid node types, see the Pega Community article on [Node Classification].

[Node types for VM-based and containerized deployments](https://docs.pega.com/bundle/platform-88/page/platform/system-administration/node-types-on-premises.html)

Example:

```yaml
nodeType: ["WebUser","bix"]
```

### Requestor configuration

Configuration related to Pega requestor settings is collected under `requestor` block.

Configuration parameters:

- `passivationTimeSec` - inactivity time after which requestor is passivated (persisted) and its resources reclaimed.

Example:

```yaml
requestor:
  passivationTimeSec: 900
```
### Security context

By default, security context for your Pega pod deployments `pegasystems/pega` image uses `pegauser`(9001) as the user and volume mounts uses `root`(0) as the group. To configure an alternative user for your custom image, set value for `runAsUser` and to configure an alternative group for volume mounts, set value for `fsGroup`. Note that pegasystems/pega image works only with pegauser(9001).
`runAsUser` and `fsGroup` must be configured in `securityContext` under each tier block and will be applied to Deployments/Statefulsets, along with these settings other allowed configuration can also be supplied here, see the [Kubernetes Documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/).

Example:

```yaml
tier:
  - name: my-tier
    securityContext:
      runAsUser: RUN_AS_USER
      fsGroup: FS_GROUP
```
### service

Specify the `service` yaml block to expose a Pega tier to other Kubernetes run services, or externally to other systems. The name of the service will be based on the tier's name, so if your tier is "web", your service name will be "pega-web". If you omit service, no Kubernetes service object is created for the tier during the deployment. For more information on services, see the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/services-networking/service).

Configuration parameters:

Parameter | Description                       | Default value
---       | ---                               | ---
`httpEnabled`    | Use this to disable the http port `80` on Pega web service. Make sure `tls` is enabled if http port is disabled. | `true`
`port`    | The port of the tier to be exposed to the cluster. For HTTP this is generally `80`. | `80`
`targetPort`    | The target port of the container to expose. The Pega container exposes web traffic on port `8080`. | `8080`
`serviceType`    | The [type of service](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) you wish to expose. | `LoadBalancer`
`annotations` | Optionally add custom annotations for advanced configuration. Specifying a custom set of annotations will result in them being used *instead of* the default configurations. | *n/a*

Example:

```yaml
service:
  port: 1234
  targetPort: 1234
  serviceType: LoadBalancer
```

### ingress

Specify the `ingress` yaml block to expose a Pega tier to access from outside Kubernetes. Pega supports the use of managing SSL certificates for HTTPS configuration using a variety of methods. Set `ingress.enabled` to true in order to deploy an ingress for the tier. For more information on services, see the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/).

Parameter | Description
---       | ---
`domain`  | Specify a domain on your network in which you create an ingress to the load balancer.
`path`  | Specify custom path to the host.
`pathType`  | Specify pathType for routing based on the Ingress controller chosen. Default is `ImplementationSpecific`
`appContextPath`  | Specify the path for access to the Pega application on a specific tier. If not specified, users will have access to the Pega application via /prweb 
`tls.enabled` | Specify the use of HTTPS for ingress connectivity. If the `tls` block is omitted, TLS will not be enabled.
`tls.secretName` | Specify the Kubernetes secret you created in which you store your SSL certificate for your deployment. For compatibility, see [provider support for SSL certificate injection](#provider-support-for-ssl-certificate-management).
`tls.useManagedCertificate` | On GKE, set to `true` to use a managed certificate; otherwise use `false`.
`tls.ssl_annotation` | On GKE or EKS, set this value to an appropriate SSL annotation for your provider.
`annotations` | Optionally add custom annotations for advanced configurations. For Kubernetes, EKS, and OpenShift deployments, including custom annotations overrides the default configuration; for GKE and AKS deployments, the deployment appends these custom annotations to the default list of annotations.

Depending on your provider or type of certificate you are using use the appropriate annotation:
  - For `EKS` - use `alb.ingress.kubernetes.io/certificate-arn: \<*certificate-arn*\>` to specify required ARN certificate.
  - For `AKS` - use `appgw.ingress.kubernetes.io/request-timeout: \<*time-out-in-seconds*\>` to configure application gateway timeout settings.

Example:

```yaml
ingress:
  domain: "tier.example.com"
  annotations:
    annotation-name-1: annotation-value-1
    annotation-name-2: annotation-value-2
```
Depending on what type of deployment you use, if there are any long-running operations such as import, append provider-specific ingress timeout annotation under each tier.

The following example shows timeout annotation overrides for an Openshift deployment:

```yaml
ingress:
  domain: "tier.example.com"
  annotations:
     haproxy.router.openshift.io/timeout: 2m
```

#### Provider support for SSL certificate management

Provider  | Kubernetes Secrets | Cloud SSL management service
---       | ---                | ---
AKS       | Supported          | None
EKS       | Not supported      | Manage certificate using Amazon Certification Manager and use ssl_annotation - see example for details.
PKS (now TKGI)       | Supported          | None
GKE       | Supported          | [Pre-shared or Google-managed certificates](#managing-certificates-in-google-cloud)

#### Managing certificates using Kubernetes secrets

In order to manage the SSL certificate using a secret, do the following:

1. Create the SSL certificate and import it into the environment using the certificate management tools of your choice.

2. Create the secret and add the certificate to the secret file.

3. Add the secret name to the pega.yaml file.

4. Pass the secret to the cluster you created in your environment before you begin the Pega Platform deployment.

Example:

```yaml
ingress:
  domain: "tier.example.com"
  tls:
    enabled: true
    secretName: web-domain-certificate
    useManagedCertificate: false
```

#### Managing certificates in AWS

Instead of Kubernetes secrets, on AWS you must manage your SSL certificates with ACM (AWS certificate manager). Using the ARN of your certificate, you configure the `ssl_annotation` in your Helm chart.

Example:

```yaml
ingress:
  domain: "tier.example.com"
  tls:
    enabled: true
    secretName:
    useManagedCertificate: false
    ssl_annotation:
      alb.ingress.kubernetes.io/certificate-arn:<certificate-arn>
```

#### Managing certificates in Google Cloud

In addition to Kubernetes secrets, on GCP you may manage your SSL certificates in GKE with two alternative methods. For more information, see the [Google Cloud documentation on SSL certificate management](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-multi-ssl).

- *Pre-shared certificate* - add the certificate to your Google Cloud project and specify the appropriate ssl annotation in the ingress section.

Example:

```yaml
ingress:
  domain: "web.dev.pega.io"
  tls:
    enabled: true
    useManagedCertificate: false
    ssl_annotation:
      ingress.gcp.kubernetes.io/pre-shared-cert: webCert
```

- *Google-managed certificate* - Pega Platform deployments can automatically generate a GKE managed certificate when you specify the appropriate SSL annotation in the ingress section. Using a static IP address is not mandatory; if you do not use it, remove the annotation. To use a static IP address, you must create the static IP address during the cluster configuration, then add it using this annotation in the pega.yaml.

Example:

```yaml
ingress:
  domain: "web.dev.pega.io"
  tls:
    enabled: true
    useManagedCertificate: true
    ssl_annotation:
      kubernetes.io/ingress.global-static-ip-name: web-ip-address
```

### Managing Resources

You can optionally configure the resource allocation and limits for a tier using the following parameters. The default value is used if you do not specify an alternative value. See [Managing Kubernetes Resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/) for more information about how Kubernetes manages resources.

Example:
```yaml
resources: 
  requests:
    memory: "12Gi"
    cpu: 3
    ephemeral-storage: 
  limits:
    memory: "12Gi"
    cpu: 4
    ephemeral-storage: 
```


Parameter                | Description                                            | Default value
---                      | ---                                                    | ---
`replicas`               | Specify the number of Pods to deploy in the tier.      | `1`
`cpuRequest`             | Deprecated, use `resources.requests.cpu`. Initial CPU request for pods in the current tier.      | `3`
`cpuLimit`               | Deprecated, use `resources.limits.cpu`. CPU limit for pods in the current tier.                | `4`
`memRequest`             | Deprecated, use `resources.requests.memory`. Initial memory request for pods in the current tier.   | `12Gi`
`memLimit`               | Deprecated, use `resources.limits.memory`. Memory limit for pods in the current tier.             | `12Gi`
`initialHeap`            | Specify the initial heap size of the JVM.              | `8192m`
`maxHeap`                | Specify the maximum heap size of the JVM.              | `8192m`
`ephemeralStorageRequest`| Deprecated, use `resources.requests.ephemeral-storage`. Ephemeral storage request for the tomcat container.    | -
`ephemeralStorageLimit`  | Deprecated, use `resources.limits.ephemeral-storage`. Ephemeral storage limit for the tomcat container.      | -

### JVM Arguments
You can optionally pass in JVM arguments to Tomcat.  Depending on the parameter/attribute used, the arguments will be placed into `JAVA_OPTS` or `CATALINA_OPTS` environmental variables.
Some of the Best-practice arguments for JVM tuning are included by default in `CATALINA_OPTS`.
Pass the required JVM parameters as `catalinaOpts` attributes in respective `values.yaml` file.  

Example:
```yaml
tier:
  - name: my-tier
    javaOpts: ""
    catalinaOpts: "-XX:SomeJVMArg=XXX"
```
Note that some JVM arguments are non-overrideable i.e. baked in the Docker image. <br>
Check [RecommendedJVMArgs.md](./RecommendedJVMArgs.md) for more details.
### nodeSelector

Pega supports configuring certain nodes in your Kubernetes cluster with a label to identify its attributes, such as persistent storage. For such configurations, use the Pega Helm chart nodeSelector property to assign pods in a tier to run on particular nodes with a specified label. For more information about assigning Pods to Nodes including how to configure your Nodes with labels, see the [Kubernetes documentation on nodeSelector](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector).

```yaml
tier:
- name: "my-tier"
  nodeType: "WebUser"

  nodeSelector:
    disktype: ssd
```

### Tolerations

Pega supports configuring tolerations for workloads. Taints are applied to nodes and tolerations are applied to pods. For more information about taints and tolerations please refer official K8S [documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/).

Example:

```yaml
tier:
- name: "my-tier"
  nodeType: "WebUser"
  
  tolerations:
  - key: "key1"
    operator: "Equal"
    value: "value1"
    effect: "NoSchedule"
  
```

### Liveness, readiness, and startup probes

Pega uses liveness, readiness, and startup probes to determine application health in your deployments. For an overview of these probes, see [Configure Liveness, Readiness and Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/). Configure a probe for *liveness* to determine if a Pod has entered a broken state; configure it for *readiness* to determine if the application is available to be exposed; configure it for *startup* to determine if a pod is ready to be checked for liveness. You can configure probes independently for each tier. If not explicitly configured, default probes are used during the deployment. Set the following parameters as part of a `livenessProbe`, `readinessProbe`, or `startupProbe` configuration.

Notes:
* Kubernetes 1.18 and later supports `startupProbe`. If your deployment uses a Kubernetes version older than 1.18, the helm charts exclude `startupProbe` and use different default values for `livenessProbe` and `readinessProbe`.
* `timeoutSeconds` cannot be greater than `periodSeconds` in some GCP environments. For details, see [this API library from Google](https://developers.google.com/resources/api-libraries/documentation/compute/v1/csharp/latest/classGoogle_1_1Apis_1_1Compute_1_1v1_1_1Data_1_1HttpHealthCheck.html#a027a3932f0681df5f198613701a83145).

#### Kubernetes pre-1.18
Parameter             | Description    | Default `livenessProbe` | Default `readinessProbe`
---                   | ---            | ---                     | ---
`initialDelaySeconds` | Number of seconds after the container has started before probes are initiated. | `200` | `30`
`timeoutSeconds`      | Number of seconds after which the probe times out. | `20` | `10`
`periodSeconds`       | How often (in seconds) to perform the probe. | `30` | `10`
`successThreshold`    | Minimum consecutive successes for the probe to be considered successful after it determines a failure. | `1` | `1`
`failureThreshold`    | The number consecutive failures for the pod to be terminated by Kubernetes. | `3` | `3`

#### Kubernetes 1.18+
Parameter             | Description    | Default `livenessProbe` | Default `readinessProbe` | Default `startupProbe`
---                   | ---            | ---                     | ---                      | ---
`initialDelaySeconds` | Number of seconds after the container has started before probes are initiated. | `0` | `0` | `10`
`timeoutSeconds`      | Number of seconds after which the probe times out. | `20` | `10` | `10`
`periodSeconds`       | How often (in seconds) to perform the probe. | `30` | `10` | `10`
`successThreshold`    | Minimum consecutive successes for the probe to be considered successful after it determines a failure. | `1` | `1` | `1`
`failureThreshold`    | The number consecutive failures for the pod to be terminated by Kubernetes. | `3` | `3` | `20`

Example:

```yaml
tier:
  - name: my-tier
      livenessProbe:
        initialDelaySeconds: 60
        timeoutSeconds: 30
        failureThreshold: 5
      readinessProbe:
        initialDelaySeconds: 400
        failureThreshold: 30
```

### Using a Kubernetes Horizontal Pod Autoscaler (HPA)

You may configure an HPA to scale your tier on a specified metric.  Only tiers that do not use volume claims are scalable with an HPA. Set `hpa.enabled` to `true` in order to deploy an HPA for the tier. For more details, see the [Kubernetes HPA documentation](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/). 

Parameter           | Description    | Default value
---                 | ---       | ---
`hpa.minReplicas`   | Minimum number of replicas that HPA can scale-down | `1` 
`hpa.maxReplicas`   | Maximum number of replicas that HPA can scale-up  | `5`
`hpa.targetAverageCPUValue` | Threshold value for scaling based on absolute CPU usage (The default value is `2.55` which represents 2.55 [Kubernetes CPU units](https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/#cpu-units)) | `2.55`
`hpa.targetAverageCPUUtilization` | Threshold value for scaling based on initial CPU request utilization (Can be set instead of `hpa.targetAverageCPUValue` to set the threshold as a percentage of the requested CPU) | 
`hpa.targetAverageMemoryUtilization` | Threshold value for scaling based on initial memory utilization (The default value is `85` which corresponds to 85% of 12Gi ) | `85`
`hpa.enableCpuTarget` | Set to true if you want to enable scaling based on CPU utilization or false if you want to disable it | true
`hpa.enableMemoryTarget` | Set to true if you want to enable scaling based on memory utilization or false if you want to disable it (Pega recommends leaving this disabled) | false

### Ensure System Availability during Voluntary Disruptions by Using a Kubernetes Pod Disruption Budget (PDB)
To limit the number of Pods running your Pega Platform application that can go down for planned disruptions, 
Pega allows you to enable a Kubernetes `PodDisruptionBudget` on a tier.  For more details on PDBs, see the Kubernetes [Pod Disruption Budgets documentation](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/#pod-disruption-budgets).

You can configure a Kubernetes `PodDisruptionBudget` on your tier by setting `pdb.enabled` to `true` in your values.yaml file.  By default, this value is
set to `false.`  You must also specify exactly one of the following parameters to configure the Pod Disruption Budget.  These parameters
are mutually exclusive, and thus only one may be set.  You may provide values that are expressed as percentages (% of pods) or integers (# of pods).  

Parameter             | Description    | Default value
---                   | ---       | ---
`pdb.minAvailable`    | The minimum number or percentage of pods in the tier that must be available.  If this minimum is reached, the Kubernetes deployment will not bring down additional pods for voluntary disruptions until more are available and healthy. | `1`
`pdb.maxUnavailable`  | The maximum number or percentage of pods in the tier that can be unavailable.  If this maximum is reached, the Kubernetes deployment will not bring down additional pods for voluntary disruptions until more are available and healthy.      | `50%` (disabled by default)

### Volume claim template

A `volumeClaimTemplate` may be configured for any tier to allow for persistent storage. This allows for stateful tiers such as `stream` to be run as a StatefulSet rather than a Deployment.  Specifying a `volumeClaimTemplate` should never be used with a custom deployment strategy for rolling updates.

### Deployment strategy

The `deploymentStrategy` can be used to optionally configure the [strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy) for any tiers deployed as a Kubernetes Deployment. This value will cannot be applied to StatefulSet deployed tiers which use the `volumeClaimTemplate` parameter.

### Environment variables

Pega supports a variety of configuration options for cluster-wide and application settings. In cases when you want to pass a specific environment variable into your deployment on a tier-by-tier basis, you specify a custom `env` block for your tier as shown in the example below.

Example:

```yaml
tier:
  - name: my-tier
    custom:
      env:
        - name: MY_ENV_NAME
          value: MY_ENV_VALUE
```

### Service Account

If the pod needs to be run with a specific [service account](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/), you can specify a custom `serviceAccountName` for your deployment tier.

Example:

```yaml
tier:
  - name: my-tier
    custom:
      serviceAccountName: MY_SERVICE_ACCOUNT_NAME
```

### Custom volumes

You can optionally specify custom `volumes` and `volumeMounts` for your deployment tier. You need to grant read and/or write permissions to the volume location to the Pega user depending on the purpose of the volume. By default, the Pega user UID is 9001.

For example:

```yaml
tier:
  - name: my-tier
    custom:
      volumeMounts:
        - name: my-volume
          mountPath: /path/to/mount
      volumes:
        - name: my-volume
          configMap:
            name: my-configmap 
```

### Sidecar Containers

Pega supports adding sidecar containers to manage requirements for your Pega application services that live outside of the primary tomcat container. This may include company policy requirements, utility images, networking containers, or other examples. For an overview of the versatility sidecar containers present, see [How Pods manage multiple containers](https://kubernetes.io/docs/concepts/workloads/pods/#how-pods-manage-multiple-containers).

You can specify custom `sidecarContainers` for your deployment tiers in the Pega Helm chart as shown in the example below. Each sidecar container definition must be a complete container definition, including a name, image, and resources.

Example:

```yaml
tier:
  - name: my-tier
    custom:
      sidecarContainers:
        - name: SIDECAR_NAME
          image: SIDECAR_IMAGE_URL
          ...
        - name: SIDECAR_NAME_2
          image: SIDECAR_IMAGE_URL_2
          ...
```

### Custom Annotations for Pods

You may optionally provide custom annotations for Pods as metadata to be consumed by other tools and libraries. Pod annotations may be specified by using the `podAnnotations` element for a given `tier`.

Example:

```yaml
tier:
  - name: my-tier
    podAnnotations:
      <annotation-key>: <annotation-value>
```

### Pod affinity

You may optionally configure the pod affinity so that it is restricted to run on particular node(s), or to prefer to run on particular nodes. Pod affinity may be specified by using the `affinity` element for a given `tier`. See the official [Kubernetes Documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/).

Example:

```yaml
tier:
  - name: my-tier
    affinity:
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
          - matchExpressions:
            - key: kubernetes.io/os
              operator: In
              values:
              - linux
```

### Pega configuration files

While Pega includes default configuration files in the Helm charts, the charts provide extension points to override the defaults with additional customizations. To change the configuration file, specify the replacement implementation to be injected into a ConfigMap.

Parameter     | Description    | Default value
---           | ---       | ---
`prconfig`    | A complete prconfig.xml file to inject.  | See [prconfig.xml](config/deploy/prconfig.xml).
`prlog4j2`    | A complete prlog4j2.xml file to inject.  | See [prlog4j2.xml](config/deploy/prlog4j2.xml).
`contextXML`  | A complete context.xml template file to inject.  | See [context.xml.tmpl](config/deploy/context.xml.tmpl).
`serverXML`   | A complete server.xml file to inject  | See [server.xml.tmpl](config/deploy/server.xml.tmpl).
`webXML`      | A complete web.xml file to inject  | No default provided, but if `config/deploy/web.xml` exists, it will be used as the default.


Example:

```yaml
tier:
  - name: my-tier
    custom:
      prconfig: |-
        ...

      prlog4j2: |-
        ...

      contextXML: |-
        ...

      serverXML: |-
        ...

      webXML: |-
        ...
```
### Pega compressed configuration files

To use [Pega configuration files](https://github.com/pegasystems/pega-helm-charts/blob/master/charts/pega/README.md#pega-configuration-files) in compressed format when deploying Pega Platform, replace each file with its compressed format file by completing the following steps:

1) Compress each configuration file using the following command in your local terminal:
```
- cat "<path_to_actual_uncompressed_file_in_local>" | gzip -c | base64 
```
Example for a prconfig.xml file:
```
cat "pega-helm-charts/charts/pega/config/deploy/prconfig.xml" | gzip -c | base64
```
2) Provide the file content with the output of the command for each file executed.
3) Set the `compressedConfigurations` in values.yaml to `true`, as in the following example:
```yaml
  compressedConfigurations: true
```

### Pega diagnostic user

While most cloud native deployments will take advantage of aggregated logging using a tool such as EFK, there may be a need to access the logs from Tomcat directly. In the event of a need to download the logs from tomcat, a username and password will be required.  You may set `pegaDiagnosticUser` and `pegaDiagnosticPassword` to set up authentication for Tomcat.

## Cassandra and Pega Customer Decision Hub deployments

If you are planning to use Cassandra (usually as a part of Pega Customer Decision Hub), you may either point to an existing deployment or deploy a new instance along with Pega.

### Using an existing Cassandra deployment

To use an existing Cassandra deployment, set `cassandra.enabled` to `false` and configure the `dds` section to reference your deployment.

Use the following parameters to configure the connection to your external Cassandra cluster

Parameter     | Tier Level Environment Variable | Description | Default value
---           |:---:| ---|:---:
`externalNodes` | N/A | A comma separated list of hosts in the Cassandra cluster. | Empty
`port` | N/A | TCP Port to connect to cassandra. | 9042
`username` | N/A | The plain text username for authentication with the Cassandra cluster.<br/>Change the value in your helm chart to the username supplied by your Cassandra cluster provider. For better security, avoid plain text usernames and leave this parameter blank; then include the username in an external secrets manager with the key CASSANDRA_USERNAME. <br/>If you make no change, Pega attempts to authenticate with the Cassandra cluster using the default username `dnode_ext`. | dnode_ext
`password` | N/A | The plain text password for authentication with the Cassandra cluster.<br/>Change the value in your helm chart to the password supplied by your Cassandra cluster provider. For better security, avoid plain text passwords and leave this parameter blank; then include the password in an external secrets manager with the key CASSANDRA_PASSWORD. <br/>If you make no change, Pega attempts to authenticate with the Cassandra cluster using the default password `dnode_ext`.| dnode_ext
`clientEncryption` | N/A | Enable (true) or disable (false) client encryption on the Cassandra connection. | false
`trustStore` | N/A | If required, provide the trustStore certificate file name.<br/>When using a trustStore certificate, you must also include a Kubernetes secret name that contains the trustStore certificate in the global.certificatesSecrets parameter.<br/>Pega deployments only support trustStores that use the Java Key Store (.jks) format. | Empty
`trustStorePassword` | N/A | If required provide trustStorePassword value in plain text. For better security leave this parameter blank and include the password in an external secrets manager with the key CASSANDRA_TRUSTSTORE_PASSWORD. | Empty
`keyStore` | N/A | If required, provide the keystore certificate file name.<br/>When using a keystore certificate, you must also include a Kubernetes secret name that contains the keystore certificate in the global.certificatesSecrets parameter.<br/>Pega deployments only support keystores that use the Java Key Store (.jks) format. | Empty
`keyStorePassword` | N/A | If required provide keyStorePassword value in plain text. For better security leave this parameter blank and include the password in an external secrets manager with the key CASSANDRA_KEYSTORE_PASSWORD. | Empty
`asyncProcessingEnabled` | CASSANDRA_ASYNC_PROCESSING_ENABLED | Enable asynchronous processing of records in DDS Dataset save operation. Failures to store individual records will not interrupt Dataset save operations. | false
`keyspacesPrefix` | N/A | Specify a prefix to use when creating Pega-managed keyspaces in Cassandra. | Empty
`extendedTokenAwarePolicy` | CASSANDRA_EXTENDED_TOKEN_AWARE_POLICY | Enable an extended token aware policy for use when a Cassandra range query runs. When enabled this policy selects a token from the token range to determine which Cassandra node to send the request. Before you can enable this policy, you must configure the token range partitioner. | false
`latencyAwarePolicy` | CASSANDRA_LATENCY_AWARE_POLICY | Enable a latency awareness policy, which collects the latencies of the queries for each Cassandra node and maintains a per-node latency score (an average). | false
`customRetryPolicy` | CASSANDRA_CUSTOM_RETRY_POLICY | Enable the use of a customized retry policy for your Pega Platform deployment for Pega Platform ’23 and earlier releases. After you enable this policy in your deployment configuration, the deployment retries Cassandra queries that time out. Configure the number of retries using the dynamic system setting (DSS): dnode/cassandra_custom_retry_policy/retryCount. The default is 1, so if you do not specify a retry count, timed out queries are retried once.| false
`customRetryPolicyEnabled` | CASSANDRA_CUSTOM_RETRY_POLICY_ENABLED | Use this parameter in Pega Platform '24 and later instead of `customRetryPolicy`. Configure the number of retries using the `customRetryPolicyCount` property.| false
`customRetryPolicyCount` | CASSANDRA_CUSTOM_RETRY_POLICY_COUNT | Specify the number of retry attempts when `customRetryPolicyEnabled` is true. For Pega Platform '23 and earlier releases use the dynamic system setting (DSS): dnode/cassandra_custom_retry_policy/retryCount. | 1
`speculativeExecutionPolicy` | CASSANDRA_SPECULATIVE_EXECUTION_POLICY | Enable the speculative execution policy for retrieving data from your Cassandra service for Pega Platform '23 and earlier releases. When enabled, Pega Platform sends a query to multiple nodes in your Cassandra service and processes the first response. This provides lower perceived latencies for your deployment, but puts greater load on your Cassandra service. Configure the speculative execution delay and max executions using the following dynamic system settings (DSS): dnode/cassandra_speculative_execution_policy/delay and dnode/cassandra_speculative_execution_policy/max_executions. | false
`speculativeExecutionPolicyEnabled` | CASSANDRA_SPECULATIVE_EXECUTION_POLICY_ENABLED | Use this parameter in Pega Platform '24 and later instead of `speculativeExecutionPolicy`. Configure the speculative execution delay and max executions using the `speculativeExecutionPolicyDelay` and `speculativeExecutionPolicyMaxExecutions` properties. | false
`speculativeExecutionPolicyDelay` | CASSANDRA_SPECULATIVE_EXECUTION_DELAY | Specify the delay in milliseconds before speculative executions are made when `speculativeExecutionPolicyEnabled` is true. For Pega Platform '23 and earlier releases use the dynamic system setting (DSS): dnode/cassandra_speculative_execution_policy/delay. | 100
`speculativeExecutionPolicyMaxExecutions` | CASSANDRA_SPECULATIVE_EXECUTION_MAX_EXECUTIONS | Specify the maximum number of speculative execution attempts when `speculativeExecutionPolicyEnabled` is true. For Pega Platform '23 and earlier releases use the dynamic system setting (DSS): dnode/cassandra_speculative_execution_policy/max_executions. | 2
`jmxMetricsEnabled` | CASSANDRA_JMX_METRICS_ENABLED | Enable reporting of DDS SDK metrics to a Java Management Extension (JMX) format for use by your organization to monitor your Cassandra service. Setting this property `false` disables metrics being exposed through the JMX interface; disabling also limits the metrics being collected using the DDS landing page. | true
`csvMetricsEnabled` | CASSANDRA_CSV_METRICS_ENABLED | Enable reporting of DDS SDK metrics to a Comma Separated Value (CSV) format for use by your organization to monitor your Cassandra service. If you enable this property, use the Pega Platform DSS: dnode/ddsclient/metrics/csv_directory to customize the filepath to which the deployment writes CSV files. By default, after you enable this property, CSV files will be written to the Pega Platform work directory. | true
`logMetricsEnabled` | CASSANDRA_LOG_METRICS_ENABLED | Enable reporting of DDS SDK metrics to your Pega Platform logs. | false


If you configured a secret in an external secrets operator, enter the secret name in `external_secret_name` parameter. For details, see [this section.](#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)

Example:

This example configuration shows the parameters required for a deployment that connects to an SSL encrypted Cassandra service. It configures an extended token aware policy for browse operations and custom retries for all queries. Metrics logging to CSV and Pega Logs are disabled. Usernames and passwords are all supplied by Kubernetes secrets.

```yaml
cassandra:
  enabled: false

dds:
  externalNodes: "CASSANDRA_NODE_IPS"
  port: "9042"
  username: ""
  password: ""
  keyspacesPrefix: "dev01"
  trustStore: "/opt/pega/certs/cass-truststore.jks"
  keyStore: "/opt/pega/certs/cass-keystore.jks"
  extendedTokenAwarePolicy: true
  customRetryPolicy: true
  csvMetricsEnabled: false
  logMetricsEnabled: false
  # The external secret below contains passwords with the following keys: CASSANDRA_USERNAME, CASSANDRA_PASSWORD, CASSANDRA_TRUSTSTORE_PASSWORD, CASSANDRA_KEYSTORE_PASSWORD
  external_secret_name: "dev02-credentials-secret"
```
In addition to being configured at the cluster level, the parameters above that have a corresponding tier level environement variable may be specified at the tier level. The following example shows how to configure the Dataflow tier to use the latency aware load balancing policy.

```yaml
  tier:
    - name: "Dataflow"
      nodeType: "BackgroundProcessing,Dataflow"

      custom:
        env:
        - name: CASSANDRA_LATENCY_AWARE_POLICY
          value: "true"
```

### Deploying Cassandra with Pega

You may deploy a Cassandra instance along with Pega.  Cassandra is a separate technology and needs to be independently managed.  When deploying Cassandra, set `cassandra.enabled` to `true` and leave the `dds` section as-is.  For more information about configuring Cassandra, see the [Cassandra Helm charts](https://github.com/helm/charts/blob/master/incubator/cassandra/values.yaml).

Pega does **not** actively update the Cassandra dependency in `requirements.yaml`. When deploying Cassandra with Pega, you should update its `version` value in `requirements.yaml`.

#### Cassandra minimum resource requirements

Deployment  | CPU     | Memory
---         | ---     | ---
Development | 2 cores | 4Gi
Production  | 4 cores | 8Gi

Example:

```yaml
cassandra:
  enabled: true
  # Set any additional Cassandra parameters. These values will be used by Cassandra's helm chart.
  persistence:
    enabled: true
  resources:
    requests:
      memory: "4Gi"
      cpu: 2
    limits:
      memory: "8Gi"
      cpu: 4

dds:
  externalNodes: ""
  port: "9042"
  username: "dnode_ext"
  password: "dnode_ext"
```
### Deploying Pega without Cassandra
To configure a Pega platform deployment without a Cassandra datastore (DDS), set `cassandra.enabled` to `false` and comment out or delete the `dds` section.

Example:

```yaml
cassandra:
  enabled: false

#dds:
#  externalNodes: ""
#  port: "9042"
#  username: "dnode_ext"
#  password: "dnode_ext"
```

## Search deployment

Use the `pegasearch` section to configure the source Elasticsearch service that the Pega Platform deployment uses for searching Rules and Work within Pega. The Elasticsearch service defined here is not related to the Elasticsearch deployment if you also define an EFK stack for logging and monitoring in your Pega Platform deployment.

### For Pega Platform 8.6 and later:

Use the chart ['backingservices'](../backingservices) to deploy the Search and Reporting Service (SRS), a Pega Platform backing service enabling the latest generation of search and reporting capabilities for your Pega applications. SRS is independent from Pega Platform and replaces the previous implementation of Elasticsearch, the legacy client-server Elasticsearch plug-in.

To use SRS, follow the deployment instructions provided at ['backingservices'](../backingservices) before you configure and deploy the Pega Helm chart. For more information, see [External Elasticsearch in your deployment](https://docs.pega.com/bundle/platform-88/page/platform/deployment/externalization-of-services/externalize-search-in-your-deployment.html).

Configure the customerDeploymentId parameter in the global section of the values.yaml to provide data isolation in SRS.  The customerDeploymentId is used as a prefix for all indexes created in ElasticSearch, and must be the value of the 'guid' claim if OAuth is used for authorization between Pega and SRS.  This parameter defaults to the name of the namespace when left empty.

You must configure the SRS URL for your Pega Platform deployment using the parameter in values.yaml as shown the following table and example:

Parameter   | Description   | Default value
---         | ---           | ---
`externalSearchService` | Set the `pegasearch.externalSearchService` as true to use Search and Reporting service as the search functionality provider to the Pega platform | false
`externalURL` | Set the `pegasearch.externalURL` value to the Search and Reporting Service endpoint url | `""`

Example:

```yaml
pegasearch:
  externalSearchService: true
  externalURL: "http://srs-service.namespace.svc.cluster.local"
```

To configure authorization for the connection between Pega Infinity and the Search and Reporting Service (SRS) use the OAuth authorization service. For more information, see ['backingservices'](../backingservices). To configure the connection to the authorization service in SRS you must configure the following authorization parameters in the Pega values.yaml as shown in the following table and example:

| Parameter             | Description                                                                                                                                                                                           | Default value      |
|-----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| `enabled`             | Set the `pegasearch.srsAuth.enabled` to 'true' to use OAuth between Infinity and SRS.                                                                                                                 | false              |
| `url`                 | Set the `pegasearch.srsAuth.url` value to the URL of the OAuth token endpoint to get the token for SRS.                                                                                             | `""`               |
| `clientId`            | Set the `pegasearch.srsAuth.clientId` value to the client id used in OAuth service.                                                                                                                   | `""`               |
| `scopes`              | Set the `pegasearch.srsAuth.scopes` value to "pega.search:full", the scope set in the OAuth service required to grant access to SRS.                                                                  | "pega.search:full" |
| `authType`            | Set the `pegasearch.srsAuth.authType` value to to authentication type use when connecting to the OAuth token endpoint. Use client_secret_basic for basic authentication or private_key_jwt to use a client assertion JWT.     | `""`               |
| `external_secret_name`| Set the `pegasearch.srsAuth.external_secret_name` value to the secret that contains the OAuth private PKCS8 key (additionally encoded with base64) used to get an authorization token for the connection between Pega tiers and SRS.  The private key should be contained in the secret key SRS_OAUTH_PRIVATE_KEY.   | `""`               |
| `privateKey`          | When not using an external secret, set the `pegasearch.srsAuth.privateKey` value to the OAuth private PKCS8 key (additionally encoded with base64) used to get an authorization token for the connection between Pega tiers and SRS.     | `""`               |
| `privateKeyAlgorithm` | Set the `pegasearch.srsAuth.privateKeyAlgorithm` value to the algorithm used to generate a private key used by the OAuth client. Allowed values: RS256 (default), RS384, RS512, ES256, ES384, ES512.  | "RS256"            |

Example:

```yaml
pegasearch:
  externalSearchService: true
  externalURL: "http://srs-service.srs-namespace.svc.cluster.local"
  srsAuth:
    enabled: true
    url: "https:/your-authorization-service-host/oauth2/v1/token"
    clientId: "your-client-id"
    authType: client_secret_basic
    scopes: "pega.search:full"
    privateKey: "LS0tLS1CRUdJTiBSU0Eg...<truncated>"
    privateKeyAlgorithm: "RS256"
```

### For Pega Platform 8.5 and earlier:

Use the following configuration to provision the legacy client-server Elasticsearch plug-in with a Pega-provided Docker image. This is a deprecated solution; as a best practice, update your deployment to Pega Platform version 8.6 or later and use SRS instead.

Parameter   | Description   | Default value
---         | ---           | ---
`image`   | Set the `pegasearch.image` parameter to a registry that can access the Pega-provided `platform/search` Docker image. Download the image from the Pega repository, tag it, and push it to your local registry. As a best practice, use the latest available image for your Pega Platform version, based on the build date specified in the tag. For example, the image tagged "8.5.6-20230829" was built on August 29, 2023. For more information, see [Pega-provided Docker images](https://docs.pega.com/bundle/platform-88/page/platform/deployment/client-managed-cloud/pega-docker-images-manage.html).| `platform/search:8.5.x-XXXXXXXX`
`imagePullPolicy` | Optionally specify an imagePullPolicy for the search container. | `""`
`replicas` | Specify the desired replica count. | `1`
`minimumMasterNodes` | To prevent data loss, you must configure the minimumMasterNodes setting so that each master-eligible node is set to the minimum number of master-eligible nodes that must be visible in order to form a cluster. Configure this value using the formula (n/2) + 1 where n is replica count or desired capacity.  For more information, see the Elasticsearch [important setting documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html) for more information. | `1`
`podSecurityContext.runAsUser`   | Elasticsearch defaults to UID 1000.  In some environments where user IDs are restricted, you may configure your own using this parameter. | `1000`
`set_vm_max_map_count`   | Elasticsearch uses a **mmapfs** directory by default to store its indices. The default operating system limits on mmap counts is likely to be too low, which may result in out of memory exceptions. An init container is provided to set the value correctly, but this action requires privileged access. If privileged access is not allowed in your environment, you may increase this setting manually by updating the `vm.max_map_count` setting in **/etc/sysctl.conf** according to the Elasticsearch documentation and can set this parameter to `false` to disable the init container. For more information, see the [Elasticsearch documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html). | `true`
`set_data_owner_on_startup`   | Set to true to enable an init container that runs a chown command on the mapped volume at startup to reset the owner of the ES data to the current user. This is needed if a random user is used to run the pod, but also requires privileges to change the ownership of files. | `false`
`podAnnotations` | Configurable annotations applied to all Elasticsearch pods. | {}
`affinity` | You may optionally configure the pod affinity so that it is restricted to run on particular node(s), or to prefer to run on particular nodes. | `""`

Additional env settings supported by Elasticsearch may be specified in a `custom.env` block as shown in the example below.

Example:

```yaml
pegasearch:
  image: "platform/search:8.5.6-20230829"
  memLimit: "3Gi"
  replicas: 1
  minimumMasterNodes: 2
  custom:
    env:
    - name: TZ
      value: "EST5EDT"
```
## Deploying Pega with externalized kafka service for stream

Beginning with Pega Platform '23, configure the default Pega Helm chart parameters that are appropriate to connect to a Kafka service that you manage in your organization to use as your stream provider for Pega Platform data integrations
Pega supports migrating existing deployments to use an externalized Kafka configuration as a stream service provider using Pega-provided Helm charts. To use your own managed Kafka infrastructure without the use of stream nodes, Pega provides instructions to run a migration with downtime and potential data loss or with minimal downtime and no downtime. For migration steps, see [Switch from embedded Stream to externalized Kafka service](MigrationToExternalStream.md).

### Stream (externalized Kafka service) settings

Example:
```yaml
# Stream (externalized Kafka service) settings.
stream:
  # Beginning with Pega Platform '23, enabled by default; when disabled, your deployment does not use a Kafka stream service configuration
  enabled: true
  # Provide externalized Kafka service broker urls.
  bootstrapServer: ""
  # Provide Security Protocol used to communicate with kafka brokers. Supported values are: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL.
  securityProtocol: PLAINTEXT
  # If required, provide trustStore certificate file name
  # When using a trustStore certificate, you must also include a Kubernetes secret name, that contains the trustStore certificate,
  # in the global.certificatesSecrets parameter.
  # Pega deployments only support trustStores using the Java Key Store (.jks) format.
  trustStore: ""
  # If required provide trustStorePassword value in plain text.
  trustStorePassword: ""
  # If required, provide keyStore certificate file name
  # When using a keyStore certificate, you must also include a Kubernetes secret name, that contains the keyStore certificate,
  # in the global.certificatesSecrets parameter.
  # Pega deployments only support keyStores using the Java Key Store (.jks) format.
  keyStore: ""
  # If required, provide keyStore value in plain text.
  keyStorePassword: ""
  # If required, provide jaasConfig value in plain text.
  jaasConfig: ""
  # If required, provide a SASL mechanism**. Supported values are: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512.
  saslMechanism: PLAIN
  # By default, topics originating from Pega Platform have the pega- prefix,
  # so that it is easy to distinguish them from topics created by other applications.
  # Pega supports customizing the name pattern for your Externalized Kafka configuration for each deployment.
  streamNamePattern: "pega-{stream.name}"
  # Your replicationFactor value cannot be more than the number of Kafka brokers and 3.
  replicationFactor: "1"
  # To avoid exposing trustStorePassword, keyStorePassword, and jaasConfig parameters, leave the values empty and
  # configure them using an External Secrets Manager, making sure you configure the keys in the secret in the order:
  # STREAM_TRUSTSTORE_PASSWORD, STREAM_KEYSTORE_PASSWORD and STREAM_JAAS_CONFIG.
  # Enter the external secret name below.
  external_secret_name: ""
```

## Pega database installation and upgrades

Pega requires a relational database that stores the rules, data, and work objects used and generated by Pega Platform. The [Pega Platform deployment guide](https://community.pega.com/knowledgebase/products/platform/deploy) provides detailed information about the requirements and instructions for installations and upgrades.  Follow the instructions for Tomcat and your environment's database server.

The Helm charts also support an automated install or upgrade with a Kubernetes Job.  The Job utilizes an installation Docker image and can be activated with the `action` parameter in the Pega Helm chart.
 
### Installations

For installations of the Pega platform, you must specify the installer Docker image and an initial default password for the `administrator@pega.com` user.

Along with this, you can configure the kubelet pull policy for the image. It is defaulted to `IfNotPresent`, meaning an image will be pulled if it is "not present". All possible options are `IfNotPresent`, `Always`, and `Never`. Always pulling an image ensures you always have the latest image at all times, even if the specific tag already exists on your machine. 

Example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
  imagePullPolicy: "PREFERRED_IMAGE_PULL_POLICY"
  adminPassword: "ADMIN_PASSWORD"
```


### Upgrades and patches

The Pega Helm charts support zero-downtime patch and upgrades processes which synchronize the required process steps to minimize downtime. With these zero-downtime processes, you and your customers can continue to access and use their applications in your environment with minimal disruption while you patch or upgrade your system.

To **upgrade Pega Platform software** deployed in a Kubernetes environment in zero-downtime, you must download the latest Pega-provided images for the version to which you are upgrading from  [Pega Digital Software Delivery](https://community.pega.com/digital-delivery) and use the Helm chart with versions 1.6.0 or later to complete the upgrade. To learn about how the upgrade process works and its requirements and the steps you must complete, see the Pega-provided runbook, [Upgrading Pega Platform in your deployment with zero-downtime](/docs/upgrading-pega-deployment-zero-downtime.md). With earlier versions of the Pega Helm charts, you must use the Pega Platform upgrade guides. To obtain the latest upgrade guide, see [Stay current with Pega](https://community.pega.com/upgrade).

To complete your Pega Infinity upgrade, after you upgrade your Pega Platform software using the Pega Helm charts and Docker images, you must use the latest Pega application software Upgrade Guide, which is separate from Pega Platform software. You can locate the appropriate upgrade guide for your installed application from the page, [All Products](https://community.pega.com/knowledgebase/products).

To **apply a Pega Platform patch** with zero-downtime to your existing Pega platform software, you use the same "zero-downtime" parameters that you use for upgrades and use the Pega-provided `platform/installer` Docker image that you downloaded for your patch version. For step-by-step guidance to apply a Pega Platform patch, see the Pega-provided runbook, [Patching Pega Platform in your deployment](/docs/patching-pega-deployment.md). The patch process applies only changes observed between the patch and your currently running version and then separately upgrades the data. For details about Pega patches, see [Pega software maintenance and extended support policy](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-software-maintenance-and-extended-support-policy).

Use the `installer` section  of the values file with the appropriate parameters to install, upgrade, or apply a patch to your Pega Platform software:

Parameter   | Description   | Default value
---         | ---           | ---
`image`   | Reference the `platform/installer` Docker image that you downloaded and pushed to your Docker registry that your deployment can access.  | `YOUR_INSTALLER_IMAGE:TAG`
`imagePullPolicy` | Specify when to pull an image. | `IfNotPresent`
`adminPassword` | Specify a temporary, initial password to log into the Pega application. This will need to be changed at first login. The adminPassword value cannot start with "@". | `"ADMIN_PASSWORD"`
`affinity` | Configures policy to assign the pods to the nodes. See the official [Kubernetes Documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/). | `""`
`upgrade.upgradeType:` |Specify the type of process, applying a patch or upgrading. | See the next table for details.
`upgrade.upgradeSteps:` |Specify the steps of a `custom` upgrade process that you want to complete. For `zero-downtime`, `out-of-place-rules`, `out-of-place-data`, or `in-place` upgrades, leave this parameter empty. | <ul>`enable_cluster_upgrade` `rules_migration` `rules_upgrade` `data_upgrade` `disable_cluster_upgrade`</ul>
`upgrade.targetRulesSchema:` |Specify the name of the schema you created the process creates for the new rules schema. | `""`
`upgrade.targetDataSchema:` | For patches to 8.4 and later or upgrades from 8.4.2 and later, specify the name of the schema the process creates for the temporary data schema. After the patch or upgrade, you must delete this temporary data schema from your database. For 8.3 Pega software patches, you can leave this value empty, as is (do not add blank text). | `""`

Upgrade type    | Description
---             | ---
`zero-downtime` |  If applying any patch or upgrading from 8.4.2 and later, use this option to minimize downtime. This patch or upgrade type migrates the rules to your designated "new rules schema", uses the temporary data schema to host some required data-specific tables, and patches or upgrades the rules to the new version with zero-downtime. With the new rules in place, the process performs a rolling reboot of your nodes, patches or upgrades any required data schema, and redeploys the application using the new rules.
`custom` |  Use this option for any upgrade in which you complete portions of the upgrade process in steps. Supported upgrade steps are: `enable_cluster_upgrade` `rules_migration` `rules_upgrade` `data_upgrade` `disable_cluster_upgrade`. To specify which steps to include in your custom upgrade, specify them in your pega.yaml file using the `upgrade.upgradeSteps` parameter.
`out-of-place-rules` | Use this option to migrate a copy of the rules to a new rules schema and run an out-of-place upgrade in that copied schema. This schema will become the rules schema after your upgrade is complete.
`out-of-place-data` |Use this option to complete an out-of-place upgrade of the data schema.  This is the final step of the out of place upgrade.
`in-place`      | Use this option to upgrade both rules and data in a single run.  This will upgrade your environment as quickly as possible but will result in application downtime.
`out-of-place` | `Deprecated and supported only with Helm charts prior to version 1.4`: For patches using Helm charts from 1.4 or earlier, you can use this process to apply a patch with zero-downtime; for upgrades from 1.4 or earlier this upgrade type minimizes downtime, but still requires some downtime. For patches or upgrades the process places the existing rules in your application into a read-only state, migrates the rules to your designated "new rules schema", and then applies the patch only to changed rules or upgrades all of the rules. With the new rules in place, the process performs a rolling reboot, patches or upgrades the data, and then redeploys your application using the new rules.

Install example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
```

Zero-downtime upgrade example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
  upgrade:
    upgradeType: "zero-downtime"
    targetRulesSchema: "new_rules_schema_name"
    targetDataSchema: "temporary_data_schema_name"
```

Custom rules upgrade without a data upgrade example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
  upgrade:
    upgradeType: "custom"
    upgradeSteps: "enable_cluster_upgrade, rules_migration, rules_upgrade, disable_cluster_upgrade"
    targetRulesSchema: "new_rules_schema_name"
    targetDataSchema: "temporary_data_schema_name"
```

Zero-downtime patch example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
  upgrade:
    upgradeType: "zero-downtime"
    targetRulesSchema: "new_rules_schema_name"
    targetDataSchema: "temporary_data_schema_name"
```

### Installer Pod Annotations

You can add annotations to the installer pod.

Example:

```yaml
installer:
  podAnnotations:
    annotation-name1: annotation-value1
    annotation-name2: annotation-value2
```


### Installer Node Selector

You can add node selector to the installer pod to launch the pod on specific node

Example:

```yaml
installer:
  nodeSelector:
    label: value
```

### Mount the custom certificates into the Tomcat container

Pega supports mounting and passing custom certificates into the tomcat container during your Pega Platform deployment. Pega supports the following certificate formats as long as they are encoded in base64: X.509 certificates such as PEM, DER, CER, CRT. To mount and pass the your custom certificates, use the `certificates` attributes as a map in the `values.yaml` file using the format in the following example.
#### (Optional) Support for providing custom certificates using External Secrets Operator
To avoid directly entering your  certificates in plain text, Pega supports Kubernetes secrets to secure certificates. Your certificates can be stored in any secrets manager provider.

• Mount secrets into your Docker containers using the External Secrets Operator([https://external-secrets.io/v0.5.3/](https://external-secrets.io/v0.5.1/)).

To support this option,

1) Create two files following the Kubernetes documentation for External Secrets Operator :
    - An external secret file that specifies what information in your secret to fetch.
    - A secret store to define access how to access the external and placing the required files in your Helm directory.
2) Copy both files into the pega-helm-charts/charts/pega/templates directory of your local Helm repository.
3) Update your local Helm repository to the latest version using the command:
    - helm repo update pega https://pegasystems.github.io/pega-helm-charts
4) Update your values.yaml file to refer to the external secret manager for certificates.
5) Add multiple custom certificates that you maintain as an externally-managed secret, each as a string, separated by a comma in the certificatesSecrets parameter.

• You can either pass certificates as external secrets or as plain text to the Pega values.yaml, but not both. If you provide both, the deployment mounts only external secrets into the tomcat container.
Example:

```yaml
certificatesSecrets: ["secret-to-be-created1","secret-to-be-created2"]
certificates:
    badssl.cer: |
      "-----BEGIN CERTIFICATE-----\n<<certificate content>>\n-----END CERTIFICATE-----\n"

```
## Deploying Hazelcast in Client-Server Model 

**For Pega Platform deployments using 8.6 and later, Pega recommends adopting a client-server model for your Hazelcast deployment.**
This deployment model introduces independent scalability for both servers and clients in Pega Platform. 
To adopt this client-server deployment model, configure the values.yaml section for `hazelcast` to use the Pega-provided `platform/clustering-service` Docker image which contains the Hazelcast clustering service image inside it. 
Using this image, your deployment starts a cluster of Hazelcast server nodes; a plugin provided by Hazelcast, the Hazelcast-Kubernetes Plugin, discovers the Hazelcast members in the cluster.

Deploying the Pega provided `platform/clustering-service` Docker image which contains the Hazelcast clustering service image inside it, 
starts a cluster of Hazelcast server nodes. For the discovery of Hazelcast members in the cluster, a plugin provided by Hazelcast, namely Hazelcast-Kubernetes Plugin is used. 
Out of the two discovery strategies that the latter plugin provides - Kubernetes API and DNS Lookup, the client-server model with Hazelcast uses DNS lookup to resolve the IP addressing of PODs running Hazelcast.
For additional information on Hazelcast member discovery, refer the plugin: [Hazelcast-Kubernetes Plugin](https://github.com/hazelcast/hazelcast-kubernetes)

For platform version 8.6 through 8.7.x, specify the `platform/clustering-service` Docker image that you downloaded in `hazelcast.image` and set `hazelcast.enabled` as `true` to deploy a Pega Platform web cluster separately from a Hazelcast cluster in a client-server deployment model.
For platform version 8.8 and later, specify the `platform/clustering-service` Docker image that you downloaded in `hazelcast.clusteringServiceImage` and set `hazelcast.clusteringServiceEnabled` as `true` to deploy a Pega Platform web cluster separately from a Hazelcast cluster in a client-server deployment model.
**Using Clustering service for client-server form of deployment is only supported from Pega Platform 8.6 or later.**

In this model, nodes running Hazelcast start independently and simultaneously with the Pega web tier nodes and create a cluster with a member count you must specify using `hazelcast.replicas` parameter. Pega web tier nodes then connect to this Hazelcast cluster in a client-sever model.

**Note:** If you are deploying Pega Platform below release 8.6, you need to set `hazelcast.enabled` as `false`, otherwise the installation will fail. 
Setting `hazelcast.enabled` as `false` deploys Pega and Hazelcast in an Embedded arrangement, in which Hazelcast and Pega Platform run on the same node. 
The default and recommended deployment strategy for Hazelcast is client-server, Embedded deployment is only being supported for backwards compatibility.
**Embedded deployment would not be supported in future platform releases.**

### Clustering Service Compatibility Matrix

Pega Infinity version   | Clustering Service version    |    Description
---                     | ---                           | ---
< 8.6                   | NA                            | Clustering Service is not supported for releases 8.5 or below 
\>= 8.6 && < 8.8         | \= 1.0.5                     | Pega Infinity 8.6.x and 8.7.x supports using a Pega-provided `platform-services/clustering-service` Docker Image that provides a clustering service version 1.0.3 or later. 
\>= 8.8                 | \= 1.3.x                     | Pega Infinity 8.8 and later supports using a Pega-provided `platform-services/clustering-service` Docker Image that provides a clustering service version 1.3.0 or later. As a best practice, use the latest available release of the clustering service. 


#### Configuration Settings
The values.yaml provides configuration options to define the deployment of Hazelcast. Apart from the below parameters when `hazelcast.enabled` is set to `true`, additional parameters are required for client-server deployment which have been documented
here: [Additional Parameters](charts/hazelcast/README.md)

Parameter   | Description   | Default value
---         | ---           | ---
`hazelcast.image` | Reference the `platform/clustering-service` Docker image that you downloaded and pushed to your Docker registry that your deployment can access. | `YOUR_HAZELCAST_IMAGE:TAG`
`hazelcast.clusteringServiceImage` | Reference the `platform/clustering-service` Docker image that you downloaded and pushed to your Docker registry that your deployment can access. | `YOUR_CLUSTERING_SERVICE_IMAGE:TAG`
`hazelcast.enabled` |  Set to `true` if client-server deployment of Pega Platform is required; otherwise leave set to `false`. Note: To avoid an installation failure, you must set this value to `false` for Pega platform deployments using versions before 8.6. | `true`
`hazelcast.clusteringServiceEnabled` |  Set to `true` if client-server deployment of Pega Platform is required; otherwise leave set to `false`. Note: Set this value to `false` for Pega platform versions below 8.8; if not set the installation will fail. | `false`
`hazelcast.migration.initiateMigration` |  Set to `true` after creating parallel cluster (new Hazelcast) to establish the connection with platform and migrate the data; Set to `false` during a deployment that removes an older Hazelcast cluster. | `false`
`hazelcast.migration.migrationJobImage` | Reference the `platform/clustering-service-kubectl` Docker image to create the migration job to run the migration script. | `YOUR_MIGRATION_JOB_IMAGE:TAG`
`hazelcast.migration.embeddedToCSMigration` |  Set to `true` while migrating the data from existing embedded Hazelcast deployment to the new c/s Hazelcast deployment. | `false`
`hazelcast.replicas` | Number of initial members to join the Hazelcast cluster. | `3`
`hazelcast.username` | Configures the username to be used in a client-server Hazelcast model for authentication between the nodes in the Pega deployment and the nodes in the Hazelcast cluster. This parameter configures the username in Hazelcast cluster and your Pega nodes so authentication occurs automatically.  | `""`
`hazelcast.password` | Configures the password to be used in a client-server Hazelcast model for authentication between the nodes in the Pega deployment and the nodes in the Hazelcast cluster. This parameter configures the password credential in Hazelcast cluster and your Pega nodes so authentication occurs automatically.  | `""`
`hazelcast.external_secret_name` | If you configured a secret in an external secrets operator, enter the secret name. For details, see [this section](#optional-support-for-providing-credentialscertificates-using-external-secrets-operator).  | `""`
`hazelcast.affinity` | Configures policy to assign the pods to the nodes. See the official [Kubernetes Documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/).  | `""`

#### Example
```yaml
hazelcast:
  image: "YOUR_HAZELCAST_IMAGE:TAG"
  clusteringServiceImage: "YOUR_CLUSTERING_SERVICE_IMAGE:TAG"
  enabled: true
  clusteringServiceEnabled: false
  migration:
    initiateMigration: false
    migrationJobImage: "YOUR_MIGRATION_JOB_IMAGE:TAG"
    embeddedToCSMigration: false
  replicas: 3
  username: ""
  password: ""
  external_secret_name: ""
```

### Enabling encryption of traffic between Ingress/LoadBalancer and Pod

Using Helm version `2.2.0`, Pega supports mounting and passing TLS certificates into the container to enable TLS between loadbalancer/ingress and pods during your Pega Platform deployment. Pega supports the keystore formats such as .jks, .keystore. To mount and pass your TLS certificates, use the `tls` section under `service` to specify the keystore content, the keystore password and the specified ports for https under 'web' tier in the `values.yaml` file using the format in the following example.

Parameter   | Description   | Default value
---         | ---           | ---
`service.tls.port` | The port of the tier to be exposed to the cluster. For HTTPS this is generally `443` | `443`
`service.tls.targetPort` | The target port of the container to expose. The TLS-enabled Pega container exposes web traffic on port `8443`. | `8443`
`service.tls.enabled` | Set as `true` if TLS is enabled for the tier, otherwise `false`. | `false`
`service.tls.external_secret_name` | If you configured a secret in an external secrets operator, enter the secret name. For details, see [this section.](#optional-support-for-providing-credentialscertificates-using-external-secrets-operator) | `""` 
`service.tls.keystore` | The keystore content for the tier. If you leave this value empty, the deployment uses the default keystore. | `""`
`service.tls.keystorepassword` | The keystore password for the tier. If you leave this value empty, the deployment uses the default password for the default keystore. | `""`
`service.tls.cacertificate` | The CA certificate for the tier. If you leave this value empty, the deployment uses the default CA certificate for the default keystore. Pass the certificateChainFile file if you are using certificateFile and certificateKeyFile. | `""`
`service.tls.certificateFile` | The content of the file that contains the server certificate, which you must provide if you do not provide a keystore and keystorepassword. The format of this file is PEM-encoded. | `""`
`service.tls.certificateKeyFile` | The content of the file that contains the server private key, which you must provide if you do not provide a keystore and keystorepassword. The format of this file is PEM-encoded. | `""`
`service.tls.traefik.enabled` | Set as `true` if you enabled Traefik for the tier and deployed the Traefik addon Helm charts; otherwise set it to `false`. | `false`
`service.tls.traefik.serverName` | The server name for the tier, SAN(Subject Alternative Name) of the certificate present inside the container | `""`
`service.tls.traefik.insecureSkipVerify` | Set to `true` to skip verifying the certificate; do this in cases where you do not need a valid root/CA certificate but want to encrypt load balancer traffic. Leave the setting to `false` to both verify the certificate and encrypt load balancer traffic. | `false`

##### Important Points to note
- By default, Pega provides a self-signed keystore and a custom root/CA certificate in Helm chart version `2.2.0`. To use the default keystore and CA certificate, leave the parameters service.tls.keystore, service.tls.keystorepassword and service.tls.cacertificate empty. The default keystore and CA certificate expire on 25/12/2025.
- To enable SSL, you must either provide a keystore with a keystorepassword or certificate, certificatekey and cacertificate files in PEM format. If you do not provide either, the deployment implements SSL by passing a Pega-provided default self-signed keystore and a custom root/CA certificate to the Pega web nodes.
- The CA certificate can be issued by any valid Certificate Authorities or you can also use a self-created CA certificate with proper chaining.
- To avoid exposing your certificates, you can use external secrets to manage your certificates. Pega also supports specifying the certificate files using the certificate parameters in the Pega values.yaml. To pass the files using these parameters, you must encode the certificate files using base64 and then enter the string output into the appropriate certificate parameter.
- To encode your keystore and certificate files use the following command:
     o	Linux:  cat ca_bundle.crt | base64 -w 0
     o	Windows: type keystore.jks | openssl base64  (needs openssl)
- Add the required, base64-encoded content in the values.yaml using either the keystore parameters (service.tls.keystore, service.tls.keystorepassword and service.tls.cacertificate) or the certificate parameters (service.tls.certificateFile, service.tls.certificateKeyFile and service.tls.cacertificate).
- Create a keystore file with the SAN(Subject Alternate Name) field present in case of Traefik ingress controller.
- You must use the latest Docker images in order to use this feature; if you use Helm chart version `2.2.0`, with outdated Docker images and set `service.tls.enabled` to `true`, the deployment logs a `Bad Gateway` error. Helm chart version `2.2.0`, you must update your Pega Platform version to the latest patch version or set `service.tls.enabled` to `false`.

#### Example:
After you enable TLS for the web tier, deploy the traefik addon for `k8s` provider, and configure the keystore file and password using the external secret operator, set the following parameters in the `values.yaml`:

```yaml
# To configure TLS between the ingress/load balancer and the backend, set the following:
tls:
   enabled: true
   external_secret_name: secret-to-be-crated
   keystore: 
   keystorepassword: 
   port: 443
   targetPort: 8443
   # set the value of CA certificate here in case of baremetal/openshift deployments - CA certificate should be in base64 format
   # pass the certificateChainFile file if you are using certificateFile and certificateKeyFile
   cacertificate:
   # provide the SSL certificate and private key as a PEM format
   certificateFile:
   certificateKeyFile:
   # if you will deploy traefik addon chart and enable traefik, set enabled=true; otherwise leave the default setting.
   traefik:
      enabled: true
      serverName: "<< SAN name of the certificate >>"
      # set insecureSkipVerify=true, if the certificate verification has to be skipped
      insecureSkipVerify: false

```

With TLS enabled for the web tier and the traefik addon deployed for `k8s` provider, you set the following parameters in the `values.yaml`:

```yaml
# To enable TLS encryption between the ingress/load balancer and the deployment backend, configure the following settings:
tls:
   enabled: true
   external_secret_name: 
   keystore: "<< encoded keystore content >>"
   keystorepassword: "<< keystore password >>"
   port: 443
   targetPort: 8443
   # set the value of CA certificate here in case of baremetal/openshift deployments - CA certificate should be in base64 format
   # pass the certificateChainFile file if you are using certificateFile and certificateKeyFile
   cacertificate: "<< encoded CA certificate >>"
  # provide the SSL certificate and private key as a PEM format
   certificateFile:
   certificateKeyFile:
  # if you will deploy traefik addon chart and enable traefik, set enabled=true; otherwise leave the default setting.
   traefik:
      enabled: true
      serverName: "<< SAN name of the certificate >>"
      # set insecureSkipVerify=true, if the certificate verification has to be skipped
      insecureSkipVerify: false

```
To enable TLS for the web tier using the certificateFile, certificateKeyFile and certificateChainFile instead of a keystore and password, you must set the following parameters in the Pega `values.yaml`:

```yaml
# To enable TLS encryption between the ingress/load balancer and the deployment backend, configure the following settings:
tls:
   enabled: true
   external_secret_name:
   keystore: 
   keystorepassword: 
   port: 443
   targetPort: 8443
   # set the value of CA certificate here in case of baremetal/openshift deployments - CA certificate should be in base64 format
   # pass the certificateChainFile file if you are using certificateFile and certificateKeyFile
   cacertificate: "<< encoded certificateChainFile certificate >>"
  # provide the SSL certificate and private key as a PEM format
   certificateFile: "<< encoded certificateFile content >>"
   certificateKeyFile: "<< encoded certificateKeyFile content >>"
   # if you will deploy traefik addon chart and enable traefik, set enabled=true; otherwise leave the default setting.
   traefik:
      enabled: false
      serverName: ""
      # set insecureSkipVerify=true, if the certificate verification has to be skipped
      insecureSkipVerify: true

```

With TLS enabled for the web tier and the traefik addon is NOT deployed for `k8s` provider, you set the following parameters in the `values.yaml`:

```yaml
# To enable TLS encryption between the ingress/load balancer and the deployment backend, configure the following settings:
tls:
   enabled: true
   external_secret_name:
   keystore: "<< encoded keystore content >>"
   keystorepassword: "<< keystore password >>"
   port: 443
   targetPort: 8443
   # set the value of CA certificate here in case of baremetal/openshift deployments - CA certificate should be in base64 format
   # pass the certificateChainFile file if you are using certificateFile and certificateKeyFile
   cacertificate: "<< encoded CA certificate >>"
  # provide the SSL certificate and private key as a PEM format
   certificateFile:
   certificateKeyFile:
   # if you will deploy traefik addon chart and enable traefik, set enabled=true; otherwise leave the default setting.
   traefik:
      enabled: false
      serverName: ""
      # set insecureSkipVerify=true, if the certificate verification has to be skipped
      insecureSkipVerify: true

```

Without TLS enabled, and no traefik addon in use, there is no reason to add and verify the certificate. You can use the following parameters in the `values.yaml`:

```yaml
# To enable TLS encryption between the ingress/load balancer and the deployment backend, configure the following settings:
tls:
   enabled: false
   external_secret_name:
   keystore: ""
   keystorepassword: ""
   port: 443
   targetPort: 8443
   # set the value of CA certificate here in case of baremetal/openshift deployments - CA certificate should be in base64 format
   # pass the certificateChainFile file if you are using certificateFile and certificateKeyFile
   cacertificate: ""
  # provide the SSL certificate and private key as a PEM format
   certificateFile:
   certificateKeyFile:
   # if you will deploy traefik addon chart and enable traefik, set enabled=true; otherwise leave the default setting.
   traefik:
      enabled: false
      serverName: ""
      # set insecureSkipVerify=true, if the certificate verification has to be skipped
      insecureSkipVerify: true

```

```yaml
# To enable HorizontalPodAutoscaler behavior specifications, configure the following settings against each tier:
behavior:
   scaleDown:
      stabilizationWindowSeconds: << provide scaleDown stabilization window in seconds >>
   scaleUp:
      stabilizationWindowSeconds: << provide scaleUp stabilization window in seconds >>
```

### Custom Ports

You can optionally specify custom ports for deployment tier. You can specify custom ports for your tiers as shown in the example below:

```yaml
tier:
  - name: my-tier
    custom:
      ports:
        - name: <name>
          containerPort: <port>
```

You can optionally specify custom ports for tier specific service. You can specify custom ports for your service as shown in the example below:
```yaml
tier:
   - name: my-tier
     service:
       customServicePorts:
       - name: <name>
         port: <port>
         targetPort: <target port>
          
```
