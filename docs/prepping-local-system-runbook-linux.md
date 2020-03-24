Preparing your local Linux system – 45 minutes
=====================================================================

To deploy Pega Platform using a local Linux system on which you can run commands with administrator privileges, prepare your system with required applications and configuration files that you need for your deployment. By preparing your system, you can complete the deployment without having to pause to obtain a Linux executable file or prepare a configuration file that is required to complete the deployment.

Before you begin
----------------

To prepare a Linux system, it is assumed:

- You have a basic familiarity with running commands from a Linux command prompt with and without root privileges.

- You use a packaging tool to install application packages. For demonstration purposes, this document refers to the Ubuntu Advanced Packaging Tool (apt) commands that is available for use on Ubuntu Linux distributions.

- You have basic familiarity with a GitHub account with which you download a Pega-managed GitHub repository that contains configuration files and scripts that support installing Pega Platform and deploying it in the Kubernetes cluster.

Creating a local folder to access all of the configuration files
----------------------------------------------------------------

Deploying with Helm requires you to run commands from a specific folder on your local system. To ensure you use the correct filepath, these instructions always use the reference `<local filepath>/<platform>-demo` folder when you extract files to a folder or run commands from a folder.

To stay consistent with the instructions, you can create a folder named `<platform>-demo` on your local system at the top level of your Windows user folder. This way, you associate the generic `<local filepath>/<platform>-demo` references to the folder `/home/<linux-username>/<platform>-demo` that is specific to your local system.

To create this folder, open a Linux command prompt and enter:

`$ mkdir /home/<linux-username>/<platform>-demo`

Where `\<platform\>-demo` is:

- For the AKS runbook, `AKS-demo`

- For the EKS runbook, `EKS-demo`

- For the GKS runbook, `GKS-demo`

- For the PKS runbook, `PKS-demo`

- For the Openshift runbook, `Openshift-demo`

You are ready to continue preparing your local system.

Installing required applications for the deployment
---------------------------------------------------

The entire deployment requires the following applications to be used during the configuration process; therefore, you should prepare your local system with all of the applications before you start your deployment:
- Helm (download the binary from the Helm GitHub repository)
- kubectl (download the binary file from the kubernetes git repository)
- Docker (Use the Linux apt or equivalent package manager)
- unzip (or an equivalent to extract files from .zip archive files.)
- az cli (only for AKS deployments)
- AWS IAM Authenticator for Kubernetes (only for EKS deployments)
- eksctl (only for EKS deployments)
- pks cli (only for PKS deployments)
- Google Cloud SDK and gcloud (only for GKE deployments)

### Installing Unzip

Install the Unzip utility with a packaging tool such as apt.

`$ sudo apt install unzip`

### Installing Helm

Pega supports using Helm version 2.1 and later. The latest runbooks use version 3.0 and it’s recommended to use this version. If you use Helm 2.x, some of the commands differ slightly for Helm 2.x.

Helm supports a variety of installation methods. To learn more, see https://helm.sh/docs/intro/install/.

Helm provides a script that you can download and then run to install the latest version:

1. To download the Helm installation script from their Git repository, from
    your home directory enter:

`$ curl
https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 >
get_helm.sh`

2. To update the permissions of the file to you can use it for installations, enter:

    `$ chmod 700 get_helm.sh`

3. To run the script, enter:

```bash
$ ./get_helm.sh
Helm v3.0.1 is available.
Downloading https://get.helm.sh/helm-v3.0.1-linux-amd64.tar.gz
Preparing to install helm into /usr/local/bin
helm installed into /usr/local/bin/helm
```

4. To review your version, enter:

```bash
$ helm version
version.BuildInfo{Version:"v3.0.1",
GitCommit:"7c22ef9ce89e0ebeb7125ba2ebf7d421f3e82ffa", GitTreeState:"clean",
GoVersion:"go1.13.4"}
```

For additional information, see [Helm documentation](https://helm.sh/docs/); for details about installation methods for previous Helm versions, see <https://v2.helm.sh/docs/using_helm/#installing-helm>.

### Installing kubectl

Kubernetes supports a variety of installation methods for the kubectl command. The organization provides a link to download the latest version of the executable file, which will be run after you move the binary file in to your PATH:

1. To download the latest binary from their git repository, enter:

```bash
$ curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s
https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
```

This command downloads and parses the \`stable.txt\` in the repository, navigates to the version listed in the file, and downloads the kubectl binary file in the /bin/linux/amd64/ directory.

2. Make the kubectl binary executable:

    `$ chmod +x ./kubectl`

3. Move the script to your executable PATH:

    `$ sudo mv ./kubectl /usr/local/bin/kubectl`

For details about installing on Linux and other supported installation methods, see [Install kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-linux).

### For PKS only: installing the PKS CLI

Install the PKS CLI binary file from the Pivotal support site as an executable that you will run as a super user with a single command with the command curl piped into your bash:

1. Use the browser of your choice to navigate to [Pivotal Network](https://network.pivotal.io/) and log in.

2. Click [Pivotal Container Service (PKS)](https://network.pivotal.io/products/pivotal-container-service).

3. Select release version 1.6 or later from the Releases dropdown.

4. Click **PKS CLI - Linux** to display the **Release Download Files**.

5. Click **PKS CLI - Linux** and in the **Save as** window, choose the `<local filepath>/<platform>-demo` folder to which you save the downloaded Linux binary file.

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

4. Extract the Google SDK distribution archive file to your system.

    `tar zxvf [ARCHIVE_FILE] google-cloud-sdk`

5. Use the install script to add Cloud SDK tools to your path.

    `./google-cloud-sdk/install.sh`

6. Restart your terminal for the changes to take effect.

7. In a Linux command prompt, initialize the Google Cloud SDK.

    `gcloud init`

    After the program initializes, you are prompted to log in.

8. Log in using your Google user account by entering **Y**:

    `To continue, you must log in. Would you like to log in (Y/n)? Y`

    You are redirected to a browser page on your system with a Google Cloud log in screen.

9. In your browser, log in to your Google Cloud user account when prompted and click **Allow** to grant permission to access Google Cloud Platform resources.

10. In your Linux command prompt, select a Cloud Platform project from the list of those where you have Owner, Editor or Viewer permissions:

```bash
Pick cloud project to use:
 [1] [my-project-1]
 [2] [my-project-2]
 ...
 Please enter your numeric choice or text value (must exactly match list item):
```

If you only have one project, `gcloud init` selects it for you. After your selection, the command confirms that you completed the setup steps successfully:

```bash
Your current project has been set to: [my-project-name].
...
Your Google Cloud SDK is configured and ready to use!

* Commands that require authentication will use [my-account-name] by default
* Commands will reference project `[my-project-name]` by default
Run `gcloud help config` to learn how to change individual settings

This gcloud configuration is called [default]. You can create additional configurations if you work with multiple accounts and/or projects.
Run `gcloud topic configurations` to learn more.
```

11. To list accounts whose credentials are stored on the local system:

    `gcloud auth list`

12. To view information about your Cloud SDK installation and the active SDK configuration:

    `gcloud info`

    These instructions were sourced from the Google document, [Quickstart for Linux](https://cloud.google.com/sdk/docs/quickstart-linux), which includes additional information.

### Installing Docker

For Linux command line users, you can follow these steps to install Docker Community Edition (CE) for the first time on a new host machine. For these instructions, you need to set up the Docker repository. Afterward, you can install and update Docker from the repository. For your convenience, the instructions available on the Docker website are included in this this section of this document.

**SET UP THE REPOSITORY**

1. Update the apt package index:

    `$ sudo apt-get update`

2. Install packages to allow apt to use a repository over HTTPS:

    `$ sudo apt-get install apt-transport-https ca-certificates curl software-properties-common`

3. Add Docker’s official GPG key:

    `$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add –`

4. Verify that you now have the key with the fingerprint `9DC8 5822 9FC7 DD38
    854A E2D8 8D81 803C 0EBF CD88` by searching for the last 8 characters of the
    fingerprint.

    `$ sudo apt-key fingerprint 0EBFCD88`

    The command should return:

```bash
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

Deploying Pega Platform using Pega-provided docker images
---------------------------------------------------------

To deploy Pega Platform, you must pull several required images from the Pega-managed Docker image repository and push them into your private Docker registry from where you reference them in the Pega Helm chart. For more information, see [Pega Helm chart](https://github.com/pegasystems/pega-helm-charts).

Pegasystems uses a standard naming practice of `hostname/product/image:tag`.  All Pega images are available from the `pega-docker.downloads.pega.com` host.  The `:tag` represents the version if Pega being deployed, for example `:8.3.1` to download Pega 8.3.1.  Pega maintains three types of required Docker images for Client-managed Cloud deployments of Pega Platform:

 Name        | Description                                           |
-------------|-------------------------------------------------------|
platform/pega  | *Download required. Deploys the Pega Platform with its customized version of the Tomcat application server |
 platform/search | *Download required. Deploys the search engine required for the Pega Platform application’s search and reporting capabilities. This Docker image contains Elasticsearch and includes all required plugins |
 platform/installer   | A utility image Pega Platform deployments use to install or upgrade all of the Pega-specific rules and database tables in the “Pega” database you have configured for your deployment.
 
You must build an installer docker image to install or upgrade all of the Pega-specific rules and database tables in the “Pega” database of your deployment. To do so, follow the tasks in the sections to download the Pega Platform distribution and build an installer image using the Pega files in the distribution.

## Downloading a Pega Platform installer docker image

Clients with appropriate licenses can log in to the image repository and download docker images that are required to install the Pega Platform onto your database.

### Requesting access to the Pega Docker repository

1. In the browser of your choice, navigate to the Pega [Digital Software Delivery](https://community.pega.com/digital-delivery) site.

2. Log into the [Pega Community](https://community.pega.com/knowledgebase/articles/pega-cloud/pega-cloud-services-patch-process-releases-83x-and-later)
    site with the credentials your Pega representative provided.

3. In the **Download and Upgrade Licensed Software** area, click **New
    request**.

4. In the right side of the page, click **Request access key**.

![Select your distribution](media/dockerimage-download.png)

5. Enter your credential details.

   After you enter valid credentials, you recieve confirmation that an email is on the way. 

![Select your distribution](media/dockerimage-download-access.png)

6. Open the email you received. It will look similar to the image shown.

![Confirmation email with access key details](media/dockerimage-access-email.png)

7. Save your access key to a text file the <local filepath>\<platform>-demo folder so you can pass it into your docker login command to ensure the it will not display in your bash history or logs.

### Downloading Pega Platform docker images to your local system

With access to the Pega-managed Docker image repository, clients log in to the image repository and download required images.

1. In a Linux bash shell with root privileges, change folders to the \home\<local filepath>\<platform>-demo directory, navigate to the <localfilepath>\<platform>-demo folder where you saved your access key, and log into the Pega-managed Docker image repository:

```bash
$ cat <localfilepath>\<platform>-demo\<access-key-filename>.txt |  docker login pega-docker.downloadsqa.pega.com --username=reg-<User ID> --password-stdin
Login Succeeded
```

2. To download the version of Pega image for your deployment, enter and specify the version tag:

```bash
$ docker pull pega-docker.downloads.pega.com/platform/pega:<version>
Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/pega:<version>
```

3. To download the version of search image image for your deployment, enter and specify the version tag:

```bash
$ docker pull pega-docker.downloads.pega.com/platform/search:<version>
Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/search:<version>
```

To build an installer docker image to install or upgrade all of the Pega-specific rules and database tables in the “Pega” database of your deployment, follow the tasks in the next sections.

## Building a Pega Platform installer docker image

These instructions require the Pega Platform distribution image to install the Pega Platform onto your database. 

Clients with appropriate licenses can download a distribution image from Pega. For additional instructions, see [Pega Digital Software Delivery User Guide](https://community.pega.com/knowledgebase/documents/pega-digital-software-delivery-user-guide).



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
