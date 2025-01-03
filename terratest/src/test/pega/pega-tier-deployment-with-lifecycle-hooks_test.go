package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithLifecycleHooks(t *testing.T) {
	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				var options = &helm.Options{
					ValuesFiles: []string{"data/values_custom_lifecycle_hooks.yaml"},
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
				assertLifecycleHook(t, yamlSplit[1], options)
			}
		}
	}
}

func assertLifecycleHook(t *testing.T, tierYaml string, options *helm.Options) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec
	require.Equal(t, 1, len(pod.Containers))
	require.Equal(t, "pega-web-tomcat", pod.Containers[0].Name)
	require.Equal(t, "pegasystems/pega", pod.Containers[0].Image)
	require.Equal(t, []string{"sleep", "15"}, pod.Containers[0].Lifecycle.PreStop.Exec.Command)
}
