Preparing your local Windows 10 system – 45 minutes
==========================================================================

In order to deploy Pega Platform using a local Windows 10 system on which you can run commands with Administrator privileges, you must prepare your system with required applications and configuration files you will use for your deployment. Pega recommends doing this first so you can complete the deployment without having to pause in order to obtain a Windows application or prepare a configuration file that is required to complete the deployment.

Assumptions and prerequisites
-----------------------------

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges.

- You use the open source packaging tool Chocolatey to install applications onto your Windows laptop. For more information, see [How Chocolatey Works](https://chocolatey.org/how-chocolatey-works).

- Basic familiarity with GitHub account with which you will download a Pega-managed GitHub repository containing configuration files and scripts that you use to install Pega Platform and then deploy it in the Kubernetes cluster.

Creating a local folder to access all of the configuration files
----------------------------------------------------------------

Deploying with Helm requires that you run commands from a specific folder on your local system. To ensure you stay oriented to the correct filepath, these instructions always use the reference \<local filepath\>\\\<platform>-demo folder when you must extract files to a folder or run commands from a folder.

In order to stay consistent with the instructions, you can create a folder called \<platform\>-demo on your local system at the top level of your Windows user folder. This way, you associate the generic \<local filepath\>\\\<platform\>-demo references to the folder C:\\Users\\\<your username\>\\\<platform\>-demo that is specific to your local system.

For Windows users: To create this folder, open a Windows PowerShell command prompt with Administrator privileges and enter:

`$ mkdir C:\Users\<Windows-username><platform>-demo`

Where \<platform\>-demo is:

- AKS-demo - for the AKS runbook

- EKS-demo - for the EKS runbook

- PKS-demo - for the PKS runbook

- Openshift-demo - for the Openshift runbook

Currently there is no runbook for running on the Google Kubernetes Engine (GKE) using the Windows 10 Google SDK. To set up a system for a deployment on GKE, see [Prepare a local Linux system – 45 minutes](prepping-local-system-runbook-linux.md). 

You are ready to continue preparing your local system.

Installing required applications for the deployment
---------------------------------------------------

The entire deployment requires the following applications to be used at some point during the process, so it's useful to prepare your local system with all of the files before you start your deployment:
- Helm
- kubectl
- Docker
- unzip (or something equivalent to extract files from .zip archive files.)
- az cli (only for AKS deployments)
- pks cli (only for PKS deployments)

Some of the required applications are binary files that you download from the organization's download area; other applications can be installed by using a Windows package manager application such as [Chocolatey](https://chocolatey.org/).

Note: In order to use the docker command in the runbooks, you install the Docker application directly from the Docker website. For your convenience, the instructions available on the Docker website are included in this document.

To install Chocolatey, follow these steps which are sourced from the [Install
Chocolatey](https://chocolatey.org/install) page.

1. Open a Windows PowerShell command prompt with Administrator privileges.

2. To ensure your PowerShell commands run without restrictions enter:

`$ Get-ExecutionPolicy`

3. To install Chocolatey and appropriate security scripts that it uses to
    ensure safety when you install applications using the Chocolatey
    application, enter:

`$ Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object
System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))`

The command returns messaging describing useful details and indicating a
successful installation. You are ready to use Chocolatey to install each
required application.

### Installing Helm and the kubernetes CLI commands:

Pega supports using Helm version 2.1 and later and the Kubernetes Command Line
Interface (CLI) 1.15 and later. The latest runbooks use Helm version 3.0 and
kubernetes-cli 1.17.0. It is recommended to use these versions. If you use Helm
2.x, some of the commands will differ slightly for Helm 2.x.

The default Helm version available in Chocolatey is 3.0.x; the default version
of kubernetes-cli is 1.17.x.

Enter the choco install command listed for each application into your PowerShell
command prompt as shown:

- To install [Helm](https://chocolatey.org/packages/kubernetes-helm): in the
    PowerShell command prompt, enter:

`$ choco install kubernetes-helm`

If, during the install process, you are prompted to run the script, reply with
**Yes**.

- To install [Kubernetes Command Line
    Interface](https://chocolatey.org/packages/kubernetes-cli): in the
    PowerShell command prompt, enter:

`$ choco install kubernetes-cli`

If, during the install process, you are prompted to run the script, reply with
**Yes**. The kubernetes-cli application includes the `kubectl` command.

### For AKS only: installing the Azure CLI

Use a Windows PowerShell command prompt with Administrator
privileges to install the Azure CLI by entering:

`$ Invoke-WebRequest -Uri https://aka.ms/installazurecliwindows -OutFile
.\\AzureCLI.msi; Start-Process msiexec.exe -Wait -ArgumentList '/I AzureCLI.msi
/quiet'`

The prompt returns when the installation is complete.

For details about installing the Azure CLI on a Windows system, see the article,
<https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-windows?view=azure-cli-latest>

### For PKS only: installing the PKS CLI

Install the PKS CLI binary executable file from the Pivotal support site that you will run with Windows Administrator permissions during the PKS deployment steps:

1. Use the browser of your choice to navigate to [Pivotal
    Network](https://network.pivotal.io/) and log in.

2. Open the [Pivotal Container Service
    (PKS)](https://network.pivotal.io/products/pivotal-container-service) page
    and select release version **1.5.1** in the pulldown.

3. Click **PKS CLI – v1.5.1**.

4. Click **PKS CLI - Windows** and for the EULA, click **AGREE**.

5. In the Windows explorer window, choose a folder in which to save the file,
    pks-windows-amd64-1.5.1-build.xx.exe, change the filename to pks.exe, and click
    **Save**.

The binary, executable file is now called “pks.exe” and should be moved to any
PowerShell PATH on your local computer so it can be run from the command line as
simply “pks”.

6. Navigate to the local path folder where you saved this file, for instance
    \<local filepath\>\\\<platform\>-demo.

7. Add this executable file to the PATH on your local computer. Add the \<platform\>-demo folder to your environment path by running the following command:

`$env:path += ";C:\Users\<Windows-username>\<platform>-demo"`

Advanced users may add the binary file to their path using their preferred method. These instructions were mostly sourced from the [Installing the PKS CLI](https://docs.pivotal.io/pks/1-6/installing-pks-cli.html).

### Installing Docker Desktop on Windows 10

In order to build a docker installation image in the section, [Prepare your Pega Platform installation Docker image – 15 minutes](#prepare-your-pega-platform-installation-docker-image--15-minutes), you must install the Community Edition (CE) of Docker for Windows. To do so, you must download, install, and log into Docker for Windows in order to complete the setup on your local system.

1. Download the installer (Docker for Windows Installer.exe) from [download.docker.com](https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe).

Locate the executable file that downloads to your local system.

2. Run the Docker for Windows Installer.exe installer using Administrator privileges.

3. Follow the install wizard to accept the license, authorize the installer to
    make changes on your computer, and proceed with the install.

As part of the initial installation, you must authorize Docker.app with your system password during the install process. The Docker application requires access to enable Hyper-V and containers, install networking components, and links to the Docker apps. The installation requires several reboots to install/enable these features. Docker Desktop takes 1-3 minutes to display on initial start.

4. In the **Configuration** window, do not select **Use Windows containers
    instead of Linux containers (this can be changed after the installation)**
    and click **OK**.

5. Click **Close**.

6. To complete your Docker Desktop setup, in the windows search bar, enter
    “docker” in the search field and select the Docker Desktop app in the search
    results.

![](media/a3caf914ced634e3b0dddada806082dc.png)

7. Log into the Docker Desktop application when you are prompted, using your
    DockerHub credentials.

![](media/5a215a2316542f28c2d66ba77ea8383a.png)

After you log in, the docker CLI is available from a Windows PowerShell command
prompt. These instructions were mostly sourced from the article,
<https://docs.docker.com/v17.09/docker-for-windows/install/>.

Adding the Pega configuration files to your Helm installation on your local system
----------------------------------------------------------------------------------

Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm, including a generic version of the following charts. After you add the repository to your local system, you can customize these Pega configuration files for your Pega Platform deployment:

- pega/pega - Use this chart to set customization parameters for your deployment. You will modify this chart later in the deployment tasks.

- pega/addons – Use this chart to install any supporting services and tools which your Kubernetes environment will require to support a Pega deployment: the required services, such as a load balancer or metrics server, that your deployment requires depend on your cloud environment. For instance you can specify whether you want to use a generic load-balancer or use one that is offered in your Kubernetes environment, such as in AKS or EKS. The runbooks provide instructions to deploy these supporting services once per Kubernetes environment when you install the addons chart, regardless of how many Pega Infinity instances are deployed.

To customize these files, you must download them from the repository to your local system, edit them with a text editor, and then save them to your local system using the same filename. In this set of tasks, you will focus on the pega/addons.yaml file; in the environment-specific runbook that you are using in the section, **Update the Helm chart values**, you will update the pega.yaml file.

1. To add the Pega repository to your Helm installation, enter:

`$ helm repo add pega https://dl.bintray.com/pegasystems/pega-helm-charts`

2. To verify the new repository, you can search it by entering:

```yaml
  $ helm search repo pega
  NAME        CHART VERSION   APP VERSION     DESCRIPTION
  pega/pega   1.2.0                           Pega installation on kubernetes
  pega/addons 1.2.0           1.0             A Helm chart for Kubernetes
```

These two charts in this /charts/pega folder of the pega-helm-charts repository, pega and addons, require customization for your deployment of Pega Platform.

Updating the Pega addons Helm chart
-----------------------------------

Update this Helm chart in order to enable the Traefik load balancer and disable
the metrics-server for deployments to the following platforms:

- AKS

- PKS

- GKE

If you are deploying to a different platform you can skip this section.

1. To download pega/addons Helm chart to the \<local filepath\>\\\<platform\>-demo, enter:

`$ helm inspect values pega/addons > addons.yaml`

2. Open the addons.yaml file from this folder in a text editor.

3. In the traefik configuration area, ensure the following two settings are
    configured to use Traefik for your deployment load-balancer:

```yaml
  traefik:
  enabled: **true**

  # Set any additional Traefik parameters. These values will be used by Traefik's Helm chart.
  # See https://github.com/Helm/charts/blob/master/stable/traefik/values.yaml
  # Set traefik.serviceType to "LoadBalancer" on gke, PKS, and pks
  serviceType: **LoadBalancer**

  Note: Do not enclose the text in quotes.
```

4. For GKE or PKS deployments, you must ensure that the Pega metrics server is disabled in the metrics-server section of this *addon* values.yaml file, since PKS deployments use the PKS metrics server

```yaml
metrics-server:

# Set this to true to install metrics-server. Follow below guidelines specific to each provider,
# open-source Kubernetes, Openshift & EKS - mandatory to set this to true if any tier as hpa.enabled is true
# GKE or PKS - set this to false since metrics-server is installed in the cluster by default.

enabled: **false**
```

5. Save the file.

Add any known, customized settings for Pega to your deployment
--------------------------------------------------------------

The Pega deployment model supports advanced configurations to fit most existing
clients' needs. If you are a Pega client and have known, required customizations
for your deployment and you already use the following files to add your known
customizations, you can copy those configurations into the configuration files
Pega added for this purpose in the [pega-helm-charts](https://github.com/pegasystems/pega-helm-charts) repository folder, pega-helm-charts/charts/pega/config/deploy:

- context.xml: add additional required data sources

- prlog4j2.xml: modify your logging configuration, if required

- prconfig.xml: adjust the standard Pega Platform configuration with known,
    required settings

Make these changes before you begin deploying Pega Platform
using Helm charts.

Downloading a Pega Platform distribution to your local system
-------------------------------------------------------------

These instructions require the Pega Platform distribution image to install the Pega Platform onto your database. To obtain a copy, you must download an image from Pega. For detailed instructions, see [Pega Digital Software Delivery User Guide](https://community.pega.com/knowledgebase/documents/pega-digital-software-delivery-user-guide).

### Requesting access to a Pega Platform distribution

1. In the browser of your choice, navigate to the Pega [Digital Software Delivery](https://community1.pega.com/digital-delivery) site.

2. Log into the [Pega Community](https://community.pega.com/knowledgebase/articles/pega-cloud/pega-cloud-services-patch-process-releases-83x-and-later)
    site with the credentials your Pega representative provided.

3. In the **Download and Upgrade Licensed Software** area, click **New
    request**.

4. In the right side of the page, click **Continue**.

If you have multiple associations with the Pega Community, the page requests you to select the organization with which you want to affiliate this request and then click **Continue**. You will receive an email with a link to your software using an email address that is associated with the organization you select on this screen.

5. In the **You're viewing products available** page, enter **Pega Platform** in the **Search**, which filters the list of products in the page.

The **Pega Platform** card should appear near the top of the card list, below
the list of all of the **Language packs for Pega Platform.**

6. In the Pega Platform card, your mouse arrow changes into a shopping cart icon, which you use to select **Pega Platform**.

The icon changes to a green check and a new cart item appears in the top right of the product list.

![](media/029c6531bd52109598047a2ee6966657.png)

7. Click **Continue**.

8. In the cart review page, in the **Pega Platform** area, select the version
    of Pega Platform for your deployment.

![](media/386d4eb20a4e2be6b767bc522cbdda91.png)

9. After your selection and review are complete, click **Finish.**

10. When the order is processed, a confirmation screen displays with details about your order.

- You receive an email with a link to the requested Pega Platform software within a few minutes. The email address used is associated with the organization you selected in this section.

![](media/748ea91e3ff43cf4544ce2f4638e86bf.png)

11. When satisfied with the order, click **Close**.

### Downloading Pega Platform to your local system

To download your Pega Platform image:

1. Open the email you received. It will look similar to the image shown.

![](media/98b1055e0e63487db7bbb2c90c9ea40c.png)

2. Click **Download now**.

3. The **Pega Licensed Software Downloads** page opens.
 
4. Under the  **My Downloads** area, click **Download software**.

Your secure **Inbox** of requested Pega software products opens. Your request for a version of Pega Platform software is listed at the top of the inbox table.

5. In the **Subject** column, click the link to your requested Pega Platform software.

The Package details window opens in the Package tab, which shows details about the Pega Platform software distribution package that you requested.

6. In the **Files:** area of the window, ensure that version of the Pega distribution image is correct.

If it is not the right version number, you must complete a new request.

7. To download the file, select the Pega distribution image checkbox and click
    **Download**.

8. In the **Save as** window, choose the \<local filepath\>\\\<platform\>-demo folder to which you save the Pega Platform distribution zip file.

9. In Windows PowerShell, change folders to the \<localfilepath\>\\\<platform\>-demo folder, where you saved the Pega Platform distribution zip and extract your files to create a new distribution image folder on your local system:

`$ Expand-Archive .\<pega-distribution-image>.zip`

![C:\\Users\\aciut\\AppData\\Local\\Temp\\SNAGHTML2318748.PNG](media/1edabbc855175f2bdca62f8a529d8c88.png)

After you expand the archive, the files in the Pega Platform distribution image are available to use in preparing your Pega Platform installation Docker image.

Prepare your Pega Platform installation Docker image – 15 minutes
-----------------------------------------------------------------

As stated previously, you must have a [DockerHub](https://hub.docker.com/) account and log into it in order to see the [Pega-installer-ready Docker image](https://hub.docker.com/r/pegasystems/pega-installer-ready). You also need the docker cli and docker-desktop installed on your system before you begin this procedure. The Pega-provided Docker image, pega-installer-ready, includes some components of a full installation image that you can use to install or upgrade the Pega Platform database. While it is built on top of a JDK, it does not contain the contents of the Pega distribution kit which are essential for installing or upgrading Pega Platform.

Pega provides this image as the primary content of the final Docker image you will use to install or upgrade Pega Platform. This section describes how you can use this Docker image in combination with a Dockerfile and the Pega Platform distribution image that you have made available on your local system. The procedure assumes you’ve downloaded the software in [Downloading Pega Platform to your local system](#downloading-a-pega-platform-distribution-to-your-local-system) and installed the required components on your local system listed in [Install required applications for the deployment](#creating-a-local-folder-to-access-all-of-the-configuration-files).

Follow these steps to create a Docker image you can use to install or upgrade
Pega Platform.

1. From your PowerShell running with Administrator privileges, ensure you are
    logged into your DockerHub account:

`$ docker login -u \<username\> -p \<username-password\>`

For details about logging into Docker from a secure password file using the `--password-stdin` option, see <https://docs.docker.com/engine/reference/commandline/login/>. From Windows, you can also ensure you are logged by using the Docker Desktop client running on your system.

2. Change your directory to the top folder of your Pega distribution, \<pega-distribution-image\>.

`$ cd .\<pega-distribution-image>`

3. Create a text file with the text editor of your choice in the \<local filepath\>\\\<platform\>-demo\\\<pega-distribution-image\> folder where you extracted the Pega distribution on your local system.

From this folder, you can list the folder content and see folders for Pega archives, Images, rules, and scripts.

![](media/152260ae774fe07d717f1b31b5560f25.png)

4. Copy the following lines of instruction into the new text file:

```yaml
FROM pegasystems/pega-installer-ready
COPY scripts /opt/pega/kit/scripts
COPY archives /opt/pega/kit/archives
COPY rules /opt/pega/kit/rules
RUN chmod -R 777 /opt/pega
```

These instructions direct a docker build function to use the Pega public Docker image, pega-install-ready, and these three folders from the Pega distribution image in order to build your Pega Platform installation image.

5. Save the text-only file with the filename, "dockerfile", without an extension, in the \<local filepath\>\\\<platform\>-demo\\\<pega-distribution-image\> folder where you extracted the Pega distribution on your local system.

6. From your PowerShell running with Administrator privileges command prompt, in your current directory, build your pega install docker image by entering:

`$ docker build -t pega-installer .`

This command uses your dockerfile to build a full Docker image with the name “pega-installer” and gives it the tag, “latest”. You can use the docker command to see that your new image exists in your Docker image inventory.

7. Tag the local version of your new image, pega-installer, with your DockerHub ID:

`$ docker tag pega-installer <your-dockerhub-ID>/pega-installer`

8. Create a private repository on your [DockerHub](https://hub.docker.com/) account that is tagged as a private repository.

9. From your default login page, click the Repositories link (at the top of the page).

10. In the Repositories view, click **Create Repository +**.

11. Enter a Name that matches the docker image you just built.

12. Provide a brief Description, that will help you remember the version of Pega
    with which you built this image, or any other useful information.

13. In the Visibility area, select the **Private**.

You should not maintain this image with Pega proprietary software as a viewable **Public** image.

14. Click **Create**.

Free DockerHub accounts support the use of a single private repository, so you may have to delete an existing private repository in order to create a new one for your Pega docker installation image.

15.  From your PowerShell running with Administrator privileges, use the docker command to push the new image to your new private repository:

`$ docker push <your-dockerhub-ID>/pega-installer`

After the command completes you will see your new image in your private
repository, similar to the image below.

![](media/9fd09158a821f828a93d6ab7c74e278a.png)
