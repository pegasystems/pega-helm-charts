# Using the Pega Platform installer-ready Docker image

Pega provides a Docker image, `pegasystems/pega-installer-ready`, that installs or upgrades the Pega Platform against your database when you run it. Pull the image from [PEGA-DOWNLOADS LINK].

## When to use this image directly

If you deploy Pega Platform with the Helm charts in this repository, the `pega` chart's installer job already runs this image for you as part of `helm install` or `helm upgrade`. In that case, you do not need to run the steps in this document; instead, follow one of the [Deploying Pega Platform](Deploying-Pega-on-AKS.md) guides for your environment, or [Patching Pega Platform in your deployment](patching-pega-deployment.md) if you are applying a patch.

Use this image when you want to install or upgrade Pega Platform outside of Kubernetes as part of a client-managed cloud environment.

## Prerequisites

- Docker is installed on the system where you will run the image.
- Your database is created and reachable, and you have a JDBC driver for it. The driver must be mounted into the `/opt/pega/lib/` directory in the container.
- You have downloaded a Pega Platform distribution kit for the license you purchased. For instructions on obtaining a kit, see [Downloading a Pega Platform distribution to your local system](building-your-own-Pega-installer-image.md#downloading-a-pega-platform-distribution-to-your-local-system).
- You can pull `pegasystems/pega-installer-ready` from [PEGA-DOWNLOADS LINK].

## Providing the Pega distribution kit

Because the `pega-installer-ready` image doesn't include a Pega distribution kit, you must supply one at container start, using one of the following methods.

1. **Mount a kit archive.** Mount a directory containing exactly one distribution kit `.zip` file at `/opt/pega/mount/kit`.
   ```bash
   $ docker run -v /some/local/directory:/opt/pega/mount/kit -e ACTION="install" pegasystems/pega-installer-ready
   ```
2. **Download the kit at startup.** Set `DISTRIBUTION_KIT_URL` to a URL the container can download the kit `.zip` from. If the URL requires authentication, also set the `CUSTOM_ARTIFACTORY_*` variables (see the environment variable tables below).
   ```bash
   $ docker run -e "DISTRIBUTION_KIT_URL=https://example.com/kit.zip" -e ACTION="install" pegasystems/pega-installer-ready
   ```

## Installing Pega Platform

Run the image with `ACTION=install` along with your database connection details:

```bash
$ docker run -it \
    -v /path/to/kit-folder:/opt/pega/mount/kit \
    -e ACTION="install" \
    -e JDBC_URL="<your-jdbc-url>" \
    -e DB_TYPE="<your-db-type>" \
    -e JDBC_CLASS="<your-jdbc-driver-class>" \
    -e DB_USERNAME="<db-username>" \
    -e DB_PASSWORD="<db-password>" \
    -e RULES_SCHEMA="rules" \
    -e DATA_SCHEMA="data" \
    -e ADMIN_PASSWORD="<temporary-admin-password>" \
    pegasystems/pega-installer-ready
```

Watch the container's console output to follow the progress of the installation. The container exits when the installation finishes.

## Upgrading Pega Platform

Run the image with `ACTION=upgrade` and an `UPGRADE_TYPE` of `in-place`, `out-of-place`, or `custom`:

```bash
$ docker run -it \
    -v /path/to/kit-folder:/opt/pega/mount/kit \
    -e ACTION="upgrade" \
    -e UPGRADE_TYPE="in-place" \
    -e JDBC_URL="<your-jdbc-url>" \
    -e DB_TYPE="<your-db-type>" \
    -e JDBC_CLASS="<your-jdbc-driver-class>" \
    -e DB_USERNAME="<db-username>" \
    -e DB_PASSWORD="<db-password>" \
    -e RULES_SCHEMA="rules" \
    -e DATA_SCHEMA="data" \
    pegasystems/pega-installer-ready
```

If you use `UPGRADE_TYPE="custom"`, you must also provide `UPGRADE_STEP` to specify which single step to run: `enable_cluster_upgrade`, `rules_migration`, `rules_upgrade`, `data_upgrade`, or `disable_cluster_upgrade`.

```bash
$ docker run -it -e ACTION="upgrade" -e UPGRADE_TYPE="custom" -e UPGRADE_STEP="enable_cluster_upgrade" pegasystems/pega-installer-ready
```

These are the same upgrade concepts used by the Helm-based patch process, so for background on choosing an upgrade type, see [Upgrading your Pega Platform deployment](upgrading-pega-deployment-zero-downtime.md). If an upgrade fails partway through, see [Resuming failed upgrades from point of failures](resuming-rules-upgrade-failures.md) before you re-run the container.

## Choosing an image version

You can pin the image to a specific release or build by adding a tag to the image name, for example `pegasystems/pega-installer-ready:3.9-jdk17`. Available tags are:

- `3.9-<jdk>` (`3.9-jdk11`, `3.9-jdk17`, `3.9-jdk21`) returns the 3.9 release built for that JDK version.
- `latest` (or no tag specified) returns the latest good build.

Every image also has a `build_timestamp` label so you can confirm exactly which build you're running:

```bash
$ docker inspect --format '{{ index .Config.Labels "build_timestamp" }}' pegasystems/pega-installer-ready[:tag]
```

The same value is exported inside the container as the `BUILD_TIMESTAMP` environment variable and is printed to the container logs on startup.

## Environment variables

You configure the image entirely through environment variables passed to `docker run`. Anywhere a table below lists a variable as sensitive, you can provide it as a mounted file at `/opt/pega/secrets/<NAME>` instead of a plain environment variable.

### Action input

| Name | Description |
| ---- | ------------ |
| ACTION | Provide `install` or `upgrade` as an action to perform a Pega installation or upgrade, respectively |
| UPGRADE_TYPE | Required only during upgrade. Provide the type of upgrade to perform: `in-place`, `out-of-place`, or `custom` |
| UPGRADE_STEP | Required only when `UPGRADE_TYPE` is `custom`. One of `enable_cluster_upgrade`, `rules_migration`, `rules_upgrade`, `data_upgrade`, `disable_cluster_upgrade` |

### Database information

| Name | Description |
| ---- | ------------ |
| JDBC_URL | The JDBC URL used to connect to the database |
| DB_TYPE | Database type |
| JDBC_CLASS | JDBC driver class |
| DB_USERNAME | Database username (can be provided as a mounted secret file) |
| DB_PASSWORD | Database password (can be provided as a mounted secret file) |

### Schema information

| Name | Description |
| ---- | ------------ |
| RULES_SCHEMA | Rules schema name |
| DATA_SCHEMA | Data schema name |
| CUSTOMERDATA_SCHEMA | Customer data schema name |

### System information

| Name | Description |
| ---- | ------------ |
| SYSTEM_NAME | System name that uniquely identifies a single system |
| PRODUCTION_LEVEL | The system production level. Range is 1-5 |
| ADMIN_PASSWORD | Sets the temporary password for administrator@pega.com |
| MT_SYSTEM | Set to enable a multitenant system, which allows organizations to act as separate Pega Platform installations |

### Customizable installation parameters

| Name | Description |
| ---- | ------------ |
| BYPASS_UDF_GENERATION | Set to `true` to skip UDF generation |
| BYPASS_PEGA_SCHEMA | Set to `true` to skip schema generation |
| ASSEMBLER | Set to `true` to run the Static Assembler |
| BYPASS_TRUNCATE_UPDATESCACHE | Set to `true` to bypass automatically truncating PR_SYS_UPDATESCACHE |
| JDBC_CUSTOM_CONNECTION | JDBC custom connection properties |
| DISTRIBUTION_KIT_URL | URL to download the distribution kit from; see [Providing the Pega distribution kit](#providing-the-pega-distribution-kit) |

### Custom artifactory authentication

Used to authenticate `DISTRIBUTION_KIT_URL` and JDBC driver downloads against a custom artifactory. Each of these can also be provided as a mounted file at `/opt/pega/secrets/<NAME>` instead of a plain environment variable.

| Name | Description |
| ---- | ------------ |
| CUSTOM_ARTIFACTORY_USERNAME | Username for Basic authentication |
| CUSTOM_ARTIFACTORY_PASSWORD | Password for Basic authentication |
| CUSTOM_ARTIFACTORY_APIKEY_HEADER | Header name to send an API key in, as an alternative to Basic authentication |
| CUSTOM_ARTIFACTORY_APIKEY | API key value sent in the `CUSTOM_ARTIFACTORY_APIKEY_HEADER` header |
| ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION | Set to `true` to enforce SSL verification when connecting to the custom artifactory |

### Thread level parameters

| Name | Description |
| ---- | ------------ |
| MAX_IDLE | Maximum idle threads |
| MAX_WAIT | Maximum wait threads |
| MAX_ACTIVE | Maximum active threads |

### Upgrade related properties

| Name | Description |
| ---- | ------------ |
| TARGET_RULES_SCHEMA | Target rules schema name |
| MIGRATION_DB_LOAD_COMMIT_RATE | The commit count to use when loading database tables |
| UPDATE_EXISITING_APPLICATIONS | Set to `true` to run Update Existing Applications |
| UPDATE_APPLICATIONS_SCHEMA | Set to `true` to run the Update Applications Schema utility, which updates the cloned Rule, Data, Work, and Work History tables with the schema changes in the latest base tables |
| RUN_RULESET_CLEANUP | Set to `true` to generate and run a SQL script that cleans old rulesets and their rules from the system |
| REBUILD_INDEXES | Set to `true` to rebuild database rules indexes after the rules load, to improve database access performance |
