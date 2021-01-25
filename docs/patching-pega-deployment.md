# Patching Pega Platform in your deployment

After you deploy Pega Platform™ on your kubernetes environment, the Pega-provided Docker images support applying a zero-downtime patch to your Pega software. These procedures are written for any level of user, from a system administrator to a development engineer who wants to use helm charts and Pega Docker images to patch the Pega software have deployed in any supported kubernetes environment.

Useful links to Pega software patching information:

- [Pega software maintenance and extended support policy](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-software-maintenance-and-extended-support-policy)
- [Pega Infinity patch calendar](https://community.pega.com/knowledgebase/articles/keeping-current-pega/pega-infinity-patch-calendar)
- [Pega Infinity patch frequently asked questions](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-infinity-patch-frequently-asked-questions)

## Deployment process overview

Pega supports client-managed cloud clients applying patches for releases 8.4 and later using a zero-time patch process to apply the latest cumulative bundle of bug and security fixes since the last minor release. For the latest Pega Community articles, see [About client managed cloud](https://community.pega.com/knowledgebase/articles/client-managed-cloud/85/about-client-managed-cloud).

The Pega zero-downtime patch process uses the out-of-place patch process so you and your customers can continue working in your application while you patch your system. Pega zero-downtime patch scripts use a temporary data schema and the patch migration script moves the rules between the appropriate schema and then performs the required rolling reboot of your deployment cluster. For a detailed overview of the process, see [Applying a patch without downtime](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/applying-patch-without-downtime).

Client-managed cloud clients use the same Pega Kubernetes tools and Helm charts in the same Pega repository that you used to install Pega Platform in a supported Kubernetes environment. The client-managed cloud patch process includes the following tasks:

1. Prepare your Docker repository by downloading the latest three Pega Platform patch release images (platform/installer, platform/pega, and platform/search) in your release stream and pushing them into your preferred Docker image repository - [Downloading your images for the patch process – 20 minutes](#downloading-your-images-for-the-patch-process--20-minutes).

2. Edit the pega Helm chart by editing parameters to specify "upgrading" your software with the software contained in your provided patch image. - [Applying a zero-downtime Pega Platform patch using Helm charts - 120-minutes](#applying-a-zero--downtime-pega-platform-patch-using-helm-charts--120-minutes).

3. **For deployments running Pega Infinity 8.2.1 through 8.2.7 only:** Create a new blank rules schema in your existing database. Leave this new schema empty. If you are running Pega Infinity 8.3 or higher, you can skip this step, since the patch scripts in your deployment automate the creation of these blank schemas.

   If you create a new schema, create each schema name in accordance with the requirements of your deployment database type:
   - Oracle/DB2 databases force unquoted identifiers to uppercase.
   - PostgreSQL databases force unquoted identifiers to lowercase.
   - MSSQL uses case sensitive identifiers; therefore you must use a consistent naming convention in order to avoid issues with your deployment.

   Pega does not support quoted identifiers in database schema names, so do not wrap your schema name with single quotes.

4. To prepare your system for the patch process running in the background, disable rule creation in your current deployment. Depending on your version of Pega Platform, the steps are slightly different, but the same steps used in on-premises systems:

   - For 8.4 and later, see [For High Availability: Preparing the cluster for patching](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/high-availability-systems-preparing-cluster-patching)
   - For 8.2 or 8.3, see [For High Availability 8.2 and 8.3 systems: Preparing the cluster for patching](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/high-availability-82-and-83-systems-preparing-cluster-patching)

5. Apply the patch by using the `helm upgrade-deploy` command as directed in the deployment section - [Deploy Pega Platform using the command line](#deploy-pega-platform-using-the-command-line).

## Assumptions and prerequisites

The process to patching your deployment assumes:

- Your Kubernetes environment has not changed and you are using the same Pega charts with which you originally deployed.

- You download the latest Pega docker images in the minor release stream with which you initially deployed.

- Your original deployment used Pega Platform 8.2.1 or later.

## Downloading your images for the patch process – 20 minutes

Patching the Pega software in your deployment requires the use of the required Docker images built from the same release stream for the software running in your deployment. To do this, you must download and make all three of the new images available in your preferred Docker registry. For details, see [Downloading Docker images for your deployment](https://github.com/pegasystems/pega-helm-charts#downloading-docker-images-for-your-deployment).

See the next section to properly reference these new images in your Pega Helm chart.

## Applying a zero-downtime Pega Platform patch using Helm charts – 120 minutes

To keep your system current with a patch, you apply a patch to the Pega Platform software in your existing deployment by using Helm, you must customize your existing `pega` Helm chart with the specific, required or settings. This includes referencing the latest Pega-provided Docker images that are available for your release.

The Pega patch application process takes at most 120 minutes total.

To complete a zero downtime patch, you must configure the following settings in your existing Pega configuration files for your Pega Platform deployment:

- Specify action.execute: upgrade-deploy to invoke the zero-downtime patch process.
- Specify the schema name or names that will be upgraded:
  - **For 8.4 and later**: specify both schema names, since the process involves migrating rules to and from each schema (dbc.connectionProperties.rulesSchema: "YOUR_RULES_SCHEMA" and jdbc.connectionProperties.dataSchema: "YOUR_DATA_SCHEMA")
  - **For 8.2 and 8.3**: specify the rules schema since the process only involves migrating rules to and from the existing rule schema (dbc.connectionProperties.rulesSchema: "YOUR_RULES_SCHEMA"); leave the existing "YOUR_RULES_SCHEMA" value (do not leave it blank).
- Ensure one of the following:
  - You pushed the images for your patch to the same repository that you used for your installation repository and the credentials for your repository account are the same as those in your `pega` Helm chart.  
  - You pushed the images for your patch to a new repository and you update the parameters with your credentials for this new repository account in your `pega` Helm chart.

- Update the tagging details, including the version and date of your latest Pega-provided `platform/pega` Docker image, that you downloaded to support your patch.

- For existing AKS and PKS deployments, provide reference details for the service host and https service port of the Kubernetes API server (this is not required for installations). For example, the Kubernetes master is running at `https://<service_host>:<https_service_port>`. For EKS and GKE deployments, you leave the existing text values (do not leave them blank).

- Update the tagging details, including the version and date of your latest Pega-provided `platform/search` Docker image, that you downloaded for your patch.

- In the installer section of the Helm chart, update the following: 
  - Update the tagging details, including the version and date of your Pega-provided `platform/installer` Docker image, that you downloaded to support your patch.
  - Specify an `out-of-place` upgrade to apply a patch using the zero-downtime patch process.
  - Specify the new target rules schema name that you created in your existing database to support the patch process within the quotes.
  - For patches to 8.4 and later, specify the new target data schema name that you created in your existing database to support the patch process within the quotes. For 8.2 or 8.3 Pega software patches, you can leave this value empty, as is (do not leave it blank).

You can leave the existing customized parameters as is; the patch process will use the remaining existing settings in your deployment.

### Updating the Pega configuration files to your Helm installation on your local system

Complete the following steps.

1. Use a text editor to open the pega.yaml file and update the following parameters in the chart based on your EKS requirements:

   | Chart parameter name    | Purpose                                   | Your setting |
   |-------------------------|-------------------------------------------|--------------|
   | actions.execute: | Specify an “upgrade” deployment type. | execute: "upgrade-deploy" |
   | jdbc.connectionProperties.rulesSchema: "YOUR_RULES_SCHEMA"  | For any patch, specify the name of the existing rules schema from which the patch process migrates the existing rules structure to your new rules schema.  | rulesSchema: "YOUR_RULES_SCHEMA" |
   | jdbc.connectionProperties.dataSchema: "YOUR_DATA_SCHEMA"  | For patches to 8.4 and later, specify the name of the temporary data schema to which the patch process migrates the existing data structure from the existing data schema; if you are applying 8.2 or 8.3 Pega software patch, you can leave this value as is (do not leave it blank).  | dataSchema: "YOUR_DATA_SCHEMA"  |
   | docker.registry.url: username:  and password: | If using a new registry since you installed Pega Platform, update the host name of a registry to an object that contains the “username” and “password” values for that registry. For more information, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). You can skip this section if the registry is the same as youyr initial installation. | <ul><li>url: “<https://index.docker.io/v1/>” </li><li>username: "\<DockerHub account username\>"</li><li> password: "\< DockerHub account password\>"</li></ul>    |
   | docker.pega.image:       | Update the tagging details, including the version and date of your latest Pega-provided `platform/pega` Docker image that you downloaded and pushed to your Docker registry. This image should match the version of the installer image with which you will apply your patch. | Image: "\<Registry host name:Port\>/my-pega:\<Pega Platform version>" |
   | <ul><li>upgrade.kube-apiserver. serviceHost</li><li>upgrade.kube-apiserver.httpsServicePort</li></ul>  | For existing AKS and PKS deployments, for the service host and https service port of the Kubernetes API server. For EKS and GKE deployments, leave the existing text values (do not leave them blank).| <ul><li>upgrade.kube-apiserver.serviceHost: "API_SERVICE_ADDRESS" </li><li>upgrade.kube-apiserver.httpsServicePort: "SERVICE_PORT_HTTPS"</li></ul> |
   | pegasearch.image: | Update the tagging details, including the version and date of your latest Pega-provided `platform/pega` Docker image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-search:\<Pega Platform version>"
   | installer.image: | Update the tagging details, including the version and date of your latest Pega-provided `platform/installer` Docker image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
   | installer.adminPassword: | Specify an initial administrator@pega.com password for your installation.  This will need to be changed at first login. The adminPassword value cannot start with "@".| adminPassword: "\<initial password\>"  |
   | installer.upgrade.upgradeType   | Specify an out-of-place upgrade to apply a patch using the zero-downtime patch process. | upgradeType: "out-of-place"  |
   | installer.upgrade.targetRulesSchema   | Specify the new target rules schema name that you created in your existing database to support the patch process within the quotes. | targetRulesSchema: ""  |
   | installer.upgrade.targetDataSchema   | For patches to 8.4 and later, specify the new target data schema name that you created in your existing database to support the patch process within the quotes. For 8.2 or 8.3 Pega software patches, you can leave this value empty, as is (do not leave it blank). | targetDataSchema: ""   |

2. Save the file.

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