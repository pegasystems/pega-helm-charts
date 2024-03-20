package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

func TestPegaTierIngressDisabled(t *testing.T) {
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					ValuesFiles: []string{"data/values_ingress_disabled.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				_, err := RenderTemplateWithErr(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})
				require.Contains(t, err.Error(), "could not find template templates/pega-tier-ingress.yaml in chart")
			}
		}
	}

}
