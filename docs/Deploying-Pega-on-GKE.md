Deploying Pega Platform on a GKE cluster
===============================

Deploy Pega Platform™ on a Google Kubernetes Engine (GKE) cluster using a PostgreSQL database you configure in Google Cloud Platform (GCP). The procedures are written for any level of user, from a system administrator to a development engineer who is interested in learning how to install and deploy Pega Platform onto a GKE cluster.

Pega helps enterprises and agencies quickly build business apps that deliver the outcomes and end-to-end customer experiences that you need. Use the procedures in this guide, to install and deploy Pega software onto a GKE cluster without much experience in either GKE configurations or Pega Platform deployments.

Create a deployment of Pega Platform on which you can implement a scalable Pega application in a GKE cluster. You can use this deployment for a Pega Platform development environment. By completing these procedures, you deploy Pega Platform on a GKE cluster with a PostgreSQL database instance and two clustered virtual machines (VMs).

Deployment process overview
------------------------

Use Kubernetes tools and the customized orchestration tools and Docker images to orchestrate a deployment in a GKE cluster that you create for the deployment:

1. Prepare your local system using [Preparing your local Linux system – 45 minutes](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md) – install required applications and configuration files.

2.  Create a GKE cluster and a Postgres database in an SQL instance in your Google Cloud Platform (GPC) account - [Prepare your GKE resources – 45 minutes](#prepare-your-resources--45-minutes).

3. Customize a configuration file with your GKE details and use the command-line tools, kubectl and Helm, to install and then deploy Pega Platform onto your GKE cluster - [Deploying Pega Platform using Helm charts – 90 minutes](#installing-and-deploying-pega-platform-using-helm-charts--90-minutes).

4. Configure your network connections in the DNS management zone of your choice so you can log in to Pega Platform - [Logging into Pega Platform – 10 minutes](#logging-into-pega-platform--10-minutes).

To understand how Pega maps Kubernetes objects with Pega applications and services, see [How Pega Platform and applications are deployed on Kubernetes](https://community.pega.com/knowledgebase/articles/cloud-choice/how-pega-platform-and-applications-are-deployed-kubernetes).

Assumptions and prerequisites
-----------------------------

This guide assumes that you use open source packaging tools on a Linux distribution to install applications onto your local system.

The following account, resources, and application versions are required for use in this document:

- A GCP account with a payment method set up to pay for the GCP resources you create and appropriate GCP account permissions and knowledge to:

  - Access an existing GCP project for your GKE resources.

  - Create an SQL Instance.

  - Select an appropriate location to deploy your database resource and GKE cluster. The cluster and PostgreSQL database into which you install Pega Platform must be in the same location.

  You are responsible for any financial costs incurred for your GCP resources.

- Pega Platform 8.3.1 or later.

- Pega Docker images – your deployment requires the use of a custom Docker image to install Pega Platform into a database that is used by your GKE cluster. After you build your image, you must make it available in a private Docker registry. To build your own image from the base image that Pega provides, you must have:

  - A DockerHub account to which you will push your final image to a private DockerHub repository. The image you build with Pega-provided components cannot be shared in a public DockerHub repository.

  - The Docker application downloaded to your local system. Log into your DockerHub account from the Docker application on your local system.

- Helm 3.0 or later. Helm is only required to use the Helm charts and not to use the Kubernetes YAML examples directly. For more information, see the [Helm documentation portal](https://helm.sh/docs/).

- kubectl – the Kubernetes command-line tool that you use to connect to and manage your Kubernetes resources.

Creating a Google Cloud Platform project - 5 minutes
-----------------------------------------
To deploy Pega Platform to a GKE cluster, you must create a Google Cloud project in which you will create your Kubernetes cluster resources.

1. Using the web browser of your choice, log in to [GCP](#https://cloud.google.com/) with your GCP account credentials.

2. Click **Console** next to your profile name to open the GCP console page.

3. In the search tool, search for "Manage resources" and select **Manage resources IAM & admin** to display the **Google Cloud Platform console > Manage resources** page.

4. Click **+CREATE PROJECT**.

5. In the New Project window, enter a unique **Project name**, select a **Location**, if appropriate, and click **CREATE**.

With the new project created, you can proceed with completing the preparation of your local system.

Preparing your GKE resources – 45 minutes
--------------------------------------

Obtain your GKE credentials so you can create a GKE cluster and configure the required PostgreSQL database in a GCP account. You can create a PostgreSQL database in any environment if the IP address of the database is available to your GKE cluster.

### Creating a GKE cluster

To deploy Pega using a GKE cluster, create the cluster in an existing  project in your Google Cloud account. During deployment the required Kubernetes configuration file is copied into the cluster. Create a multi-zonal cluster with two VMs with sufficient memory and CPU resources to support a deployment of Pega Platform that can perform under high workloads.

You can create this cluster using gcloud or the Google Cloud Console. This demo provides steps using the Google Cloud Console using the web browser of your choice; for steps to create the cluster using gcloud in the Google Cloud SDK, see the **gcloud** tab on the page [Creating a multi-zonal cluster]( https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-multi-zonal-cluster).

To log in to your demo cluster, you must have the following information:

- The name of your GKE cluster

- The login credentials for your Google account: username and password

- Whether any SSL information is required to authenticate your login and if so, the appropriate authentication certificates.

To use the Google Cloud Console:

1. In a web browser, log in to <https://cloud.google.com/> and navigate to your
    **Console** in the upper right corner.

2. In your **Google Cloud Platform** console, use the **Navigation Menu** to go
    to the **Kubernetes Engine > Clusters** page.

3. On the **Kubernetes Clusters** page, click **+CREATE CLUSTER**.

4. Choose the **Standard cluster** template.

5. On the Standard cluster page, enter the following details:

    a. **Name** - enter the permanent name you use to refer to your GKE cluster.

    b. **Location type** - select **Zonal** or **Regional**.

    c. **Zone** - select an appropriate zone or region from the list.

    d. **Master version** - select an appropriate version. (The default version is most appropriate.)

    e. **Node pools: default-pool - Number of Nodes** -  enter "2".

    f. **Node pools: default-pool - Machine configuration > Machine family** - select the "General-purpose" tab.

    g. **Node pools: default-pool - Machine configuration > Series** - select "N1".

    h. **Node pools: default-pool - Machine configuration > Machine type** - select n1-highmem-4 (4 vCPU **Cores** and 26 GB **Memory**) for a minimum deployment; however using n1-highmem-8 (8 vCPU **Cores** and 52 GB **Memory**) is suitable for deployments that will process heavier workloads.

    The remaining fields can be left to their default values; however, if you have specific cluster requirements, update the template with your changes before proceeding.

6. Scroll to the bottom of the page and click **Create**.

### Creating a database resource

Pega Platform deployments on GKE require you to install Pega Platform software in an SQL database. After you create an SQL instance that is available to your GKE cluster, you must create a postSQL database in which to install Pega Platform. When you are finished, you will need the database name and the SQL instance IP address that you create in this procedure in order to add this information to your pega.yaml Helm chart.

#### Creating an SQL instance

Create an SQL instance that is available to your GKE cluster. In this example, the SQL instance is created in GCP; however, you can create or use a database resource that is available to the GKE cluster.

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

    c. Select an appropriate **Region** and **Zone** for your database server. Select the same zone or region that you used to create your GKE cluster.

    d. **Database version**, select **PostgreSQL 11**.

    e. **Configuration options \> Connectivity**, select **Public IP**, click **+ Add Network**, enter a **Name** and **Network** of one or more IP addresses to whitelist for this PostgreSQL database, and click **Done**.
    
    As a best practice, add the following IP addresses: your local system from where you install helm, the worker nodes of the cluster. One method for finding the IP address for worker nodes of the cluster is to view the nodes in your GKE cluster with kubectl command-line tool and then use the command options, `kubectl describe nodes | grep ExternalIP`.

6. In **Configuration options \> Machine type and storage**:

    a. **Machine type**, select n1-highmem-2 (2 vCPU **Cores** and 13 GB **Memory**).

    b. **Network throughput**, select **SSD (Recommended)**.

    c. **Storage capacity**, enter **20 GB** and select **Enable automatic storage increases**.

7. Configure the remaining setting using the default values:

    a. **Auto backups and high availability** - select an automated backup time.

    b. **Flags** - no flags are required.

    c. **For Maintenance** - any preference is supported.

    d.  **Labels**, no labels are required. Labels can help clarify billing details for your GKE resources.

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

Installing and deploying Pega Platform using Helm charts – 90 minutes
---------------------------------------------------------------------

To deploy Pega Platform by using Helm, customize the pega.yaml Helm chart that holds the specific settings for your deployment needs and then run a series of Helm commands to complete the deployment.

An installation with deployment will take about 90 minutes total, because a Pega Platform installation in your PostgreSQL database takes up to an hour.

### Updating the pega.yaml Helm chart values

To deploy Pega Platform, configure the parameters in the pega.yaml Helm chart to your deployment resource. Pega maintains a repository of Helm charts that are required to deploy Pega Platform by using Helm, including a generic version of this chart. To configure parameters this file, download it from the repository to your local system, edit it with a text editor, and then save it with the same filename. To simplify the instruction, you can download the file to the \gke-demo folder you have already created on your local system. 

Configure the parameters so the pega.yaml Helm chart matches your deployment resources in these areas:

- Specify that this is an GKE deployment.

- Credentials for your DockerHub account in order to access the required Docker images.

- Access your GCP SQL database.

- Install the version of Pega Platform that you built into your Docker installation image.

- Specify host names for your web and stream tiers.

1. To download the pega.yaml Helm chart to the \<local filepath\>/gke-demo, enter:

`$ helm inspect values pega/pega > pega.yaml`

2. Use a text editor to open the pega.yaml file and update the following parameters in the chart based on your GKE requirements:

| Chart parameter name    | Purpose                                   | Your setting |
|-------------------------|-------------------------------------------|--------------|
| provider:               | Specify a GKE deployment.                 | provider:"gke"|
| actions.execute:        | Specify a “deploy” deployment type.       | execute: "deploy"   |
| Jdbc.url:               | Specify the database IP address and database name for your Pega Platform installation.        | <ul><li>url: "jdbc:postgresql://**localhost**:5432/**dbName**"</li><li>where **localhost** is the public IP address you configured for your database connectivity and **dbName** is the name you entered for your PostgreSQL database in [Creating a database resource](#creating-a-database-resource).</li></ul>  |
| Jdbc.driverClass:       | Specify the driver class for a PostgreSQL database. | driverClass: "org.postgresql.Driver"                |
| Jdbc.dbType:            | Specify PostgreSQL database type.         | dbType: "postgres”   |
| Jdbc.driverUri:         | Specify the database driver Pega Platform uses during the deployment.| <ul><li>driverUri: "latest jar file available” </li><li>For PostgreSQL databases, use the URL of the latest PostgreSQL driver file that is publicly available at <https://jdbc.postgresql.org/download.html>.</li></ul>||
| Jdbc: username: password: | Set the security credentials for your database server to allow installation of Pega Platform into your database.   |<ul><li>username: "\<name of your database user\>" </li><li>password: "\<password for your database user\>"</li><li>-- For GCP PostgreSQL databases, the default user is “postgres”.</li></ul> |
| jdbc.rulesSchema: jdbc.dataSchema:  | Set the names of both your rules and the data schema to the values that Pega Platform uses for these two schemas.      | rulesSchema: "rules" dataSchema: "data" |
| docker.registry.url: username: password: | Map the host name of a registry to an object that contains the “username” and “password” values for that registry. For more information, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). | <ul><li>url: “<https://index.docker.io/v1/>” </li><li>username: "\<DockerHub account username\>"</li><li> password: "\< DockerHub account password\>"</li></ul> |
| docker.pega.image:       | Refer to the latest Pega Platform deployment image on DockerHub.  | <ul><li>Image: "pegasystems/pega:latest" </li><li>For a list of default images that Pega provides: <https://hub.docker.com/r/pegasystems/pega-ready/tags></li></ul> |
| upgrade:    | Do not set for installations or deployments. | upgrade: for non-upgrade, keep the default value. |
| tier.name: ”web” tier.service.domain:| Set a host name for the pega-web service of the DNS zone. | <ul><li>domain: "\<the host name for your web service tier\>" </li><li>Assign this host name with an external IP address and log into Pega Platform with this host name in the URL. Your web tier host name must comply with your networking standards and be available as an external IP address.</li></ul>|
| tier.name: ”stream” tier.service.domain: | Set the host name for the pega-stream service of the DNS zone.   | <ul><li>domain: "\<the host name for your stream service tier\>" </li><li>Your stream tier host name should comply with your networking standards. </li></ul>|
| installer.image:        | Specify the Docker image you built to install Pega Platform. | <ul><li>Image: "\<your installation Docker image :your tag\>" </li><li>You created this image in  [Preparing your local Linux system](docs/prepping-local-system-runbook-linux.md)</li></ul>|
| installer. adminPassword:                | Specify a password for your initial log in to Pega Platform.    | adminPassword: "\<initial password\>"  |

3. Save the file.

### Deploy Pega Platform using the command line

A Helm installation and a Pega Platform installation are separate processes. The Helm install command uses Helm to install your deployment as directed in the Helm charts, one in the **charts\\addons** folder and one in the **charts\\pega** folder.

In this document, you specify that the Helm chart always “deploys” by using the setting, actions.execute: “deploy”. In the following tasks, you overwrite this function on your *initial* Helm install by specifying `--set global.actions.execute:install-deploy`, which invokes an installation of Pega Platform using your installation Docker image and then
automatically followed by a deploy. In subsequent Helm deployments, you should not use the override argument, `--set global.actions.execute=`, since Pega Platform is already installed in your database.

1. Open a Linux bash shell and change the location to the top folder of your gke-demo directory that you created in [Preparing your local Linux system](docs/prepping-local-system-runbook-linux.md).

`$ cd /home/<local filepath>/gke-demo`

2. Use the gcloud command to ensure you are logged into your account.

```yaml
$ gcloud info
Google Cloud SDK [274.x.y]
...
Account: [your Google account]
Project: [your Google project]

Current Properties:
  [core]
    project: [your Google project]
    custom_ca_certs_file: [any files you used for SSL certification]
    account: [your Google account]
    disable_usage_reporting: [your response during authorization]
  [compute]
    zone: [your selected zone]
```

3. View the status of all of your GKE clusters and verify the name of the cluster for the Pega Platform deployment.

`$ gcloud container clusters list`

4. Download the cluster Kubeconfig access credential file, which is specific to your cluster, into your \<local filepath\>/.kube directory.

    If your gcloud configuration includes the zone you chose for your cluster, you can skip adding the `-z <zone-name>` option to the command.

```yaml
$ gcloud container clusters get-credentials <cluster-name> -z <zone-name>
Fetching cluster endpoint and auth data.
kubeconfig entry generated for <cluster-name>.
```

5. View the nodes in your GKE cluster, including cluster names and status.

```yaml
$ kubectl get nodes
NAME                                             STATUS   ROLES    AGE    VERSION
gke-demo-default-pool-abc   Ready    <none>   3d2h   v1.13.11-gke.14
gke-demo-default-pool-def   Ready    <none>   3d2h   v1.13.11-gke.14
```

6. Install the Kubernetes dashboard by applying an available version on Github.

```yaml
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
```

7. To ensure security while using this dashboard, generate a new secret for use as an access token for the Kubernetes dashboard.

```yaml
$ kubectl -n kube-system get secret
NAME                               TYPE                                  DATA   AGE
...
kubernetes-dashboard-certs         Opaque                                0      3m48s
kubernetes-dashboard-csrf          Opaque                                1      3m47s
kubernetes-dashboard-key-holder    Opaque                                2      3m47s
kubernetes-dashboard-token-XYZ   kubernetes.io/service-account-token   3      3m48s
...
```

8. Set one of the service-account secret tokens as the dashboard token.

```yaml
$ kubectl -n kube-system describe secrets kubernetes-dashboard-token-XYZ
Name:         kubernetes-dashboard-token-XYZ
Namespace:    kubernetes-dashboard
Labels:       <none>
Annotations:  kubernetes.io/service-account.name: kubernetes-dashboard
              kubernetes.io/service-account.uid: <unique ID>

Type:  kubernetes.io/service-account-token

Data
====
ca.crt:     1115 bytes
namespace:  20 bytes
token:      <Unique encryption details>

```

9. Establish a required cluster role binding setting so that you can launch the Kubernetes dashboard.

```yaml
$ kubectl create clusterrolebinding dashboard-admin -n kube-system --clusterrole=cluster-admin --serviceaccount=kube-system:kubernetes-dashboard
```

10. Start the proxy server for the Kubernetes dashboard.

`$ kubectl proxy`

11. To access the Dashboard UI, open a web browser and navigate to the following URL:

<http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/>

12. In the **Kubernetes Dashboard** sign in window, choose the appropriate authentication method:

- To use a cluster Kubeconfig access credential file: select **Kubeconfig**, navigate to your \<local filepath\>/.kube directory and select the config file. Click **SIGN IN**.

- To use a cluster a Kubeconfig token: select **Token** and paste your Kubeconfig token into the **Enter token** area. Click **SIGN IN**.

    You can now view your deployment details using the Kubernetes dashboard. Use this dashboard to review the status of your deployment. Without a deployment, only GKE resources display. The dashboard does not display your GKE cluster name or your resource name, which is expected behavior.

    To continue using the Kubernetes dashboard to see the progress of your deployment, keep this  Linux shell open.


13. Open a new Linux bash shell and change the location to the top folder of your gke-demo directory.

`$ cd /home/<local filepath>/gke-demo`

14. Create namespaces in preparation for the pega.yaml and addons.yaml deployments.

```yaml
$ kubectl create namespace mypega-gke-demo
namespace/mypega-gke-demo created
$ kubectl create namespace pegaaddons
namespace/pegaaddons created
```

15. Install the addons chart, which you updated in [Preparing your local system](#Prepare-your-local-system-–-45-minutes).

```yaml
$ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
NAME: addons
LAST DEPLOYED: Fri Jan  3 18:58:28 2020
NAMESPACE: pegaaddons
STATUS: deployed
REVISION: 1
```

The `pegaddons` namespace contains the deployment’s load balancer and the metric server configurations that you configured in the addons.yaml Helm chart. A successful pegaaddons deployment returns details of deployment progress. For further verification of your deployment progress, you can refresh the Kubernetes dashboard and look in the `pegaaddons` **Namespace** view.

16. Deploy Pega Platform for the first time by specifying to install Pega Platform into the database specified in the Helm chart when you install the pega.yaml Helm chart.

```yaml
helm install mypega-gke-demo pega/pega --namespace mypega-gke-demo --values pega.yaml --set global.actions.execute=install-deploy
NAME: mypega-gke-demo
LAST DEPLOYED: Fri Jan  3 19:00:19 2020
NAMESPACE: mypega-gke-demo
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

For subsequent Helm installs, use the command `helm install mypega-gke-demo pega/pega --namespace mypega-gke-demo` to deploy Pega Platform and avoid another Pega Platform installation.

A successful Pega deployment immediately returns details that show progress for your `mypega-gke-demo` deployment.

17. Refresh the Kubernetes dashboard that you opened in Step 11. If you closed the dashboard, start the proxy server for the Kubernetes dashboard as directed in Step 10, and relaunch the web browser as directed in Step 11.

18. In the dashboard, in **Namespace** select the `mypega-pks-demo` view and then click on the **Pods** view. Initially, some pods will have a red status, which means they are initializing:

![Initial view of pods during deploying](media/dashboard-mypega-pks-demo-install-initial.png)

    Note: A deployment takes about 15 minutes for all resource configurations to initialize; however a full Pega Platform installation into the database can take up to an hour.

    To follow the progress of an installation, use the dashboard. For subsequent deployments, you do not need to do this. Initially, while the resources make requests to complete the configuration, you will see red warnings while the configuration is finishing, which is expected behavior.

19. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-db-install** job, and click the logs icon on the right side of that row.

    After you open the logs view, you can click the icon for automatic refresh to see current updates to the install log.

20.  To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the `mypega-gke-demo` namespace pods.

A successful deployment does not show errors across the various workloads. The `mypega-pks-demo` Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment has 100% complete **Workloads**.

Logging in to Pega Platform – 10 minutes
---------------------------------------

After you complete your deployment, as a best practice, associate the host name of the pega-web tier ingress with the IP address that the deployment load balancer assigned to the tier during deployment. The host name of the pega-web tier ingress used in this demo, **gke.web.dev.pega.io**, is set in the pega.yaml file in the following lines:

```yaml
tier:
  - name: "web"

    service:
      # Enter the domain name to access web nodes via a load balancer.
      #  e.g. web.mypega.example.com
      domain: "**gke.web.dev.pega.io**"
```

To log in to Pega Platform with this host name, assign the host name with the same IP address that the deployment load balancer assigned to the web tier. This final step ensures that you can log in to Pega Platform with your host name, on which you can independently manage security protocols that match your networking infrastructure standards.

You can view the networking endpoints associated with your GKE deployment by using the Google Cloud Platform console. From the Navigation menu, go to the **Kubernetes Engine > Clusters > Services & Ingresses** page to display the IP address of this tier and the pega-web tier ingress host name. Use the page filter to look at the pega-web resources in your cluster.

![Initial view of pods during deploying](media/pega-web-networking.png)

To manually associate the host name of the pega-web tier ingress with the tier endpoint, use the DNS lookup management system of your choice. As an example, if your organization has a GCP **Cloud DNS** that is configured to manage your DNS lookups, create a record set that specifies the pega-web tier the host name and add the IP address of the pega-web tier.

For GCP **Cloud DNS** documentation details, see [Quickstart](https://cloud.google.com/dns/docs/quickstart).

### Logging in by using the domain name of the web tier

With the ingress host name name associated with this IP address in your DNS service, you can log in to Pega Platform with a web browser using the URL: http://\<pega-web tier ingress host name>/prweb.

![](media/25b18c61607e4e979a13f3cfc1b64f5c.png)
