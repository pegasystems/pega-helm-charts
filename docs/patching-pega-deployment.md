# Patching Pega Platform in your deployment

After you deploy Pega Platform™ on your kubernetes environment, the Pega-provided Docker images support applying a zero-downtime patch to your Pega software. These procedures are written for any level of user, from a system administrator to a development engineer who wants to use helm charts and Pega Docker images to patch the Pega software have deployed in any supported kubernetes environment.

Useful links to Pega software patching information:

- [Pega software maintenance and extended support policy](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-software-maintenance-and-extended-support-policy)
- [Pega Infinity patch calendar](https://community.pega.com/knowledgebase/articles/keeping-current-pega/pega-infinity-patch-calendar)
- [Pega Infinity patch frequently asked questions](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-infinity-patch-frequently-asked-questions)

## Deployment process overview

Pega supports client-managed cloud clients applying patches for releases 8.4 and later using a zero-time patch process to apply the latest cumulative bundle of bug and security fixes since the last minor release. For the latest Pega Community articles, see [About client managed cloud](https://community.pega.com/knowledgebase/articles/client-managed-cloud/85/about-client-managed-cloud).

The Pega zero-downtime patch process uses the out-of-place patch process so you and your customers can continue working in your application while you patch your system. Pega zero-downtime patch scripts use a temporary data schema and the patch migration script moves the rules between the appropriate schema and then performs the required rolling reboot of your deployment cluster. For a detailed overview of the process, see [Applying a patch without downtime](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/applying-patch-without-downtime).

For releases 8.4 and later, client-managed cloud clients use the same Pega Kubernetes tools and Helm charts in the same Pega repository that you used to install Pega Platform in a supported Kubernetes environment. The client-managed cloud patch process includes the following tasks:

1. Prepare your Docker repository by downloading the latest three Pega Platform patch release images (platform/installer, platform/pega, and platform/search) in your release stream and pushing them into your preferred Docker image repository - [Downloading your images for the patch process – 20 minutes](#downloading-your-images-for-the-patch-process--20-minutes).

2. Edit the pega Helm chart by editing parameters to specify "upgrading" your software with the software contained in your provided patch image. - [Applying a zero-downtime Pega Platform patch using Helm charts - 120-minutes](#applying-a-zero--downtime-pega-platform-patch-using-helm-charts--120-minutes).

3. For deployments running Pega Infinity 8.3.1 through 8.3.4, create a new blank rules schema and a new temporary database schema in your existing database. Leave these new schemas empty. If you are running Pega Infinity 8.4 or higher, you can skip this step, since the patch scripts in your deployment automate the creation of these blank schemas.

   Pega does not support quoted identifiers, so do not wrap your schema name with single quotes. Create each schema name in accordance with the requirements of your deployment database type:
   - Oracle/DB2 databases force unquoted identifiers to uppercase.
   - PostgreSQL databases force unquoted identifiers to lowercase.
   - MSSQL uses case sensitive identifiers; therefore you must use a consistent naming convention in order to avoid issues with your deployment.

4. Apply the patch by using the `helm update` command as directed in the deployment section - [Deploy Pega Platform using the command line](#deploy-pega-platform-using-the-command-line).

## Assumptions and prerequisites

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges or a Linux command prompt with root privileges.

- You use open source packaging tools on Windows or Linux to install applications onto your local system.

- Your Kubernetes environment has not changed since you originally deployed, and you still have the appropriate account permissions and knowledge to run Helm commands on your environment.

- You can download the latest Pega distribution in the release stream with which you initially deployed.

- Your original deployment used Pega Platform 8.3.1 or later.

## Downloading your images for the patch process – 20 minutes

Patching the Pega software in your deployment requires the use of the required Docker images built from the same release stream for the software running in your deployment. To do this, you must download and make all three of the new images available in your preferred Docker registry. For details, see [Downloading Docker images for your deployment](https://github.com/pegasystems/pega-helm-charts#downloading-docker-images-for-your-deployment).

See the next section to make properly reference these images in your Pega Helm chart.

## Applying a zero-downtime Pega Platform patch using Helm charts – at least 120 minutes

To apply a patch to the Pega Platform software in your existing deployment by using Helm, you must customize your existing `pega` Helm chart with the specific, required settings. To keep your system current with a patch, you must use the latest Pega-provided Docker images that are available for your release stream.

The Pega patch application process takes about 15 minutes total.

### Updating the Pega configuration files to your Helm installation on your local system

Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm. To complete a zero downtime patch, you must configure the following settings in your existing Pega configuration files for your Pega Platform deployment:

- Specify action.execute: upgrade-deploy to invoke the zero-downtime patch process.
- Ensure one of the following:
  - You pushed the images for your patch to the same repository that you used for your installation repository and the credentials for your repository account are the same as those in your `pega` Helm chart.  
  - You pushed the images for your patch to a new repository and you update the parameters with your credentials for this new repository account in your `pega` Helm chart.

- Specifying the repository and Docker installation image details 

- .

You can leave the existing customized paramters as is; the patch process will use the remaining existing settings in your deployment.

1. Use a text editor to open the pega.yaml file and update the following parameters in the chart based on your EKS requirements:

   | Chart parameter name    | Purpose                                   | Your setting |
   |-------------------------|-------------------------------------------|--------------|
   | actions.execute: | Specify an “upgrade” deployment type. | execute: "upgrade" |
   | jdbc.connectionProperties.rulesSchema: "YOUR_RULES_SCHEMA"  | Specify the name of the new rules schema to which the patch process migrates the existing rules structure from the existing rules schema  | rulesSchema: "YOUR_RULES_SCHEMA" |
   | jdbc.connectionProperties.dataSchema: "YOUR_DATA_SCHEMA"  | Specify the name of the temporary data schema to which the patch process migrates the existing data structure from the existing data schema  | dataSchema: "YOUR_DATA_SCHEMA"  |
   | docker.registry.url: username: password: | Map the host name of a registry to an object that contains the “username” and “password” values for that registry. For more information, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). | <ul><li>url: “<https://index.docker.io/v1/>” </li><li>username: "\<DockerHub account username\>"</li><li> password: "\< DockerHub account password\>"</li></ul>    |
   | docker.pega.image:       | Specify the new Pega-provided `Pega` image that matches the version of the installer image with which you will apply your patch. You downloaded and pushed this image when you pushed your new installer image to your Docker registry.  | Image: "\<Registry host name:Port\>/my-pega:\<Pega Platform version>" |
   | upgrade.kube-apiserver. serviceHost & upgrade.kube-apiserver.httpsServicePort  | For the Helm charts 1.2 through 1.4, you must specify upgrade to invoke the zero-downtime patch process | upgrade.kube-apiserver.serviceHost: "API_SERVICE_ADDRESS" upgrade.kube-apiserver.httpsServicePort: "SERVICE_PORT_HTTPS" |
   | tier.name: ”web” tier.service.domain:| Set a host name for the pega-web service of the DNS zone. To support the use of HTTPS for ingress connectivity enable SSL/TLS termination protocols on the tier ingress and provide your ARN certificate.| <ul><li>domain: "\<the host name for your web service tier\>" </li><li>ingress.tls.enabled: "true"</li><li>ingress.ssl_annotation: "alb.ingress.kubernetes.io/certificate-arn: \<certificate-arn>\"</li><li>Assign this host name with the DNS host name that the load balancer associates with the web tier; after the deployment is complete, you can log into Pega Platform with your host name in the URL. Your web tier host name must comply with your networking standards and be available on an external network.</li></ul>|
   | tier.name: ”stream” tier.service.domain: | Set the host name for the pega-stream service of the DNS zone.   | <ul><li>domain: "\<the host name for your stream service tier\>" </li><li>Your stream tier host name should comply with your networking standards. </li></ul>|
   | pegasearch.image: | Specify the Pega-provided Docker `search` image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-search:\<Pega Platform version>"
   | installer.image:        | Specify the Pega-provided Docker `installer` image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
   | installer. adminPassword:                | Specify an initial administrator@pega.com password for your installation.  This will need to be changed at first login. The adminPassword value cannot start with "@".     | adminPassword: "\<initial password\>"  |

3. Save the file.

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

3. Install the addons Helm chart, which you updated in [Updating the addons Helm chart values](#Updating-the-addons-Helm-chart-values).

   ```yaml
   $ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
   ```

   The `pegaaddons` namespace contains the deployment’s load balancer and the metric server configurations that you configured in the addons.yaml Helm chart. A successful pegaaddons deployment returns details of deployment progress. For further verification of your deployment progress, you can refresh the Kubernetes dashboard and look in the `pegaaddons` **Namespace** view.

4. Deploy Pega Platform for the first time by 
installing the pega Helm chart, which you updated in [Updating the pega Helm chart values](#Updating-the-pega-Helm-chart-values). This installs Pega Platform software into the database you specified in the pega chart.

   ```yaml
   helm install mypega-EKS-demo pega/pega --namespace mypega-EKS-demo --values pega.yaml --set global.actions.execute=install-deploy
   ```

   For subsequent Helm installs, use the command `helm install mypega-EKS-demo pega/pega --namespace mypega-EKS-demo` to deploy Pega Platform and avoid another Pega Platform installation.

   A successful Pega deployment immediately returns details that show progress for your `mypega-EKS-demo` deployment.

5. Refresh the Kubernetes dashboard that you opened in the previous section. If you closed the dashboard, start the proxy server for the Kubernetes dashboard and then relaunch the web browser.

6. In the dashboard, in **Namespace** select the `mypega-EKS-demo` view and then click on the **Pods** view. Initially, you can some pods have a red status, which means they are initializing:

    ! [](media/dashboard-mypega-pks-demo-install-initial.png)

    Note: A deployment takes about 15 minutes for all resource configurations to initialize; however a full Pega Platform installation into the database can take up to an hour.

    To follow the progress of an installation, use the dashboard. For subsequent deployments, you do not need to do this. Initially, while the resources make requests to complete the configuration, you will see red warnings while the configuration is finishing, which is expected behavior.

7. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-db-install** job, and click the logs icon on the right side of that row.

    After you open the logs view, you can click the icon for automatic refresh to see current updates to the install log.

8. To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the `mypega-EKS-demo` namespace pods.

    ! [](media/f7779bd94bdf3160ca1856cdafb32f2b.png)

    A successful deployment does not show errors across the various workloads. The `mypega-EKS-demo` Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment has 100% complete **Workloads**.

    ! [](media/0fb2d07a5a8113a9725b704e686fbfe6.png)

## Logging in to Pega Platform – 10 minutes

After you complete your deployment, as a best practice, associate the host name of the pega-web tier ingress with the DNS host name that the deployment load balancer assigned to the tier during deployment. The host name of the pega-web tier ingress used in this demo, **eks.web.dev.pega.io**, is set in the pega.yaml file in the following lines:

```yaml
tier:
  - name: "web"

    service:
      # Enter the domain name to access web nodes via a load balancer.
      #  e.g. web.mypega.example.com
      domain: "**eks.web.dev.pega.io**"
```

To log in to Pega Platform with this host name, you can log into your ingress load balancer and note the DNS host name that the load balancer associates with web tier; after you copy the DNS host name, you can assign the host name you gave to the web tier with the DNS host name that the deployment load balancer assigned to the web tier. This final step ensures that you can log in to Pega Platform with the host name you configured for your deployment in the pega Helm chart, so you can independently manage security protocols that match your networking infrastructure standards.

To manually associate the host name of the pega-web tier ingress with the tier endpoint, use the DNS lookup management system of your choice. If your organization has an AWS Route 53 DNS lookup service already established to manage your DNS lookups, use the Route 53 Dashboard to create a record set that specifies the pega-web tier the host name and add the DNS host name you found when you log on the load balancer.

For AWS Route53 Cloud DNS lookup service documentation details, see [What is Amazon Route 53?](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/Welcome.html). If not using AWS Route53 Cloud DNS lookup service, see the documentation for your DNS lookup service.

With the ingress host name name associated with this DNS host host in your DNS service, you can log in to Pega Platform with a web browser using the URL: http://\<pega-web tier ingress host name>/prweb.

! [](media/25b18c61607e4e979a13f3cfc1b64f5c.png)