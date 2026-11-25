# Migrating from Helm 3 to Helm 4

The Pega Helm charts are validated with both Helm 3 and Helm 4. The charts themselves do not change between the two versions, so you do not need a different chart or different `values.yaml` files. Only the Helm client changes.

Use this page when you move an existing Helm 3 based deployment or automation to the Helm 4 client. For the full list of Helm 4 changes, see the [Helm 4 release announcement](https://helm.sh/blog/helm-4-released/) and the [Helm 4 overview](https://helm.sh/docs/overview/).

## Supported versions

| Helm client | Status |
| --- | --- |
| Helm 3.x | Supported and validated in CI. |
| Helm 4.x | Supported and validated in CI. |

The Helm project extended the Helm 3 security fix period, and security fixes now end on February 10, 2027. Check the [Helm 4 release announcement](https://helm.sh/blog/helm-4-released/) for the current support dates, and plan to move your automation to Helm 4 before security fixes end.

## Before you migrate

1. Install the Helm 4 client next to your Helm 3 client, for example under a different file name, so that you can fall back if needed.
2. Run `helm list --all-namespaces` and note the releases that you manage.
3. Back up your values files and, if needed, the release secrets (`kubectl get secrets -l owner=helm`). Helm 4 is expected to read the release records that Helm 3 created, so no conversion should be needed.
4. Update the chart dependencies with the Helm 4 client:

    ```bash
    helm dependency update ./charts/pega
    ```

## Migration steps

1. Run `helm version` and confirm that it returns Helm 4.x.
2. Render the chart and review the output against your current deployment:

    ```bash
    helm template <release> pega/pega --namespace <namespace> --values values.yaml > helm4.yaml
    ```

3. Upgrade the existing release with the Helm 4 client. Use the same chart version and values files that you use today:

    ```bash
    helm upgrade <release> pega/pega --namespace <namespace> --values values.yaml --set global.actions.execute=deploy
    ```

4. Verify the release with `helm status <release> --namespace <namespace>` and `kubectl get pods --namespace <namespace>`.

Switching the same release back to the Helm 3 client is expected to work, but test it in a non-production environment before you rely on it as a rollback plan.

## What changes in Helm 4

| Area | Helm 3 | Helm 4 | What to do |
| --- | --- | --- | --- |
| Apply method | Client-side apply | Server-side apply is the default for new installs. According to the Helm 4 announcement, upgrades keep the apply method that the release already uses. | Review the Helm 4 documentation if you depend on field ownership or conflict behavior. |
| Rollback flag | `--atomic` | `--rollback-on-failure` | Update scripts. The old flag is expected to keep working with a deprecation warning. |
| Replace flag | `--force` | `--force-replace` | Update scripts. The old flag is expected to keep working with a deprecation warning. |
| Waiting for resources | Polls resource status | Uses kstatus-based readiness checks when you use `--wait`. | If your deployment is slow to become ready, review the `--timeout` value and the `helm upgrade --help` output for the available wait options. |
| Post-renderers | Executable passed with `--post-renderer` | Post-renderers are plugins. | Package your post-renderer as a Helm plugin and pass its plugin name. |
| Registry login | `helm registry login <registry-url>` | `helm registry login <domain>` (the domain only, without a path) | Update scripts that log in to an OCI registry. |
| Chart metadata | Unknown `Chart.yaml` fields are ignored | Unknown fields are rejected by `helm lint --strict` | Remove unsupported fields from charts that you maintain. The Pega charts have no unsupported fields. |

## Charts that you maintain

- The Pega `addons`, `backingservices`, and `pega` charts, including their subcharts, use `apiVersion: v2`, and dependencies are declared in `Chart.yaml`. Helm 4 still installs `apiVersion: v1` charts, but you should move any of your own v1 charts to v2.
- If you use the optional Cassandra dependency, the `pega` chart now uses Cassandra chart `0.15.3`, which deploys Cassandra 3.11.6. When you upgrade an existing release that has `cassandra.enabled=true`, the Cassandra image tag changes, which can cause the Cassandra pods to be recreated. Plan this upgrade for a maintenance window.

## Validate your migration

After the upgrade, confirm the following:

- `helm lint --strict` passes for your chart and values files.
- `helm get manifest <release>` shows the expected resources.
- The Pega web and batch pods reach the `Running` and `Ready` state.
- `helm uninstall <release> --namespace <namespace>` removes the release when you test in a non-production namespace.
