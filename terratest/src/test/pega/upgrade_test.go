package pega

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
)

// set action execute to install
var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "upgrade",
		"cassandra.enabled":      "false",
		"global.provider":        "k8s",
	},
}

// VerifyUpgradeActionShouldNotRenderDeployments - Tests all the skipped templates for action upgrade. These templates not supposed to be rendered for upgrade action.
func VerifyUpgradeActionSkippedTemplates(t *testing.T) {
	output := helm.RenderTemplate(t, options, PegaHelmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-installer-role.yaml",
		"templates/pega-environment-config.yaml",
		"charts/installer/templates/pega-installer-status-rolebinding.yaml",
		"charts/pegasearch/templates/pega-search-deployment.yaml",
		"charts/pegasearch/templates/pega-search-service.yaml",
		"charts/pegasearch/templates/pega-search-transport-service.yaml",
		"charts/installer/templates/pega-install-environment-config.yaml",
		"templates/pega-tier-config.yaml",
		"templates/pega-tier-deployment.yaml",
		"templates/pega-tier-hpa.yaml",
		"templates/pega-tier-ingress.yaml",
		"templates/pega-tier-service.yaml",
	})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	// assert that above templates are not rendered
	require.Empty(t, deployment)
}

// VerifyUpgradeActionInstallJob - Tests upgrade job yaml rendered with the values as provided in default values.yaml
func VerifyUpgradeActionInstallJob(t *testing.T) {
	var upgradeJobObj k8sbatch.Job
	var upgradeSlice = ReturnJobSlices(t, PegaHelmChartPath, options)
	helm.UnmarshalK8SYaml(t, upgradeSlice[1], &upgradeJobObj)
	VerifyPegaJob(t, options, &upgradeJobObj, pegaJob{"pega-db-upgrade", []string{}, "pega-upgrade-environment-config"})
}

//TestUpgradeActions - Test all objects deployed for upgrade action with the values as provided in default values.yaml
func TestUpgradeActions(t *testing.T) {
	VerifyUpgradeActionSkippedTemplates(t)
	VerifyUpgradeActionInstallJob(t)
	VerifyUpgradeEnvConfig(t, options, PegaHelmChartPath)
	VerfiyRegistrySecret(t, PegaHelmChartPath, options)
	VerifyCredentialsSecret(t, PegaHelmChartPath, options)
	VerifyInstallerConfigMaps(t, options, PegaHelmChartPath)
}
