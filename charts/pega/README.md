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
pks         | VMware Tanzu Kubernetes Grid Integrated Edition (TKGI), which used to be Pivotal Container Service (PKS)
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

When you download Docker images, it is recommended that you specify the absolute image version and the image name instead of using the `latest` tag; for example: `pegasystems/pega:8.4.4` or `platform-services/search-n-reporting-service:1.12.0`. When you download these images with these details from the Pega repository, you pull the latest available image. If you pull images only specifying `latest`, you may not get the image you wanted.

For this reason, it is also recommended that you specify the `docker.pega.imagePullPolicy: "IfNotPresent"` option in production, since it will ensure that a new generic tagged image will not overwrite the locally cached version.

Example:

 ```yaml
docker:
  registry:
    url: "YOUR_DOCKER_REGISTRY"
    username: "YOUR_DOCKER_REGISTRY_USERNAME"
    password: "YOUR_DOCKER_REGISTRY_PASSWORD"
  pega:
    image: "pegasystems/pega:8.4.4"
    imagePullPolicy: "Always"
```

## Deploying with Pega-provided busybox and k8s-wait-for utility images
To deploy Pega Platform, the Pega helm chart requires the use of the busybox and k8s-wait-for images. For clients who want to pull these images from a registry other than Docker Hub, they must tag and push these images to another registry, and then pull these images by specifying `busybox` and `k8s-wait-for` values as described below.

Example:

 ```yaml
utilityImages:
  busybox:
    image: "busybox:1.31.0"
    imagePullPolicy: "IfNotPresent"
  k8s_wait_for:
    image: "dcasavant/k8s-wait-for"
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

Pega supports deployments using a multi-tier architecture model that separates application processing from k8s functions. Isolating processing in its own tier supports unique deployment configurations, including the Pega application prconfig, resource allocations, and scaling characteristics. Use the tier section in the helm chart to specify into which tiers you wish to deploy the tier with nodes dedicated to the logical tasks of the tier.

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
### Security context

By default, security context for your Pega pod deployments `pegasystems/pega` image uses `pegauser`(9001) as the user. To configure an alternative user for your custom image, set value for `runAsUser`. Note that pegasystems/pega image works only with pegauser(9001).
`runAsUser` must be configured in `securityContext` under each tier block and will be applied to Deployments/Statefulsets, see the [Kubernetes Documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/).

Example:

```yaml
tier:
  - name: my-tier
    securityContext:
      runAsUser: RUN_AS_USER
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

Parameter       | Description                                            | Default value
---             | ---                                                    | ---
`replicas`      | Specify the number of Pods to deploy in the tier.      | `1`
`cpuRequest`    | Initial CPU request for pods in the current tier.      | `3`
`cpuLimit`      | CPU limit for pods in the current tier.                | `4`
`memRequest`    | Initial memory request for pods in the current tier.   | `12Gi`
`memLimit`      | Memory limit for pods in the current tier.             | `12Gi`
`initialHeap`   | Specify the initial heap size of the JVM.              | `4096m`
`maxHeap`       | Specify the maximum heap size of the JVM.              | `8192m`

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
- name: "stream"
  nodeType: "Stream"

  nodeSelector:
    disktype: ssd
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

### Pega configuration files

While Pega includes default configuration files in the Helm charts, the charts provide extension points to override the defaults with additional customizations. To change the configuration file, specify the replacement implementation to be injected into a ConfigMap.

Parameter     | Description    | Default value
---           | ---       | ---
`prconfig`    | A complete prconfig.xml file to inject.  | See [prconfig.xml](config/deploy/prconfig.xml).
`prlog4j2`    | A complete prlog4j2.xml file to inject.  | See [prlog4j2.xml](config/deploy/prlog4j2.xml).
`contextXML`  | A complete context.xml template file to inject.  | See [context.xml.tmpl](config/deploy/context.xml.tmpl).
`serverXML`   | A complete server.xml file to inject  | See [server.xml](config/deploy/server.xml).
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

### Pega diagnostic user

While most cloud native deployments will take advantage of aggregated logging using a tool such as EFK, there may be a need to access the logs from Tomcat directly. In the event of a need to download the logs from tomcat, a username and password will be required.  You may set `pegaDiagnosticUser` and `pegaDiagnosticPassword` to set up authentication for Tomcat.

## Cassandra and Pega Customer Decision Hub deployments

If you are planning to use Cassandra (usually as a part of Pega Customer Decision Hub), you may either point to an existing deployment or deploy a new instance along with Pega.

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

## Search deployment

Use the `pegasearch` section to configure the source ElasticSearch service that the Pega Platform deployment uses for searching Rules and Work within Pega. The ElasticSearch service defined here is not related to the ElasticSearch deployment if you also define an EFK stack for logging and monitoring in your Pega Platform deployment.
### For Pega Platform 8.6 and later:

Pega recommends using the chart ['backingservices'](../backingservices) to enable Pega Infinity backing service and to deploy the latest generation of search and reporting capabilities to your Pega applications that run independently on nodes provisioned exclusively to run these services.
This is an independently manageable search solution that replaces the previous implementation of Elasticsearch. The SRS supports, but does not require you to enable, Elasticsearch for your Pega Infinity deployment searching capability.

To use this search and reporting service, follow the deployment instructions provided at ['backingservices'](../backingservices) before you configure and deploy `pega` helm chart. You must configure the SRS URL for your Pega Infinity deployment using the parameter in values.yaml as shown the the following table and exmple:

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

Use the below configuration to provision an internally deployed instance of elasticsearch for search functionality within the platform:

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
 
### Installations

For installations of the Pega platform, you must specify the installer Docker image and an initial default password for the `administrator@pega.com` user.

Example:

```yaml
installer:
  image: "YOUR_INSTALLER_IMAGE:TAG"
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
`adminPassword` | Specify a temporary, initial password to log into teh Pega application. This will need to be changed at first login. The adminPassword value cannot start with "@". | `"ADMIN_PASSWORD"`
`upgrade.upgradeType:` |Specify the type of process, applying a patch or upgrading. | See the next table for details.
`upgrade.upgradeSteps:` |Specify the steps of a `custom` upgrade process that you want to complete. For `zero-downtime`, `out-of-place-rules`, `out-of-place-data`, or `in-place` upgrades, leave this parameter empty. | <ul>`enable_cluster_upgrade` `rules_migration` `rules_upgrade` `data_upgrade` `disable_cluster_upgrade`</ul>
`upgrade.targetRulesSchema:` |Specify the name of the schema you created the process creates for the new rules schema. | `""`
`upgrade.targetDataSchema:` | For patches to 8.4 and later or upgrades from 8.4.2 and later, specify the name of the schema the process creates for the temporary data schema. After the patch or upgrade, you must delete this temporary data schema from your database. For 8.2 or 8.3 Pega software patches, you can leave this value empty, as is (do not add blank text). | `""`

Upgrade type    | Description
---             | ---
`zero-downtime` |  If applying any patch or upgrading from 8.4.2 and later, use this option to minimize downtime. This patch or upgrade type migrates the rules to your designated "new rules schema", uses the temporary data schema to host some required data-specific tables, and patches (only changed rules) or upgrades (all) the rules to the new version with zero-downtime. With the new rules in place, the process performs a rolling reboot of your nodes, patches or upgrades any required data schema, and redeploys the application using the new rules.
`custom` |  Use this option for any upgrade in which you complete portions of the upgrade process in steps. Supported upgrade steps are: `enable_cluster_upgrade` `rules_migration` `rules_upgrade` `data_upgrade` `disable_cluster_upgrade`. To specify which steps to include in your custom upgrade, specify them in your pega.yaml file using the `upgrade.upgradeSteps` parameter.
`out-of-place-rules` | Use this option to complete an out-of-place upgrade of rules that the upgrade process migrates to the new rules schema. This schema will become the rules schema after your upgrade is complete.
`out-of-place-data` |Use this option to complete an out-of-place upgrade of data the upgrade process migrates to a new, temporary data schema. The upgrade process removes this temporary schema after your application is running with updated data.
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
### Mount the custom certificates into the Tomcat container

Pega supports mounting and passing custom certificates into the tomcat container during your Pega Platform deployment. Pega supports the following certificate formats as long as they are encoded in base64: X.509 certificates such as PEM, DER, CER, CRT. To mount and pass the your custom certificates, use the `certificates` attributes as a map in the `values.yaml` file using the format in the following example.

Example:

```yaml
certificates:
    badssl.cer: |
      "-----BEGIN CERTIFICATE-----\n<<certificate content>>\n-----END CERTIFICATE-----\n"

```