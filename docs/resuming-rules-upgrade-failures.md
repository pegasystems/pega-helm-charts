# Resuming failed upgrades from point of failures

## You can only resume rules_upgrade failures from point of failure when you set the upgradeType parameter to “custom” when you upgrade.

## Client-required steps

### Steps to enable resume functionality
- Specify `action.execute: upgrade` to upgrade your application using the software version contained in your Pega-provided "installer" image.
- Specify `installer.upgrade.upgradeType: custom`
- Specify `installer.upgrade.upgradeSteps: rules_upgrade` to run rules_upgrade
- Provide a Persistent Volume Claim name in the `installer.installerMountVolumeClaimName` parameter. This is a client-managed PVC for mounting upgrade artifacts.
- You must create the PVC manually in the same namespace where you deploy Pega Platform before the volume will be bound.
- Set `installer.upgrade.automaticResumeEnabled` to `true` to enable this functionality
- Run the upgrade process by using the `helm upgrade release --namespace mypega` command. For more information, see  [Upgrading your Pega Platform deployment using the command line](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/upgrading-pega-deployment-zero-downtime.md#upgrading-your-pega-platform-deployment-using-the-command-line).

### Steps to resume failed rules_upgrade
- To debug upgrade failure, follow the instructions in the README section  - [Debugging failed upgrades using helm commands](https://github.com/pegasystems/pega-helm-charts/blob/master/README.md#debugging-failed-upgrades-using-helm-commands)
- Identify the root cause and make necessary changes to fix the issue before resuming the upgrade
- Resume the upgrade process by using the `helm upgrade release --namespace mypega` command. For more information, see - [Upgrading your Pega Platform deployment using the command line](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/upgrading-pega-deployment-zero-downtime.md#upgrading-your-pega-platform-deployment-using-the-command-line).
- The upgrade skips steps that were completed successfully in the previous run.