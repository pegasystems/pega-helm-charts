# Resuming failed upgrades from point of failures

## Pega provides ability for clients to resume only rules_upgrade failures from point of failure , when upgraded with upgradeType custom

## Client-required steps

### Steps to enable resume functionality
- Specify `action.execute: upgrade` to upgrade your application using the software version contained in your Pega-provided "installer" image. 
- Specify `installer.upgrade.upgradeType: custom`
- Specify `installer.upgrade.upgradeSteps: rules_upgrade` to run rules_upgrade
- Input pvc name to `installer.installerMountVolumeClaimName` , which is a client managed Persistent Volume Claim for mounting upgrade artifacts. 
- Given PVC must be created manually in the same namespace where pega is deployed before volume will be bound.
- Set `installer.upgrade.automaticResumeEnabled` to `true` to enable this functionality 
- Invoke the upgrade process by using the `helm upgrade release --namespace mypega` command as directed in the deployment section - [Upgrading your Pega Platform deployment using the command line](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/upgrading-pega-deployment-zero-downtime.md#upgrading-your-pega-platform-deployment-using-the-command-line).

### Steps to resume failed rules_upgrade
- Debug upgrade failure as directed in README section - [Debugging failed upgrades using helm commands](https://github.com/pegasystems/pega-helm-charts/blob/master/README.md#debugging-failed-upgrades-using-helm-commands)
- Identify the root cause , and make necessary changes to fix the issue before resuming the upgrade
- Invoke the upgrade process by using the `helm upgrade release --namespace mypega` command as directed in the deployment section - [Upgrading your Pega Platform deployment using the command line](https://github.com/pegasystems/pega-helm-charts/blob/master/docs/upgrading-pega-deployment-zero-downtime.md#upgrading-your-pega-platform-deployment-using-the-command-line).
- Above invocation skips steps that are completed successfully in the previous run




