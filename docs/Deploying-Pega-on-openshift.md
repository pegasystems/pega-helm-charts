# Deploying Pega Platform on a Red Hat OpenShift container platform cluster

Deploy Pega Platform™ on a Red Hat OpenShift container platform cluster using a PostgreSQL database. The following procedures are written for any level of user, from a system administrator to a development engineer who is interested in learning how to install and deploy Pega Platform onto an OpenShift cluster.

Pega helps enterprises and agencies quickly build business apps that deliver the outcomes and end-to-end customer experiences that you need. Use the procedures in this guide, to install and deploy Pega software onto an OpenShift cluster without much experience in either OpenShift configurations or Pega Platform deployments.

Create a deployment of Pega Platform on which you can implement a scalable Pega application in a OpenShift cluster using Red Hat OpenShift Container Platform (Self Managed). You can use this deployment for a Pega Platform development environment. By completing these procedures, you deploy Pega Platform on a OpenShift cluster with a PostgreSQL database instance.

## Supported Products

Pegasystems supports deployment of Pega Platform using [OpenShift Container Platform](https://www.redhat.com/en/technologies/cloud-computing/openshift/container-platform) as a self managed solution. Managed solutions, other OpenShift products, and the IBM Cloud are not currently supported.

## Deployment process overview

Use Kubernetes tools and the customized orchestration tools and Docker images to orchestrate a deployment in a OpenShift cluster that you create for the deployment:

1. Prepare your local system:

   - To prepare a local Linux system, install required applications and configuration files - [Preparing your local Linux system – 45 minutes](prepping-local-system-runbook-linux.md).

   - To prepare a local Windows system, install required applications and configuration files -
   [Preparing your local Windows 10 system – 45 minutes](prepping-local-system-runbook-windows.md).

2. Create a cluster using your preferred cloud resources and install Red Hat OpenShift Container Platform (Self Managed) on the cluster. For details, see [OpenShift Documentation](https://docs.openshift.com/).

3. Customize two configuration files with which you configure your OpenShift cluster to run Pega in a web application cluster and run backing services, including search, in a Pega search and reporting service cluster. In this section, you will use the command-line tools, kubectl and Helm, to install and then deploy Pega Platform onto your existing OpenShift cluster - [Deploying Pega Platform using Helm charts – 90 minutes](#installing-and-deploying-pega-platform-using-helm-charts--90-minutes).

4. Configure your network connections in the DNS management zone of your choice so you can log in to Pega Platform - [Logging in to Pega Platform – 10 minutes](#logging-in-to-pega-platform--10-minutes).

To understand how Pega maps Kubernetes objects with Pega applications and services, see [Understanding the Pega deployment architecture](https://community.pega.com/knowledgebase/articles/client-managed-cloud/cloud/understanding-pega-deployment-architecture).

## Assumptions and prerequisites

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges or a Linux command prompt with root privileges.

- You use open source packaging tools on Windows or Linux to install applications onto your local system.

The following account, resources, and application versions are required for use in this document:

- An  account with a payment method set up to pay for the cloud infrastructure resources you create and appropriate account permissions and knowledge to:
  - Select an appropriate location to deploy your database resource and OpenShift cluster. The cluster and PostgreSQL database into which you install Pega Platform must be in the same location.
  - Create an Openshift cluster or use an existing cluster you can access.
  - Create an SQL instance using a cloud resource that is accessible to the pods in the Openshift deployment using a supported database. For support details, see [URL and Driver Class](https://github.com/pegasystems/pega-helm-charts/tree/master/charts/pega).

  You are responsible for any financial costs incurred for your cloud infrastructure resources.

- Pega Platform 8.6 or later, which supports the latest Pega software feature support, including the Search and Reporting Service (SRS).

- Pega Docker images – your deployment requires the use of several Docker images that you download and make available in a private Docker registry. For step-by-step details, see [Downloading and managing Pega Platform docker images (linux)](prepping-local-system-runbook-linux.md#downloading-and-managing-pega-platform-docker-images) or [Downloading and managing Pega Platform docker images (windows)](prepping-local-system-runbook-windows.md#downloading-and-managing-pega-platform-docker-images).

- Helm 3.0 or later. Helm is only required to use the Helm charts and not to use the Kubernetes Yaml examples directly. For more information, see the [Helm documentation portal](https://helm.sh/docs/).

- kubectl – the Kubernetes command-line tool that you use to connect to and manage your Kubernetes resources.

- openshift-cli - the OpenShift command-line tool that you use to connect to your OpenShift cluster.

## Creating a Red Hat OpenShift project - One hour

There are multiple ways of creating an Openshift cluster. For an overview, see [OpenShift Container Platform installation overview](https://docs.openshift.com/container-platform/4.7/installing/index.html#ocp-installation-overview).

### Creating a database resource

Pega Platform deployments on OpenShift require you to install Pega Platform software in an SQL database. After you create an SQL instance that is available to your OpenShift cluster, you must create a PostgreSQL database in which to install Pega Platform. When you are finished, you will need the database name and the SQL instance IP address that you create in this procedure in order to add this information to your pega.yaml Helm chart.

#### Creating an SQL instance

Create an SQL instance that is available to your OpenShift cluster. In this example, the SQL instance is created in GCP; however, you can create or use any database resource that is available to the OpenShift cluster.

1. Use a web browser to log in to <https://cloud.google.com/> and navigate to your
    **Console** in the upper right corner.

2. In your **Google Cloud Platform** console, use the **Navigation Menu** to go
    to **Storage** grouping and select **SQL**.

3. On the **SQL** page, click **+Create Instance**.

4. On the **Choose your database engine** window, click **Choose PostgreSQL**.

5. On the **Create PostgreSQL instance** page, add details to the following
    required fields for this database server:

    a. **Instance ID**, enter a database server ID. (demo-runbook-sql1)

    b. **Default user password**, enter a “postgres” user password.

    c. Select an appropriate **Region** and **Zone** for your database server. Select the same zone or region that you used to create your OpenShift cluster.

    d. **Database version**, select **PostgreSQL 11 or 12**.

    e. **Configuration options \> Connectivity**, select **Public IP**, click **+ Add Network**, enter a **Name** and **Network** of one or more IP addresses to allow access to this PostgreSQL database, and click **Done**.

    As a best practice, add the following IP addresses: your local system from where you install helm, the worker nodes of the cluster. One method for finding the IP address for worker nodes of the cluster is to view the nodes in your OpenShift cluster with kubectl command-line tool and then use the command options, `kubectl describe nodes | grep ExternalIP`.

6. In **Configuration options \> Machine type and storage**:

    a. **Machine type**, select n1-highmem-2 (2 vCPU **Cores** and 13 GB **Memory**).

    b. **Network throughput**, select **SSD (Recommended)**.

    c. **Storage capacity**, enter **20 GB** and select **Enable automatic storage increases**.

7. Configure the remaining setting using the default values:

    a. **Auto backups and high availability** - select an automated backup time.

    b. **Flags** - no flags are required.

    c. **For Maintenance** - any preference is supported.

    d.  **Labels**, no labels are required. Labels can help clarify billing details for your OpenShift resources.

8. Click **Create**.

  A deployment progress page displays the status of your deployment until it is complete, which takes up to 5 minutes. When complete, the GCP UI displays all of the SQL resources in your account, which includes your newly created SQL instance:

![cid:image007.png\@01D5A3B1.62812F70](media/9aa072ea703232c2f6651fe95854e8dc.62812f70)

#### Creating a database in your SQL instance

Create a PostgreSQL database in your new SQL instance for the Pega Platform installation. Use the database editing tool of your choice to log into your SQL instance and create this new PostgreSQL database. The following example was completed using pgAdmin 4.

1. Log into your SQL instance.

    You can find your access information and login credentials, by selecting the SQL instance in the GCP console.

2. In the database editor tool, navigate to Databases and create a new database.

   No additional configuration is required.

With your SQL service IP address and your new database name, you are ready to continue to the next section.

## Installing and deploying Pega Platform using Helm charts – 90 minutes

To deploy Pega Platform by using Helm, customize the backing-services.yaml and pega.yaml Helm chart that hold the specific settings for your deployment needs and then run a series of Helm commands to complete the deployment.

An installation with deployment will take about 90 minutes total, because a Pega Platform installation in your PostgreSQL database takes up to an hour.

### Adding the Pega configuration files to your Helm installation on your local system

Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm, including a generic version of the following charts. After you add the repository to your local system, you can customize these Pega configuration files for your Pega Platform deployment:

- pega/pega - Use this chart to set customization parameters for your deployment. You will modify this chart later in the deployment tasks.

- pega/backingservices - Use this chart to set customization parameters for the Pega-provided Search and Reporting Service (SRS) your deployment. You will modify this chart later in the deployment tasks.

To customize these files, you must download them from the source github repository to your local system, edit them with a text editor, and then save them to your local system using the same filename.

1. To add the Pega repository to your Helm installation, enter:

    `$ helm repo add pega https://pegasystems.github.io/pega-helm-charts`

2. To verify the new repository, you can search it by entering:

```bash
  $ helm search repo pega
  NAME                  CHART VERSION   APP VERSION     DESCRIPTION
  pega/pega             2.2.0                           Helm chart to configure required installation and deployment configuration settings in your environment for your deployment.
  pega/addons           2.2.0           1.0             Helm chart to configure supporting services and tools in your environment for your deployment.
  pega/backingservices  2.2.0                           Helm Chart to provision the latest Search and Reporting Service (SRS) for your Pega Infinity deployment.
```

The addons charts is not required for OpenShift deployments of Pega Platform. Use of the backingservices chart is optional, but recommended for delivering your ElasticSearch service for Pega Infinity 8.6 and later.

### Updating the backingservices.yaml Helm chart values for the SRS (Supported when installing or upgrading to Pega Infinity 8.6 and later)

To configure the parameters in the backingservices.yaml file, download the file in the charts/backingservices folder, edit it with a text editor, and then save it to your local system using the same filename.

Configure the following parameters so the backingservices.yaml Helm chart matches your Elastic resources in these areas:

- Your docker registry credentials.

- Your deployment name with which the deployment appends the names of your SRS pods.

- Your Pega-provided SRS Docker image in the repository to which you pushed the image.

- Your choice to enable either:
  - An ElasticSearch service that is defined at [ElasticSearch Helm Chart](https://github.com/elastic/helm-charts/tree/master/elasticsearch) which the SRS can provide to your deployment.
  - Your ElasticSearch cluster for which you specify the required authentication details.
- Add additional, recommended environmental variables for ElasticSearch basic authentication. For details, see [Examples/openshift](https://github.com/elastic/helm-charts/blob/master/elasticsearch/examples/openshift/values.yaml).

- Three additional parameters that ElasticSearch recommends for OpenShift deployments.

1. To download the backingservices.yaml Helm chart to the \<local filepath>\openshift-demo, enter:

   `$ helm inspect values pega/backingservices > <local filepath>/openshift-demo/backingservices.yaml`

2. Use a text editor to open the backingservices.yaml file and update the following parameters in the chart based on your OpenShift requirements:

| Chart parameter name    | Purpose                                   | Your setting |
|-------------------------|-------------------------------------------|--------------|
|global.imageCredentials.registry: username: password:  | Include the URL of your Docker registry along with the registry “username” and “password” credentials. | <ul><li>url: “\<URL of your registry>” </li><li>username: "\<Registry account username\>"</li><li> password: "\<Registry account password\>"</li></ul>|
|srs.enabled:      | Enable the SRS to provision an internal ElasticSearch service within the SRS cluster that is defined at [ElasticSearch Helm Chart](https://github.com/elastic/helm-charts/tree/master/elasticsearch).    | enabled: "true"|
|srs.deploymentName:      | Specify unique name for the deployment based on org app and/or SRS applicable environment name.      | deploymentName: "acme-demo-dev-srs"|
|srs.srsRuntime.srsImage: | Specify the Pega-provided SRS Docker image that you downloaded and pushed to your Docker registry. | srs.srsRuntime.srsImage: "\<Registry host name:Port>my-pega-srs:\<srs-version>". For `<srs-version>` tag details, see [SRS Version compatibility matrix](../charts/backingservices/charts/srs/README.md#srs-version-compatibility-matrix). |
| srs.srsStorage.provisionInternalESCluster: | Enabled by default to provision an Elasticsearch cluster. | <ul><li>Set srs.srsStorage.provisionInternalESCluster:`true` and run `$ make es-prerequisite NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT>`</li><li>Set srs.srsStorage.provisionInternalESCluster:`false` if you want to use an existing, externally provisioned ElasticSearch cluster. </li></ul> |
|srs.srsStorage.domain: port: protocol: basicAuthentication: awsIAM: requireInternetAccess: | Disabled by default. Enable only when srs.srsStorage.provisionInternalESCluster is false and you want to configure SRS to use an existing, externally provisioned ElasticSearch cluster. For an ElasticSearch cluster secured with Basic Authentication, use `srs.srsStorage.basicAuthentication` section to provide access credentials. For an AWS ElasticSearch cluster secured with IAM role based authentication, use `srs.srsStorage.awsIAM` section to set the aws region where AWS ElasticSearch cluster is hosted. For unsecured managed ElasticSearch cluster do not configure these options. | <ul><li>srs.srsStorage.domain: "\<external-es domain name\>"</li><li>srs.srsStorage.port: "\<external es port\>"</li><li>srs.srsStorage.protocol: "\<external es http protocol, `http` or `https`\>"</li><li>srs.srsStorage.basicAuthentication.username: "\<external es `basic Authentication username`\>"</li><li>srs.srsStorage.basicAuthentication.password: "\<external es `basic Authentication password`\>"</li>     <li>srs.srsStorage.awsIAM.region: "\<external AWS es cluster hosted `region`\>"</li><li> srs.srsStorage.requireInternetAccess: "\<set to `true` if you host your external ElasticSearch cluster is outside of the current network and the deployment must access it over the internet.\>"</li></ul> |
|elasticsearch: <ul><li>esConfig: elasticsearch.yml: xpack.security.enabled: "true”</li><li>esConfig: elasticsearch.yml: xpack.security.transport.ssl.enabled: "true”</li><li>volumeClaimTemplate: resources: requests: storage:  "\<30Gi>” </li></ul> | To run an internal ElasticSearch service (within your deployment), specify details for both the ElasticSearch security configuration and a cluster disk volume size. While the security pack should be enabled, the transport configuration should match whether an SSL certificate is available; by default it is enabled. For the disk volume size, the default is 30Gi; set this value to at least three times the size of your estimated search data size. | <ul><li>elasticsearch: esConfig: elasticsearch.yml: <ul><li>xpack.security.enabled: "true” </li><li>xpack.security.transport.ssl.enabled: "true” </li></ul></li><li>elasticsearch: volumeClaimTemplate: resources: requests: storage:  "\<30Gi>” </li></ul> |
|Additional OpenShift-required parameters:<ul><li>securityContext:</li><li>podSecurityContext:</li><li>sysctlInitContainer:</li></ul> | Manually add ElasticSearch-recommended parameters to ensure that your SRS pods do not have a RunAsUser parameter and do not require an initialization container. For details, see [Examples/OpenShift](https://github.com/elastic/helm-charts/tree/master/elasticsearch/examples/openshift). | <ul><li>securityContext.runAsUser: null</li><li>podSecurityContext.fsGroup: null</li><li>    podSecurityContext.runAsUser: null </li><li>sysctlInitContainer.enabled: false</li></ul> |

3. Save the file.

4. To use an internal Elasticsearch cluster (srs.srsStorage.provisionInternalESCluster:true) for your deployment, you must run `$ make es-prerequisite NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT>`.

### Adding customized settings for Pega to your deployment

The Pega deployment model supports advanced configurations to fit most existing
clients' needs. If you are a Pega client and have known, required customizations
for your deployment and you already use the following files to add your known
customizations, you can copy those configurations into the configuration files
Pega added for this purpose in the [pega-helm-charts](https://github.com/pegasystems/pega-helm-charts) repository folder, pega-helm-charts/charts/pega/config/deploy:

- context.xml: add additional required data sources

- prlog4j2.xml: modify your logging configuration, if required

- prconfig.xml: adjust the standard Pega Platform configuration with known,
    required settings

Make these changes before you begin deploying Pega Platform using Helm charts.

### Updating the pega.yaml Helm chart values

To configure the parameters in the pega.yaml Helm, download the file in the charts/pega folder, edit it with a text editor, and then save it to your local system using the same filename.

Configure the following parameters so the pega.yaml Helm chart matches your deployment resources in these areas:

- Specify that this is an OpenShift deployment.

- Specify to deploy your configuration to ensure that when you invoke this helm configuration file, it invokes a deployment, not an install or update.

- Credentials for your DockerHub account in order to access the required Docker images.

- Access your SQL database (in this example one configured in GCP).

- Access your ElasticSearch service (For 8.6 and later, Pega recommends deploying your service using an SRS cluster).

- Install the version of Pega Platform that you built into your Docker installation image.

- Specify host names for your web and stream tiers.

- Enable encryption of traffic between the ingress/load balancer and the pods by specifying SSL certificates for your web tiers.

- Enable Hazelcast client-server model for Pega Platform 8.6 and later.

- For new deployments, Pega recommends deploying Pega Platform using an externalized Kafka configuration as a stream service provider to use your own managed Kafka infrastructure. Deployment of stream with externalized Kafka configuration requires Pega Infinity 8.4 or later.

1. To download the pega.yaml Helm chart to the \<local filepath\>/openshift-demo, enter:

   `$ helm inspect values pega/pega > /<local filepath>/openshift-demo/pega.yaml`

   You can compare this file to the Pega-provided example Pega helm chart that was used to test an OpenShift deployment, [pega-openshift.yaml](./resources/pega-openshift.yaml) to help you see what parameters can be further customized for your OpenShift cluster.

2. Use a text editor to open the pega.yaml file and update the following parameters in the chart based on your OpenShift requirements:

    | Chart parameter name    | Purpose                                   | Your setting |
    |-------------------------|-------------------------------------------|--------------|
    | provider:               | Specify a OpenShift deployment.           | provider:"openshift" |
    | actions.execute:        | Specify a “deploy” deployment type.       | execute: "deploy"   |
    | Jdbc.url:               | Specify the database IP address and database name for your Pega Platform installation. | <ul><li>url: "jdbc:postgresql://**localhost**:5432/**dbName**"</li><li>where **localhost** is the public IP address you configured for your database connectivity and **dbName** is the name you entered for your PostgreSQL database in [Creating a database resource](#creating-a-database-resource).</li></ul>  |
    | Jdbc.driverClass:       | Specify the driver class for a PostgreSQL database. | driverClass: "org.postgresql.Driver"                |
    | Jdbc.dbType:            | Specify PostgreSQL database type.         | dbType: "postgres”   |
    | Jdbc.driverUri:         | Specify the database driver Pega Platform uses during the deployment.| <ul><li>driverUri: "latest jar file available” </li><li>For PostgreSQL databases, use the URL of the latest PostgreSQL driver file that is publicly available at <https://jdbc.postgresql.org/download.html>.</li></ul> |
    | Jdbc: username: password: | Set the security credentials for your database server to allow installation of Pega Platform into your database.   |<ul><li>username: "\<name of your database user\>" </li><li>password: "\<password for your database user\>"</li><li>For GCP PostgreSQL databases, the default user is “postgres”.</li></ul> |
    | jdbc.rulesSchema: jdbc.dataSchema:  | Set the names of both your rules and the data schema to the values that Pega Platform uses for these two schemas.      | rulesSchema: "rules" dataSchema: "data" |
    | customArtifactory.authentication: basic.username: basic.password: apiKey.headerName: apiKey.value: | To download a JDBC driver from your custom artifactory which is secured with Basic or APIKey Authentication. Use `customArtifactory.authentication.basic` section to provide access credentials or use `customArtifactory.authentication.apiKey` section to provide APIKey value and dedicated APIKey header details. | <ul><li>basic.username: "\<Custom artifactory basic Authentication username\>"</li><li>basic.password: "\<Custom artifactory basic Authentication password\>"</li><li>apiKey.headerName: "\<Custom artifactory dedicated APIKey Authentication header name\>"</li><li>apiKey.value: "\<Custom artifactory APIKey value for APIKey authentication\>"</li> </ul> |
    | customArtifactory.certificate: | Custom artifactory SSL certificate verification is enabled by default. If your custom artifactory domain has a self-signed SSL certificate, provide the certificate. You can disable SSL certificate verification by setting `customArtifactory.enableSSLVerification` to `false`;however, this setting establishes an insecure connection. | <ul><li> certificate: "\<custom artifactory SSL certificate to be verified\>"</li></ul> |
    | docker.registry.url: username: password: | Include the URL of your Docker registry along with the registry “username” and “password” credentials. | <ul><li>url: “\<URL of your registry>” </li><li>username: "\<Registry account username\>"</li><li> password: "\<Registry account password\>"</li></ul> |
    | docker.pega.image:       | Specify the Pega-provided `Pega` image you downloaded and pushed to your Docker registry.  | Image: "<Registry host name:Port\>/my-pega:\<Pega Platform version>" |
    | tier.name: ”web” tier.ingress.domain:| Set a host name for the pega-web service of the DNS zone. Pega supports specifying certificates for an ingress using the same methods that OpenShift supports. Note that if you configure both secrets and pre-shared certificates on the ingress, the load balancer ignores the secrets and uses the list of pre-shared certificates. For details, see [Using multiple SSL certificates in HTTP(s) load balancing with Ingress](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-multi-ssl).  | <ul><li>tier.name: "\<the host name for your web service tier\>" </li><li>Assign this host name with an external IP address and log into Pega Platform with this host name in the URL. Your web tier host name must comply with your networking standards and be available as an external IP address.</li><li>tier.ingress.tls: set to `true` to support HTTPS in the ingress. See step 12 to support the management of the certificates in your deployment.</li></ul> |
    | tier.name: ”web” tier.service.tls:| Set this parameter as `true` to encrypt the traffic between the load balancer/ingress and pods. Beginning with Helm version `2.2.0` Pega provides a default self-signed certificate; Pega also supports specifying your own CA certificate.  | <ul><li>tier.service.tls.enabled: set to `true` to enable the traffic encryption </li><li>tier.service.tls.port: 443</li><li>tier.service.tls.targetPort: 8443</li><li>tier.service.tls.keystore: The base64 encoded content of the keystore file. Leave this value empty to use the default, Pega-provided self-signed certificate.</li><li>tier.service.tls.keystorepassword: the password of the keystore file</li><li>tier.service.tls.cacertificate: the base64 encrypted content of the destination CA certificate. Leave this empty in case of using default Pega shipped certificate or use this parameter to specify the CA certificate chain for the certificate which is present in the keystore. The openshift uses this to validate the endpoint certificate in the keystore, securing the connection from the loadbalancer/ingress to the destination pods</li><li>tier.service.traefik.enabled: set to `false` as this option is for `k8s` provider not for `openshift`</li></ul> |
    | tier.name: ”stream” (Deprecated) tier.ingress.domain: | The "Stream tier" is deprecated, please enable externalized Kafka service configuration under External Services.Set the host name for the pega-stream service of the DNS zone.   | <ul><li>domain: "\<the host name for your stream service tier\>" </li><li>Your stream tier host name should comply with your networking standards.</li><li>tier.ingress.tls: set to `true` to support HTTPS in the ingress and pass the SSL certificate in the cluster using a secret. For details, see step 12 in the section, **Deploying Pega Platform using the command line**.</li><li>To remove the exposure of a stream from external network traffic, delete the `service` and `ingress` blocks in the tier.</li></ul> |
    | pegasearch: | For Pega Platform 8.6 and later, Pega recommends using the chart 'backingservices' to enable Pega SRS. To deploy this service, you must configure your SRS cluster using the backingservices Helm charts and provide the SRS URL for your Pega Infinity deployment. | <ul><li>externalSearchService: true</li><li>externalURL: pegasearch.externalURL For example, http://srs-service.mypega-pks-demo.svc.cluster.local </li></ul> |
    | installer.image:  | Specify the Docker `installer` image for installing Pega Platform that you pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
    | installer. adminPassword:                | Specify an initial administrator@pega.com password for your installation.  This will need to be changed at first login. The adminPassword value cannot start with "@".  | adminPassword: "\<initial password\>"  |
    | hazelcast: | For Pega Platform 8.6 and later, Pega recommends using Hazelcast in client-server model. Embedded deployment would not be supported in future platform releases.| |
    | hazelcast.image:        | Specify the Pega-provided `clustering-service` Docker image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
    | hazelcast.enabled: hazelcast.replicas: hazelcast.username: hazelcast.password: | Either to enable Hazelcast in client-server model and configure the number of replicas and username & passowrd for authentication | <ul><li>enabled: true/false <br/> Set to true if you want to deploy pega platform in client-server Hazelcast model, otherwise false. *Note: Set this value as false for Pega platform versions below 8.6, if not set the installation will fail.* </li><li>replicas: <No. of initial server members to join(3 or more based on deployment)> </li><li>username: "\<UserName for authentication\>" </li><li> password: "\<Password for authentication\>" </li></ul> |
    | stream.enabled: stream.bootstrapServer: stream.securityProtocol: stream.trustStore: stream.trustStorePassword: stream.keyStore: stream.keyStorePassword: stream.saslMechanism: stream.jaasConfig: stream.streamNamePattern: stream.replicationFactor: stream.external_secret_name:| Enable an externalized kafka configuration to connect to your existing stream service, by configuring these required settings | <ul><li>enabled: true/false <br/> Set to true if you want to deploy Pega Platform to use an externalized Kafka configuration, otherwise leave set to false. Note: Pega recommends enabling an externalized Kafka configuration and has deprecated using a stream tier configuration starting at version 8.7.</li><li>bootstrapServer: Provide your existing Kafka broker URLs separated by commas.</li><li>securityProtocol: Provide the required security protocol that your deployment will use to communicate with your existing brokers. Valid values are: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL.</li><li> trustStore: When using a trustStore certificate, you must also include a Kubernetes secret name, that contains the trustStore certificate, in the global.certificatesSecrets parameter. Pega deployments only support trustStores using the Java Key Store (.jks) format.</li><li> trustStorePassword: If required, provide keyStore value in plain text.</li><li> keyStore: When using a keyStore certificate, you must also include a Kubernetes secret name, that contains the keyStore certificate, in the global.certificatesSecrets parameter. Pega deployments only support keyStores using the Java Key Store (.jks) format.<li> keyStorePassword: If required, provide keyStore value in plain text.</li><li> jaasConfig: If required, provide jaasConfig value in plain text.</li><li> saslMechanism: If required, provide SASL Mechanism. Supported values are: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512.</li> <li> streamNamePattern: By default, topics originating from Pega Platform have the "pega-" prefix, so that it is easy to distinguish them from topics created by other applications. Pega supports customizing the name pattern for your Externalized Kafka configuration for each deployment.</li> <li> replicationFactor: Your replicationFactor value cannot be more than the number of Kafka brokers and 3.</li> <li> external_secret_name: To avoid exposing trustStorePassword, keyStorePassword, and jaasConfig parameters, leave the values empty and configure them using an External Secrets Manager, making sure that you configure the keys in the secret in the order:STREAM_TRUSTSTORE_PASSWORD, STREAM_KEYSTORE_PASSWORD and STREAM_JAAS_CONFIG. Enter the external secret name.</li></ul> |
   
3. Save the file.

#### (Optional) Add Support for providing DB credentials using External Secrets Operator

Create two files following the Kubernetes documentation for External Secrets Operator [External Secrets Operator](https://external-secrets.io/v0.5.3/) :
•	An external secret file that specifies what information in your secret to fetch.
•	A secret store to define access how to access the external and placing the required files in your Helm directory.

- Copy both files into the pega-helm-charts/charts/pega/templates directory of your Helm
- Update repo to the latest-> helm repo update pega https://pegasystems.github.io/pega-helm-charts
- Update Pega.yaml file to refer to the external secret manager for DB password.

### Deploying Pega Platform using the command line

A Helm installation and a Pega Platform installation are separate processes. The Helm install command uses Helm to install your deployment as directed in the Helm charts, one in the **charts\\addons** folder and one in the **charts\\pega** folder.

In this document, you specify that the Helm chart always “deploys” by using the setting, actions.execute: “deploy”. In the following tasks, you overwrite this function on your *initial* Helm install by specifying `--set global.actions.execute:install-deploy`, which invokes an installation of Pega Platform using your installation Docker image and then
automatically followed by a deploy. In subsequent Helm deployments, you should not use the override argument, `--set global.actions.execute=`, since Pega Platform is already installed in your database.

1. Open a Linux bash shell and change the location to the top folder of your openshift-demo directory that you created in [Preparing your local Linux system](prepping-local-system-runbook-linux.md).

   `$ cd /home/<local filepath>/openshift-demo`

2. To use the oc command to establish authentication with your OpenShift cluster, login to your cluster, specifying the cluster IP address:

   ```bash
   $ oc login (https:<cluster-IP-address>:<port>)
   Login successful.
   ```

3. To view the nodes in your OpenShift cluster, including cluster names and status, enter:

  ```bash
   $ oc get nodes
   NAME                              STATUS    ROLES     AGE       VERSION
   os-runbook-zcbxd-master-0         Ready     master    47h       v1.20.0+87cc9a4
   os-runbook-zcbxd-master-1         Ready     master    47h       v1.20.0+87cc9a4
   os-runbook-zcbxd-master-2         Ready     master    47h       v1.20.0+87cc9a4
   os-runbook-zcbxd-worker-b-d8tsq   Ready     worker    47h       v1.20.0+87cc9a4
   os-runbook-zcbxd-worker-c-5jlct   Ready     worker    47h       v1.20.0+87cc9a4
   os-runbook-zcbxd-worker-d-rmrsz   Ready     worker    47h       v1.20.0+87cc9a4
   ```

4. To create namespaces in preparation for the pega.yaml and backingservices.yaml deployments, enter:

```bash
   $ kubectl create namespace mypega-openshift-demo
```

5. For Pega Platform 8.6 and later installations, to install the backingservices chart that you updated in [Updating the backingservices.yaml Helm chart values (Supported when installing or upgrading to Pega Infinity 8.6 and later)](#Updating the backingservices.yaml Helm chart values (Supported when installing or upgrading to Pega Infinity 8.6 and later)), enter:

 ```bash
   $ helm install backingservices pega/backingservices --namespace mypega-openshift-demo --values backingservices.yaml
   NAME: backingservices
   LAST DEPLOYED: Fri Jul 16 15:31:58 2021
   NAMESPACE: mypega-openshift-demo
   STATUS: deployed
   REVISION: 1
```

The `mypega-openshift-demo` namespace used for pega deployment can also be used for backingservice deployment that you configured in backingservices.yaml helm chart.

6. To deploy Pega Platform for the first time by specifying to install Pega Platform into the database specified in the Helm chart when you install the pega.yaml Helm chart, enter:

```bash
    $ helm install mypega-openshift-demo pega/pega --namespace mypega-openshift-demo --values pega.yaml --set global.actions.execute=install-deploy
    NAME: mypega-openshift-demo
    LAST DEPLOYED: Fri Jan  3 19:00:19 2020
    NAMESPACE: mypega-openshift-demo
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None
```

For subsequent Helm installs, use the command `helm install mypega-openshift-demo pega/pega --namespace mypega-openshift-demo` to deploy Pega Platform and avoid another Pega Platform installation.

A successful Pega deployment immediately returns details that show progress for your `mypega-openshift-demo` deployment.

7. Refer to your OpenShift dashboard to follow the progress of or troubleshoot an installation. For subsequent deployments, you do not need to do this. Initially, while the resources make requests to complete the configuration, you will see red warnings while the configuration is finishing, which is expected behavior.

8. To view the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the `mypega-openshift-demo` namespace pods.

A successful deployment will not show errors across any various workloads.

## Logging in to Pega Platform – 10 minutes

After you complete your deployment, as a best practice, associate the host name of the pega-web tier ingress with the IP address that the deployment load balancer assigned to the tier during deployment. The host name of the pega-web tier ingress used in this demo, **openshift.web.dev.pega.io**, is set in the pega.yaml file in the following lines:

```yaml
tier:
  - name: "web"

    service:
      # Enter the domain name to access web nodes via a load balancer.
      #  e.g. web.mypega.example.com
      domain: "**openshift.web.dev.pega.io**"
```

To log in to Pega Platform with this host name, assign the host name with the same IP address that the deployment load balancer assigned to the web tier. This final step ensures that you can log in to Pega Platform with your host name, on which you can independently manage security protocols that match your networking infrastructure standards.

You can view the networking endpoints associated with your OpenShift deployment by using the OpenShift dashboard. From the Navigation menu, go to the **Networking > Routes** and in the `pega-web` page, in the **Router canonical hostname** area, you can see the pega-web tier ingress host name that Openshift creates based on the pega-web tier ingress host domain name that you set in the pega.yaml file.

![Initial view of the pega-web network](media/openshift-pega-web-route.png)

To manually associate the host name of the pega-web tier ingress with the tier endpoint, use the DNS lookup management system of your choice. As an example, if your organization has a GCP **Cloud DNS** that is configured to manage your DNS lookups, create a record set that specifies the pega-web tier the host name and add the IP address of the pega-web tier.

For GCP **Cloud DNS** documentation details, see [Quickstart](https://cloud.google.com/dns/docs/quickstart).

For details from Openshift to use this generated canonical hostname, click the link, "Do you need to set up custom DNS?"

### Logging in by using the domain name of the web tier

With the ingress host name name associated with this IP address in your DNS service, you can log in to Pega Platform with a web browser using the URL: `http://\<pega-web tier ingress host name>/prweb`.

![](media/25b18c61607e4e979a13f3cfc1b64f5c.png)
