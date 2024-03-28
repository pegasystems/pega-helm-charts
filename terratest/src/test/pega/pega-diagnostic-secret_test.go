package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

const diagnosticWebUser = "webuser"
const diagnosticWebPassword = "webpass"
const diagnosticGlobalUser = "globaluser"
const diagnosticGlobalPassword = "globalpass"

func TestWebTierPegaDiagnosticSecret(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					ValuesFiles: []string{"data/values_with_tier_diagnostic_user.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": getUpgradeTypeForUpgradeAction(operation),
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-diagnostic-secret.yaml"})
				verifyDiagnosticSecret(t, yamlContent, diagnosticWebUser, diagnosticWebPassword)
			}
		}
	}
}

func TestGlobalPegaDiagnosticSecret(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					ValuesFiles: []string{"data/values_with_global_diagnostic_user.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": getUpgradeTypeForUpgradeAction(operation),
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-diagnostic-secret.yaml"})
				verifyDiagnosticSecret(t, yamlContent, diagnosticGlobalUser, diagnosticGlobalPassword)
			}
		}
	}
}

func TestNoPegaDiagnosticSecret(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					ValuesFiles: []string{"data/values_with_no_diagnostic_user.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": getUpgradeTypeForUpgradeAction(operation),
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-diagnostic-secret.yaml"})
				var secretobj k8score.Secret
				UnmarshalK8SYaml(t, yamlContent, &secretobj)
				secretData := secretobj.Data
				require.Nil(t, secretData["PEGA_DIAGNOSTIC_USER"])
				require.Nil(t, secretData["PEGA_DIAGNOSTIC_PASSWORD"])
			}
		}
	}
}

func verifyDiagnosticSecret(t *testing.T, yamlContent string, user string, password string) {
	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, user, string(secretData["PEGA_DIAGNOSTIC_USER"]))
	require.Equal(t, password, string(secretData["PEGA_DIAGNOSTIC_PASSWORD"]))
}
