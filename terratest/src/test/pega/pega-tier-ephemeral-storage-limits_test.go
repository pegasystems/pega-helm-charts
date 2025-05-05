package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestPegaDeploymentEphemeralStorageLimits(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				var options = &helm.Options{
					ValuesFiles: []string{"data/values_ephemeral_storage.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}
				deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(deploymentYaml, "---")
				assertWeb(t, yamlSplit[1], options)
				assertEphemeralStorage(t, yamlSplit[1])
			}
		}
	}
}

func assertEphemeralStorage(t *testing.T, tierYaml string) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec
	require.Equal(t, "30G", pod.Containers[0].Resources.Limits.StorageEphemeral().String())
	require.Equal(t, "20G", pod.Containers[0].Resources.Requests.StorageEphemeral().String())
}
