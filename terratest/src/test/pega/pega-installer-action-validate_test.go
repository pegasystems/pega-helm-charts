package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

func TestPegaInstallerAction_WithValidUpgradeType(t *testing.T) {

	var supportedActions = []string{"upgrade"}
	var supportedUpgradeTypes = []string{"in-place" ,"out-of-place", "custom", "out-of-place-rules" ,"out-of-place-data"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, action := range supportedActions {
		for _, upgradeType := range supportedUpgradeTypes {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":               "k8s",
					"global.actions.execute":        action,
					"installer.upgrade.upgradeType": upgradeType,
				},
			}

			yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-action-validate.yaml"})

			require.Contains(t, yamlContent, "could not find template charts/installer/templates/pega-installer-action-validate.yaml")
			require.Contains(t, err.Error(), "could not find template charts/installer/templates/pega-installer-action-validate.yaml")
		}
	}
}
func TestPegaInstallerAction_WithValidUpgradedeployType(t *testing.T) {

	var supportedActions = []string{"upgrade-deploy"}
	var supportedUpgradeTypes = []string{"zero-downtime"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, action := range supportedActions {
		for _, upgradeType := range supportedUpgradeTypes {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":               "k8s",
					"global.actions.execute":        action,
					"installer.upgrade.upgradeType": upgradeType,
				},
			}

			yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-action-validate.yaml"})

			require.Contains(t, yamlContent, "could not find template charts/installer/templates/pega-installer-action-validate.yaml")
			require.Contains(t, err.Error(), "could not find template charts/installer/templates/pega-installer-action-validate.yaml")
		}
	}
}

func TestPegaInstallerAction_WithInValidUpgradeType(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		SetValues: map[string]string{
			"global.provider":               "k8s",
			"global.actions.execute":        "upgrade",
			"installer.upgrade.upgradeType": "invalidValue",
		},
	}

	yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-action-validate.yaml"})
	require.Contains(t, yamlContent, "Upgrade Type value is not correct")
	require.Contains(t, err.Error(), "Upgrade Type value is not correct")

}
