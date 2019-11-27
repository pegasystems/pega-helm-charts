Preparing your local Linux system to deploy Pega Platform– 45 minutes
=====================================================================

In order to deploy using a local Linux system on which you can run commands with
Administrator privileges, you must prepare your system with required
applications and configuration files you will use for your deployment of Pega
Platform. Pega recommends doing this first so you can complete the deployment
without having to pause in order to obtain a Linus executable files or prepare a
configuration file that is required to complete the deployment.

Assumptions and Prerequisites
-----------------------------

This guide assumes:

-   You have a basic familiarity with running commands from a Linux command
    prompt with and without root privileges.

-   You use a packaging tool to install application packages. For demonstration
    purposes, this document refers to the Ubuntu Advanced Packaging Tool (apt)
    commands that is available for use on Ubuntu Linux distributions.

-   Basic familiarity with GitHub account with which you will download a
    Pega-managed GitHub repository containing configuration files and scripts
    that you use to install Pega Platform and then deploy it in a Kubernetes
    cluster.

Create a local folder to access all of the configuration files
--------------------------------------------------------------

Deploying with Helm requires that you run commands from a specific folder on
your local system. To ensure you stay oriented to the correct filepath, these
instructions always use the reference \<local filepath\>/\<platform\>-demo
directory when you must extract files to a folder or run commands from a folder.

In order to stay consistent with the instructions, it is recommended that you
create a directory called \<platform\>-demo on your local system in your /home
directory. This way, you associate the generic \<local
filepath\>/\<platform\>-demo references to the folder
/home/\<*linux-username*\>/\<platform\>-demo that is specific to your local
system.

To create this folder, open a Linux command prompt and enter:

\$ mkdir /home/\<*linux-username*\>/\<platform\>-demo

Where \<platform\>-demo is:

-   AKS-demo -for the AKS runbook

-   EKS-demo -for the EKS runbook

-   GKE-demo -for the GKE runbook

-   PKS-demo -for the PKS runbook

-   Openshift-demo - for the Openshift runbook

You are ready to continue preparing your local system.

Installing required applications for the deployment 
----------------------------------------------------

Some of the required applications are binary files that you download from the
organizations, download area. In order to use the docker command, the
instructions available on the docker website are included in this document if
you’d prefer to just doc the installation and not navigate to the docker
documentation support website.

### Installing unzip

Install the unzip utility with a packaging tool such as apt.

\$ sudo apt install unzip

### Installing Helm

Helm supports a variety of installation methods. Pega supports using v2.16.1.
Helm provides a script that you can download and then run to install the latest
version:

1.  To download the helm installation script from their git repository, enter:

\$ curl -LO https://git.io/get_helm.sh

1.  To update the permissions of the file to you can use it for installations,
    enter:

\$ chmod 700 get_helm.sh

1.  To run the script, enter:

\$ ./get_helm.sh

For details about this method and other supported installation methods, see
<https://v2.helm.sh/docs/using_helm/#installing-helm>.

1.  To review your version, enter:

\$ helm version

Client: &version.Version{SemVer:"v2.16.1",
GitCommit:"bbdfe5e7803a12bbdf97e94cd847859890cf4050", GitTreeState:"clean"}

### Installing kubectl

Kubernetes supports a variety of installation methods for the kubectl command.
The organization provides a link to download the latest version of the
executable file, which will be run after you move the binary file in to your
PATH:

1.  To download the latest binary from their git repository, enter:

\$ curl -LO https://storage.googleapis.com/kubernetes-release/release/\`curl -s
https://storage.googleapis.com/kubernetes-release/release/stable.txt\`/bin/linux/amd64/kubectl

This command downloads and parses the \`stable.txt\` in the repository,
navigates to the version listed in the file, and downloads the kubectl binary
file in the /bin/linux/amd64/ directory.

1.  Make the kubectl binary executable:

\$ chmod +x ./kubectl

1.  Move the script to your executable PATH:

\$ sudo mv ./kubectl /usr/local/bin/kubectl

For details about installing on Linux and other supported installation methods,
see
<https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-linux>.

### For PKS only: installing the PKS CLI

Install the PKS CLI binary file from the Pivotal support site as an executable
that you as a super user with a single command with the command curl piped into
your bash:

1.  Use the browser of your choice to navigate to [Pivotal
    Network](https://network.pivotal.io/) and log in.

2.  Click [Pivotal Container Service
    (PKS)](https://network.pivotal.io/products/pivotal-container-service).

3.  Select release version 1.6 or later from the Releases dropdown.

4.  Click **PKS CLI - Linux** to display the **Release Download Files**.

5.  Click **PKS CLI - Linux** and in the **Save as** window, choose the \<local
    filepath\>/\<platform\>-demo folder to which you save the downloaded Linux
    binary file.

6.  Rename the downloaded binary file to pks.

7.  On the command line, run the following command to make the PKS CLI binary
    executable:

\$ chmod +x pks

1.  Move the binary file into your \$PATH directory.

\$ sudo mv ./pks /bin

These instructions were mostly sourced from the [Installing the PKS
CLI](https://docs.pivotal.io/pks/1-6/installing-pks-cli.html).

### For AKS only: installing the Azure CLI

Install the Azure CLI as a super user with a single command with the command
curl piped into your bash:

\$ curl -sL https://aka.ms/InstallAzureCLIDeb \| sudo bash

For details, see the article,
<https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-apt?view=azure-cli-latest>.

### Installing Docker 

For Linux command line users, you can follow these steps to install Docker CE
for the first time on a new host machine. For these instructions, you need to
set up the Docker repository. Afterward, you can install and update Docker from
the repository.

**SET UP THE REPOSITORY**

1.  Update the apt package index:

\$ sudo apt-get update

1.  Install packages to allow apt to use a repository over HTTPS:

\$ sudo apt-get install apt-transport-https ca-certificates curl
software-properties-common

1.  Add Docker’s official GPG key:

\$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg \| sudo apt-key add –

1.  Verify that you now have the key with the fingerprint 9DC8 5822 9FC7 DD38
    854A E2D8 8D81 803C 0EBF CD88, by searching for the last 8 characters of the
    fingerprint.

\$ sudo apt-key fingerprint 0EBFCD88

The command should return:

pub 4096R/0EBFCD88 2017-02-22

9DC8 5822 9FC7 DD38 854A E2D8 8D81 803C 0EBF CD88

uid Docker Release (CE deb) \<docker\@docker.com\>

sub 4096R/F273FCD8 2017-02-22

1.  Use the following command to set up the stable repository. You always need
    the stable repository, even if you want to install builds from the edge or
    test repositories as well. To add the edge or test repository, add the word
    edge or test (or both) after the word stable in the commands below.

\$ sudo add-apt-repository "deb [arch=amd64]
https://download.docker.com/linux/ubuntu \$(lsb_release -cs) stable"

Note: Starting with Docker 17.06, stable releases are also pushed to the edge
and test repositories.

**INSTALL DOCKER CE**

1.  Update the apt package index.

\$ sudo apt-get update

1.  Install the latest version of Docker CE, or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

\$ sudo apt-get install docker-ce docker-ce-cli containerd.io

1.  Verify that Docker Engine - Community is installed correctly by running the
    hello-world image.

\$ sudo docker run hello-world

These instructions are sourced from the article
<https://docs.docker.com/v17.09/engine/installation/linux/docker-ce/ubuntu/>.
For additional information and links to other Linux flavors, you can find
instructions for other supported flavors of Linux.

Cloning the pega-helm-charts github repository to your local system
-------------------------------------------------------------------

Pega maintains a Github repository that contains the Helm charts that are
required to deploy Pega Platform using Helm. You must clone the repository to
your local system from which you will complete the deployment. You can either
clone by downloading a zip file of the repository or using Github desktop.

To access the GitHub website and begin the cloning process:

1.  Sign in to [GitHub](https://github.com/) with your GitHub credentials using
    the browser of your choice.

2.  Navigate to the main page of the
    [pega-helm-charts](https://github.com/pegasystems/pega-helm-charts)
    repository and choose the **Master** branch.

3.  In the repository heading, click **Clone or download**.

4.  In the **Clone with HTTPS** popup window, click **Download Zip**.

5.  In the explorer window, navigate to the local path, \<local
    filepath\>/\<platform\>-demo, ensure that the file name is
    pega-helm-charts-master.zip and click **Save**.

6.  In a Linux bash shell, change folders to the \<local
    filepath\>/\<platform\>-demo directory, where you saved the Pega Platform
    distribution zip and extract your files to create a new distribution image
    folder on your local system:

\$ unzip ./\<pega-distribution-image\>.zip

After you extract the files from the archive, you will run the deployment
commands from several of the folders in the \<local
filepath\>/\<platform\>-demo/pega-helm-charts-master/charts folder.

These instructions were mostly sourced from the [GitHub
help](https://help.github.com/en/desktop/contributing-to-projects/cloning-a-repository-from-github-to-github-desktop).

Updating the Pega addons Helm chart
-----------------------------------

Update this Helm chart in order to enable the Traefik load balancer and disable
the metrics-server for deployments to the following platforms:

-   AKS

-   PKS

-   GKE

If you are deploying to a different platform you can skip this section.

1.  In your pega-helm-chart repository, navigate to the \<local
    filepath\>/\<platform\>-demo /pega-helm-charts-master/charts/addons folder.

2.  Open the values.yaml file from this folder in a text editor

3.  In the traefik configuration area, ensure the following two settings are
    configured to use Traefik for your deployment load-balancer:

traefik:

enabled: **true**

\# Set any additional Traefik parameters. These values will be used by Traefik's
Helm chart.

\# See https://github.com/Helm/charts/blob/master/stable/traefik/values.yaml

\# Set traefik.serviceType to "LoadBalancer" on gke, PKS, and pks

serviceType: **LoadBalancer**

Note: Do not enclose the text in quotes.

1.  For PKS deployments, you must ensure that the Pega metrics server is
    disabled in the metrics-server section of this *addon* values.yaml file,
    since PKS deployments use the PKS metrics server

metrics-server:

\# Set this to true to install metrics-server. Follow below guidelines specific
to each provider,

\# open-source Kubernetes, Openshift & EKS - mandatory to set this to true if
any tier as hpa.enabled is true

\# GKE or PKS - set this to false since metrics-server is installed in the
cluster by default.

enabled: **false**

1.  Save the file.

Add any known, customized settings for Pega to your deployment
--------------------------------------------------------------

The Pega deployment model supports advanced configurations to fit most existing
client’s needs. If you are a Pega client and have known, required customizations
for your deployment and you already use the following files to add your known
customizations, you can copy those configurations into the configuration files
Pega added for this purpose in the pega-helm-charts repository folder, \<local
filepath\>/\<platform\>-demo/pega-helm-charts-master/charts/pega/config/deploy:

-   context.xml: add additional required data sources

-   prlog4j2.xml modify your logging configuration, if required

-   prconfig.xml: adjust the standard Pega Platform configuration with known,
    required settings

Make these changes before you begin in the section, [Deploying Pega Platform
using Helm charts – 30 minutes](#_Deploying_Pega_Platform).

Downloading a Pega Platform distribution to your local system
-------------------------------------------------------------

These instructions require the Pega Platform distribution image to install the
Pega Platform onto your database. To obtain a copy, you must download an image
from Pega. For detailed instructions, see [Pega Digital Software Delivery User
Guide](https://community.pega.com/knowledgebase/documents/pega-digital-software-delivery-user-guide).

### Requesting access to a Pega Platform distribution

1.  In the browser of your choice, navigate to the Pega [Digital Software
    Delivery](https://community1.pega.com/digital-delivery) site.

2.  Log into the [Pega
    Community](https://community.pega.com/knowledgebase/articles/pega-cloud/pega-cloud-services-patch-process-releases-83x-and-later)
    site with the credentials your Pega representative provided.

3.  In the **Download and Upgrade Licensed Software** area, click **New
    request**.

4.  In the right side of the page click **Continue**.

If you have multiple associations with the Pega Community, the page requests you
to first select the organization with which you want to affiliate this request
and then click **Continue**. You will receive an email with a link to your
software using an email address that is associated with the organization you
select on this screen.

1.  In the **You're viewing products available** page, enter **Pega Platform**
    in the **Search**, which will filter the list of products in the page.

The **Pega Platform** card should appear near the top of the card list, below
the list of all of the **Language packs for Pega Platform.**

1.  In the Pega Platform card, use your mouse to activate the icon into a
    shopping cart and click the shopping cart.

The icon changes to a green check and a new cart item appears in the top right
of the product list.

![](media/029c6531bd52109598047a2ee6966657.png)

1.  Click **Continue**.

2.  In the cart review page, in the **Pega Platform** area, select the version
    of Pega Platform for your deployment.

    ![](media/386d4eb20a4e2be6b767bc522cbdda91.png)

3.  After your selection and review are complete, click **Finish.**

4.  When the order is processed, a confirmation screen displays with details
    about your order.

    An email with a link to the requested Pega Platform software is sent within
    a few minutes. The email address used is associated with the organization
    you selected in this section.

![](media/748ea91e3ff43cf4544ce2f4638e86bf.png)

1.  When satisfied with the order, click **Close**.

### Downloading Pega Platform to your local system

To download your Pega Platform image,

1.  Open the email you received. It will look similar to the image shown.

![](media/98b1055e0e63487db7bbb2c90c9ea40c.png)

1.  Click **Download now**.

2.  The **Pega Licensed Software Downloads** page opens.

You can download your requested Pega Platform software using the link under **My
Downloads**.

1.  Click **Download software**.

Your secure **Inbox** of requested Pega software products opens. Your request
for a version of Pega Platform software is listed at the top of the inbox table.

1.  In the **Subject** column, click the link to your requested Pega Platform
    software.

The Package details window opens in the Package tab, which shows details about
the Pega Platform software distribution package that you requested.

1.  In the **Files:** area of the window, ensure that version of the Pega
    distribution image is correct.

If it is not the right version number, you must complete a new request.

1.  To download the file, select the Pega distribution image checkbox and click
    **Download**.

2.  In the **Save as** window, choose the \<local filepath\>/\<platform\>-demo
    folder to which you save the Pega Platform distribution zip file.

3.  In a Linux bash shell, change folders to the /home/\<local
    filepath\>/\<platform\>-demo directory, where you saved the Pega Platform
    distribution zip and extract your files to create a new distribution image
    folder on your local system:

\$ unzip .\\\<pega-distribution-image\>.zip

After you extract the archive, the files in the Pega Platform distribution image
are available to use in preparing your Pega Platform installation Docker image.

Prepare your Pega Platform installation Docker image
----------------------------------------------------

As stated previously, you are required to have a
[DockerHub](https://hub.docker.com/) account and log into it in order to see the
[pega-installer-ready Docker
image](https://hub.docker.com/r/pegasystems/pega-installer-ready). You also need
the docker cli and docker-desktop installed on your system before you begin this
procedure. The Pega-provided Docker image, pega-installer-ready, includes some
components of a full installation image that you can use to install or upgrade
the Pega Platform database. While it is built on top of a JDK, it does not
contain the contents of the Pega distribution kit which are essential for
installing or upgrading Pega Platform.

Pega provides this image as the primary content of the final Docker image you
will use to install or upgrade Pega Platform. This section describes how you can
use this Docker image in combination with a Dockerfile and the Pega Platform
distribution image that you have made available on your local system. The
procedure assumes you’ve downloaded the software in [Downloading Pega Platform
to your local system](#_Downloading_Pega_Platform) and installed the required
components on your local system listed in [Install required applications for the
deployment](#_Install_required_applications).

Follow these steps to create a Docker image you can use to install or upgrade
Pega Platform.

1.  From a Linux bash shell, ensure you are logged into your DockerHub account:

\$ docker login -u \<username\> --p \<username-password\>

For details about logging into Docker from a secure password file, see
<https://docs.docker.com/engine/reference/commandline/login/>. From Windows, you
can also ensure you are logged by used the Docker Desktop client running on your
system.

1.  Change your directory to the top folder of your Pega distribution,
    \<pega-distribution-image\>.

\$ cd ./\<pega-distribution-image\>/

1.  Create a text file with the text editor of your choice in the \<local
    filepath\>/\<platform\>-demo/\<pega-distribution-image\> folder where you
    extracted the Pega distribution on your local system.

You can list the folder content and see folders for Pega archives, Images,
rules, and scripts.

![](media/152260ae774fe07d717f1b31b5560f25.png)

1.  Copy the following lines into the file to build your docker image using the
    public image on DockerHub that Pega provides to build install images,
    pegasystems/pega-installer-ready:

>   FROM pegasystems/pega-installer-ready

>   COPY scripts /opt/pega/kit/scripts

>   COPY archives /opt/pega/kit/archives

>   COPY rules /opt/pega/kit/rules

1.  Save the text-only file as dockerfile, without an extension.

2.  From your command prompt (Linux or Powershell running with Administrator
    privileges), in your current directory, build your pega install docker image
    by entering:

`$ docker build -t pega-installer .`

This command uses your dockerfile to build a full Docker image with the name
“pega-installer” and gives it the tag, “latest”. You can use the docker command
to see that your new image exists in your Docker image inventory.

1.  Tag the local version of your new image, pega-installer, with your DockerHub
    ID:

`$ docker tag pega-installer <your-dockerhub-ID>/pega-installer`

1.  Create a private repository on your [DockerHub](https://hub.docker.com/)
    account that is tagged as a private repository.

2.  From your default login page, click the Repositories link (at the top of the
    page).

3.  In the Repositories view, click **Create Repository +**.

4.  Enter a Name that matches the docker image you just built.

5.  Provide a brief Description, that will help you remember the version of Pega
    with which you built this image, or any other useful information.

6.  In the Visibility area, select the **Private**.

    You should not mage this image with Pega proprietary software a viewable
    **Public** image.

7.  Click **Create**.

Free DockerHub accounts support the use of a single private repository, so you
may have to delete an existing private repository in order to create a new one
for your Pega docker installation image.

1.  From a Linux bash shell, use the docker command to push the new image to
    your new private repository:

`$ docker push <your-dockerbug-ID>/pega-installer`

After the command completes you will see your new image in your private
repository, similar to the image below.

![](media/9fd09158a821f828a93d6ab7c74e278a.png)
