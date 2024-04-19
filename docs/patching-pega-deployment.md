# Patching Pega Platform in your deployment

After you deploy Pega Platform™ on your kubernetes environment, the Pega-provided Docker images support applying a zero-downtime patch to your Pega software. The following procedures are written for any level of user, from a system administrator to a development engineer who wants to use helm charts and Pega Docker images to patch the Pega software have deployed in any supported kubernetes environment.

Useful links to Pega software patching information:

- [Pega software maintenance and extended support policy](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-software-maintenance-and-extended-support-policy)
- [Pega Infinity patch calendar](https://community.pega.com/knowledgebase/articles/keeping-current-pega/pega-infinity-patch-calendar)
- [Pega Infinity patch frequently asked questions](https://community.pega.com/knowledgebase/articles/keeping-current-pega/85/pega-infinity-patch-frequently-asked-questions)

## Kubernetes-based patching process overview

Pega supports client-managed cloud clients applying patches for releases 8.4 and later using a zero-downtime patch process to apply the latest cumulative bundle of bug and security fixes since the last minor release. For the latest Pega Community articles, see [About client managed cloud](https://community.pega.com/knowledgebase/articles/client-managed-cloud/85/about-client-managed-cloud).

The Pega zero-downtime patch process uses the zero-downtime patch process so you and your customers can continue working in your application while you patch your system. Pega zero-downtime patch scripts use a temporary data schema and the patch migration script moves the rules between the appropriate schema and then performs the required rolling reboot of your deployment cluster. For a detailed overview of the process, see [Applying a patch without downtime](https://community.pega.com/knowledgebase/articles/keeping-current-pega/86/applying-patch-without-downtime).

## Client-required steps
Client-managed cloud clients use the same Pega Kubernetes tools and Helm charts in the same Pega repository that you used to install Pega Platform in a supported Kubernetes environment. The client-managed cloud patch process includes the following tasks:

1. Prepare your Docker repository by downloading the latest three Pega Platform patch release images (platform/installer, platform/pega, and platform/search) in your release stream and pushing them into your preferred Docker image repository For step-by-step details, see [Downloading and managing Pega Platform docker images (linux)](prepping-local-system-runbook-linux.md#downloading-and-managing-pega-platform-docker-images) or [Downloading and managing Pega Platform docker images (windows)](prepping-local-system-runbook-windows.md#downloading-and-managing-pega-platform-docker-images).

2. Edit the pega Helm chart by editing parameters to specify "upgrade-deploy" your software with the software contained in your provided patch image. - [Applying a zero-downtime Pega Platform patch using Helm charts - 120-minutes](#applying-a-zero-downtime-pega-platform-patch-using-helm-charts--120-minutes).

3. Apply the patch by using the `helm upgrade release --namespace mypega` command as directed in the deployment section - [Patching your Pega Platform deployment using the command line](#patching-your-pega-platform-deployment-using-the-command-line).

## Assumptions and prerequisites

The process to patch your deployment assumes:

- Your Kubernetes environment has not changed, and you are using the same Pega charts with which you originally deployed.

- Your original deployment used Pega Platform 8.3.0 or later.

## Applying a zero-downtime Pega Platform patch using Helm charts – 120 minutes

To keep your system current with a patch, you apply a patch to the Pega Platform software in your existing deployment by using Helm, you must customize your existing `pega` Helm chart with the specific, required or settings. This includes referencing the latest Pega-provided Docker images that are available for your release.

The Pega patch application process takes at most 120 minutes total.

To complete a zero downtime patch, you must configure the following settings in your existing Pega configuration files for your Pega Platform deployment:

- Specify action.execute: upgrade-deploy to invoke the zero-downtime patch process.
- Specify the schema name or names that will be upgraded:
  - **For 8.4 and later**: specify both schema names, since the process involves migrating rules to and from each schema (jdbc.rulesSchema: "YOUR_RULES_SCHEMA" and jdbc.dataSchema: "YOUR_DATA_SCHEMA").
  - **For 8.3**: specify the rules schema since the process only involves migrating rules to and from the existing rule schema (jdbc.rulesSchema: "YOUR_RULES_SCHEMA"); leave the existing "YOUR_RULES_SCHEMA" value (do not leave blank text).
- Ensure one of the following:
  - You pushed the images for your patch to the same repository that you used for your installation repository and the credentials for your repository account are the same as those in your `pega` Helm chart.  
  - You pushed the images for your patch to a new repository and you update the parameters with your credentials for this new repository account in your `pega` Helm chart.

- Update the tagging details, including the version and date of your latest Pega-provided `platform/pega` Docker image, that you downloaded to support your patch.

- For existing AKS and PKS deployments, provide reference details for the service host and https service port of the Kubernetes API server (this is not required for installations). For example, the Kubernetes master is running at `https://<service_host>:<https_service_port>`. For EKS and GKE deployments, you leave the existing text values (do not leave them blank).

- Update the tagging details, including the version and date of your latest Pega-provided `platform/search` Docker image, that you downloaded for your patch.

- In the installer section of the Helm chart, update the following:
  - Update the tagging details, including the version and date of your Pega-provided `platform/installer` Docker image, that you downloaded to support your patch.
  - Specify an `zero-downtime` upgrade to apply a patch using the zero-downtime patch process.
  - **For patches to 8.4 and later**, specify the new target new rules and temporary data schema names that the process creates in your existing database to support the patch process within the quotes. **For 8.3 Pega software patches**, you can leave this value empty, as is (do not leave blank text).

You can leave the existing customized parameters as is; the patch process will use the remaining existing settings in your deployment.

### Updating the Pega configuration files to your Helm installation on your local system

Complete the following steps.

1. Use a text editor to open the pega.yaml file and update the following parameters in the chart based on your Kubernetes environment requirements:

   | Chart parameter name    | Purpose                                   | Your setting |
   |-------------------------|-------------------------------------------|--------------|
   | actions.execute: | To apply a patch using the zero-downtime patch process, specify an “upgrade-deploy” deployment type. | execute: "upgrade-deploy" |
   | jdbc.rulesSchema: "YOUR_RULES_SCHEMA"  | For any patch, specify the name of the existing rules schema from which the patch process migrates the existing rules structure to your new rules schema.  | rulesSchema: "YOUR_RULES_SCHEMA" |
   | jdbc.dataSchema: "YOUR_DATA_SCHEMA"  | Specify the name of the existing data schema to which the patch process migrates the existing data structure from the existing data schema  | dataSchema: "YOUR_DATA_SCHEMA"  |
   | docker.registry.url: username:  and password: | If using a new registry since you installed Pega Platform, update the host name of a registry to an object that contains the “username” and “password” values for that registry. For more information, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). You can skip this section if the registry is the same as your initial installation. | <ul><li>url: “<https://index.docker.io/v1/>” </li><li>username: "\<DockerHub account username\>"</li><li> password: "\< DockerHub account password\>"</li></ul>    |
   | docker.pega.image:       | Update the tagging details, including the version and date of your latest Pega-provided `platform/pega` Docker image that you downloaded and pushed to your Docker registry. This image should match the version of the installer image with which you will apply your patch. | Image: "\<Registry host name:Port\>/my-pega:\<Pega Platform version>" |
   | <ul><li>upgrade.kube-apiserver. serviceHost</li><li>upgrade.kube-apiserver.httpsServicePort</li></ul>  | For existing AKS and PKS deployments, for the service host and https service port of the Kubernetes API server. For EKS and GKE deployments, leave the existing text values (do not leave them blank).| <ul><li>upgrade.kube-apiserver.serviceHost: "API_SERVICE_ADDRESS" </li><li>upgrade.kube-apiserver.httpsServicePort: "SERVICE_PORT_HTTPS"</li></ul> |
   | pegasearch.image: | Update the tagging details, including the version and date of your latest Pega-provided `platform/pega` Docker image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-search:\<Pega Platform version>"
   | installer.image: | Update the tagging details, including the version and date of your latest Pega-provided `platform/installer` Docker image that you downloaded and pushed to your Docker registry. | Image: "\<Registry host name:Port>/my-pega-installer:\<Pega Platform version>" |
   | installer.adminPassword: | Specify an initial administrator@pega.com password for your installation.  This will need to be changed at first login. The adminPassword value cannot start with "@".| adminPassword: "\<initial password\>"  |
   | installer.upgrade.upgradeType   | Specify an zero-downtime upgrade to apply a patch using the zero-downtime patch process. | upgradeType: "zero-downtime"  |
   | installer.upgrade.targetRulesSchema   | Specify a new rules schema name that the process creates in your existing database to support the patch process within the quotes. | targetRulesSchema: ""  |
   | installer.upgrade.targetDataSchema   | For patches to 8.4 and later, specify a new target data schema name that the process creates in your existing database to support the patch process within the quotes. For 8.3 Pega software patches, you can leave this value empty, as is (do not leave blank text). | targetDataSchema: "" |
2. Save the file.

### Patching your Pega Platform deployment using the command line

In this document, you specify that the Helm chart always “deploys” by using the setting, actions.execute: "upgrade-deploy argument". After you have your customizations saved in your pega Helm chart, you are ready to apply the patch.

1. Do one of the following:

   - Open Windows PowerShell running as Administrator on your local system and change the location to the top folder of your `\<platform\>-demo` folder that you created in [Preparing your local Windows 10 system](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md).

   `$ cd <local filepath>\<platform>-demo`

   - Open a Linux bash shell and change the location to the top folder of your `\<platform\>-demo` directory that you created in [Preparing your local Linux system](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md).

   `$ cd /home/<local filepath>/<platform>-demo`

2. Patch Pega Platform by upgrading using your updated `pega` Helm chart.

   ```yaml
   helm upgrade mypega-<platform>-demo pega/pega --namespace mypega-<platform>-demo --values pega.yaml
   ```

   A successful upgrade immediately returns details that shows progress for your `mypega-<platform>-demo` deployment.

3. Refresh the Kubernetes dashboard that you opened in the previous section. If you closed the dashboard, start the proxy server for the Kubernetes dashboard and then relaunch the web browser.

4. In the dashboard, in **Namespace** select the `mypega-<platform>-demo` view and then click on the **Pods** view. Initially, you can some pods have a red status, which means they are initializing:

    You can follow the progress of your patch using the dashboard. Initially, while the resources make requests to complete the configuration, you will see red warnings while the configuration is finishing, which is expected behavior.

5. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-zdt-upgrade** job, and click the logs icon on the right side of that row.

   After you open the logs view, you can click the icon for automatic refresh to see current updates to the upgrade (patch) log.

6. To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the `mypega-<platform>-demo` namespace pods.

   A successful deployment does not show errors across the various workloads. The `mypega-<platform>-demo` Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment has 100% complete **Workloads**.

   It takes a little over an hour for the patch process to patch the applicable rules and then perform a rolling reboot of your nodes.
