package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
)

// Path to the helm chart we will test
const pegaHelmChartPath = "../../../charts/pega"

// set action execute to install
var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "upgrade",
		"cassandra.enabled":      "false",
	},
}

func VerifyUpgradeActionShouldNotRenderDeployments(t *testing.T) {
	// with action as 'install' below templates should not be rendered
	output := helm.RenderTemplate(t, options, pegaHelmChartPath, []string{
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

func VerifyUpgradeActionInstallJob(t *testing.T) {
	var upgradeJobObj k8sbatch.Job
	var upgradeSlice = returnJobSlices(t, pegaHelmChartPath, options)
	helm.UnmarshalK8SYaml(t, upgradeSlice[1], &upgradeJobObj)
	VerifyJob(t, options, &upgradeJobObj, pegaJob{"pega-db-upgrade", []string{}, "pega-upgrade-environment-config"})
}

func TestUpgradeActions(t *testing.T) {
	VerifyUpgradeActionShouldNotRenderDeployments(t)
	VerifyUpgradeActionInstallJob(t)
	VerifyUpgradeEnvConfig(t, options, pegaHelmChartPath)
	VerfiyRegistrySecret(t, pegaHelmChartPath, options)
	VerifyCredentialsSecret(t, pegaHelmChartPath, options)
}
