Deploying Pega Platform on PKS
===============================

Use this set of procedures to deploy Pega Platform™ on a Pivotal Container
Service (PKS) cluster using a PostgreSQL database you configure in Google Cloud
Platform (GCP). They are written for any level of user, from a System
Administrator to a development engineer who is interested in learning how to
install and deploy Pega Platform onto PKS.

Pega helps enterprises and agencies quickly build business apps that deliver the
outcomes and end-to-end customer experiences that you need. Use the procedures in
this guide, to install and deploy Pega software onto PKS without much experience in either PKS configurations or Pega Platform deployments.

Create a deployment of Pega Platform on which you can implement a scalable Pega application in a PKS cluster. You can use this deployment for a Pega Platform development environment. By completing these procedure, you will have a Pega Platform cluster in PKS with a PostgreSQL database instance and two clustered virtual machines (VMs).

Deployment process overview
------------------------

Use Kubernetes tools and the customized orchestration tools and Docker images to orchestrate a deployment in a PKS cluster that you create for the deployment:

1. Prepare your local system:

    - To prepare a local Linux system, install required applications and configuration files - [Preparing your local Linux system – 45 minutes](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md).

    - To prepare a local Windows system, install required applications and configuration files -
    [Preparing your local Windows 10 system – 45 minutes](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md).

2. Request a PKS cluster form Pivotal and create an SQL instance in an account such as Google Cloud Platform (GPC) - [Prepare your PKS resources – 45 minutes](#prepare-your-resources--45-minutes).

3. Customize a configuration file with your PKS details and use kubectl and Helm to install and then deploy Pega Platform onto your PKS cluster - [Deploying Pega Platform using Helm charts – 90 minutes](#installing-and-deploying-pega-platform-using-helm-charts--90-minutes).

4. Configure your network connections in the DNS management zone of your choice so you can log in to Pega Platform - [Logging into Pega Platform – 10 minutes](#logging-into-pega-platform--10-minutes).

To understand how Pega maps Kubernetes objects with Pega applications and services, see [How Pega Platform and applications are deployed on Kubernetes](https://community.pega.com/knowledgebase/articles/cloud-choice/how-pega-platform-and-applications-are-deployed-kubernetes).

Assumptions and prerequisites
-----------------------------

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges or a Linux command prompt with root privileges.

- You use open source packaging tools on Windows or Linux to install applications onto your local system.

The following account, resources, and application versions are required for use in this document:

- A GCP account with a payment method set up to pay for the GCP resources you create. You should also have sufficient GCP account permissions and knowledge to:

  - Create a PostgreSQL Database.

  - Select an appropriate location in which to deploy your database resource; the document assumes your location is US East.

  You are responsible for any financial costs incurred for your GCP resources.

- Pega Platform 8.3.1 or later.

- Pega Docker images – your deployment requires the use of a custom Docker image to install Pega Platform into a database that is used by your PKS cluster. After you build your image, you must make it available in a private Docker registry. In order to construct your own image from the base image that Pega provides, you must have:

  - A DockerHub account to which you will push your final image to a private DockerHub repository. The image you build with Pega-provided components cannot be shared in a public DockerHub repository.

  - The Docker application downloaded to your local system, either Linux- or Windows-based. Log into your DockerHub account from the Docker application on your local system.

- Helm 3.0 or later. Helm is only required to use the Helm charts and not to use the Kubernetes YAML examples directly. For detailed usage, refer to the [Helm documentation portal](https://helm.sh/docs/).

- Kubectl –the Kubernetes command line tool you must use to connect and to manage your Kubernetes resources.

Prepare your resources – 45 minutes
-----------------------------------

This section covers the details necessary to obtain your PKS credentials and
configure the required PostgreSQL database in a GCP account. Pega supports
creating a PostgreSQL database in any environment if the IP address of
the database is available to your PKS cluster.

### Requesting a PKS cluster

In many cases a Pega client can request a demo cluster from Pivotal running PKS
1.5 or later. If your organization is already familiar with the Pivotal
Container Service, you can create your own local PKS environment.

At a minimum, your cluster should be provisioned with at least two worker nodes
that have 32GB of RAM in order to support the typical processing loads in a Pega
Platform deployment. Pivotal supports SSL authentication, which you can request
if your organization requires it.

In order to login to your demo cluster, you must have the following information:

- The target IP address of your PKS API

- The login credentials: username and password

- Whether any SSL information is required to authenticate your login.

During deployment the required Kubernetes configuration file is copied into the cluster.

### Creating a database resource

PKS deployments require you to install Pega Platform software in an SQL database. After you create an SQL instance that is available to your PKS cluster, you must create a PostgreSQL database in which to install Pega Platform. When you are finished, you will need the database name and the SQL instance IP address which you create in this procedure in order to add this information to your pega.yaml Helm chart.

#### Creating an SQL instance

To begin, create an SQL Instance that is available to your PKS cluster. In this example, we create an SQL instance in GCP; however, you can create or use an database resource that is available to the PKS cluster.

1. Use a web browser to log onto <https://cloud.google.com/> and navigate to your **Console** in the upper right corner.

2. In your **Google Cloud Platform** console, use the **Navigation Menu** to go to **Storage** grouping and select **SQL**.

3. In the SQL page, click **+Create Instance**.

4. In the **Choose your database engine** screen, click **Choose PostgreSQL**.

5. In the **Create PostgreSQL instance** page, add details to the following
    required fields for this database server:

    a. In **Instance ID**, enter a database server ID.

    b. In **Default user password**, enter a “postgres” user password.

    c. Select an appropriate **Region** and **Zone** for your database server. Contact your Pivotal representative so you can select the same zone or region that Pivotal used to create your PKS demo cluster.

    d. In **Database version**, select **PostgreSQL 11**.

    e. In **Configuration options \> Connectivity**, select **Public IP**, click **+ Add Network**, enter a **Name** and **Network** of one or more IP address to whitelist for this PostgreSQL database, and click **Done**.

    For clusters that are provisioned by Pivotal: you can launch the Kubernetes dashboard to view the external IP address of the nodes in your cluster; to add that IP network to the database whitelist, enter the first three sets of number,and use 0/24 for the final set in this IP range. For example: 123.123.123.0/24.

6. In **Configuration options \> Machine type and storage**:

    a. For **Machine type**, select 4 vCPU **Cores** and 15 GB **Memory**.

    b. For **Network throughput**, select **SSD (Recommended)**.

    c. For **Storage capacity**, enter **20 GB** and select **Enable automatic storage increases**.

7. Configure the remaining setting using the default values:

    a. For **Auto backups and high availability**, select backups can be automated 1AM – 5AM in a single zone.

    b. For **Flags**, no flags are required.

    c. **For Maintenance**, any preference is supported.

    d. For **Labels**, no labels are required.
    
    Labels can help clarify billing details for your PKS resources.

8. Click **Create**.

    A deployment progress page displays the status of your deployment until it is complete, which takes up to 5 minutes. When complete, the GCP UI displays all of the SQL resources in your account, which includes your newly created SQL instance:

![cid:image007.png\@01D5A3B1.62812F70](media/9aa072ea703232c2f6651fe95854e8dc.62812f70)

#### Creating a database in your SQL instance

Create a PostgreSQL database in your new SQL instance for the Pega Platform installation. Use the database editing tool of your choice to log into your SQL instance and create this new PostgreSQL database. The following example was completed using pgAdmin4.

1. Use a database editor tool, such as pgadmin4, to log into your SQL instance.

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

- Specify that this is an PKS deployment.

- Access your DockerHub account in order to access the required Docker images.

- Access your GCP SQL database you created in [Create an SQL database resource](#_Create_an_SQL).

- Install the version of Pega Platform that you built into your Docker install image in [Build your Docker image from the pega-installer-ready](#_Build_your_Docker) into your SQL database.

- Access the Pega Docker images you specify for your deployments.

- Specify host names for your web and stream tiers.

1. To download the pega.yaml Helm chart to the \<local filepath\>/pks-demo, enter:

`$ helm inspect values pega/pega > pega.yaml`

2. Open the pega.yaml file in a text editor and update the following parameters in the chart based on your PKS requirements:

| Chart parameter name    | Purpose                                   | Your setting |
|-------------------------|-------------------------------------------|--------------|
| provider:               | Specify a PKS deployment.                 | provider:"pks"|
| actions.execute:        | Specify a “deploy” deployment type.       | execute: "deploy"   |
| Jdbc.url:               | Specify the database IP address and database name for your Pega Platform installation.        | <ul><li>url: "jdbc:postgresql://**localhost**:5432/**dbName**"</li><li>where **localhost** is the public IP address you configured for your database connectivity and **dbName** is the name you entered for your PostgreSQL database in [Creating a database resource](#creating-a-database-resource).</li></ul>|
| Jdbc.driverClass:       | Specify the driver class for a PostgreSQL database. | driverClass: "org.postgresql.Driver"                |
| Jdbc.dbType:            | Specify PostgreSQL database type.         | dbType: "postgres”   |
| Jdbc.driverUri:         | Specify the database driver Pega Platform uses during the deployment.| <ul><li>driverUri: "latest jar file available” </li><li>For PostgreSQL databases, use the URL of the latest PostgreSQL driver file that is publicly available at <https://jdbc.postgresql.org/download.html>.</li></ul>|
| Jdbc: username: password: | Set the security credentials for your database server to allow installation of Pega Platform into your database.   | <ul><li>username: "\<name of your database user\>" </li><li>password: "\<password for your database user\>"</li><li>-- For GCP PostgreSQL databases, the default user is “postgres”.</li></ul>|
| jdbc.rulesSchema: jdbc.dataSchema:  | Set the names of both your rules and the data schema to the values that Pega Platform uses for these two schemas.      | rulesSchema: "rules" dataSchema: "data" |
| docker.registry.url: username: password: | Map the hostname of a registry to an object that contains the “username” and “password” values for that registry. For more information, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). | <ul><li>url: “<https://index.docker.io/v1/>” </li><li>username: "\<DockerHub account username\>"</li><li> password: "\< DockerHub account password\>"</li></ul> |
| docker.pega.image:       | Refer to the latest Pega Platform deployment image on DockerHub.  | <ul><li>Image: "pegasystems/pega:latest" </li><li>For a list of default images that Pega provides: <https://hub.docker.com/r/pegasystems/pega-ready/tags></li></ul> |
| upgrade:    | Do not set for installations or deployments | upgrade: for non-upgrade, keep the default value. |
| tier.name: ”web” tier.service.domain:| Set a hostname for the pega-web service of the DNS zone. | <ul><li>domain: "\<the hostname for your web service tier\>" </li><li>Assign this hostname with an external IP address and log into Pega Platform with this hostname in the URL. Your web tier hostname must comply with your networking standards and be available as an external IP address.</li></ul>|
| tier.name: ”stream” tier.service.domain: | Set the hostname for the pega-stream service of the DNS zone.   | <ul><li>domain: "\<the hostname for your stream service tier\>" </li><li>Your stream tier hostname should comply with your networking standards. </li></ul>|
| installer.image:        | Specify the Docker image you built to install Pega Platform. | <ul><li>Image: "\<your installation Docker image :your tag\>" </li><li>You created this image in  [Preparing your local Linux system](docs/prepping-local-system-runbook-linux.md)</li></ul>|
| installer. adminPassword:                | Specify a password for your initial log in to Pega Platform.    | adminPassword: "\<initial password\>"  |

3. Save the file.

### Deploy Pega Platform using the command line

These steps walk you through:

- Connecting your local system to your PKS cluster.
- Enabling the use of a web browser-based Kubernetes dashboard you can use to monitor
your deployment.
- Performing the Helm commands required to complete your deployment of Pega Platform on to your PKS environment.

It's easy to confuse a Helm install with a Pega Platform install, but they are
separate processes. The Helm install command uses Helm to install your
deployment as directed in the Helm charts, one in the charts\\addons folder and
one in the charts\\pega folder. 

In this document, you specify that the Helm chart always “deploys” by using the setting, actions.execute: “deploy”. In the following tasks, you overwrite this function on your *initial* Helm install by specifying `--set global.actions.execute:install-deploy`, which invokes an installation of Pega Platform using your installation Docker image and then
automatically followed by a deploy. In subsequent Helm deployments, you should not use the override argument, `--set global.actions.execute=`, since Pega Platform is already installed in your database.

1. Do one of the following:

- Open Windows PowerShell running as Administrator on your local system and change the location to the top folder of your pks-demo folder that you created in [Preparing your local Windows 10 system](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md).

`$ cd <local filepath>\pks-demo`

- Open a Linux bash shell and change the location to the top folder of your pks-demo directory that you created in [Preparing your local Linux system](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md).

`$ cd /home/<local filepath>/pks-demo`

2. Use the PKS CLI to log into your account using the Pivotal-provided API and login credentials and skip SSL validation.

`$ pks login -a <API> -u <USERNAME> -p <PASSWORD> -k`

If you need to validate with SSL, replace the -k with --ca-cert \<PATH TO CERT\>.

3. To use the PKS CLI to view your all of your PKS cluster names and status of each, enter:

`$ pks clusters`

Your cluster name is displayed in the **Name** field.

4. To use the PKS CLI to download the cluster Kubeconfig access credential file, which is specific to your cluster, into your \<local filepath\>/.kube directory, enter:

```yaml
$ pks get-credentials <cluster-name>`
Fetching credentials for cluster pega-platform.
Context set for cluster pega-platform.
```

If you need to use a Bearer Token Access Credentials instead of this credential file, see the Pivotal document, [Accessing Dashboard](https://docs.pivotal.io/pks/1-3/access-dashboard.html).

5. To use the kubectl command to view the VM nodes, including cluster names and status, enter:

`$ kubectl get nodes`

6. Establish a required cluster role binding setting so you can launch the Kubernetes dashboard.

`$ kubectl create clusterrolebinding dashboard-admin -n kube-system
--clusterrole=cluster-admin --serviceaccount=kube-system:kubernetes-dashboard`

7. To start the proxy server for the Kubernetes dashboard, enter:

`$ kubectl proxy`

8. To access the Dashboard UI, open a web browser and navigate to the following:

<http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/>

9. In the **Kubernetes Dashboard** sign in window, choose the appropriate sign in method:

- To use a cluster Kubeconfig access credential file: select **Kubeconfig**, navigate to your \<local filepath\>/.kube directory and select the config file. Click **SIGN IN**.

- To use a cluster a Kubeconfig token: select **Token** and paste your
    Kubeconfig token into the **Enter token** area. Click **SIGN IN**.

    You can now view your deployment details visually using the Kubernetes dashboard. You will use this dashboard to review the status of your deployment as you continue. At this point, with no deployment, you only see the PKS resources. Note that the Kubernetes dashboard does not display your PKS cluster name or your resource name. This is expected behavior.

    In order to continue using the Kubernetes dashboard to see the progress of your deployment, keep this PowerShell or Linux shell prompt open and open a new one for the remaining steps.

10. Do one of the following:

- Open a new Windows PowerShell running as Administrator on your local system and change the location to the top folder of your pks-demo folder.

`$ cd <local filepath>\pks-demo`

- Open a new Linux bash shell and change the location to the top folder of your pks-demo directory.

`$ cd /home/<local filepath>/pks-demo`

11. Create namespaces for both your Pega deployment and the addons:

```yaml
$ kubectl create namespace mypega-pks-demo
namespace/mypega-pks-demo created
$ kubectl create namespace pegaaddons
namespace/pegaaddons created
```

12. To install the addons chart, which you updated in [Preparing your local system](#Prepare-your-local-system-–-45-minutes), enter:

```yaml
$ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
```

The `pegaddons` namespace contains the deployment’s load balancer and disables the metric server. A successful pegaaddons deployment returns details of deployment progress. For further verification of your deployment progress, you can refresh the Kubernetes dashboard and look in the `pegaaddons` Namespace view.

13. To deploy Pega Platform for the first time by specifying to install Pega Platform into the database you specified in the Helm chart, install the pega.yaml Helm chart:

```yaml
helm install mypega-pks-demo pega/pega --namespace mypega-pks-demo --values pega.yaml --set global.actions.execute=install-deploy
```

For subsequent Helm installs, use the command `helm install mypega-pks-demo pega/pega --namespace mypega-pks-demo` to deploy Pega Platform and avoid another Pega Platform installation.

A successful Pega deployment immediately returns details that show progress for your `mypega-pks-demo` deployment.

14. Refresh the Kubernetes dashboard you opened in step 8. If you closed the dashboard, open a new command prompt running as Administrator and relaunch the web browser as directed in Step 10.

15. In the dashboard, use the **Namespace** pulldown to change the view to `mypega-pks-demo` and click on the **Pods** view. Initially, you can some pods have a red status, which means they are initializing:

![](media/dashboard-mypega-pks-demo-install-initial.png)

Note: A deployment takes about 15 minutes for all of the resource configurations to complete; however a full Pega Platform installation into the database can take up to an hour.

To follow the progress of an installation, use the dashboard. For subsequent deployments, you will not need to do this. Initially, some of the resources are making requests to complete the configuration; therefore, you will see red warnings while the configuration is finishing. This is expected behavior.

16. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-db-install** job, and click the logs icon located on the right side of that row.

    After you open the logs view, you can click the icon for automatic refresh to see current updates to the install log.

17. To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the `mypega-pks-demo` namespace pods.

![](media/f7779bd94bdf3160ca1856cdafb32f2b.png)

A successful deployment will not show errors across the various workloads. The `mypega-pks-demo` Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment will have 100% complete **Workloads**.

![](media/0fb2d07a5a8113a9725b704e686fbfe6.png)

Logging into Pega Platform – 10 minutes
---------------------------------------

After you complete your deployment, it is a best practice to associate the hostname of the pega-web tier ingress with the IP address that the deployment load balancer assigned to the tier during deployment. The hostname of the pega-web tier ingress used in this demo is **pks.web.dev.pega.io**, is set in the pega.yaml file in the lines:

```yaml
tier:
  - name: "web"

    service:
      # Enter the domain name to access web nodes via a load balancer.
      #  e.g. web.mypega.example.com
      domain: "**pks.web.dev.pega.io**"
```

In order to sign into Pega Platform using this hostname, you must assign it with the same IP address that the deployment load balancer has assigned to the web tier. This final step ensures that you can log onto Pega Platform using your hostname, on which you can independently manage security protocols that match your networking infrastructure standards.

To manually associate the hostname of the pega-web tier ingress with the tier endpoint, use the DNS lookup management system of your choice. As an example, if your organization has a GCP **Cloud DNS** configured to manage your DNS lookups, you can create a new record set with the pega-web tier the hostname and add the IP address of the pega-web tier.

For GCP **Cloud DNS** documentation details, see [Quickstart](https://cloud.google.com/dns/docs/quickstart). If not using the GCP **Cloud DNS**, for configuration details, see the documentation for your DNS lookup management system.

With the ingress hostname name associated with this IP address in your DNS service, you can log into Pega Platform with a web browser using the URL: http://\<pega-web tier ingress hostname>/prweb.

![](media/25b18c61607e4e979a13f3cfc1b64f5c.png)
