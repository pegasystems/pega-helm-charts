package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/networking/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestConstellation(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				var constellationOptions = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"constellation.enabled":         "true",
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}
				deploymentYaml := RenderTemplate(t, constellationOptions, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(deploymentYaml, "---")
				assertWeb(t, yamlSplit[1], constellationOptions)
				assertBatch(t, yamlSplit[2], constellationOptions)
				assertStream(t, yamlSplit[3], constellationOptions)
				ingressYaml := RenderTemplate(t, constellationOptions, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})
				assertPegaTierIngress(t, ingressYaml, constellationOptions)
			}
		}
	}
}

func assertPegaTierIngress(t *testing.T, ingressYaml string, options *helm.Options) {
	var tierIngressObj v1.Ingress
	UnmarshalK8SYaml(t, ingressYaml, &tierIngressObj)
	require.Equal(t, "/c11n", tierIngressObj.Spec.Rules[0].HTTP.Paths[0].Path)
	require.Equal(t, "constellation", tierIngressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
	require.Equal(t, int32(3000), tierIngressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number)
}
