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

To stay consistent with the instructions, you can create a folder named `<platform>-demo` on your local system at the top level of your home directory. This way, you associate the generic `<local filepath>/<platform>-demo` references to the folder `/home/<linux-username>/<platform>-demo` that is specific to your local system.

To create this folder, open a Linux command prompt and enter:

`$ mkdir /home/<linux-username>/<platform>-demo`

Where `\<platform\>-demo` is:

- For the AKS runbook, `AKS-demo`

- For the EKS runbook, `EKS-demo`

- For the GKS runbook, `GKS-demo`

- For the TKGI runbook, `TKGI-demo`

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
- AWS cli (only for EKS deployments)
- pks cli (only for TKGI deployments)
- Google Cloud SDK and gcloud (only for GKE deployments)

### Installing Unzip

Install the Unzip utility with a packaging tool such as apt.

`$ sudo apt install unzip`

### Installing Helm

Pega supports using Helm version 2.1 and later. The latest runbooks use version 3.0 and it’s recommended to use this version. If you use Helm 2.x, some of the commands differ slightly for Helm 2.x.

Helm supports a variety of installation methods. To learn more, see https://helm.sh/docs/intro/install/.

Helm provides a script that you can download and then run to install the latest version:

1. To download the Helm installation script from the Helm Git repository, from your home directory, enter:

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

### For TKGI only: installing the PKS CLI

Install the PKS CLI binary file from the VMware Enterprise PKS support site as an executable that you will run as a super user with a single command with the command curl piped into your bash:

1. Use the browser of your choice to navigate to [Installing the PKS CLI](https://docs.pivotal.io/pks/1-7/installing-pks-cli.html) and log in.

2. Select release version 1.7 or later from the Releases dropdown.

3. Click **PKS CLI - Linux** to display the **Release Download Files**.

4. Click **PKS CLI - Linux** and in the **Save as** window, choose the `<local filepath>/TKGI-demo` folder to which you save the downloaded Linux binary file.

5. Rename the downloaded binary file to "pks".

6. On the command line, to make the PKS CLI binary executable, enter:

    `$ chmod +x pks`

7. To move the binary file into your $PATH directory, enter:

    `$ sudo mv ./pks /bin`

These instructions were mostly sourced from the [Installing the PKS CLI](https://docs.pivotal.io/pks/1-7/installing-pks-cli.html).

### For AKS only: installing the Azure CLI

To install the Azure CLI as a super user with a single command with the command curl piped into your bash, enter:

`$ curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash`

For details, see the article, <https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-apt?view=azure-cli-latest>.

### For GKE only: installing and initializing the Google Cloud SDK 

In order to use the Google Cloud SDK for your deployment, you must install the software and then initialize its use by referencing the Google Cloud project in which you will deploy Pega Platform.

To install the Google Cloud SDK:

1. In your browser, log in to your Google user account.

2. Ensure that Python 2.7 is installed.

3. Download the appropriate Google SDK distribution for your Linux version. See [Installing from versioned archives](https://cloud.google.com/sdk/docs/downloads-versioned-archives). There are also versions available for Debian and Ubuntu or Red Hat and CentOS Linux.

4. To extract the Google SDK distribution archive file to your system, enter:

    `tar zxvf [ARCHIVE_FILE] google-cloud-sdk`

5. To use the install script to add Cloud SDK tools to your path, enter:

    `./google-cloud-sdk/install.sh`

6. Restart your terminal for the changes to take effect.

7. From your command prompt, to initialize the Google Cloud SDK, enter:

    `gcloud init`

    After the program initializes, you are prompted to log in.

8. Log in using your Google user account by entering **Y**:

    `To continue, you must log in. Would you like to log in (Y/n)? Y`

    You are redirected to a browser page on your system with a Google Cloud log in screen.

9. In your browser, log in to your Google Cloud user account when prompted and click **Allow** to grant permission to access Google Cloud Platform resources.

10. In your command prompt, select a Cloud Platform project from the list of those where you have Owner, Editor or Viewer permissions:

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

11. To list accounts whose credentials are stored on the local system, enter:

    `gcloud auth list`

12. To view information about your Cloud SDK installation and the active SDK configuration, enter:

    `gcloud info`

These instructions were sourced from the Google document, [Quickstart for Linux](https://cloud.google.com/sdk/docs/quickstart-linux), which includes additional information.

### For EKS only: installing AWS IAM authenticator

Pega recommends using the AWS IAM Authenticator for Kubernetes to authenticate with your Kubernetes cluster using your AWS credentials. You must download the binary file from AWS and then install the CLI on your local computer.

1. To download the latest binary from the AWS site, from your home directory, enter:

```bash
   $ curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.15.10/2020-02-22/bin/linux/amd64/aws-iam-authenticator
   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                  Dload  Upload   Total   Spent    Left  Speed
   100 33.6M  100 33.6M    0     0  1996k      0  0:00:17  0:00:17 --:--:-- 4786k
```

2. To make the AWS IAM Authenticator binary executable, enter:

    `$ chmod +x ./aws-iam-authenticator`

3. To use the AWS recommended command to copy the AWS IAM Authenticator file to a newly created folder, `$HOME/bin/aws-iam-authenticator`, and ensure that $HOME/bin comes first in your $PATH, enter:

   `$ mkdir -p $HOME/bin && cp ./aws-iam-authenticator $HOME/bin/aws-iam-authenticator && export PATH=$PATH:$HOME/bin`

4. To add $HOME/bin to your PATH environment variable, enter:

   `$ echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc`

5.  To verify the AWS IAM Authenticator is working, enter:

   ```bash
   $ ws-iam-authenticator help
   A tool to authenticate to Kubernetes using AWS IAM credentials

   Usage:
     aws-iam-authenticator [command]
```

These instructions were sourced from the AWS document, [Installing aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html).

### For EKS only: installing the eksctl utility

Pega recommends deploying your EKS cluster for your Pega Platform deployment using the AWS `eksctl` command line utility for creating and managing clusters on Amazon EKS.

1. To download and extract the latest `eksctl` binary from the AWS site, from your home directory, enter:

   `$ curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp`

2. To move the extracted binary to /usr/local/bin, which allows the command to be run without the "sudo" prefix, enter:

   `$ sudo mv /tmp/eksctl /usr/local/bin`

3. To verify the `eksctl` version, enter:

   ```bash
   $ eksctl version
   0.17.0
   ```

These instructions were sourced from the AWS document, [Install eksctl](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html), which includes additional information about how to use the 'eksctl' utility in combination with the 'kubectl' command-line utility, which you should have already installed earlier in this document.

### For EKS only: installing AWS CLI version 2

Install the AWS CLI utility with a packaging tool such as apt-get. Depending on the tool you use, youmay need a have SAML profile for your AWS account on your system.

`$ sudo apt-get install awscli`

To confirm the installation or upgrade, enter:

```bash
$ aws --version
aws-cli/2.0.10 Python/3.7.3 Linux/5.3.0-51-generic botocore/2.0.0dev14
```

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

After you install Docker CE on your Linux system, you can run Docker commands without using the 'sudo' preface. To learn how, see [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user). The recommended steps to change ownership of the Unix socket that binds to the Docker daemon on your Linux system (from the default, `root`, to your system user) originate in that document and are added here for your convenience. You can skip steps 3 - 6 if you run the docker command with the `sudo` preface.

1. Update the apt package index.

    `$ sudo apt-get update`

2. Install the latest version of Docker CE or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

    `$ sudo apt-get install docker-ce docker-ce-cli containerd.io`

3. To create a `docker` group of users, enter:

    `$ sudo groupadd docker`

4. To add your user name to the docker group, enter:

    `$ sudo usermod -aG docker $USER`

5. Log out and log back in so that your group membership is re-evaluated.

6. To verify that the dpocker command runs without the `sudo` preface and that Docker CE is installed correctly by running the hello-world image, enter:

    `$ docker run hello-world`

    These installation instructions are sourced from the Docker documentation,  [Install Docker Engine on Ubuntu](https://docs.docker.com/v17.09/engine/installation/linux/docker-ce/ubuntu/). For additional information and links to other Linux distributions, see the instructions provided by other supported Linux distributions.

Creating SSL certificates for HTTPS load balancer on a Kubernetes ingress
-----------------------------------

When you configure an HTTPS load balancer through a Kubernetes ingress, you can configure the load balancer to present authentication certificates to the client. After you generate your certificates and corresponding private keys using any SSL Certificate creation tool, you pass the certificates and private keys to your Kubernetes ingress using the method supported for the specific environment. For Pega support details, see the [Addons Helm chart](https://github.com/pegasystems/pega-helm-charts/tree/master/charts/addons).

Use the generation method best suited to your company standards. The following steps create one or more certificates and private keys for a host name using a manually verified domain.

1. Navigate to the certificate generator, such as [SSL For Free](https://www.sslforfree.com/).

2. Enter your host name for which you are creating a certificate and click **Create FreeSSL Certificate**.

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

With the host name associated with the verification value, wait several minutes to ensure the configuration is established.

6. Click the link for name/host record.

   A new page displays the status of the connection.

7. When successful, you will see that the link returned the correct value and you can generate the SSL certificate files.

   If it is not successful, you may have to wait longer for the DNS lookup to correctly associate the verification value with the host name; if it continues to not work, you may need to update your DNS settings.

8. In the main page of **SSL for Free**, click **Download All SSL Certificate Files** and save the certificate file and private key file in the \<platform\>-demo folder you have already created on your local system.

You will manage these certificates in the environment to which you deploy Pega Platform. For environment-specific details, see the **Deploying Pega Platform using the command line** section in the runbook for that environment.

Deploying Pega Platform using Pega-provided docker images
---------------------------------------------------------

To deploy Pega Platform, you must pull several required images from the Pega-managed Docker image repository and push them into your private Docker registry from where you reference them in the Pega Helm chart. For more information, see [Pega Helm chart](https://github.com/pegasystems/pega-helm-charts).

Pegasystems uses a standard naming practice of hostname/product/image:tag. Pega images are available from the host site, pega-docker.downloads.pega.com. Pega maintains three types of required Docker images for Client-managed Cloud deployments of Pega Platform:

 Name        | Description                                           | Tags     |
-------------|-------------------------------------------------------|----------|
`platform/installer`   | A utility image with which you install all of the Pega-specific rules and database tables in the “Pega” database that you have configured for your deployment. This installation is required before a deployment can take place.| `<version>` |
`platform/pega`  | (Download required) Deploys Pega Platform with its customized version of the Tomcat application server.| `<version>` or `<version>-YYYYMMDD` |
`platform/search` | (Download required) Deploys the required search engine for Pega Platform search and reporting capabilities. This Docker image contains Elasticsearch and includes all required plugins.| `<version>` or `<version>-YYYYMMDD` |

When you decide on a Pega Platform version for your downloaded Docker images, you should use the same version tag for each of the three images you download.

For the `platform/installer` image, the :tag represents the version of Pega you want to install, for example the tag :8.5.1 will install Pega Platform version 8.5.1.

For `platform/pega` and `platform/search` images, Pega also offers an image with a version tag appended with a datestamp using the pattern `pegaVersion-YYYYMMDD` to indicate the version and the date that Pega built the image. For example, if you pull the `platform/pega` with a tag, `pega:8.5.1-20201026`, the tag indicates that Pega built this 8.5.1 image on 26 October 2020. Using the version tag without the datestamp will always point to the most recently built image for that version.

The datestamp ensures that the image you download includes the changes that Pega engineering commits to the repository using pull requests by a certain date. While Pega builds the most current patch version of each minor release one time each day, Pega makes the last five daily-built images available for client downloads.  After Pega releases a new patch version, the prior patch version no longer receives daily builds with a datestamp tag.

After you obtain access to the Pega-provided host repository and pull each image, you can re-tag and push each of the three Pega-provided images to your preferred Docker registry to make them available to the deployment as described in the next section. You then provide your registry URL, credentials, and reference each image appropriately in the Pega Helm chart. You can find example usage details for referencing the three images in a repository in the appropriate runbook for your type of deployment.

These images do not expire, and you can keep them in your repository for as long as you require.

Pega also supports experienced client's ability to build your own installer docker image using components of a full Pega Platform 8.3 or later distribution image to install or upgrade the Pega Platform database. If you build your own installation image, you do not have to download the Pega-provided docker image listed in the table. For details on building your own image, see [Building a Pega Platform installer docker image](building-your-own-Pega-installer-image.md).

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

### Downloading and managing Pega Platform docker images

With access credentials to the Pega Docker image repository, you log in, download each required image, retag each image appropriately, and finally upload each image to your own Docker registry. For an overview of tagging and managing Docker images, see the Docker article, [Deploy a registry server](https://docs.docker.com/registry/deploying/).

Important: Use the same Docker image tag to ensure you download compatible images.

Pega supports any of the following Docker image registries from which your deployment will access the three Pega-provided Docker images. For details about setting up your choice of Docker registry, click the link for that registry's documentation:

- [DockerHub](https://docs.docker.com/docker-hub/repos/)
- [Amazon elstic Container Registry (ECR)](https://docs.aws.amazon.com/AmazonECR/latest/userguide/get-set-up-for-amazon-ecr.html)
- [Azure Container Registry](https://azure.microsoft.com/en-us/services/container-registry/)
- [Google Cloud Container Registry](https://cloud.google.com/container-registry/)

With a Docker registry configured, clients push their three Pega-provided images to their registry so it is available to the deployment. Clients must also provide their registry URL, credentials, and then reference each image appropriately in the Pega Helm chart.

Examples:

```yaml
# If using a custom Docker registry, supply the credentials here to pull Docker images.
  docker:
    registry:
      url: "YOUR_DOCKER_REGISTRY"
      username: "YOUR_DOCKER_REGISTRY_USERNAME"
      password: "YOUR_DOCKER_REGISTRY_PASSWORD"
      # Docker image information for the Pega docker image, containing the application server.
    pega:
      image: "<Registry host name:Port>/my-pega:<Pega Platform version>"

# Elasticsearch deployment settings.
pegasearch:
  image: "<Registry host name:Port>/my-pega-search:<Pega Platform version>"
  memLimit: "3Gi"
  replicas: 1

# Pega Installer settings.
installer:
  image: "<Registry host name:Port>/my-pega-installer:<Pega Platform version>"
  adminPassword: "ADMIN_PASSWORD"
```

Example usage details for referencing your three images in a repository are included in the appropriate runbook for your type of deployment.

It is a best practice to retag each of your Pega Docker images by including your registry host name and port; if this information is not included, the docker tag command uses the Docker public registry located at [registry-1.docker.io](https://registry-1.docker.io/) by default. For more details and naming convention guidance, see the [docker tag](https://docs.docker.com/engine/reference/commandline/tag/) documentation.

1. In a Linux bash shell with root privileges, change folders to the \home\<local filepath>\<platform>-demo directory, navigate to the <localfilepath>\<platform>-demo folder where you saved the file that contains your access key and log into the Pega-managed Docker image repository:

```bash
$ cat <localfilepath>\<platform>-demo\<access-key-filename>.txt | docker login pega-docker.downloads.pega.com --username=<reg-XXXXX> --password-stdin
Login Succeeded
```

2. To download your preferred version of the `Pega` image to your local system, specify the version tag when you enter:

```bash
$ docker pull pega-docker.downloads.pega.com/platform/pega:<version>
Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/pega:<version>
```

3. Retag the `Pega` image for your deployment with a customized tag that includes your Docker registry host name and a name that is useful to your organization, such as `<Registry host name:Port>/my-pega:<Pega Platform version>`.

   `$ docker tag pega-docker.downloads.pega.com/platform/pega:8.4.0 <Registry host name:Port>/my-pega:8.4.0`

4. To push the retagged `my-pega` image to your registry, enter:

   `$ docker push <Registry host name:Port>/my-pega:8.4.0`

5. To download your preferred version of the `search` image to your local system, specify the version tag when you enter:
 
```bash
$ docker pull pega-docker.downloads.pega.com/platform/search:<version>
Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/search:<version>
```

6. To retag the `search` image for your deployment with a customized tag that includes your Docker registry host name and a name that is useful to your organization, such as `<Registry host name:Port>/my-pega-search:<Pega Platform version>`, enter:

   `$ docker tag pega-docker.downloads.pega.com/platform/search:8.4.0 <Registry host name:Port>/my-pega-search:8.4.0`

7. To push the retagged `my-pega-search` image to your registry, enter:

   `$ docker push <Registry host name:Port>/my-pega-search:8.4.0`

8. To download your preferred version of the `installer` image to your local system, specify the version tag when you enter:

```bash
$ docker pull pega-docker.downloads.pega.com/platform/installer:<version>
Digest: <encryption verification>
Status: Downloaded pega-docker.downloads.pega.com/platform/installer:<version>
```

9. To retag the `installer` image for your deployment with a customized tag that includes your Docker registry host name and a name that is useful to your organization, such as `<Registry host name:Port>/my-pega-installer:<Pega Platform version>`, enter:

   `$ docker tag pega-docker.downloads.pega.com/platform/installer:8.4.0 <Registry host name:Port>/my-pega-installer:8.4.0`

10. To push the retagged `my-pega-installer` image to your registry, enter:

   `$ docker push <Registry host name:Port>/my-pega-installer:8.4.0`

After you push these three downloaded images to your private Docker registry, you are ready to begin deploying Pega Platform to a support Kubernetes environment. Use the runbook in this Github directory for your deployment.
