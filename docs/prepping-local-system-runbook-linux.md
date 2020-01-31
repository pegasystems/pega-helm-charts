Preparing your local Linux system – 45 minutes
=====================================================================

In order to deploy Pega Platform using a local Linux system on which you can run commands with
Administrator privileges, you must prepare your system with required applications and configuration files you will use for your deployment. Pega recommends doing this first so you can complete the deployment without having to pause in order to obtain a Linux executable file or prepare a configuration file that is required to complete the deployment.

Assumptions and prerequisites
-----------------------------

This guide assumes:

- You have a basic familiarity with running commands from a Linux command prompt with and without root privileges.

- You use a packaging tool to install application packages. For demonstration purposes, this document refers to the Ubuntu Advanced Packaging Tool (apt) commands that is available for use on Ubuntu Linux distributions.

- Basic familiarity with GitHub account with which you will download a Pega-managed GitHub repository containing configuration files and scripts that you use to install Pega Platform and then deploy it in the Kubernetes cluster.

Creating a local folder to access all of the configuration files
----------------------------------------------------------------

Deploying with Helm requires that you run commands from a specific folder on your local system. To ensure you stay oriented to the correct filepath, these instructions always use the reference \<local filepath\>/\<platform\>-demo directory when you must extract files to a folder or run commands from a folder.

In order to stay consistent with the instructions, you can create a directory called \<platform\>-demo on your local system in your /home directory. This way, you associate the generic \<local filepath\>/\<platform\>-demo references to the folder /home/\<*linux-username*\>/\<platform\>-demo that is specific to your local
system.

To create this folder, open a Linux command prompt and enter:

`$ mkdir /home/<linux-username>/<platform>-demo`

Where \<platform\>-demo is:

- AKS-demo - for the AKS runbook

- EKS-demo - for the EKS runbook

- GKE-demo - for the GKE runbook

- PKS-demo - for the PKS runbook

- Openshift-demo - for the Openshift runbook

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
- Google Cloud SDK and gcloud (only for GKE deployments)

Some of the required applications are binary files that you download from the
organization's download area; other applications can be installed by using a Linux package manager. 

In order to use the Docker command in the runbooks, you must have the Docker application installed; however, you must install the application directly from the Docker website. For your convenience, the instructions available on the Docker website are included in this document.

### Installing Unzip

Install the Unzip utility with a packaging tool such as apt.

`$ sudo apt install unzip`

### Installing Helm

Pega supports using Helm version 2.1 and later. The latest runbooks use version 3.0 and it’s recommended to use this version. If you use Helm 2.x, some of the commands will differ slightly for Helm 2.x.

Helm supports a variety of installation methods. To learn more, see https://helm.sh/docs/intro/install/.

Helm provides a script that you can download and then run to install the latest version:

1. To download the Helm installation script from their Git repository, from
    your home directory enter:

`$ sudo curl
https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 >
get_helm.sh | bash`

2. To update the permissions of the file to you can use it for installations,
    enter:

`$ chmod 700 get_helm.sh`

3. To run the script, enter:

```yaml
$ ./get_helm.sh
Helm v3.0.1 is available.
Downloading https://get.helm.sh/helm-v3.0.1-linux-amd64.tar.gz
Preparing to install helm into /usr/local/bin
helm installed into /usr/local/bin/helm
```

4. To review your version, enter:

```yaml
$ helm version
version.BuildInfo{Version:"v3.0.1",
GitCommit:"7c22ef9ce89e0ebeb7125ba2ebf7d421f3e82ffa", GitTreeState:"clean",
GoVersion:"go1.13.4"}
```

For details about installation methods for previous Helm versions, see
<https://v2.helm.sh/docs/using_helm/#installing-helm>.

### Installing kubectl

Kubernetes supports a variety of installation methods for the kubectl command.
The organization provides a link to download the latest version of the
executable file, which will be run after you move the binary file in to your
PATH:

1. To download the latest binary from their git repository, enter:

```yaml
$ curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s
https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
```

This command downloads and parses the \`stable.txt\` in the repository, navigates to the version listed in the file, and downloads the kubectl binary file in the /bin/linux/amd64/ directory.

2. Make the kubectl binary executable:

`$ chmod +x ./kubectl`

3. Move the script to your executable PATH:

`$ sudo mv ./kubectl /usr/local/bin/kubectl`

For details about installing on Linux and other supported installation methods,
see <https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-linux>.

### For PKS only: installing the PKS CLI

Install the PKS CLI binary file from the Pivotal support site as an executable that you will run as a super user with a single command with the command curl piped into your bash:

1. Use the browser of your choice to navigate to [Pivotal Network](https://network.pivotal.io/) and log in.

2. Click [Pivotal Container Service (PKS)](https://network.pivotal.io/products/pivotal-container-service).

3. Select release version 1.6 or later from the Releases dropdown.

4. Click **PKS CLI - Linux** to display the **Release Download Files**.

5. Click **PKS CLI - Linux** and in the **Save as** window, choose the \<local filepath\>/\<platform\>-demo folder to which you save the downloaded Linux binary file.

6. Rename the downloaded binary file to "pks".

7. On the command line, run the following command to make the PKS CLI binary executable:

`$ chmod +x pks`

8. Move the binary file into your $PATH directory.

`$ sudo mv ./pks /bin`

These instructions were mostly sourced from the [Installing the PKS CLI](https://docs.pivotal.io/pks/1-6/installing-pks-cli.html).

### For AKS only: installing the Azure CLI

Install the Azure CLI as a super user with a single command with the command curl piped into your bash:

`$ curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash`

For details, see the article, <https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-apt?view=azure-cli-latest>.

### For GKE only: installing and initializing the Google Cloud SDK 

In order to use the Google Cloud SDK for your deployment, you must install the software and then initialize its use by referencing the Google Cloud project in which you will deploy Pega Platform.

To install the Google Cloud SDK:

1. In your browser, log in to your Google user account.

2. Ensure that Python 2.7 is installed.

3. Download the appropriate Google SDK distribution for your Linux version. See [Installing from versioned archives](https://cloud.google.com/sdk/docs/downloads-versioned-archives). There are also versions available for Debian and Ubuntu or Red Hat and CentOS Linux.

3. Extract the Google SDK distribution archive file to your system.

`tar zxvf [ARCHIVE_FILE] google-cloud-sdk`

4. Use the install script to add Cloud SDK tools to your path.

`./google-cloud-sdk/install.sh`

5. Restart your terminal for the changes to take effect.

6. In a Linux command prompt, initialize the Google Cloud SDK.

`gcloud init`

After the program initializes, you are prompted to log in.

7. Log in using your Google user account by entering **Y**:

`To continue, you must log in. Would you like to log in (Y/n)? Y`

You are redirected to a browser page on your system with a Google Cloud log in screen.

8. In your browser, log in to your Google Cloud user account when prompted and click **Allow** to grant permission to access Google Cloud Platform resources.

9. In your Linux command prompt, select a Cloud Platform project from the list of those where you have Owner, Editor or Viewer permissions:

```yaml
Pick cloud project to use:
 [1] [my-project-1]
 [2] [my-project-2]
 ...
 Please enter your numeric choice or text value (must exactly match list item):
```

If you only have one project, `gcloud init` selects it for you. After your selection, the command confirms that you completed the setup steps successfully:

```yaml
Your current project has been set to: [my-project-name].
...
Your Google Cloud SDK is configured and ready to use!

* Commands that require authentication will use [my-account-name] by default
* Commands will reference project `[my-project-name]` by default
Run `gcloud help config` to learn how to change individual settings

This gcloud configuration is called [default]. You can create additional configurations if you work with multiple accounts and/or projects.
Run `gcloud topic configurations` to learn more.
```

10. To list accounts whose credentials are stored on the local system:

`gcloud auth list`

11. To view information about your Cloud SDK installation and the active SDK configuration:

`gcloud info`

These instructions were sourced from the Google document, [Quickstart for Linux](https://cloud.google.com/sdk/docs/quickstart-linux), which includes additional information.

### Installing Docker

For Linux command line users, you can follow these steps to install Docker Community Edition (CE) for the first time on a new host machine. For these instructions, you need to set up the Docker repository. Afterward, you can install and update Docker from the repository.

**SET UP THE REPOSITORY**

1. Update the apt package index:

`$ sudo apt-get update`

2. Install packages to allow apt to use a repository over HTTPS:

`$ sudo apt-get install apt-transport-https ca-certificates curl
software-properties-common`

3. Add Docker’s official GPG key:

`$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add –`

4. Verify that you now have the key with the fingerprint `9DC8 5822 9FC7 DD38
    854A E2D8 8D81 803C 0EBF CD88` by searching for the last 8 characters of the
    fingerprint.

`$ sudo apt-key fingerprint 0EBFCD88`

The command should return:

```yaml
pub 4096R/0EBFCD88 2017-02-22
9DC8 5822 9FC7 DD38 854A E2D8 8D81 803C 0EBF CD88
uid Docker Release (CE deb) <docker@docker.com>
sub 4096R/F273FCD8 2017-02-22
```

5. Use the following command to set up the stable repository. You always need the stable repository, even if you want to install builds from the edge or test repositories as well. To add the edge or test repository, add the word "edge" or "test" (or both) after the word stable in the commands below.

`$ sudo add-apt-repository "deb [arch=amd64]
https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`

Note: Starting with Docker 17.06, stable releases are also pushed to the edge
and test repositories.


**INSTALL DOCKER CE**

1. Update the apt package index.

`$ sudo apt-get update`

2. Install the latest version of Docker CE or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

`$ sudo apt-get install docker-ce docker-ce-cli containerd.io`

3. Verify that Docker Engine - Community is installed correctly by running the
    hello-world image.

`$ sudo docker run hello-world`

These instructions are sourced from the article <https://docs.docker.com/v17.09/engine/installation/linux/docker-ce/ubuntu/>. For additional information and links to other Linux distributions, see the instructions provided by other supported Linux distributions.

Adding the Pega configuration files to your Helm installation on your local system
----------------------------------------------------------------------------------

Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm, including a generic version of the following charts. After you add the repository to your local system, you can customize these Pega configuration files for your Pega Platform deployment:

- pega/pega - Use this chart to set customization parameters for your deployment. You will modify this chart later in the deployment tasks.

- pega/addons – Use this chart to install any supporting services and tools which your Kubernetes environment will require to support a Pega deployment. The required services, such as a load balancer or metrics server, that your deployment requires depend on your cloud environment. For instance you can specify whether you want to use a generic load-balancer or use one that is offered in your Kubernetes environment, such as in AKS or EKS. The runbooks provide instructions to deploy these supporting services once per Kubernetes environment, when you install the addons chart, regardless of how many Pega Infinity instances are deployed.

To customize these files, you must download them from the repository to your local system, edit them with a text editor, and then save them to your local system using the same filename. In this set of tasks, you will focus on the pega/addons.yaml file; in the environment-specific runbook that you are using in the section, **Update the Helm chart values**, you will update the pega.yaml file.

To simplify the instruction, you can download the file to the \<platform\>-demo folder you have already created on your local system. 

1. To add the Pega repository to your Helm installation, enter:

`$ helm repo add pega https://dl.bintray.com/pegasystems/pega-helm-charts`

2. To verify the new repository, you can search it by entering:

```yaml
$ helm search repo pega
NAME        CHART VERSION   APP VERSION     DESCRIPTION
pega/pega   1.2.0                           Pega installation on kubernetes
pega/addons 1.2.0           1.0             A Helm chart for Kubernetes
```

These two charts in this repository, pega and addons, require customization for your deployment of Pega Platform.

Updating the Pega/addons Helm chart
-----------------------------------

Update this Helm chart in order to enable the Traefik load balancer and disable
the metrics-server for deployments to the following platforms:

- AKS

- PKS

- GKE

If you are deploying to a different platform you can skip this section.

1. To download pega/addons Helm chart to the \<local filepath\>/\<platform\>-demo, enter:

`$ helm inspect values pega/addons > addons.yaml`

2. Open the addons.yaml file from this folder in a text editor

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

4. For GKE or PKS deployments, you must ensure that the Pega metrics server is disabled in the metrics-server section of this *addon* values.yaml file, since PKS deployments use the PKS metrics server.

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

The Pega deployment model supports advanced configurations to fit most existing clients' needs. If you are a Pega client and have known, required customizations for your deployment and you already use the following files to add your known customizations, you can copy those configurations into the configuration files Pega added for this purpose in the pega-helm-charts repository folder, \<local filepath\>/\<platform\>-demo/pega-helm-charts-master/charts/pega/config/deploy:

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

1. In the browser of your choice, navigate to the Pega [Digital Software Delivery](https://community.pega.com/digital-delivery) site.

2. Log into the [Pega Community](https://community.pega.com/knowledgebase/articles/pega-cloud/pega-cloud-services-patch-process-releases-83x-and-later)
    site with the credentials your Pega representative provided.

3. In the **Download and Upgrade Licensed Software** area, click **New
    request**.

4. In the right side of the page click **Continue**.

If you have multiple associations with the Pega Community, the page requests you to select the organization with which you want to affiliate this request and then click **Continue**. You will receive an email with a link to your software using an email address that is associated with the organization you select on this screen.

5. In the **You're viewing products available** page, enter **Pega Platform** in the **Search**, which will filter the list of products in the page.

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

- An email with a link to the requested Pega Platform software is sent within a few minutes. The email address used is associated with the organization you selected in this section.

![](media/748ea91e3ff43cf4544ce2f4638e86bf.png)

11. When satisfied with the order, click **Close**.

### Downloading Pega Platform to your local system

To download your Pega Platform image,

1. Open the email you received. It will look similar to the image shown.

![](media/98b1055e0e63487db7bbb2c90c9ea40c.png)

2. Click **Download now**.

3. The **Pega Licensed Software Downloads** page opens.

You can download your requested Pega Platform software using the link under **My Downloads**.

4. Click **Download software**.

Your secure **Inbox** of requested Pega software products opens. Your request for a version of Pega Platform software is listed at the top of the inbox table.

5. In the **Subject** column, click the link to your requested Pega Platform software.

The Package details window opens in the Package tab, which shows details about the Pega Platform software distribution package that you requested.

6. In the **Files:** area of the window, ensure that version of the Pega distribution image is correct.

If it is not the right version number, you must complete a new request.

7. To download the file, select the Pega distribution image checkbox and click
    **Download**.

8. In the **Save as** window, choose the \<local filepath\>/\<platform\>-demo folder to which you save the Pega Platform distribution zip file.

9. In a Linux bash shell, change folders to the /home/\<local filepath\>/\<platform\>-demo directory, where you saved the Pega Platform distribution zip and extract your files to create a new distribution image folder on your local system:

`$ unzip ./<pega-distribution-image>.zip`

After you extract the archive, the files in the Pega Platform distribution image are available to use in preparing your Pega Platform installation Docker image.

Prepare your Pega Platform installation Docker image – 15 minutes
-----------------------------------------------------------------

As stated previously, you must have a [DockerHub](https://hub.docker.com/) account and log into it in order to see the [pega-installer-ready Docker image](https://hub.docker.com/r/pegasystems/pega-installer-ready). You also need the Docker cli and Docker-desktop installed on your system before you begin this procedure. The Pega-provided Docker image, pega-installer-ready, includes some components of a full installation image that you can use to install or upgrade the Pega Platform database. While it is built on top of a JDK, it does not contain the contents of the Pega distribution kit which are essential for installing or upgrading Pega Platform.

Pega provides this image as the primary content of the final Docker image you will use to install or upgrade Pega Platform. This section describes how you can use this Docker image in combination with a Dockerfile and the Pega Platform distribution image that you have made available on your local system. The procedure assumes you’ve downloaded the software in [Downloading Pega Platform to your local system](#downloading-a-pega-platform-distribution-to-your-local-system) and installed the required components on your local system listed in [Install required applications for the deployment](#creating-a-local-folder-to-access-all-of-the-configuration-files).

Follow these steps to create a Docker image you can use to install or upgrade Pega Platform.

1. From a Linux bash shell, ensure you are logged into your DockerHub account:

`$ docker login -u <username> -p <username-password>`

For details about logging into Docker from a secure password file using the `--password-stdin` option, see <https://docs.docker.com/engine/reference/commandline/login/>.

2. Change your directory to the top folder of your Pega distribution <pega-distribution-image\>.

`$ cd ./<pega-distribution-image>/`

3. Create a text file with the text editor of your choice in the \<local filepath\>/\<platform\>-demo/\<pega-distribution-image\> folder where you extracted the Pega distribution on your local system.

From this folder, you can list the folder content and see folders for Pega archives, Images, rules, and scripts.

4. Copy the following lines of instruction into the new text file:

```yaml
FROM pegasystems/pega-installer-ready
COPY scripts /opt/pega/kit/scripts
COPY archives /opt/pega/kit/archives
COPY rules /opt/pega/kit/rules
RUN chmod -R 777 /opt/pega
```

These instructions direct a docker build function to use the Pega public Docker image, pega-install-ready, and these three folders from the Pega distribution image in order to build your Pega Platform installation image.

5. Save the text-only file with the filename, "dockerfile", without an extension, in the \<local filepath\>/\<platform\>-demo/\<pega-distribution-image\> folder where you extracted the Pega distribution on your local system.

6. From your Linux command prompt, in your current directory, build your pega install Docker image by entering:

`$ docker build -t pega-installer .`

This command uses your dockerfile to build a full Docker image with the name “pega-installer” and gives it the tag, “latest”. You can use the Docker command to see that your new image exists in your Docker image inventory.

7. Tag the local version of your new image, pega-installer, with your DockerHub ID:

`$ docker tag pega-installer <your-dockerhub-ID>/pega-installer`

8. Create a private repository on your [DockerHub](https://hub.docker.com/) account that is tagged as a private repository.

9. From your default login page, click the Repositories link (at the top of the page).

10. In the Repositories view, click **Create Repository +**.

11. Enter a Name that matches the Docker image you just built.

12. Provide a brief Description, that will help you remember the version of Pega
    with which you built this image, or any other useful information.

13. In the Visibility area, select the **Private**.

You should not maintain this image with Pega proprietary software as a viewable **Public** image.

14. Click **Create**.

Free DockerHub accounts support the use of a single private repository, so you may have to delete an existing private repository in order to create a new one for your Pega Docker installation image.

15. From a Linux bash shell, use the Docker command to push the new image to
    your new private repository:

`$ docker push <your-dockerhub-ID>/pega-installer`

After the command completes you will see your new image in your private repository, similar to the image below.

![](media/9fd09158a821f828a93d6ab7c74e278a.png)
