package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
)

const pegaHelmChartPath = "../../../charts/pega"

var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "install-deploy",
	},
}

func VerifyInstallDeployActionSkippedTemplates(t *testing.T) {
	// with action as 'install-deploy' below templates should not be rendered
	output := helm.RenderTemplate(t, options, pegaHelmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-upgrade-environment-config.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	// assert that above templates are not rendered
	require.Empty(t, deployment)
}
func VerifyInstallDeployActionInstallerJob(t *testing.T) {
	var installerJobObj k8sbatch.Job
	var installerSlice = returnJobSlices(t, pegaHelmChartPath, options)
	helm.UnmarshalK8SYaml(t, installerSlice[1], &installerJobObj)
	VerifyJob(t, options, &installerJobObj, pegaJob{"pega-db-install", []string{}, "pega-install-environment-config"})
}

func TestInstallDeployActions(t *testing.T) {
	VerifyInstallDeployActionSkippedTemplates(t)
	VerifyInstallDeployActionInstallerJob(t)
	VerifyPegaStandardTierDeployment(t, pegaHelmChartPath, options, []string{"wait-for-pegainstall", "wait-for-pegasearch", "wait-for-cassandra"})
	VerifyInstallerConfigMaps(t, options, pegaHelmChartPath)
	VerifyInstallEnvConfig(t, options, pegaHelmChartPath)
	VerifyInstallerRoleBinding(t, options, pegaHelmChartPath)
	VerifyInstallerRole(t, options, pegaHelmChartPath)

}
