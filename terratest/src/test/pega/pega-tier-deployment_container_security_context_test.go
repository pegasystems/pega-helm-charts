package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestPegaTierDeploymentContainerSecurityContext(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy"}
	var deploymentNames = []string{"myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		var depObj appsv1.Deployment

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"global.deployment.name": depName,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				require.Nil(t, depObj.Spec.Template.Spec.Containers[0].SecurityContext)

			}
		}
	}
}

func TestPegaTierDeploymentSecurityContextForPegaContainer(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift"}
	var supportedOperations = []string{"deploy"}
	var deploymentNames = []string{"myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		var depObj appsv1.Deployment

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                                   vendor,
						"global.actions.execute":                            operation,
						"global.deployment.name":                            depName,
						"global.tier[0].containerSecurityContext.runAsUser": "7009",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

				require.Equal(t, int64(7009), *depObj.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser)
			}
		}
	}
}
