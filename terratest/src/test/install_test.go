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
	},
}

func VerifyInstallActionSkippedTemplates(t *testing.T) {
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-action-validate.yaml",
		//		"templates/pega-credentials-secret.yaml",
		"templates/pega-environment-config.yaml",
		//		"templates/pega-registry-secret.yaml",
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

func VerifyInstallActionInstallJob(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	var upgradeJobObj k8sbatch.Job
	var upgradeSlice = returnJobSlices(t, helmChartPath, options)
	helm.UnmarshalK8SYaml(t, upgradeSlice[1], &upgradeJobObj)
	VerifyJob(t, options, &upgradeJobObj, pegaJob{"pega-db-install", []string{}, "pega-install-environment-config"})
}

func TestInstallActions(t *testing.T) {
	VerifyInstallActionSkippedTemplates(t)
	VerifyInstallActionSkippedTemplates(t)
	VerifyInstallEnvConfig(t, options, pegaHelmChartPath)
	VerfiyRegistrySecret(t, pegaHelmChartPath, options)
	VerifyCredentialsSecret(t, pegaHelmChartPath, options)
	VerifyInstallerConfigMaps(t, options, pegaHelmChartPath)
}
