Preparing your local Windows 10 system – 45 minutes
==========================================================================

To deploy Pega Platform using a local Windows 10 system on which you can run commands with administrator privileges, prepare your system with required applications and configuration files that you need for your deployment. By preparing your system, you can complete the deployment without having to pause to install a Windows application or prepare a configuration file.

Before you begin
----------------

To prepare a Wndows 10 system, it is assumed:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with administrator privileges.

- You use the open source packaging tool Chocolatey to install applications onto your Windows laptop. For more information, see [How Chocolatey Works](https://chocolatey.org/how-chocolatey-works).

- You have basic familiarity with a GitHub account with which you download a Pega-managed GitHub repository that contains configuration files and scripts that support installing Pega Platform and deploying it in the Kubernetes cluster.

Creating a local folder to access all of the configuration files
----------------------------------------------------------------

Deploying with Helm requires you to run commands from a specific folder on your local system. To ensure you use the correct filepath, these instructions always use the reference `<local filepath>\<platform>-demo` folder when you extract files to a folder or run commands from a folder.

To stay consistent with the instructions, you can create a folder named `<platform>-demo` on your local system at the top level of your Windows user folder. This way, you associate the generic `<local filepath>\<platform>-demo` references to the folder `C:\\Users\\\<your username\>\\\<platform\>-demo` that is specific to your local system.

To create this folder, open a Windows PowerShell command prompt with administrator privileges and enter:

    `$ mkdir C:\Users\<Windows-username><platform>-demo`

Where `<platform>-demo` is:

- For the AKS runbook, `AKS-demo`

- For the EKS runbook, `EKS-demo`

- For the PKS runbook, `PKS-demo`

- For the Openshift runbook, `Openshift-demo`

Currently there is no runbook for running on the Google Kubernetes Engine (GKE) using the Windows 10 Google SDK. To set up a system for a deployment on GKE, see [Prepare a local Linux system – 45 minutes](prepping-local-system-runbook-linux.md). 

You are ready to continue preparing your local system.

Installing required applications for the deployment
---------------------------------------------------

The entire deployment requires the following applications to be used during the configuration process; therefore, you should prepare your local system with all of the applications before you start your deployment:
- Helm
- kubectl
- Docker
- unzip (or an equivalent to extract files from .zip archive files.)
- az cli (only for AKS deployments)
- AWS IAM Authenticator for Kubernetes (only for EKS deployments)
- eksctl (only for EKS deployments)
- pks cli (only for PKS deployments)

Some of the required applications are binary files that you download from the organization's download area; you can install other applications by using a Windows package manager application such as [Chocolatey](https://chocolatey.org/).

Note: To use the docker command in the runbooks, you install the Docker application directly from the Docker website. For your convenience, the instructions available on the Docker website are included in this document.

To install Chocolatey, follow these steps which are sourced from the [Install
Chocolatey](https://chocolatey.org/install) page.

1. Open a Windows PowerShell command prompt with administrator privileges.

2. To ensure your PowerShell commands run without restrictions enter:

    `$ Get-ExecutionPolicy`

3. To install Chocolatey and appropriate security scripts that it uses to
    ensure safety when you install applications using the Chocolatey
    application, enter:

    `$ Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-ObjectSystem.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))`

    The command returns messaging describing useful details and indicating a successful installation. You are ready to use Chocolatey to install each required application.

### Installing Helm and the kubernetes CLI commands:

Pega supports using Helm version 2.1 and later and the Kubernetes Command Line Interface (CLI) 1.15 and later. The latest runbooks use Helm version 3.0 and
kubernetes-cli 1.17.0. It is recommended to use these versions. If you use Helm 2.x, some of the commands differ slightly for Helm 2.x.

The default Helm version available in Chocolatey is 3.0.x; the version of kubernetes-cli used in the runbooks is version 1.17.x.

Enter the choco install command listed for each application into your PowerShell command prompt as shown:

- To install [Helm](https://chocolatey.org/packages/kubernetes-helm): in the PowerShell command prompt, enter:

    `$ choco install kubernetes-helm`

If a prompt to run the script appears, enter **Yes**.

For additional information, see [Helm documentation](https://helm.sh/docs/).

- To install [Kubernetes Command Line Interface](https://chocolatey.org/packages/kubernetes-cli): in the PowerShell command prompt, enter:

    `$ choco install kubernetes-cli`

If a prompt to run the script appears, enter **Yes**. 

The kubernetes-cli application includes the `kubectl` command.

For details about installing on Windows 10 and supported installation methods, see [Install kubectl on Windows](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-windows).

For Amazon EKS only: - To install [eksctl](https://chocolatey.org/packages/eksctl): in the PowerShell command prompt, enter:

`$ choco install eksctl`

If a prompt to run the script appears, enter **Yes**.

For Amazon EKS only: - To install [AWS IAM Authenticator for Kubernetes](https://chocolatey.org/packages/aws-iam-authenticator#files): in the PowerShell command prompt, enter:

`$ choco install aws-iam-authenticator`

If a prompt to run the script appears, enter **Yes**. 

Confirm the AWS CLI that comes withe the `AWS IAM Authenticator for Kubernetes` installation.

```yaml
$ aws --version
aws-cli/1.16.272 Python/3.6.0 Windows/10 botocore/1.13.8
```

For additional AWS tool information, see [Install the AWS CLI version 1 on Windows](https://docs.aws.amazon.com/cli/latest/userguide/install-windows.html).

### For AKS only: installing the Azure CLI

To install the Azure CLI using a Windows PowerShell command prompt with administrator privileges, enter:

`$ Invoke-WebRequest -Uri https://aka.ms/installazurecliwindows -OutFile .\\AzureCLI.msi; Start-Process msiexec.exe -Wait -ArgumentList '/I AzureCLI.msi /quiet'`

The prompt returns when the installation is completed.

For details about installing the Azure CLI on a Windows system, see [Install Azure CLI on Windows](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-windows?view=azure-cli-latest).

### For PKS only: installing the PKS CLI

Install the PKS CLI binary executable file from the Pivotal support site that you run with Windows administrator permissions during the PKS deployment steps:

1. Use the browser of your choice to navigate to [Pivotal Network](https://network.pivotal.io/) and log in.

2. Open the [Pivotal Container Service (PKS)](https://network.pivotal.io/products/pivotal-container-service) page and release version **1.5.1**.

3. Click **PKS CLI – v1.5.1**.

4. Click **PKS CLI - Windows**.

5. In the EULA tab, click **AGREE**.

6. In the Windows explorer window, choose the `<platform>-demo` folder to save the `pks-windows-amd64-1.5.1-build.xx.exe` file, change the filename to `pks.exe`, and then click **Save**.

7. Add this executable file to the PATH on your local computer so you can run `pks` from the command-line by entering:

    `$env:path += ";C:\Users\<Windows-username>\<platform>-demo"`

    Advanced users may add the binary file to their path using their preferred method. These instructions were mostly sourced from the [Installing the PKS CLI](https://docs.pivotal.io/pks/1-6/installing-pks-cli.html).

### Installing Docker Desktop on Windows 10

In order to build a docker installation image in the section, [Prepare your Pega Platform installation Docker image – 15 minutes](#prepare-your-pega-platform-installation-docker-image--15-minutes), you must install the Community Edition (CE) of Docker for Windows. To do so, you must download, install, and log into Docker for Windows in order to complete the setup on your local system.

1. Download the installer (Docker for Windows Installer.exe) from [download.docker.com](https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe).

Locate the executable file that downloads to your local system.

2. Run the Docker for Windows Installer.exe installer using administrator privileges.

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

Creating SSL certificates for HTTPS load balancer on a Kubernetes ingress
-----------------------------------

When you configure an HTTPS load balancer through a Kubernetes ingress, you can configure the load balancer to present authentication certificates to the client. After you generate your certificates and corresponding private keys using any SSL Certificate creation tool, you pass the certificates and private keys to your Kubernetes ingress using the method supported for the specific environment. For Pega support details, see the [Addons Helm chart](https://github.com/pegasystems/pega-helm-charts/tree/master/charts/addons). 

Use the generation method best suited to your company standards. The following steps create one or more certificates and private keys for a hostname using a manually verified domain.

1. Navigate to the certificate generator, such as [SSL For Free](https://www.sslforfree.com/).

2. Enter your hostname for which you are creating a certificate and click **Create FreeSSL Certificate**.

3. To validate the certificate manually, click **Manual Verification (DNS)**.

   The **Manually Verify Domain (DNS)** option displays.

4. Click **Manually Verify Domain**.

   The **Update DNS Records** screen displays a name/host record, an associated verification value.

5. Using the DNS lookup management system of your choice, associate the generated name/host record with the verification value. For instance:

```java
Name          :   _acme-challenge.demotest.dev.pega.io
Type          :   TXT
TTL (Seconds) :   1
Value         :  "ezEiD0Lkvvzlgfaqdohe3ZcX7s4vVF6hHlBBKI3sL38"
```

With the hostname associated with the verification value, wait several minutes to ensure the configuration is established.

6. Click the link for name/host record.

   A new page displays the status of the connection.

7. When successful, you will see that the link returned the correct value and you can generate the SSL certificate files.

   If it is not successful, you may have to wait longer for the DNS lookup to correctly associate the verification value with the hostname; if it continues to not work, you may need to update your DNS settings.

8. In the main page of **SSL for Free**, click **Download All SSL Certificate Files** and save the certificate file and private key file in the \<platform\>-demo folder you have already created on your local system.

You will manage these certificates in the environment to which you deploy Pega Platform. For environment-specific details, see the **Deploying Pega Platform using the command line** section in the runbook for that environment.

Downloading a Pega Platform installer docker image
-----------------------------------------------

To download the docker images Pega Platform distribution image to install the Pega Platform onto your database. 

.

## Requesting access to a Pega Platform distribution

1. In the browser of your choice, navigate to the Pega [Digital Software Delivery](https://community.pega.com/digital-delivery) site.

2. Log into the [Pega Community](https://community.pega.com/knowledgebase/articles/pega-cloud/pega-cloud-services-patch-process-releases-83x-and-later)
    site with the credentials your Pega representative provided.

3. In the **Download and Upgrade Licensed Software** area, click **New
    request**.

4. In the right side of the page, click **Continue**.

![](media/029c6531bd52109598047a2ee6966657.png)

Open the email you received. It will look similar to the image shown.

![](media/dockerimage-access-email.png)






Building a Pega Platform installer docker image
-----------------------------------------------

These instructions require the Pega Platform distribution image to install the Pega Platform onto your database. To obtain a copy, you must download an image from Pega. For detailed instructions, see [Pega Digital Software Delivery User Guide](https://community.pega.com/knowledgebase/documents/pega-digital-software-delivery-user-guide).

## Requesting access to a Pega Platform distribution

1. In the browser of your choice, navigate to the Pega [Digital Software Delivery](https://community.pega.com/digital-delivery) site.

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

## Downloading Pega Platform to your local system

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

## Prepare your Pega Platform installation Docker image

As stated previously, you must have a [DockerHub](https://hub.docker.com/) account and log into it in order to see the [Pega-installer-ready Docker image](https://hub.docker.com/r/pegasystems/pega-installer-ready). You also need the docker cli and docker-desktop installed on your system before you begin this procedure. The Pega-provided Docker image, pega-installer-ready, includes some components of a full installation image that you can use to install or upgrade the Pega Platform database. While it is built on top of a JDK, it does not contain the contents of the Pega distribution kit which are essential for installing or upgrading Pega Platform.

Pega provides this image as the primary content of the final Docker image you will use to install or upgrade Pega Platform. This section describes how you can use this Docker image in combination with a Dockerfile and the Pega Platform distribution image that you have made available on your local system. The procedure assumes you’ve downloaded the software in [Downloading Pega Platform to your local system](#downloading-a-pega-platform-distribution-to-your-local-system) and installed the required components on your local system listed in [Install required applications for the deployment](#creating-a-local-folder-to-access-all-of-the-configuration-files).

Follow these steps to create a Docker image you can use to install or upgrade
Pega Platform.

1. From your PowerShell running with administrator privileges, ensure you are logged into your DockerHub account:

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

6. From your PowerShell running with administrator privileges command prompt, in your current directory, build your pega install docker image by entering:

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

15.  From your PowerShell running with administrator privileges, use the docker command to push the new image to your new private repository:

    `$ docker push <your-dockerhub-ID>/pega-installer`

After the command completes you will see your new image in your private
repository, similar to the image below.

![](media/9fd09158a821f828a93d6ab7c74e278a.png)
