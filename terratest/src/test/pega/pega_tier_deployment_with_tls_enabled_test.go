package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierDeploymentWithTlsEnabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"global.deployment.name":        depName,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.storageClassName":       "storage-class",
					},
					ValuesFiles: []string{"data/pega-tier-service-override_values.yaml"},
					SetStrValues: map[string]string{
						"service.tls.enabled": "true",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				assertWeb(t, yamlSplit[1], options)
				assertBatch(t, yamlSplit[2], options)
				assertStream(t, yamlSplit[3], options)
				assertStreamWithSorageClass(t, yamlSplit[3], options)

			}
		}
	}
}
