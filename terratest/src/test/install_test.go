package test

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
)

const pegaHelmChartPath = "../../../charts/pega"

// set action execute to install
var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "install",
		"cassandra.enabled":      "false",
		"global.provider":        "k8s",
	},
}

// VerifyInstallActionSkippedTemplates - Tests all the skipped templates for action install. These templates not supposed to be rendered for install action.
func VerifyInstallActionSkippedTemplates(t *testing.T) {
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"templates/pega-environment-config.yaml",
		"templates/pega-tier-config.yaml",
		"templates/pega-tier-deployment.yaml",
		"templates/pega-tier-hpa.yaml",
		"templates/pega-tier-ingress.yaml",
		"templates/pega-tier-service.yaml",
		"charts/installer/templates/pega-installer-role.yaml",
		"charts/installer/templates/pega-installer-status-rolebinding.yaml",
	})

	var emptyObjects appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &emptyObjects)

	// assert that above templates are not rendered
	require.Empty(t, emptyObjects)
}

// VerifyInstallActionInstallJob - Tests Install job yaml rendered with the values as provided in default values.yaml
func VerifyInstallActionInstallJob(t *testing.T) {
	var upgradeJobObj k8sbatch.Job
	var upgradeSlice = ReturnJobSlices(t, pegaHelmChartPath, options)
	helm.UnmarshalK8SYaml(t, upgradeSlice[1], &upgradeJobObj)
	VerifyPegaJob(t, options, &upgradeJobObj, pegaJob{"pega-db-install", []string{}, "pega-install-environment-config"})
}

//TestInstallActions - Test all objects deployed for install action with the values as provided in default values.yaml
func TestInstallActions(t *testing.T) {
	VerifyInstallActionSkippedTemplates(t)
	VerifyInstallActionInstallJob(t)
	VerifyInstallEnvConfig(t, options, pegaHelmChartPath)
	VerfiyRegistrySecret(t, pegaHelmChartPath, options)
	VerifyCredentialsSecret(t, pegaHelmChartPath, options)
	VerifyInstallerConfigMaps(t, options, pegaHelmChartPath)
}
