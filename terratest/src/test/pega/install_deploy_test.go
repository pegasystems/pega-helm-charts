package pega

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
)

// Sets the the action to install-deploy, all test cases present in this file uses this action
var installDeployOptions = &helm.Options{
	SetValues: map[string]string{
		"global.provider":        "k8s",
		"global.actions.execute": "install-deploy",
	},
}

// VerifyInstallDeployActionSkippedTemplates - Tests all the skipped templates for action install-deploy. These templates not supposed to be rendered for install-deploy action.
func VerifyInstallDeployActionSkippedTemplates(t *testing.T) {
	output := helm.RenderTemplate(t, installDeployOptions, PegaHelmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-upgrade-environment-config.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	// assert that above templates are not rendered
	require.Empty(t, deployment)
}

// VerifyInstallDeployActionInstallerJob - Tests Install job yaml rendered with the values as provided in default values.yaml for action install-deploy
func VerifyInstallDeployActionInstallerJob(t *testing.T) {
	var installerJobObj k8sbatch.Job
	var installerSlice = ReturnJobSlices(t, PegaHelmChartPath, installDeployOptions)
	helm.UnmarshalK8SYaml(t, installerSlice[1], &installerJobObj)
	VerifyPegaJob(t, installDeployOptions, &installerJobObj, pegaJob{"pega-db-install", []string{}, "pega-install-environment-config"})
}

// TestInstallDeployActions - Test all objects deployed for install-deploy action with the values as provided in default values.yaml
func TestInstallDeployActions(t *testing.T) {
	VerifyInstallDeployActionSkippedTemplates(t)
	VerifyInstallDeployActionInstallerJob(t)
	VerifyInstallerConfigMaps(t, installDeployOptions, PegaHelmChartPath)
	VerifyInstallEnvConfig(t, installDeployOptions, PegaHelmChartPath)
	VerifyInstallerRoleBinding(t, installDeployOptions, PegaHelmChartPath)
	VerifyInstallerRole(t, installDeployOptions, PegaHelmChartPath)
	VerifyPegaStandardTierDeployment(t, PegaHelmChartPath, installDeployOptions, []string{"wait-for-pegainstall", "wait-for-pegasearch", "wait-for-cassandra"})
}
