# Using the Pega Platform installer-ready Docker image

Pega provides a Docker image, `pegasystems/pega-installer-ready`, that installs or upgrades Pega Platform against your database when the `pega` Helm chart's installer job runs it. Pull the image from [PEGA-DOWNLOADS LINK].

## When to use pega-installer-ready instead of the combined installer image

Pega distributes a combined installer image that already contains a Pega distribution kit, and pointing `installer.image` at that image is the standard way client-managed cloud clients install or upgrade Pega Platform.

Use `pega-installer-ready` instead when you need to supply your own kit rather than use the one bundled into Pega's combined image, for example if that bundled kit has a known vulnerability you need to avoid. In this case, you point `installer.image` at `pega-installer-ready` itself, and supply the distribution kit to the installer job separately, either from a URL or a mounted volume. You still deploy the same way as any other Pega Platform deployment, by editing values.yaml and running `helm install` or `helm upgrade`, as described in the [Deploying Pega Platform](Deploying-Pega-on-AKS.md) guide for your environment. This document only covers the settings that differ when you use `pega-installer-ready`.

## Prerequisites

- You have completed the general Kubernetes and Helm setup described in the [Deploying Pega Platform](Deploying-Pega-on-AKS.md) guide for your environment.
- Your database is created and reachable.
- You can pull `pegasystems/pega-installer-ready` from [PEGA-DOWNLOADS LINK].
- You have a URL the installer job can use to download your Pega Platform distribution kit; the chart only supports providing the kit this way, not by mounting it from a local file. For instructions on obtaining a kit and hosting it somewhere the installer job can reach, see [Downloading a Pega Platform distribution to your local system](building-your-own-Pega-installer-image.md#downloading-a-pega-platform-distribution-to-your-local-system); you can stop after downloading, since you won't be building a combined image.

## Configuring the pega chart to use pega-installer-ready

1. Set `installer.image` to `pega-installer-ready` instead of a combined image:

   ```yaml
   installer:
     image: "pegasystems/pega-installer-ready:<tag>"
   ```

2. Set `installer.distributionKit.url` to a URL the installer job can download the kit `.zip` from. If the URL requires authentication, also set the `global.customArtifactory` values.

   ```yaml
   installer:
     distributionKit:
       url: "https://example.com/kit.zip"
   ```

   By default, the `pega-installer-ready` container downloads the kit itself at startup. If you'd rather have a separate init container download the kit before the installer container starts, for example to avoid giving the installer container direct network access, set `installer.distributionKit.downloadContainer.image` to an image that has `curl` available:

   ```yaml
   installer:
     distributionKit:
       url: "https://example.com/kit.zip"
       downloadContainer:
         image: "curlimages/curl:<tag>"
   ```

## Installing Pega Platform

Set `global.actions.execute` to `install` or `install-deploy`, and fill in your database connection details, then deploy as usual with `helm install`.

```yaml
global:
  actions:
    execute: "install-deploy"
  jdbc:
    url: "<your-jdbc-url>"
    dbType: "<your-db-type>"
    driverClass: "<your-jdbc-driver-class>"
    username: "<db-username>"
    password: "<db-password>"
    rulesSchema: "rules"
    dataSchema: "data"
installer:
  image: "pegasystems/pega-installer-ready:<tag>"
  adminPassword: "<temporary-admin-password>"
  distributionKit:
    url: "https://example.com/kit.zip"
```

```bash
$ helm install mypega-demo pega/pega --namespace mypega --values pega.yaml
```

Follow the progress of the installer job the same way you would for any Pega Platform deployment, by watching the installer pod and job logs from your Kubernetes dashboard or `kubectl`.

## Upgrading Pega Platform

Set `global.actions.execute` to `upgrade` or `upgrade-deploy`, and set `installer.upgrade.upgradeType` to `in-place`, `out-of-place-rules`, `out-of-place-data`, `zero-downtime`, or `custom`.

```yaml
global:
  actions:
    execute: "upgrade-deploy"
installer:
  image: "pegasystems/pega-installer-ready:<tag>"
  upgrade:
    upgradeType: "in-place"
```

If you use `upgradeType: "custom"`, also set `installer.upgrade.upgradeSteps` to a comma-separated list of the steps to run, in order: `enable_cluster_upgrade`, `rules_migration`, `rules_upgrade`, `data_upgrade`, `disable_cluster_upgrade`.

For background on choosing an upgrade type, see [Upgrading your Pega Platform deployment](upgrading-pega-deployment-zero-downtime.md); for applying a zero-downtime patch, see [Patching Pega Platform in your deployment](patching-pega-deployment.md). If an upgrade fails partway through, see [Resuming failed upgrades from point of failures](resuming-rules-upgrade-failures.md) before you retry.

## Choosing an image version

You can pin `installer.image` to a specific release or build by adding a tag, for example `pegasystems/pega-installer-ready:3.9-jdk17`. Available tags are:

- `3.9-<jdk>` (`3.9-jdk11`, `3.9-jdk17`, `3.9-jdk21`) returns the 3.9 release built for that JDK version.
- `latest` (or no tag specified) returns the latest good build.