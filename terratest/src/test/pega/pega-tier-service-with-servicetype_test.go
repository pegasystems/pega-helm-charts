package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaServiceWithServiceType(t *testing.T) {

	var supportedVendors = []string{"openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)
				var options = &helm.Options{
					ValuesFiles: []string{"data/values_with_servicetype.yaml"},
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
				serviceyamlContent := strings.Split(yamlContent, "---")
				var pegaServiceObj k8score.Service
				UnmarshalK8SYaml(t, serviceyamlContent[1], &pegaServiceObj)
				serviceType := pegaServiceObj.Spec.Type
				require.Equal(t, k8score.ServiceType("LoadBalancer"), serviceType)
			}
		}
	}
}
