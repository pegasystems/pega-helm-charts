Deploying Pega Platform on AKS
===============================

About this document
-------------------

Use this set of procedures to deploy Pega Platform™ on an Azure Kubernetes
Service (AKS) cluster set up in your Microsoft Azure account. It is designed for
any level of user, from a System Administrator to a development engineer who's
interested in learning how to install and deploy Pega Platform onto AKS.

Pega helps enterprises and agencies quickly build business apps that deliver the
outcomes and end-to-end customer experiences they need. Using the procedures in
this guide, a user will be able to install and deploy Pega software onto AKS
without much experience in either AKS configurations or Pega Platform
deployments.

How the deployment works
------------------------

Pega provides customized orchestration tools and docker images required to
orchestrate a deployment in an AKS cluster you create for the deployment using
the following stepped tasks:

1. Prepare your local system using the appropriate instruction set:

- [Prepare a local Linux system – 45 minutes](#https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md) – install required applications and configuration files.

- [Prepare a local Windows 10 system – 45 minutes](#https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md) – install required applications and configuration files.

2. [Prepare your AKS resources – 45 minutes](#prepare-your-aks-resources--45-minutes) – create an AKS cluster, an SQL database, and a storage resource in your Azure account.

4. [Deploying Pega Platform using Helm charts – 90 minutes](#installing-and-deploying-pega-platform-using-helm-charts--90-minutes) – customize a configuration file with your AKS details and use kubectl and helm to install and then deploy Pega Platform onto your AKS cluster.

5. [Logging into Pega Platform – 10 minutes](#logging-into-pega-platform--10-minutes) – configure your network connections in your DNS zone so you can log onto Pega Platform.

You can review the Pega architecture overview to understand how Pega maps
Kubernetes objects with Pega applications and services: [How Pega Platform and
applications are deployed on Kubernetes](https://community.pega.com/knowledgebase/articles/cloud-choice/how-pega-platform-and-applications-are-deployed-kubernetes).

### Document goals

By following the procedures in this document, you will create a deployment of
Pega Platform on which you can implement a scalable Pega application in an AKS
cluster running in your Microsoft Azure account. This deployment can be used for
a Pega Platform development environment. The procedures cover these required
sections:

- Instructions for creating an AKS cluster that has suitable resources to
    support a scalable dev/test environment for your Pega application.

- Instructions for creating a Microsoft Azure SQL Database that hosts Pega
    Platform data and rules in your Azure account.

- Instructions for creating a Pega Platform installation Docker image you use
    to install and deploy Pega Platform onto your AKS cluster.

- Guidance for logging onto a Kubernetes dashboard and exploring areas in the
    dashboard to help you see the status of your deployment and troubleshoot
    deployment and installation problems.

    After you follow this document to its completion, you'll have a Pega
    Platform cluster with a single Microsoft Azure SQL Database instance and two
    clustered virtual machines (VMs).

### Assumptions and Prerequisites

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges or a Linux command prompt with root privileges.

- You use opensource packaging tools on Windows or Linux to install applications onto your local system.

The following account and application versions are required for use in this
document:

- A Microsoft Azure account with a payment method set up to pay for the Azure resources you create using this document. You should also have sufficient Microsoft Azure account permissions and knowledge to:

  - Create an AKS service, a SQL Database, and a storage resource.

  - Select an appropriate Location in which to deploy Microsoft Azure resources;
    the document assumes your Location is US East.

  - You are responsible for any financial costs incurred for your AKS resources.

- Pega Platform 8.3.1 or later

- Pega Docker images – your deployment requires the use of a custom Docker image to install Pega Platform into a database that is used by your AKS cluster. After you build your image, you must make it available in a private Docker registry. In order to construct your own image from the base image that Pega provides, you must have:

  - A DockerHub account to which you will push your final image to a private DockerHub repository. The image you build with Pega-provided components cannot be shared in a public DockerHub repository.

  - The Docker application downloaded to your local system, either Linux- or Windows-based. Log into your DockerHub account from docker on your local system.

- A Pega Platform distribution kit downloaded and extracted on your local system.

- Helm – install Helm 3.0 or later. Helm is only required to use the Helm charts and not to use the Kubernetes yaml examples directly. For detailed usage, refer to the Helm documentation portal.

- Kubectl –the kubernetes command line tool you must use to connect and to manage your Kubernetes resources.

Prepare your local system – 45 minutes
--------------------------------------

This document requires that you prepare a local Linux or Windows 10 system on
which you can run commands with root or Administrator privileges. To prepare
your system use the guide appropriate for your type of system:

- [Preparing your local Linux system – 45 minutes](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md)

- [Preparing your local Windows 10 system – 45 minutes](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md)

After you have your local system prepared, you are ready to prepare your PKS
resources.

### Create an AKS cluster

Create an AKS cluster in your Azure account for your deployment. The deployment
will place configuration files into your cluster during the deployment steps in
the section [Deploy Pega Platform using the command
line](#deploy-pega-platform-using-the-command-line). For now, you will create a
simple cluster with two VMs with sufficient memory and CPU resources to support
a deployment of Pega Platform that can perform under high workloads.

1. In a browser, login to Microsoft Azure Portal (<https://portal.azure.com/>)
    with your credentials.

2. Search for **kubernetes services** and select it in the dropdown list.

![](media/aeba40691a7035ab1e87c8426764c488.png)

3. Click **+Add**.

1. In the **Create an AKS Kubernetes service page**, in the **Basics** tab
    page, add details to the following required fields for this Kubernetes
    service:

1. In **Project Details**, select your Azure **Subscription** (the default is
    **Microsoft Azure**).

1. In **Resource Group** area, select the appropriate group or click **Create
    new** and provide a name.

1. In **Cluster details**, enter a **Kubernetes cluster name**.

1. Select a **Region** that is appropriate for your deployment.

1. Select the default **Kubernetes version**. Pega requires that you select
    version 1.13.10 or later.

1. For **DNS name prefix**, entire a prefix your organization requires or, if
    not required, use the default, which Azure creates based on the **Kubernetes
    cluster name** you provided.

1. In Primary node pool specify the size of your VMs for this service.

    Pega recommends standard deployments use VMs of type **D8Sv3** (8 vcpus, 32 GiB
memory). Change the filters if you cannot find this type using the search
function.

13. Specify a **Node count**. Note that the cost varies with different VM
    configurations.

     Pega recommends standard deployments use at least two nodes.

1. Click **Next: Scale**.

1. In the **Scale** tab page, keep the default, so **Virtual nodes** and VM
    Scale sets are both set to **Disabled**.

1. Click **Next: Authentication**.

1. In the **Authentication** tab page, add details to the following required
    fields for this Kubernetes service:

1. In **Cluster Infrastructure** area/**Service principal**, do one of the
    following:

- If you don’t have an existing service principal leave the field at the
    default to create a new service principal.

- If you have an existing service principal, click **Configure service
    principal**, and in the **Configure service principal** window, select **Use
    existing** and enter a **Service principal client ID** and **client
    secret**.

18. In **Kubernetes authentication and authorization**, leave Enable RBAC
    toggled to Yes.

1. Click **Next: Networking**.

1. In the **Networking** tab page, add details to the following required fields
    for this Kubernetes service:

1. Ensure that **HTTP application routing** is **Enabled**.

1. Ensure that **Networking configuration** is set to **Advanced**.

1. In **Configure virtual networks**, leave the network and address settings to
    defaults unless you have specific, advanced networking requirements.

1. Click **Next: Monitoring**.

1. In the **Monitoring** tab page, in **Azure Monitor**, accept the defaults to
    **Enable container monitoring** and to use the default analytics workspace
    to store monitoring data.

1. Click **Next: Tags**.

1. In the **Tags** tab page, add any tags with which you want to identify
    resources such as owner, user, organization name using the **Name** and
    **Value** tags.

   Tags can help clarify billing details for your AKS resources.

1. Click **Next: Review + create**.

   Azure runs a validation and when validated, your service is ready to create.

1. Check your configurations on the **Create and Review** tab.

   Azure validates your configuration for your parameters and marks **Validation
passed** on successful validation.

1. Click **Create**.

    You’ll see a deployment progress page displayed until your deployment is
complete, which takes about 15 minutes. When complete, the Azure UI displays all
of the resources created for your AKS cluster in **Deployment details**:

![](media/a7339875194a80ba4c4c7d45299f0c99.png)

### Create an SQL database resource

Create an SQL database in which to install Pega Platform. AKS deployments
require you to install in an SQL database. When you are finished, you will need
to obtain the database name and the server which you are creating in this
procedure so you can customize your “pega” helm chart with this information.

1. In a browser, login to Microsoft Azure Portal (https://portal.azure.com/)
    with your credentials.

2. Search for **SQL Databases** and select it in the dropdown list.

3. Click **+Add**.

4. In the **Create SQL Database** page, in the **Basics** tab page, add details
    to the following required fields for this database server:

5. In **Project details**, select your Azure Subscription.

6. In **Resource Group**, select the same group you specified or created for
    your AKS cluster.

7. In **Database details**, enter a database server name.

8. In **Server**, click **Create new**.

9. In the **New server** window, enter a name, security details, select the
    same **Region** you used for your AKS cluster, and click **OK**.

10. For **Want to use SQL elastic pool?**, leave the default set to **No**.

11. In **Compute + storage server**, click configure database and ensure you
    choose Gen5 compute hardware and a minimum of **4** vCores and **32 GB**
    storage.

12. Click **Next: Networking \>**.

13. In the **Networking** tab page, add details to the following required fields
    for this database server:

14. In the **Network connectivity** area, for a **Connectivity method**, select
    **Public endpoint**.

15. In the **Firewall rules** area, for **Allow Azure services and resources to
    access this server**, select **Yes**.

16. In the **Firewall rules** area, for **Add current client IP address**,
    select **Yes**.

17. Click **Next: Additional settings \>**.

18. In the **Additional settings** tab page, use the following settings for this
    database server:

    - **Data source**: **None**

    - **Database collation**: use the default

    - **Enable Advanced data security**: Not now

1. Click **Next: Tags**.

1. In the **Tags** tab page, add any tags with which you want to identify
    resources such as owner, user, organization name using the **Name** and
    **Value** tags.

Tags can help clarify billing details for your AKS resources.

21. Click **Next: Review + create**.

1. Check your configurations on the **Create and Review** tab.

Azure validates your configuration for your parameters.

23. Click **Create**.

    You’ll see a deployment progress page displayed until your database server
deployment is complete, which takes about 5 minutes. When complete, the Azure UI
displays all of the resources created for your SQL database deployment details:

![](media/ec3afaa6f4e3e9224dec832be51c7dc5.png)

### Locate your SQL database details

After you finalize your SQL database configuration and you are ready to deploy,
you will need the following database details to complete your deployment
configuration by updating this database URL string in the “pega” helm chart:

`jdbc:sqlserver://**YOUR_DB_HOST**:1433;databaseName=**YOUR_DB_NAME**;selectMethod=cursor;sendStringParametersAsUnicode=false:`

To complete this configuration, you’ll use the SQL Server name for
**YOUR_DB_HOST** and the SQL database name for **YOUR_DB_NAME**. In order to
locate these details in your Azure portal:

1. Click the Portal menu in the upper left corner of you Azure portal and in
    the **Favorites** section, select **SQL databases**.

2. Under the **Name** column, click the database you just created to display
    the SQL database details page for your database.

    See the screenshot below to locate **YOUR_DB_NAME** and **YOUR_DB_HOST**. You
will use these names when you update your values Helm chart.

![](media/35a5b26c419b985b6a25b19a62de07a8.png)

Installing and deploying Pega Platform using Helm charts – 90 minutes
---------------------------------------------------------------------

In order to deploy Pega Platform using Helm, you must customize the “pega” helm
chart that holds the specific settings for your deployment needs and then run a series of helm commands to complete the deployment.

An installation followed by a deployment will take about 90 minutes total, since
it takes about an hour for Pega Platform to completely install in your SQL database.

### Update the Helm chart values

To deploy Pega Platform, you must finalize a number of parameters in the pega.yaml Helm chart. Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm, including a generic version of this chart. You will use the pega.yaml file to set customization parameters that are specific to your deployment. 

To customize this file, you must download it from the repository to your local system, edit it with a text editor, and then save it with the same filename. To simplify the instruction, you can download the file to the aks-demo folder you have already created on your local system. By customizing these parameters, you are configuring these required conditions:

- Specify that this is an AKS deployment.

- Access your DockerHub account in order to access the required Docker images.

- Access your Azure SQL database you created in [Create an SQL database
    resource](#create-an-sql-database-resource) and located in [Locate your SQL
    database details](#locate-your-sql-database-details).

- Install the version of Pega Platform that you built into your docker install
    image in [Build your Docker image from the
    pega-installer-ready](#_Build_your_Docker) into your Azure SQL database.

- Access the Pega Docker images you specify for your deployments.

- Specify host names for your web and stream tiers.

To finalize these details, follow these steps:

1. To download pega/pega Helm chat to the \<local filepath>\aks-demo, enter:

`$ helm inspect values pega/pega > pega.yaml`

2. Open the pega.yaml file from this folder in a text editor and update the following
    settings in the chart based on your AKS requirements:

| Chart parameter name    | Purpose                                   | Your setting |
|-------------------------|-------------------------------------------|--------------|
| provider:               | Specify a AKS deployment.                 | provider:"aks"|
| actions.execute:        | Specify a “deploy” deployment type.       | execute: "deploy"   |
| Jdbc.url:               | Specify the server and database name for your Pega Platform installation.   | url: “jdbc:sqlserver://**YOUR_DB_HOST_NAME**:1433; databaseName=**YOUR_DB_NAME**; selectMethod=cursor; sendStringParametersAsUnicode=false” You can locate **YOUR_DB_HOST_NAME** and **YOUR_DB_NAME** of your Azure SQL database, see [Locate your SQL database details](#locate-your-sql-database-details). |
| Jdbc.driverClass:       | Specify the driver class for this type of database. | driverClass: "com.microsoft.sqlserver.jdbc.SQLServerDriver"                                              |
| Jdbc.dbType:            | Specify the database type.                | dbType: " mssql”    |
| Jdbc.driverUri:         | Specify the database driver Pega Platform uses during the deployment. For AKS, we can obtain the URL of the required 7.4.1. driver file that is publicly available in the referenced Maven repository.                              | driverUri: "https://repo1.maven.org/maven2/com/microsoft/sqlserver/mssql-jdbc/7.4.1.jre11/mssql-jdbc-7.4.1.jre11.jar" |
| Jdbc: username: password: | Set the security credentials for your database server to allow installation of Pega Platform into your database.           | username: "\<name of your database user\>" password: "\<password for your database user\>"     |
| jdbc.rulesSchema: jdbc.dataSchema:       | Set the names of both your rules and the data schema to the values that Pega uses for these two schemas.         | rulesSchema: "rules" dataSchema: "data"      |
| docker.registry.url: username: password: | This object maps the hostname of a registry to an object containing the “username” and “password” for that registry. For details, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). | url: “<https://index.docker.io/v1/>” username: "\<DockerHub account username\>" password: "\< DockerHub account password\>"     |
| docker.pega.image: | Refer to the latest Page Platform deployment image on Dockerhub. | Image: "pegasystems/pega:latest" Pega provides these images: <https://hub.docker.com/r/pegasystems/pega-ready/tags>     |
| tier.name: ”web” tier.service.domain:    | Set the hostname for the pega-web service of the DNS zone.  | domain: "\<the hostname for your web service tier\>" You will assign this hostname with an external IP address and log into Pega Platform using this hostname in the URL. Your web tier hostname should comply with your networking standards and be available as an external IP address.|
| tier.name: ”stream” tier.service.domain: | Set the hostname for the pega-stream service of the DNS zone.  | domain: "\<the hostname for your stream service tier\>" Your stream tier hostname should comply with your networking standards |
| installer.image:   | Specify the Docker image you built to install Pega Platform.   | Image: "\<your installation Docker image :your tag\>" You created this image in [Constructing your Docker image from the pega-installer-ready](#_Constructing_your_Docker)   |
| installer. adminPassword: | Specify a password for your initial logon to Pega Platform. | adminPassword: "\<initial password\>" |

3. Save the file.

### Deploy Pega Platform using the command line

These steps walk you through connecting your local system to your AKS cluster;
enabling the use of a browser-based Kubernetes dashboard you can use to monitor
your deployment; and performing the helm commands required to complete your
deployment of Pega Platform on to your AKS environment.

It's easy to confuse a helm install with a Pega Platform install, but they are
separate processes. The helm install command uses helm to install your
deployment as directed in the helm charts, one in the charts\\addons folder and
one in the charts\\pega folder. In this document, you specify that the helm
chart always “deploys” by using the setting, actions.execute: “deploy”. For the
following steps, you overwrite this function on your initial helm install by
specifying “--set global.actions.execute:install-deploy”, which invokes an
installation of Pega Platform using your installation docker image and then
automatically followed by a deploy. In subsequent helm deployments, you should
not use the override argument, “--set global.actions.execute=”, since Pega
Platform is already installed in your database.

1. Do one of the following:

- Open a Windows PowerShell running as Administrator on your local system and change the location to the top folder of your aks-demo folder that you created in [Create a local folder to access all of the configuration file](#create-a-local-folder-to-access-all-of-the-configuration-files).

`$ cd <local filepath>\aks-demo`

- Open a Linux bash shell and change the location to the top folder of your aks-demo directory that you created in [Create a local folder to access all of the configuration file](#create-a-local-folder-to-access-all-of-the-configuration-files).

`$ cd /home/<local filepath>/aks-demo`

2. Use the Azure CLI to log into your account.

`$ az login`

A browser window opens and you are asked to log into your Azure account.

3. Log into the Azure account you will use to deploy Pega Platform.

    After you log into your Azure home page, ensure that you see confirmation
information in your command prompt. For example, in a Windows PowerShell you'll see confirmation of your credentials similar to this screen shot.

![](media/0a2d6fd5f1c8174de2e5f9a97b8801a6.png)

If you cannot log into your Azure home page or see that the Azure CLI recognizes your account correctly, contact your Microsoft account representative.

4. In the upper right corner of your browser, click **Portal**.

    You are brought to your Azure home page.

1. In the **Recent resources** area, click your AKS cluster to review these
    settings:

- Name: listed in the page header

- Resource group: listed in the main details section

6. Prepare your environment using the Azure CLI, in your Windows PowerShell:

```yaml
$ az aks get-credentials --resource-group <resource-group-name> --name <cluster-name>
Merged "runbook-demo" as current context in <local filepath>\<cluster-name>.kube\config
```

7.  Use the kubectl command to view the VM nodes created when you created the
    AKS cluster:

`$ kubectl get nodes`

![](media/3a59b1ecf9d827e0003d46880029cdd8.png)

8. Establish a required cluster role binding setting so you can launch the
    kubernetes dashboard.

`$ kubectl create clusterrolebinding dashboard-admin -n kube-system --clusterrole=cluster-admin --serviceaccount=kube-system:kubernetes-dashboard`

9. Launch the kubernetes dashboard to view your AKS resources before you deploy
    Pega Platform by replacing the names specific to your AKS cluster.

`$ az aks browse --resource-group <resource-group-name> --name <cluster-name>`

![](media/81a7ae961cabc463381869e3bae5c722.png)

You can now view your deployment details visually using the Kubernetes dashboard. You’ll use this dashboard to review the status of your deployment as you continue. At this point, with no deployment, you only see the AKS resources. Note that the Kubernetes dashboard does not display your AKS cluster name or your resource name. This is expected behavior.

In order to continue using the Kubernetes dashboard to see the progress of your deployment, keep this PowerShell or Linux command prompt open and open a new one for the remaining steps.

10. Do one of the following:

- Open a new Windows PowerShell running as Administrator on your local system and change the location to the top folder of your aks-demo folder that you created in [Create a local folder to access all of the configuration file](#create-a-local-folder-to-access-all-of-the-configuration-files).

`$ cd \<local filepath>\aks-demo`

- Open a new Linux bash shell and change the location to the top folder of your aks-demo directory that you created in [Create a local folder to access all of the configuration file](#create-a-local-folder-to-access-all-of-the-configuration-files).

`$ cd /home/<local filepath>/aks-demo`

11. Create namespaces for both your Pega deployment and the addons:

```yaml
$ kubectl create namespace mypega
namespace/mypega created
$ kubectl create namespace pegaaddons
namespace/pegaaddons created
```

12. To install the addons chart, which enables the deployment’s load balancer and disables the metric server and you havealready configured onyour local system, enter:

```yaml
$ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
```

A successful pegaaddon deployment returns details of deployment progress. For further verification of your deployment progress, you can refresh the Kubernetes dashboard and look in the pegaaddons Namespace view.

13. To deploy Pega Platform for the first time by specifying to install Pega Platform into the database you specified in the Helm chart, install the pega.yaml Helm chart:

```yaml
helm install mypega pega/pega --namespace mypega --values pega.yaml --set
 global.actions.execute=install-deploy
```

For subsequent Helm installs, use the command 'helm install mypega pega/pega --namespace mypega --values pega.yaml' to deploy Pega Platform and avoid another Pega Platform installation.

A successful Pega deployment immediately returns details that show progress for your deployment.

14. Refresh the Kubernetes dashboard you opened in step 9; if you closed it, open a new command prompt running as Administrator and relaunch the browser as directed in Step 9.

In the dashboard, use the **Namespace** pulldown to change the view to **pega**
and click on the **Pods** view.

![](media/055d24b4ac0c0dfcb9c68cec334ce42a.png)

Initially, some of the resources are making requests to complete the configuration; therefore, you will see red warnings while the configuration is finishing. This is expected behavior.

A deployment takes about 15 minutes for all of the resource configurations to complete; however a full Pega Platform installation into the database can take up to an hour. To follow the progress of an installation, use the dashboard; for subsequent deployments, you will not need to do this.

15. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-db-install** job, and click the logs icon located on the right side of that row.

    After you open the logs view, you can click the icon for automatic refresh to see current updates to the install log.

16. To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the **mypega** namespace pods.

![](media/f7779bd94bdf3160ca1856cdafb32f2b.png)

A successful deployment will not show errors across the various workloads. The **mypega** Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment will have 100% complete **Workloads**.

![](media/0fb2d07a5a8113a9725b704e686fbfe6.png)

Logging into Pega Platform – 10 minutes
---------------------------------------

After you have completed your deployment, you must associate the hostname of the
pega-web tier with the IP address that the deployment load balancer gave to the
tier. This final step ensures that you can log onto Pega Platform using your
hostname, on which you can independently manage security protocols that match
your networking infrastructure standards.

### Logging in using the IP address

To view the pega deployment components, enter:

`$ kubectl get services --namespace pega`

![](media/f329e9f92feed8cb5959d91db246aa84.png)

The pega-web tier external IP address and port number are displayed. Port 80 is
used for http traffic, which means you can’t use https encryption when access
the web-tier in a browser with

Instead, Pega recommends using the domain name of the web tier.

### Logging in using the domain name of the web tier

You must manually set the IP address of the web tier domain name in order to log
into Pega Platform domain name that is set during the deployment.

The example of a domain name of the web tier used in this demo is
**aks.web.dev.pega.io**, which you set in the values.yaml file here:

```yaml
tier:
   -name: "web"
    # Enter the domain name to access web nodes via a load balancer.
    # e.g. web.mypega.example.com

    service:
    # Enter the domain name to access web nodes via a load balancer.
    #  e.g. web.mypega.example.com
    domain: "**aks.web.dev.pega.io**"
```

When you set this to be "\<the hostname for your web service tier\>" as directed
in [Update the Helm chart values](#update-the-helm-chart-values), you will
manually associate "\<the hostname for your web service tier\>" with the IP
address of the web tier domain name.

In order to sign into Pega using "\<the hostname for your web service tier\>",
you must assign the domain name with the same IP address that the deployment
load balancer has assigned to the web tier.

1. From your command prompt, review the IP addresses that are in the pega service

`$ kubectl get services --namespace pega`

![](media/f329e9f92feed8cb5959d91db246aa84.png)

2. In a browser, login to Microsoft Azure Portal (https://portal.azure.com/)
    with your credentials.

3.  Search for **DNS zones** and select it in the dropdown list.

    You are brought to the DNS zones for your AKS cluster.

![](media/3c7f4a5c3c21c6577a4dbab5f5cfa79d.png)

4. In the Name column, click on the In the DNS zone for your deployment.

The DNS zone page displays.

5. To associate the IP address of the web tier (in step 1) with the domain name
    you configured during your deployment, you must add a **Record set** to your
    DNS zone for **"\<the hostname for your web service tier\>"**:
      
    a. Click **+Record** set

    b.  In the **Name** field, enter **"\<the hostname for your web service
    tier\>"**.

    c.  In the **Type** field, select **A**.

    d.  In **Alias record set** the configuration remains **No**.

    e.  Set **TTL** to 5 and **TTL Unit** to Minutes.

    f.  In the IP Address field, enter the IP from the pega-web tier.

    g.  Click **OK**.

    The new record set appears in the list of record sets for this Azure DNS zone.

![](media/ccb6329a621c6f11970e25531cfa1857.png)

With the domain name set to this IP address, you can log into Pega Platform with a browser using the URL: http://\<the hostname for your web service tier>/prweb

![](media/25b18c61607e4e979a13f3cfc1b64f5c.png)
