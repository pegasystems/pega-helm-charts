# Using the Pega Platform installer-ready Docker image

Pega provides a Docker image, `platform/installer-ready`, that installs or upgrades Pega Platform against your database when the `pega` Helm chart's installer job runs it. For instructions on how to pull this image, refer to the [Installing with the Pega Installer Docker Image](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/install-process-use-docker-image.html) guide.

## When to use installer-ready instead of the combined installer image

Pega distributes a combined installer image that already contains a Pega distribution kit, and pointing `installer.image` at that image is the standard way client-managed cloud clients install or upgrade Pega Platform.

Use `installer-ready` instead when you need to use a standalone kit rather than use the one bundled into Pega's combined image, for example if that bundled kit has a known vulnerability that cannot be a part of a Docker image. In this case, you point `installer.image` at `installer-ready` itself, and supply the distribution kit to the installer job separately, either from a URL or a mounted volume. You still deploy the same way as any other Pega Platform deployment, by editing values.yaml and running `helm install` or `helm upgrade`, as described in the [Installing with the Pega Installer Docker Image](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/install-process-use-docker-image.html) guide for your environment. This document only covers the settings that differ when you use `installer-ready`.

## Prerequisites

- You have completed the general Kubernetes and Helm setup described in the [Preparing the environment](https://docs.pega.com/bundle/platform/page/platform/deployment/client-managed-cloud/preparing-the-deployment-install.html) guide for your environment.
- Your database is created and reachable.
- The `platform/installer-ready` image has been pulled.
- You have either a URL the installer job can use to download your Pega Platform distribution kit or a manually created PersistentVolumeClaim (PVC) containing the kit `.zip` file. For instructions on obtaining a kit and hosting it somewhere the installer job can reach, see [Downloading the Pega Platform distribution package](https://docs.pega.com/bundle/platform/page/platform/deployment/downloading-pega-distribution-package.html).

## Configuring the pega chart to use installer-ready

1. Set `installer.image` to `installer-ready` instead of a combined image:

  ```yaml
  installer:
    image: "platform/installer-ready:<tag>"
  ```

2. Provide the distribution kit to the `installer-ready` image using one of the following options:

   1. Provide `installer.distributionKit.url` as a URL the installer job can download the kit `.zip` from. If the URL requires authentication, also set the `global.customArtifactory` values. Additionally, set `installer.distributionKit.downloadContainer.image` to an image that has `curl` available to be used as an init container.

      ```yaml
      installer:
        distributionKit:
          url: "https://example.com/kit.zip"
          downloadContainer:
            image: "curlimages/curl:<tag>"
      ```

   2. Provide the distribution kit through a persistent volume. Create the PVC in the same namespace as the Pega deployment and ensure it contains only the distribution kit `.zip` file. Set `installer.distributionKitVolumeClaimName` to the claim name and leave `installer.distributionKit.url` empty.

      ```yaml
      installer:
        distributionKit:
          url: ""
        distributionKitVolumeClaimName: "pega-distribution-kit"
      ```

      The installer job mounts the PVC at `/opt/pega/mount/kit`.

   If both `installer.distributionKit.url` and `installer.distributionKitVolumeClaimName` are specified, the URL takes precedence.

## Choosing an image version

You can pin `installer.image` to a specific release or build by adding a tag, for example `platform/installer-ready:3.9-jdk17`. Select the proper image to match your distribution's JDK. Reference the "Java Versions" section of the [Platform Support Guide](https://docs.pega.com/bundle/platform/page/platform/deployment/platform-support-guide/platform-support-guide.html) for details.

Available tags are:
- `3.9-<jdk>` (`3.9-jdk11`, `3.9-jdk17`, `3.9-jdk21`) returns the 3.9 release built for that JDK version.