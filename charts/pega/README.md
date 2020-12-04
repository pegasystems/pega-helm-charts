# Pega Helm chart

The Pega Helm chart is used to deploy an instance of Pega Infinity into a Kubernetes environment.  This readme provides a detailed description of possible configurations and their default values as applicable.

## Supported providers

Enter your Kubernetes provider which will allow the Helm charts to configure to any differences between deployment environments. These values are case-sensitive and must be lowercase.

Value       | Deployment target
---         | ---
k8s         | Open-source Kubernetes
openshift   | Red Hat Openshift
eks         | Amazon Elastic Kubernetes Service (EKS)
gke         | Google Kubernetes Engine (GKE)
pks         | VMware Tanzu Kubernetes Grid Integrated Edition (TKGI), which used to be Pivotal Container Service (PKS)
aks         | Microsoft Azure Kubernetes Service (AKS)

Example for a kubernetes environment:

```yaml
provider: "k8s"
```

## Actions

Use the `action` section in the helm chart to specify a deployment action.  The standard actions is to deploy against an already installed database, but you can also install a Pega system. These values are case-sensitive and must be lowercase.

For additional, required installation parameters, see the [Installer section](#install).

Value             | Action
---               | ---
deploy            | Start the Pega containers using an existing Pega database installation.
install           | Install Pega Platform into your database without deploying.
install-deploy    | Install Pega Platform into your database and then deploy.

<!--upgrade           | Upgrade the Pega Platform installation in your database.
upgrade-deploy    | Upgrade the Pega Platform installation in your database, and then deploy.
-->
Example:

```yaml
action: "deploy"
```

## JDBC Configuration

Use the `jdbc` section  of the values file to specify how to connect to the Pega database. *Pega must be installed to this database before deploying on Kubernetes*.  

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

### Driver URI

Pega requires a database driver JAR to be provided for connecting to the relational database.  This JAR may either be baked into your image by extending the Pega provided Docker image, or it may be pulled in dynamically when the container is deployed.  If you want to pull in the driver during deployment, you will need to specify a URL to the driver using the `jdbc.driverUri` parameter.  This address must be visible and accessible from the process running inside the container.

The Pega Docker images use Java 11, which requires that the JDBC driver that you specify is compatible with Java 11.

### Authentication

The simplest way to provide database authorization is via the `jdbc.username` and `jdbc.password` parameters. These values will create a Kubernetes Secret and at runtime will be obfuscated and stored in a secrets file.

### Connection Properties

You may optionally set your connection properties that will be sent to our JDBC driver when establishing new connections.  The format of the string is `[propertyName=property;]`.

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

## Docker

Specify the location for the Pega Docker image.  This image is available on DockerHub, but can also be mirrored and/or extended with the use of a private registry.  Specify the url of the image with `docker.pega.image`. You may optionally specify an imagePullPolicy with `docker.pega.imagePullPolicy`.

When using a private registry that requires a username and password, specify them using the `docker.registry.username` and `docker.registry.password` parameters.

Example:

 ```yaml
docker:
  registry:
    url: "YOUR_DOCKER_REGISTRY"
    username: "YOUR_DOCKER_REGISTRY_USERNAME"
    password: "YOUR_DOCKER_REGISTRY_PASSWORD"
  pega:
    image: "pegasystems/pega"
    imagePullPolicy: "Always"
```

## Tiers of a Pega deployment

Pega supports deployment using a multi-tier architecture to separate processing and functions. Isolating processing in its own tier also allows for unique deployment configuration such as its own prconfig, resource allocations, or scaling characteristics.  Use the `tier` section in the helm chart to specify which tiers you wish to deploy and their logical tasks.  

### Tier examples

Three values.yaml files are provided to showcase real world deployment examples.  These examples can be used as a starting point for customization and are not expected to deployed as-is.

For more information about the architecture for how Pega Platform runs in a Pega cluster, see [How Pega Platform and applications are deployed on Kubernetes](https://community.pega.com/knowledgebase/articles/cloud-choice/how-pega-platform-and-applications-are-deployed-kubernetes).

#### Standard deployment using three tiers

To provision a three tier Pega cluster, use the default example in the in the helm chart, which is a good starting point for most deployments:

Tier name     | Description
---           |---
web           | Interactive, foreground processing nodes that are exposed to the load balancer. Pega recommends that these node use the node classification “WebUser” `nodetype`.
batch         | Background processing nodes which handle workloads for non-interactive processing. Pega recommends that these node use the node classification “BackgroundProcessing” `nodetype`. These nodes should not be exposed to the load balancer.
stream        | Nodes that run an embedded deployment of Kafka and are exposed to the load balancer. Pega recommends that these node use the node classification “Stream” `nodetype`.

#### Small deployment with a single tier

To get started running a personal deployment of Pega on kubernetes, you can handle all processing on a single tier.  This configuration provides the most resource utilization efficiency when the characteristics of a production deployment are not necessary.  The [values-minimal.yaml](./values-minimal.yaml) configuration provides a starting point for this simple model.

Tier Name   | Description
---         | ---
pega        | One tier handles all foreground and background processing and is given a `nodeType` of "WebUser,BackgroundProcessing,search,Stream".

#### Large deployment for production isolation of processing

To run a larger scale Pega deployment in production, you can split additional processing out to dedicated tiers.  The [values-large.yaml](./values-large.yaml) configuration provides an example of a multi-tier deployment that Pega recommends as a good starting point for larger deployments.

Tier Name   | Description
---         | ---
web         | Interactive, foreground processing nodes that are exposed to the load balancer. Pega recommends that these node use the node classification “WebUser” `nodetype`.
batch       | Background processing nodes which handle some of the non-interactive processing.  Pega recommends that these node use the node classification   “BackgroundProcessing,Search,Batch” `nodetype`. These nodes should not be exposed to the load balancer.
stream      | Nodes that run an embedded deployment of Kafka and are exposed to the load balancer. Pega recommends that these node use the node classification “Stream” `nodetype`.
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

[Node types for client-managed cloud environments](https://community.pega.com/knowledgebase/articles/performance/node-classification)

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

### service

Specify the `service` yaml block to expose a Pega tier to other Kubernetes run services, or externally to other systems. The name of the service will be based on the tier's name, so if your tier is "web", your service name will be "pega-web". If you omit service, no Kubernetes service object is created for the tier during the deployment. For more information on services, see the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/services-networking/service).

Configuration parameters:

Parameter | Description                       | Default value
---       | ---                               | ---
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

Specify the `ingress` yaml block to expose a Pega tier to access from outside Kubernetes. Pega supports the use of managing SSL certificates for HTTPS configuration using a variety of methods. For more information on services, see the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/).

Parameter | Description
---       | ---
`domain`  | Specify a domain on your network in which you create an ingress to the load balancer.
`appContextPath`  | Specify the path for access to the Pega application on a specific tier. If not specified, users will have access to the Pega application via /prweb 
`tls.enabled` | Specify the use of HTTPS for ingress connectivity. If the `tls` block is omitted, TLS will not be enabled.
`tls.secretName` | Specify the Kubernetes secret you created in which you store your SSL certificate for your deployment. For compatibility, see [provider support for SSL certificate injection](#provider-support-for-ssl-certificate-management).
`tls.useManagedCertificate` | On GKE, set to `true` to use a managed certificate; otherwise use `false`.
`tls.ssl_annotation` | On GKE or EKS, set this value to an appropriate SSL annotation for your provider.
`annotations` | Optionally add custom annotations for advanced configurations. For Kubernetes and EKS deployments, including custom annotations overrides the default configuration; for GKE and AKS deployments, the deployment appends these custom annotations to the default list of annotations.

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

Parameter       | Description    | Default value
---             | ---       | ---
`replicas`      | Specify the number of Pods to deploy in the tier. | `1`
`cpuRequest`    | Initial CPU request for pods in the current tier.  | `2`
`cpuLimit`      | CPU limit for pods in the current tier.  | `4`
`memRequest`    | Initial memory request for pods in the current tier. | `6Gi`
`memLimit`      | Memory limit for pods in the current tier. | `8Gi`
`initialHeap`   | This specifies the initial heap size of the JVM.  | `4096m`
`maxHeap`       | This specifies the maximum heap size of the JVM.  | `7168m`

### nodeSelector

Pega supports configuring certain nodes in your Kubernetes cluster with a label to identify its attributes, such as persistent storage. For such configurations, use the Pega Helm chart nodeSelector property to assign pods in a tier to run on particular nodes with a specified label. For more information about assigning Pods to Nodes including how to configure your Nodes with labels, see the [Kubernetes documentation on nodeSelector](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector).

```yaml
tier:
- name: "stream"
  nodeType: "Stream"

  nodeSelector:
    disktype: ssd
```

### Liveness and readiness probes

[Probes are used by Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) to determine application health.  Configure a probe for *liveness* to determine if a Pod has entered a broken state; configure it for *readiness* to determine if the application is available to be exposed.  You can configure probes independently for each tier.  If not explicitly configured, default probes are used during the deployment.  Set the following parameters as part of a `livenessProbe` or `readinessProbe` configuration.

Parameter           | Description    | Default value
---                 | ---            | ---
`initialDelaySeconds` | Number of seconds after the container has started before liveness or readiness probes are initiated. | `300`
`timeoutSeconds`      | Number of seconds after which the probe times out. | `20`
`periodSeconds`       | How often (in seconds) to perform the probe. Some providers such as GCP require this value to be greater than the timeout value. | `30`
`successThreshold`    | Minimum consecutive successes for the probe to be considered successful after it determines a failure. | `1`
`failureThreshold`    | The number consecutive failures for the pod to be terminated by Kubernetes. | `3`

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
`hpa.targetAverageCPUUtilization` | Threshold value for scaling based on initial CPU request utilization (The default value is `70` which corresponds to 70% of 2) | `70`
`hpa.targetAverageMemoryUtilization` | Threshold value for scaling based on initial memory utilization (The default value is `85` which corresponds to 85% of 6Gi ) | `85`

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

### Custom Annotations for Pods

You may optionally provide custom annotations for Pods as metadata to be consumed by other tools and libraries. Pod annotations may be specified by using the `podAnnotations` element for a given `tier`.

Example:

```yaml
tier:
  - name: my-tier
    podAnnotations:
      <annotation-key>: <annotation-value>
```

### Pega configuration files

While Pega includes default configuration files in the Helm charts, the charts provide extension points to override the defaults with additional customizations. To change the configuration file, specify the replacement implementation to be injected into a ConfigMap.

Parameter     | Description    | Default value
---           | ---       | ---
`prconfig`    | A complete prconfig.xml file to inject.  | See [prconfig.xml](config/deploy/prconfig.xml).
`prlog4j2`    | A complete prlog4j2.xml file to inject.  | See [prlog4j2.xml](config/deploy/prlog4j2.xml).
`contextXML`  | A complete context.xml template file to inject.  | See [context.xml.tmpl](config/deploy/context.xml.tmpl).


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
```

### Pega diagnostic user

While most cloud native deployments will take advantage of aggregated logging using a tool such as EFK, there may be a need to access the logs from Tomcat directly. In the event of a need to download the logs from tomcat, a username and password will be required.  You may set `pegaDiagnosticUser` and `pegaDiagnosticPassword` to set up authentication for Tomcat.

## Cassandra and DDS deployment

If you are planning to use Cassandra (usually as a part of Pega Decisioning), you may either point to an existing deployment or deploy a new instance along with Pega. 

### Using an existing Cassandra deployment

To use an existing Cassandra deployment, set `cassandra.enabled` to `false` and configure the `dds` section to reference your deployment.

Example:

```yaml
cassandra:
  enabled: false

dds:
  externalNodes: "CASSANDRA_NODE_IPS"
  port: "9042"
  username: "cassandra_username"
  password: "cassandra_password"
```

### Deploying Cassandra with Pega

You may deploy a Cassandra instance along with Pega.  Cassandra is a separate technology and needs to be independently managed.  When deploying Cassandra, set `cassandra.enabled` to `true` and leave the `dds` section as-is.  For more information about configuring Cassandra, see the [Cassandra Helm charts](https://github.com/helm/charts/blob/master/incubator/cassandra/values.yaml).

*Cassandra minimum resource requirements*

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

## Search deployment

Use the `pegasearch` section to configure a deployment of ElasticSearch for searching Rules and Work within Pega.  This deployment is used exclusively for Pega search, and is not the same ElasticSearch deployment used by the EFK stack or any other dedicated service such as Pega BI.

Parameter   | Description   | Default value
---         | ---           | ---
`image`   | Set the `pegasearch.image` location to a registry that can access the Pega search Docker image. The image is [available on DockerHub](https://hub.docker.com/r/pegasystems/search), and you may choose to mirror it in a private Docker repository. | `pegasystems/search:latest`
`imagePullPolicy` | Optionally specify an imagePullPolicy for the search container. | `""`
`replicas` | Specify the desired replica count. | `1`
`minimumMasterNodes` | To prevent data loss, you must configure the minimumMasterNodes setting so that each master-eligible node is set to the minimum number of master-eligible nodes that must be visible in order to form a cluster. Configure this value using the formula (n/2) + 1 where n is replica count or desired capacity.  For more information, see the ElasticSearch [important setting documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html) for more information. | `1`
`podSecurityContext.runAsUser`   | ElasticSearch defaults to UID 1000.  In some environments where user IDs are restricted, you may configure your own using this parameter. | `1000`
`set_vm_max_map_count`   | Elasticsearch uses a **mmapfs** directory by default to store its indices. The default operating system limits on mmap counts is likely to be too low, which may result in out of memory exceptions. An init container is provided to set the value correctly, but this action requires privileged access. If privileged access is not allowed in your environment, you may increase this setting manually by updating the `vm.max_map_count` setting in **/etc/sysctl.conf** according to the ElasticSearch documentation and can set this parameter to `false` to disable the init container. For more information, see the [ElasticSearch documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html). | `true`
`set_data_owner_on_startup`   | Set to true to enable an init container that runs a chown command on the mapped volume at startup to reset the owner of the ES data to the current user. This is needed if a random user is used to run the pod, but also requires privileges to change the ownership of files. | `false`
`podAnnotations` | Configurable annotations applied to all Elasticsearch pods. | {}

Additional env settings supported by ElasticSearch may be specified in a `custom.env` block as shown in the example below.

Example:

```yaml
pegasearch:
  image: "pegasystems/search:8.3"
  memLimit: "3Gi"
  replicas: 1
  minimumMasterNodes: 2
  custom:
    env:
    - name: TZ
      value: "EST5EDT"
```
## Pega database installation and upgrades

Pega requires a relational database that stores the rules, data, and work objects used and generated by Pega Platform. The [Pega Platform deployment guide](https://community.pega.com/knowledgebase/products/platform/deploy) provides detailed information about the requirements and instructions for installations and upgrades.  Follow the instructions for Tomcat and your environment's database server.

The Helm charts also support an automated install or upgrade with a Kubernetes Job.  The Job utilizes an installation Docker image and can be activated with the `action` parameter in the Pega Helm chart.
 
### Install
 
For installations of the Pega platform, you must specify the installer Docker image and an initial default password for the `administrator@pega.com` user.

Example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
  adminPassword: "ADMIN_PASSWORD"
```

### Upgrade

For upgrades of the Pega platform, you must specify the installer Docker image and the type of upgrade to execute.

Upgrade type    | Description
---             | ---
`in-place`      | An in-place upgrade will upgrade both rules and data in a single run.  This will upgrade your environment as quickly as possible but will result in downtime.
`out-of-place`  | An out-of-place upgrade involves more steps to minimize downtime.  It will place the rules into a read-only state, then migrate the rules to a new schema. Next it will upgrade the rules to the new version. Lastly it will separately upgrade the data.

Example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
  upgrade:
    upgradeType: "out-of-place"
    targetRulesSchema: "rules_upgrade"
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
