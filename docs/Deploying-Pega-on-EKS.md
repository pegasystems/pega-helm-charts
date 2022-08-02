# Deploying Pega Platform on an Amazon EKS cluster

Deploy Pega Platform™ on an Amazon Elastic Kubernetes Service (Amazon EKS) cluster using an Amazon Relational Database Service (Amazon RDS). The following procedures are written for any level of user, from a system administrator to a development engineer who is interested in learning how to install and deploy Pega Platform onto a EKS cluster.

Pega helps enterprises and agencies quickly build business apps that deliver the outcomes and end-to-end customer experiences that you need. Use the procedures in this guide, to install and deploy Pega software onto a EKS cluster without much experience in either EKS configurations or Pega Platform deployments.

Create a deployment of Pega Platform on which you can implement a scalable Pega application in a EKS cluster. You can use this deployment for a Pega Platform development environment. By completing these procedures, you deploy Pega Platform on a EKS cluster with a Amazon RDS database instance and two clustered virtual machines (VMs).

## Deployment process overview

Use Kubernetes tools and the customized orchestration tools and Docker images to orchestrate a deployment in a EKS cluster that you create for the deployment:

1. Prepare your local system:

   - To prepare a local Linux system, install required applications and configuration files - [Preparing your local Linux system – 45 minutes](prepping-local-system-runbook-linux.md).

   - To prepare a local Windows system, install required applications and configuration files -
   [Preparing your local Windows 10 system – 45 minutes](prepping-local-system-runbook-windows.md).

2. Create an Amazon EKS cluster and create an Amazon RDS instance in your AWS account - [Prepare your Amazon EKS resources – 45 minutes](#prepare-your-Amazon-EKS-resources--45-minutes).

3. Customize a configuration file with your Amazon EKS cluster details and use the command-line tools, AWS CLI, eksctl, kubectl and Helm, to install and then deploy Pega Platform onto your Amazon EKS cluster - [Deploying Pega Platform using Helm charts – 90 minutes](#installing-and-deploying-pega-platform-using-helm-charts--90-minutes).

4. Configure your network connections in the DNS management zone of your choice so you can log in to Pega Platform - [Logging in to Pega Platform – 10 minutes](#logging-in-to-pega-platform--10-minutes).

To understand how Pega maps Kubernetes objects with Pega applications and services, see [Understanding the Pega deployment architecture](https://community.pega.com/knowledgebase/articles/client-managed-cloud/cloud/understanding-pega-deployment-architecture).

## Assumptions and prerequisites

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges or a Linux command prompt with root privileges.

- You use open source packaging tools on Windows or Linux to install applications onto your local system.

The following account, resources, and application versions are required for use in this document:

- An Amazon AWS account with a payment method set up to pay for the Amazon cluster and RDS resources you create and appropriate AWS account permissions and knowledge to:

  - Create an Amazon RDS DB instance.

  - Select an appropriate location in which to deploy your database resource; the document assumes your location is US East.

  You are responsible for any financial costs incurred for your AWS resources.

- Pega Platform 8.3.1 or later.

- Pega Docker images – your deployment requires the use of several Docker images that you download and make available in a private Docker registry. For step-by-step details, see [Downloading and managing Pega Platform docker images (linux)](prepping-local-system-runbook-linux.md#downloading-and-managing-pega-platform-docker-images) or [Downloading and managing Pega Platform docker images (windows)](prepping-local-system-runbook-windows.md#downloading-and-managing-pega-platform-docker-images).

- Helm 3.0 or later. Helm is only required to use the Helm charts and not to use the Kubernetes YAML examples directly. For more information, see the [Helm documentation portal](https://helm.sh/docs/).

- kubectl – the Kubernetes command-line tool that you use to connect to and manage your Kubernetes resources.

- AWS IAM Authenticator for Kubernetes - the AWS command-line tool that you use to configure the required AWS CLI credentials for deploying your Amazon EKS cluster: access key, secret access key, AWS Region, and output format.

- eksctl - the Amazon EKS command-line tool that you use for creating and managing clusters on Amazon EKS. If you have to update the version of kubernetes on an existing Amazon EKS cluster, enter the following command, replacing `<cluster-name>` with your existing cluster name: `eksctl update cluster --name <cluster-name> --approve`.

## Prepare your Amazon EKS resources – 45 minutes

This section covers the details necessary to obtain your AWS credentials and configure the required Amazon RDS database in an AWS account.

### Creating your IAM user access keys

In order to create an EKS cluster, Pega recommends using IAM user access keys for deployment authentications. If your organization does not support using IAM credentials, refer to your organization's guidance for how to manage authentication between your deployment in your EKS cluster and your organization.

Use the following steps, which are sourced from the AWS documentation, **Access Key and Secret Access Key** on the page, [Quickly Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#cli-quick-configuration).

To create access keys for an IAM user:

1. Sign in to the AWS Management Console and open the IAM console at https://console.aws.amazon.com/iam/.

2. In the navigation pane, choose Users.

3. Choose the name of the user whose access keys you want to create, and then choose the Security credentials tab.

4. In the Access keys section, choose Create access key.

5. To view the new access key pair, choose Show. You will not have access to the secret access key again after this dialog box closes. Your credentials will look something like this:

   - Access key ID: AKIAIOSFODNN7EXAMPLE
   - Secret access key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

   You need these two details when you configure the load balancer for your deployment in the Helm addons configuration section, [Updating the addons.yaml Helm chart values](#Updating-the-addonsyaml-Helm-chart-values).

6. To download the key pair, choose Download .csv file. Store the keys in a secure location. You will not have access to the secret access key again after this dialog box closes.

   Keep the keys confidential in order to protect your AWS account and never email them. Do not share them outside your organization, even if an inquiry appears to come from AWS or Amazon.com. No one who legitimately represents Amazon will ever ask you for your secret key.

7. After you download the .csv file, choose Close.

When you create an access key, the key pair is active by default, and you can use the pair right away.

### Configuring your AWS CLI credentials

To save your IAM user access keys and other preferences for your EKS deployment to a configuration file on your local system, use the `aws configure` command. this command prompts you for four pieces of information you must specify in order to deploy an EKS cluster in your AWS account (access key, secret access key, AWS Region, and output format). For complete details about what this information includes, see the overview article, [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html). To see details that may be useful to customize your stored credentials to meet your organization's business needs, see [Configuration and Credential File Settings](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).

To setup your local system and save your AWS credentials and profile to the $USER/.aws file, enter:

`$ aws configure`

You are prompted for your AWS access credentials and details. Enter your own values. For guidance on completing each value, see [Quickly Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#cli-quick-configuration).

```yaml
AWS Access Key ID [None]: your-key-ID
AWS Secret Access Key [None]: your-secrete-access-key
Default region name [None]: your-region-preference
Default output format [None]:  specify your preference for a result format.
```

With your credentials saved locally, you must push your Pega-provided Docker images to your Docker registry. For details on where the AWS CLI stores your credentials locally, see [Where Are Configuration Settings Stored?](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-where)

### Managing SSL certificates in AWS certificate manager

Pega supports the use of an HTTPS load balancer through a Kubernetes ingress, which requires you to configure the load balancer to present authentication certificates to the client. In EKS clusters, Pega requires that you use an AWS Load Balancer Controller (formerly named AWS ALB Ingress Controller). For an overview, see [Application load balancing on Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html).

To configure this ingress controller, Pega allows your deployment to use the AWS Load Balancer Controller without a certificate for testing purposes; however, for the best security practice, it is a recommend practice to specify an SSL certificate that you create or import into AWS Credential Manager. After you have one your Amazon Resource Name (ARN) credential using this certificate or multiple certificates uploaded to AWS Credential Manager, use one of the following choices for the `annotation` parameter of the web tier ingress configuration:

- Leave it blank so that the deployment automatically associates the existing certificate with your ingress controller.
- Specify an ARN certificate out of multiple certificates or want to provide secondary certificates that the deployment associates with your ingress controller.

To import or create a new SSL certificate to support HTTPS, see [Importing Certificates into AWS Certificate Manager](https://docs.aws.amazon.com/acm/latest/userguide/import-certificate.html). After you have your ARN certificate, include a reference to the required ARN certificate using an appropriate "web tier" `ingress` annotation as in the example:

```yaml
ingress:
  tls:
    # Enable TLS encryption
    enabled: true
    # secretName:
    # useManagedCertificate: false
    ssl_annotation: 
      alb.ingress.kubernetes.io/certificate-arn: <certificate-arn>
```

Where `alb.ingress.kubernetes.io/certificate-arn` is the required annotation name and `<certificate-arn>` takes the form, `arn:aws:acm:<region>:<AWS account>:certificate/xxxxxxx`, which you copy from the AWS console view of the load balancer configuration. You add these parameters when you complete the configuration of your Helm page chart. For details, see, [Updating the addons.yaml Helm chart values](#Updating-the-addonsyaml-Helm-chart-values).

### Creating an Amazon EKS cluster

You can create your EKS cluster using the `eksctl` command line utility. This example shows how to define your configuration in a yaml file that you pass to the `eksctl` command. For more details and options available to advanced EKS users, see [Create your Amazon EKS cluster and worker nodes](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html#eksctl-create-cluster).

At a minimum, for if you are creating a simple Pega demo, your cluster must be provisioned with at least two worker nodes that have 32GB of RAM in order to support the typical processing loads in a Pega Platform deployment; for this option, you can use a minimum of two m5.xlarge nodes for your deployment worker nodes. Pega has not tested this method using Windows worker nodes. Larger deployments may required additional processing capacity.

To create an Amazon EKS cluster with Linux worker nodes:

1. Save the following text to a file named similar to `my-EKS-cluster-spec.yaml` in your EKS-demo folder:

```yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
name: pega-85-demo
region: us-east-1
version: '1.15'

nodeGroups:
- name: linux-ng
instanceType: m5.xlarge
minSize: 2
```
You should use a name and region here to reflect the version of Pega Platform you want to deploy and specify the region in which your cluster will run.

2. To create your Amazon EKS cluster and Windows and Linux worker nodes, from your /EKS-demo folder, enter.

`eksctl create cluster -f ./my-EKS-cluster-spec.yaml`

 It takes 10 to 15 minutes for the cluster provisioning to complete. During deployment this command copies the required Kubernetes configuration file to the cluster and into your default ~/.kube directory.
 
 3. After provisioning is complete, verify that the worker nodes joined the cluster and are in Ready state, by entering:

`kubectl get nodes`

With your cluster created and running as expected, you must create a database resource for your Pega Platform installation. If you have to delete your cluster for some reason before you have namespaces deployed in it, use the command, `eksctl delete cluster --name <cluster-name>'.

4. To deploy the Kubernetes dashboard and see your EKS cluster, enter:

kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml

5. Create a Service Account and Cluster Role Binding, each named `eks-admin` to securely connect to the dashboard with administrator-level permissions, since by default, the Kubernetes dashboard limits permissions.

For client convenience, Pega suggests saving the following text, to define the service account and cluster role binding for your deployment, to a file named similar to `eks-admin-service-account.yaml` in your EKS-demo folder.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: eks-admin
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: eks-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
-kind: ServiceAccount
 name: eks-admin
 namespace: kube-system
```

6. Apply the service account and cluster role binding to your cluster.

```yaml
$ kubectl apply -f eks-admin-service-account.yaml
service account/eks-admin created
```

7. To generate an authentication token for the eks-admin service account required to connect to the dashboard of your administrative service account, enter:

`kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '{print $1}')`

8. Copy the <authentication_token> value to connect to the Kubernetes dashboard.

9. To start the proxy server for the Kubernetes dashboard, enter:

    `$ kubectl proxy`

10. To access the Dashboard UI, open a web browser and navigate to the following URL:

    `http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/`

11. In the **Kubernetes Dashboard** sign in window, choose to use a cluster Kubeconfig token: select **Token** and paste the <authentication_token> value that you just generated into the **Enter token** area.

12. Click **SIGN IN**.

    You can now view your EKS cluster details using the Kubernetes dashboard. After you install Pega software, you can use this dashboard to review the status of all of the related Kubernetes objects used in your deployment. Without a deployment, only EKS resources display. The dashboard does not display your EKS cluster name or your resource name, which is expected behavior.

    To continue using the Kubernetes dashboard to see the progress of your deployment, keep this PowerShell or Linux shell open.

At this point, you must create a database instance into which you must install Pega Platform.

### Creating a database resource

Amazon EKS deployments require that you install Pega Platform software in an Amazon RDS instance that contains a PostgreSQL database. In general, it is recommended, but not required to name this database `pega`; the idea is to ensure that your organization recognizes the AWS database resource in which you install Pega Platform rules and the data your Pega deployment manages.

During the configuration screens with which you create an RDS instance for your EKS deployment, you can create a database and name it `pega` and add a password during the DB instance provisioning in the **Additional configuration** section. After the DB instance is created, this runbook section includes an optional section to review and optimize the database for any specific organization requirements before you install and deploy Pega using this database.

#### Creating an Amazon RDS instance

Create a database that is available to your EKS cluster. At a minimum, for a simple Pega demo, your Amazon RDS instance must be provisioned with at least 20GB of storage to support the installation of a Pega Platform software and a limited space for processing data; to provision a Pega Platform deployment to store typical Pega workload data, the size varies by workload, but should be over 70GB of total storage.

Pega Platform deployments require that the region you select is located in the same region where your EKS worker nodes are located.

1. Use a web browser to log in to <https://aws.amazon.com/console/>, which logs you into your default region.

2. In the AWS Management Console, use the search tool to navigate to **Amazon RDS** main page.

3. In **Amazon RDS** main page, in the Create Database section, click **Create database**.

    This creates an RDS database in your default region. If this region is different from the region where your EKS worker nodes are located, you must reset your region to the region when you created your EKS cluster.

4. In a new AWS console, create an Amazon RDS parameter group for this new database that it can use to limit the number of database connections.

   a.  Open the Amazon RDS navigation (menu in upper left of screen) and right-click **Parameter groups** to open the menu in a new screen.

   b. Click **Create parameter group**.

   c. Enter the following group details:

   - in **Parameter group family** select the database family to which this parameter group will apply; this family is based on the PostgreSQL engine and version you chose for your database instance.

   - in **Type** select `DB parameter group`.

   - in **Group name** enter an applicable name for this database parameter group.

   - in **Description** enter any distinguishing details for this parameter group, such as `DB connection limitation settings`.

   d. Click **Create**.

   e. Select the new parameter group created and click on it to open.

   f. In the search area, enter `max_connections`, select the parameter and click **Edit parameters**.

   g. Update the value to a minimum of 5000 and then click **Save changes**.

   With your changes saved, you can close the parameter group window in your browser.

5. To create your database, return to the **Create database** page and select the following options:

   a. For **Choose a database creation method**, select **Standard create**.

   b. For **Engine Options**, select **PostgreSQL**.

   c. For **Version**, select the latest available PostgreSQL version 11 series. For example, select `PostgreSQL 11.9-R1`.

   d. For **Templates**, select **Production**.

   e. In the **Settings** section, specify these details:

   - In the **DB instance identifier**, enter a unique \<*databasename*\>.
   - In **Credentials Settings**, leave **Master username** set to the default.
   - Create a master password for your database that meets your organization standards and then confirm it.

   f. In the **DB instance size** section, add these details:

   - Select **Standard classes**.
   - Select **db.m5.large** or greater. The **db.m5.large** selection provides the minimum hardware requirement for Pega Platform installations (a minimum of **4** vCores and **32 GB** storage).
RAM
   g. In the **Storage** section, for details, accept the default values:

   - **Storage type** is **Provisioned IOPS**.
   - **Allocated storage** is **100GiB**. While the size can be a 20GiB, this is not enough storage for processing. Pega recommends a minimal storage size of 100GiB to store typical Pega application workload data. Provisioning a smaller disk than 100GiB can impact performance.
   - **Provisioned IOPS** is **3000**.
   - **Storage autoscaling** and a **Maximum storage threshold** can be selected to allow the storage to increase if your deployment exceeds the threshold you define. Clear this setting if you do not want storage autoscaling enabled on the database instance.

   h. In the **Availability & durability** > **Multi-AZ deployment** section, for details, accept the default value, **Create a standby instance (recommended for production usage)**.

   i. In the **Connectivity** section, choose the existing VPC cluster network, so that you create the database in the same VPC as cluster:

   - **Virtual private cloud (VPC)** select your cluster VCP which should be accessible by your organization.
   - **Subnet group** select **default** or specify a subnet your deployment requires.
   - **Public access** select **no** unless your organization requires public access and you have appropriate security protocols for your subnet in place.
   - **VPC security group** select **Create new** and in **New VPC security group name**, specify an appropriate `security-group-name`.
   - If required, click **Additional configuration** and specify a non-default **Database port** if your access requires a port other than port `5432`.

    j. In **Database authentication** > **Database authentication options** select backups `Password authentication`.

    k. In the **Additional configuration** add the following settings and configurations:

    - **Database options** - create and name a database in this RDS DB instance and add a secure password. When the instance is provisioned, Amazon RDS creates this database within this DB instance. If you do not specify a database name, Amazon RDS does not create a database when it creates the DB instance and you must do it later before you install Pega Platform.

    - **DB parameter group** - select the parameter group you created your organization's deployment in step 4.

    - **Backup** - clear `Enable automatic backups` unless it is required by your organization's deployment.

    - **Encryption** - clear `Enable encryption` unless it is required by your organization's deployment.

    - **Performance Insights** - clear `Enable Performance Insights` unless your organization's deployment can take advantage of this RDS database performance monitoring feature that makes it easy to diagnose and solve performance challenges on Amazon RDS databases.

    - **Monitoring** - clear `Enable Enhanced monitoring` unless your organization can take advantage of this AWS  enhanced monitoring metrics service.

    - **Log exports** - skip this section unless your organization requires integration with Amazon CloudWatch Logs.

    - **Maintenance** - clear `Enable auto minor version upgrade` unless your organization's deployment requires automated version upgrades; In the **Maintenance window section, leave the default `No preference`selection.

    - **Deletion protection** - clear `Enable deletion protection` unless your organization's deployment requires it.

    l. After reviewing your choices and approving the **Estimated monthly costs** for these choices,  click **Create database**.

    A deployment progress page displays the status of your deployment until it is complete, which takes up to 5 minutes. When complete, the EKS UI displays all of the database resources in your account, which includes your newly created RDS DB instance, which should be in an **Available** state.

6. Find and note the Security group ID of your Linux worker nodes in your deployment to later add each of these security groups to your RDS instance security group (in the next step).

   a. Open an EC2 dashboard, click **Instances** in the left navigation pane and select one of your worker nodes in the **Instances** page.

   b. In the instance details area below the primary list of instances, click the **Security** tab, copy the Security groups to a text editor for use in the next step.

   c. Repeat steps a. and b. for each of your deployment's worker nodes.

After you note the security groups for each of your deployment's worker nodes, proceed to the next step.

7. Update the RDS security group to provide inbound access from each of your deployment's worker nodes.

   a. Open the Amazon RDS navigation (menu in upper left of screen), select **Databases** and open the page for your new PostgreSQL DB instance.

   b. In the **Connectivity & Security** tab, in the **Security** area, click the RDS security group to open it.

   c. In the *SG-security-group-name* page, in the **Inbound rules** tab, click **Edit inbound rules**.

   d. In the **Edit inbound rules** page, click **Add rule** with the following parameters:

   - **Type** - select **PostgreSQL**.
   - **Protocol** and **Port range** are set automatically.
   - **Source** - select **Custom** and use the search tool to find a security group associated with the worker nodes in your organization's deployment and then select that security group.

   e. Repeat step d. to add an inbound rule that is associated with each of your organization's deployment worker node security groups.

   f. After you add all of the required inbound rules for all of the security groups click **Save**.

After you complete these steps, your EKS cluster has the  appropriate access to the RDS DB instance and you can proceed to the next section.

#### Optional: Reviewing or creating a database in your RDS DB instance

The RDS DB instance must contain a database in order to install Pega Platform in your EKS deployment. If you did not create a PostgreSQL database in your new RDS DB instance in step 5k, you must do it now. If you did, you can also review the database and add any organizational requirements other than a password. Use the database editing tool of your choice to log into your RDS DB instance and create this new PostgreSQL database. The following example used pgAdmin4.

1. Use a database editor tool, such as pgadmin4, to log into your RDS DB instance.

    You can find your access information and login credentials, by selecting the RDS DB instance in the EKS console.

2. If you have to create a new database, in the database editor tool, navigate to Databases and create a new database; if you already have a database created, you can review the settings to ensure it meets your organization's database requirements.

   Pega does not require any additional configuration.

With new database available in your RDS DB instance, you are ready to continue to the next section.

## Installing and deploying Pega Platform using Helm charts – 90 minutes

To deploy Pega Platform by using Helm, customize the default Helm charts that holds the specific settings for your deployment needs and then run a series of Helm commands to complete the deployment.

An installation with deployment will take about 90 minutes total, because a Pega Platform installation in your PostgreSQL database takes up to an hour.

### Adding the Pega configuration files to your Helm installation on your local system

Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm, including a generic version of the following charts. After you add the repository to your local system, you can customize these Pega configuration files for your Pega Platform deployment:

- pega/addons – Use this chart to install any supporting services and tools which your Kubernetes environment will require to support a Pega deployment: the required services, such as a load balancer or metrics server, that your deployment requires depend on your cloud environment. For instance you can specify whether you want to use a generic load-balancer or use one that is offered in your Kubernetes environment. With the instructions in this runbook, you deploy these supporting services once per Kubernetes environment when you install the addons chart, regardless of how many Pega Infinity instances are deployed.

- pega/pega - Use this chart to set customization parameters for your deployment. You will modify this chart later in the deployment tasks.

- pega/backingservices - Use this chart to set customization parameters for the Pega-provided Search and Reporting Service (SRS) your deployment. You will modify this chart later in the deployment tasks.

To customize these files, you must download them from the source github repository to your local system, edit them with a text editor, and then save them to your local system using the same filename.

1. To add the Pega repository to your Helm installation, enter:

    `$ helm repo add pega https://pegasystems.github.io/pega-helm-charts`

2. To verify the new repository, you can search it by entering:

```
  $ helm search repo pega
  NAME                  CHART VERSION   APP VERSION     DESCRIPTION
  pega/pega             2.2.0                           Helm chart to configure required installation and deployment configuration settings in your environment for your deployment.
  pega/addons           2.2.0           1.0             Helm chart to configure supporting services and tools in your environment for your deployment.
  pega/backingservices  2.2.0                           Helm Chart to provision the latest Search and Reporting Service (SRS) for your Pega Infinity deployment.
```

The pega and addons charts in the /charts/pega folder require customization for your organization's EKS deployment of Pega Platform. The backingservices chart is optional, but recommended for Pega Infinity 8.6 and later.

#### Updating the addons.yaml Helm chart values

To configure the use of an Amazon AWS ALB ingress controller in the addons.yaml file, download the file in the charts/addons folder, edit it with a text editor, and then save it to your local system using the same filename.

1. Download the example [addons-eks.yaml](./resources/addons-eks.yaml) to the \<local filepath\>/EKS-demo that will over-write the default pega/addons file.

   This example addons file specifies the use of an Amazon AWS ALB ingress controller and contains several parameters that will specify details from your EKS environment so your Pega Platform deployment can use the load balancer in your EKS environment.

2. Use a text editor to open the addons-eks.yaml file and update the following parameters in the chart based on your EKS requirements:

- Specify `aws-load-balancer-controller:> enabled: true` to install an AWS load balance -controller for your deployment.
- Specify your EKS cluster name in the `clusterName: <YOUR_EKS_CLUSTER_NAME>` parameter.
- Specify the region of your EKS cluster name in the `region: <YOUR_EKS_CLUSTER_REGION>` parameter. Resources created by the ALB Ingress controller will be prefixed with this string.
- Specify the the AWS VPC ID of your EKS cluster name in the `VpcID: <YOUR_EKS_CLUSTER_VPC_ID>` parameter. You must enter your VPC ID here if ec2metadata is unavailable from the controller pod.
- Uncomment and specify the Amazon EKS Amazon ECR image repository in the image.repository: <Amazon EKS Amazon ECR image repository> parameter. This is required for AWS GovCloud deployments
- Specify complete  required required annotation to specify the role that you associate with the primary IAM user who is responsible for your EKS deployment in the `serviceAccount.annotations.eks.amazonaws.com/role-arn: <YOUR_IAM_ROLE_ARN>` parameter.

To ensure logging for your deployment is properly configured to take advantage of the built-in EFK logging tools in EKS deployments, refer to the [Amazon EKS Workshop](https://eksworkshop.com/logging/).

#### Updating the backingservices.yaml Helm chart values for the SRS (Supported when installing or upgrading to Pega Infinity 8.6 and later)

To configure the parameters in the backingservices.yaml file, download the file in the charts/backingservices folder, edit it with a text editor, and then save it to your local system using the same filename.

1. To download the backingservices.yaml Helm chart to the \<local filepath>\eks-demo, enter:

  `$ helm inspect values pega/backingservices > <local filepath>/eks-demo/backingservices.yaml`

2. Use a text editor to open the backingservices.yaml file and update the following parameters in the chart based on your EKS requirements:

| Chart parameter name              | Purpose                                   | Your setting |
|:---------------------------------|:-------------------------------------------|:--------------|
| global.imageCredentials.registry: username: password:  | Include the URL of your Docker registry along with the registry “username” and “password” credentials. | <ul><li>url: “\<URL of your registry>” </li><li>username: "\<Registry account username\>"</li><li> password: "\<Registry account password\>"</li></ul> |
| srs.deploymentName:        | Specify unique name for the deployment based on org app and/or SRS applicable environment name.      | deploymentName: "acme-demo-dev-srs"   |
| srs.srsRuntime.srsImage: | Specify the Pega-provided SRS Docker image that you downloaded and pushed to your Docker registry. | srs.srsRuntime.srsImage: "\<Registry host name:Port>my-pega-srs:\<srs-version>". For `<srs-version>` tag details, see [SRS Version compatibility matrix](../charts/backingservices/README.md#srs-version-compatibility-matrix).    |
| srs.srsStorage.provisionInternalESCluster: | Enabled by default to provision an Elasticsearch cluster. | <ul><li>Set srs.srsStorage.provisionInternalESCluster:`true` and run `$ make es-prerequisite NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT>`</li><li>Set srs.srsStorage.provisionInternalESCluster:`false` if you want to use an existing, externally provisioned ElasticSearch cluster. </li></ul> |
| srs.srsStorage.domain: port: protocol: basicAuthentication: awsIAM: requireInternetAccess: | Disabled by default. Enable only when srs.srsStorage.provisionInternalESCluster is false and you want to configure SRS to use an existing, externally provisioned Elasticsearch cluster. For an Elasticsearch cluster secured with Basic Authentication, use `srs.srsStorage.basicAuthentication` section to provide access credentials. For an AWS Elasticsearch cluster secured with IAM role based authentication, use `srs.srsStorage.awsIAM` section to set the aws region where AWS Elasticsearch cluster is hosted. For unsecured managed ElasticSearch cluster do not configure these options. | <ul><li>srs.srsStorage.domain: "\<external-es domain name\>"</li> <li>srs.srsStorage.port: "\<external es port\>"</li> <li>srs.srsStorage.protocol: "\<external es http protocol, `http` or `https`\>"</li>     <li>srs.srsStorage.basicAuthentication.username: "\<external es `basic Authentication username`\>"</li>     <li>srs.srsStorage.basicAuthentication.password: "\<external es `basic Authentication password`\>"</li>     <li>srs.srsStorage.awsIAM.region: "\<external AWS es cluster hosted `region`\>"</li><li> srs.srsStorage.requireInternetAccess: "\<set to `true` if you host your external Elasticsearch cluster outside of the current network and the deployment must access it over the internet.\>"</li></ul>     |
| elasticsearch: volumeClaimTemplate: resources: requests: storage: | Specify the Elasticsearch cluster disk volume size. Default is 30Gi, set this value to at least three times the size of your estimated search data size | <ul><li>elasticsearch: volumeClaimTemplate: resources: requests: storage:  "\<30Gi>” </li></ul> |

3. Save the file.

4. To use an internal Elasticsearch cluster (srs.srsStorage.provisionInternalESCluster:true) for your deployment, you must run `$ make es-prerequisite NAMESPACE=<NAMESPACE_USED_FOR_DEPLOYMENT>`.

#### Add supported custom settings for Pega to your deployment

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

#### (Optional) Add Support for providing DB credentials using External Secrets Operator

Create two files following the Kubernetes documentation for External Secrets Operator [External Secrets Operator](https://external-secrets.io/v0.5.3/) :
  •	An external secret file that specifies what information in your secret to fetch.
  •	A secret store to define access how to access the external and placing the required files in your Helm directory.

- Copy both files into the pega-helm-charts/charts/pega/templates directory of your Helm
- Update repo to the latest-> helm repo update pega https://pegasystems.github.io/pega-helm-charts
- Update Pega.yaml file to refer to the external secret manager for DB password.

#### Updating the pega.yaml Helm chart values

To configure the parameters in the pega.yaml fie, download the file in the charts/pega folder, edit it with a text editor, and then save it to your local system using the same filename.

Configure the parameters so the pega.yaml Helm chart matches your deployment resources in these areas:

- Specify that this is an EKS deployment.

- Credentials for your DockerHub account in order to access the required Docker images.

- Access your RDS database.

- Access your ElasticSearch service (For 8.6 and later, Pega recommends deploying your service using an SRS cluster).

- Install the version of Pega Platform built into your Docker installation image.

- Specify host names for your web and stream tiers and import and use any required SSL certificates for your web tiers.

- Enable encryption of traffic between the ingress/load balancer and the pods by specifying SSL certificates for your web tiers.

- Enable Hazelcast client-server model for Pega Platform 8.6 and later.

1. To download the pega.yaml Helm chart to the \<local filepath\>/EKS-demo, enter:

   `$ helm inspect values pega/pega > pega.yaml`

2. Use a text editor to open the pega.yaml file and update the following parameters in the chart based on your EKS requirements:

   | Chart parameter name    | Purpose                                   | Your setting |
   |-------------------------|-------------------------------------------|--------------|
   | provider:               | Specify an EKS deployment.                 | provider:"eks"|
   | actions.execute:        | Specify a “deploy” deployment type.       | execute: "deploy"   |
   | jdbc.url:               | Specify the database IP address and database name for your Pega Platform installation.        | <ul><li>url: "jdbc:postgresql://**localhost**:5432/`pega`"</li><li>where **localhost** is the public IP address you configured for your database connectivity and `pega` is the recommended name of the database you created in your RDS DB instance.</li></ul>  |
   | jdbc.driverClass:       | Specify the driver class for a PostgreSQL database. | driverClass: "org.postgresql.Driver"|
   | jdbc.dbType:            | Specify PostgreSQL database type.         | dbType: "postgres”   |
   | jdbc.driverUri:         | Specify the database driver Pega Platform uses during the deployment.| <ul> <li>driverUri: "latest jar file available” </li> <li>For PostgreSQL databases, use the URL of the latest PostgreSQL driver file that is publicly available at <https://jdbc.postgresql.org/download.html>.</li></ul>    |
   | jdbc: username: password: | Set the security credentials for your database server to allow installation of Pega Platform into your database.   | <ul><li>username: "\<name of your database user\>" </li><li>password: "\<password for your database user\>"</li><li>-- For RDS postgreSQL databases, previously set default <username>.</li></ul>   |
   | jdbc.rulesSchema: jdbc.dataSchema:  | Set the names of both your rules and the data schema to the values that Pega Platform uses for these two schemas.      | <ul><li>rulesSchema: "rules" </li><li>dataSchema: "data"</li></ul>  |
   | customArtifactory.authentication: basic.username: basic.password: apiKey.headerName: apiKey.value: | To download a JDBC driver from your custom artifactory which is secured with Basic or APIKey Authentication. Use `customArtifactory.authentication.basic` section to provide access credentials or use `customArtifactory.authentication.apiKey` section to provide APIKey value and dedicated APIKey header details. | <ul><li>basic.username: "\<Custom artifactory basic Authentication username\>"</li><li>basic.password: "\<Custom artifactory basic Authentication password\>"</li><li>apiKey.headerName: "\<Custom artifactory dedicated APIKey Authentication header name\>"</li><li>apiKey.value: "\<Custom artifactory APIKey value for APIKey authentication\>"</li> </ul> |
   | customArtifactory.certificate: | Custom artifactory SSL certificate verification is enabled by default. If your custom artifactory domain has a self-signed SSL certificate, provide the certificate. You can disable SSL certificate verification by setting `customArtifactory.enableSSLVerification` to `false`;however, this setting establishes an insecure connection. | <ul><li> certificate: "\<custom artifactory SSL certificate to be verified\>"</li></ul> |
   | docker.registry.url: username: password: | Map the host name of a registry to an object that contains the “username” and “password” values for that registry. For more information, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). | <ul><li>url: “<https://index.docker.io/v1/>” </li><li>username: "\<DockerHub account username\>"</li><li> password: "\< DockerHub account password\>"</li></ul>    |
   | docker.pega.image:       | Specify the Pega-provided `Pega` image that you downloaded and pushed to your Docker registry.  | Image: "\<Registry host name:Port\>/my-pega:\<Pega Platform version>" |
   | tier.name: ”web” tier.service.domain:| Set a host name for the pega-web service of the DNS zone. To support the use of HTTPS for ingress connectivity enable SSL/TLS termination protocols on the tier ingress and provide your ARN certificate, where `alb.ingress.kubernetes.io/certificate-arn` is the required annotation name and `<certificate-arn>` takes the form, `arn:aws:acm:<region>:<AWS account>:certificate/xxxxxxx` which you copy from the AWS console view of the load balancer configuration.| <ul><li>domain: "\<the host name for your web service tier\>" </li><li>ingress.tls.enabled: "true"</li><li>ingress.ssl_annotation: alb.ingress.kubernetes.io/certificate-arn: \<certificate-arn></li><li>Assign this host name with the DNS host name that the load balancer associates with the web tier; after the deployment is complete, you can log into Pega Platform with your host name in the URL. Your web tier host name must comply with your networking standards and be available on an external network.</li></ul> |
   | tier.name: ”web” tier.service.tls:| Set this parameter as `true` to encrypt the traffic between the load balancer/ingress and pods. Beginning with Helm version `2.2.0` Pega provides a default self-signed certificate; Pega also supports specifying your own CA certificate. | <ul><li>tier.service.tls.enabled: set to `true` to enable the traffic encryption </li><li>tier.service.tls.port: 443</li><li>tier.service.tls.targetPort: 8443</li><li>tier.service.tls.external_secret_name: The external secret name for fetching certificates from a secure location. For details, see [this option in the Pega Helm chart](../charts/pega#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)</li><li>tier.service.tls.keystore: The base64 encoded content of the keystore file. Leave this value empty to use the default, Pega-provided self-signed certificate.</li><li>tier.service.tls.keystorepassword: the password of the keystore file</li><li>tier.service.tls.cacertificate: the base64 encrypted content of the root CA certificate.  You can leave this value empty for EKS deployments.</li><li>tier.service.traefik.enabled:  set to `false` as this option is for `k8s` provider not for `EKS`</li></ul> |
   | tier.name: ”stream” tier.service.domain: | Set the host name for the pega-stream service of the DNS zone.   | <ul><li>domain: "\<the host name for your stream service tier\>" </li><li>Your stream tier host name should comply with your networking standards. </li></ul> |
   | pegasearch: | For Pega Platform 8.6 and later, Pega recommends using the chart 'backingservices' to enable Pega SRS. To deploy this service, you must configure your SRS cluster using the backingservices Helm charts and provide the SRS URL for your Pega Infinity deployment. | <ul><li>externalSearchService: true</li><li>externalURL: pegasearch.externalURL For example, http://srs-service.mypega-pks-demo.svc.cluster.local </li></ul> |
   | installer.image:        | Specify the Pega-provided Docker `installer` image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
   | installer. adminPassword:                | Specify an initial administrator@pega.com password for your installation.  This will need to be changed at first login. The adminPassword value cannot start with "@".     | adminPassword: "\<initial password\>"  |
   | hazelcast: | For Pega Platform 8.6 and later, Pega recommends using Hazelcast in client-server model. Embedded deployment would not be supported in future platform releases.| |
   | hazelcast.image:        | Specify the Pega-provided `clustering-service` Docker image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
   | hazelcast.enabled: hazelcast.replicas: hazelcast.username: hazelcast.password: | Either to enable Hazelcast in client-server model and configure the number of replicas and username & passowrd for authentication | <ul><li>enabled: true/false <br/> Set to true if you want to deploy pega platform in client-server Hazelcast model, otherwise false. *Note: Set this value as false for Pega platform versions below 8.6, if not set the installation will fail.* </li><li>replicas: <No. of initial server members to join(3 or more based on deployment)> </li><li>username: "\<UserName for authentication\>" </li><li> password: "\<Password for authentication\>" </li></ul> |

3. Save the file.

4. To use an internal Elasticsearch cluster (srs.srsStorage.requireInternetAccess:true) for your deployment, you must download the Makefile file from the repository ( path from root : pega-helm-charts/charts/backingservices/Makefile) and replace <YOUR_NAMESPACE> with the namespace you used for the deployment, then run `$ es-prerequisite`.

### Deploy Pega Platform using the command line

A Helm installation and a Pega Platform installation are separate processes. The Helm install command uses Helm to install your deployment as directed in the Helm charts, one in the **charts\\addons** folder and one in the **charts\\pega** folder.

In this document, you specify that the Helm chart always “deploys” by using the setting, actions.execute: “deploy”. In the following tasks, you overwrite this function on your *initial* Helm install by specifying `--set global.actions.execute:install-deploy`, which invokes an installation of Pega Platform using your installation Docker image and then
automatically followed by a deploy. In subsequent Helm deployments, you should not use the override argument, `--set global.actions.execute=`, since Pega Platform is already installed in your database.

1. Do one of the following:

   - Open Windows PowerShell running as Administrator on your local system and change the location to the top folder of your EKS-demo folder that you created in [Preparing your local Windows 10 system](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md).

   `$ cd <local filepath>\EKS-demo`

   - Open a Linux bash shell and change the location to the top folder of your EKS-demo directory that you created in [Preparing your local Linux system](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md).

   `$ cd /home/<local filepath>/EKS-demo`

2. Create namespaces in preparation for the pega.yaml and addons.yaml deployments.

   ```yaml
   $ kubectl create namespace mypega-EKS-demo
   namespace/mypega-EKS-demo created
   $ kubectl create namespace pegaaddons
   namespace/pegaaddons created
   ```

3. Install the addons Helm chart, which you updated in [Updating the addons.yaml Helm chart values](#Updating-the-addonsyaml-Helm-chart-values).

   ```yaml
   $ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
   ```

   The `pegaaddons` namespace contains the deployment’s load balancer and the metric server configurations that you configured in the addons.yaml Helm chart. A successful pegaaddons deployment returns details of deployment progress. For further verification of your deployment progress, you can refresh the Kubernetes dashboard and look in the `pegaaddons` **Namespace** view.

4. For Pega Platform 8.6 and later installations, to install the backingservices chart that you updated in [Updating the backingservices.yaml Helm chart values (Supported when installing or upgrading to Pega Infinity 8.6 and later)](#Updating the backingservices.yaml Helm chart values (Supported when installing or upgrading to Pega Infinity 8.6 and later)), enter:

   ```yaml
   $ helm install backingservices pega/backingservices --namespace mypega-EKS-demo --values backingservices.yaml
   ```
   The `mypega-EKS-demo` namespace used for pega deployment can also be used for backingservice deployment that you configured in backingservices.yaml helm chart.

5. Deploy Pega Platform for the first time by installing the pega Helm chart, which you updated in [Updating the addons.yaml Helm chart values](#Updating-the-addonsyaml-Helm-chart-values). This installs Pega Platform software into the database you specified in the pega chart.

   ```yaml
   helm install mypega-EKS-demo pega/pega --namespace mypega-EKS-demo --values pega.yaml --set global.actions.execute=install-deploy
   ```

   For subsequent Helm installs, use the command `helm install mypega-EKS-demo pega/pega --namespace mypega-EKS-demo` to deploy Pega Platform and avoid another Pega Platform installation.

   A successful Pega deployment immediately returns details that show progress for your `mypega-EKS-demo` deployment.

6. Refresh the Kubernetes dashboard that you opened in the previous section. If you closed the dashboard, start the proxy server for the Kubernetes dashboard and then relaunch the web browser.

7. In the dashboard, in **Namespace** select the `mypega-EKS-demo` view and then click on the **Pods** view. Initially, you can some pods have a red status, which means they are initializing.

    Note: A deployment takes about 15 minutes for all resource configurations to initialize; however a full Pega Platform installation into the database can take up to an hour.

    To follow the progress of an installation, use the dashboard. For subsequent deployments, you do not need to do this. Initially, while the resources make requests to complete the configuration, you will see red warnings while the configuration is finishing, which is expected behavior.

8. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-db-install** job, and click the logs icon on the right side of that row.

    After you open the logs view, you can click the icon for automatic refresh to see current updates to the install log.

9. To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the `mypega-EKS-demo` namespace pods.

    A successful deployment does not show errors across the various workloads. The `mypega-EKS-demo` Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment has 100% complete **Workloads**.

## Logging in to Pega Platform – 10 minutes

After you complete your deployment, as a best practice, associate the host name of the pega-web tier ingress with the DNS host name that the deployment load balancer assigned to the tier during deployment. The host name of the pega-web tier ingress used in this demo, **eks.web.dev.pega.io**, is set in the pega.yaml file in the following lines:

```yaml
tier:
  - name: "web"

    service:
      # Enter the domain name to access web nodes via a load balancer.
      #  e.g. web.mypega.example.com
      domain: "eks.web.dev.pega.io"
```

To log in to Pega Platform with this host name, you can log into your ingress load balancer and note the DNS host name that the load balancer associates with web tier; after you copy the DNS host name, you can assign the host name you gave to the web tier with the DNS host name that the deployment load balancer assigned to the web tier. This final step ensures that you can log in to Pega Platform with the host name you configured for your deployment in the pega Helm chart, so you can independently manage security protocols that match your networking infrastructure standards.

To manually associate the host name of the pega-web tier ingress with the tier endpoint, use the DNS lookup management system of your choice. If your organization has an AWS Route 53 DNS lookup service already established to manage your DNS lookups, use the Route 53 Dashboard to create a record set that specifies the pega-web tier the host name and add the DNS host name you found when you log on the load balancer.

For AWS Route53 Cloud DNS lookup service documentation details, see [What is Amazon Route 53?](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/Welcome.html). If not using AWS Route53 Cloud DNS lookup service, see the documentation for your DNS lookup service.

With the ingress host name name associated with this DNS host host in your DNS service, you can log in to Pega Platform with a web browser using the URL: http://\<pega-web tier ingress host name>/prweb.
