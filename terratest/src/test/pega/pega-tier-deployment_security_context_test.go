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

func TestPegaTierDeploymentSecurityContextDefaults(t *testing.T) {
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

				if vendor != "openshift" {
					// for openshift and for others
					require.Equal(t, int64(0), *depObj.Spec.Template.Spec.SecurityContext.FSGroup)
					require.Equal(t, int64(9001), *depObj.Spec.Template.Spec.SecurityContext.RunAsUser)
				} else {
					require.Nil(t, depObj.Spec.Template.Spec.SecurityContext)
				}

			}
		}
	}
}

func TestPegaTierDeploymentSecurityContextAdditional(t *testing.T) {
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
						"global.provider":                                      vendor,
						"global.actions.execute":                               operation,
						"global.deployment.name":                               depName,
						"global.tier[0].securityContext.runAsNonRoot":          "true",
						"global.tier[0].securityContext.supplementalGroups[0]": "2000",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

				require.Equal(t, bool(true), *depObj.Spec.Template.Spec.SecurityContext.RunAsNonRoot)
				require.Equal(t, int64(2000), *&depObj.Spec.Template.Spec.SecurityContext.SupplementalGroups[0])

				if vendor != "openshift" {
					// for openshift and for others
					require.Equal(t, int64(0), *depObj.Spec.Template.Spec.SecurityContext.FSGroup)
					require.Equal(t, int64(9001), *depObj.Spec.Template.Spec.SecurityContext.RunAsUser)
				} else {
					require.Nil(t, depObj.Spec.Template.Spec.SecurityContext.FSGroup)
					require.Nil(t, depObj.Spec.Template.Spec.SecurityContext.RunAsUser)
				}

			}
		}
	}
}

func TestPegaTierDeploymentSecurityContextOverrideDefault(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke"}
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
						"global.provider":                        vendor,
						"global.actions.execute":                 operation,
						"global.deployment.name":                 depName,
						"global.tier[0].securityContext.fsGroup": "2",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

				if vendor != "openshift" {
					// for openshift and for others
					require.Equal(t, int64(2), *depObj.Spec.Template.Spec.SecurityContext.FSGroup)
					require.Equal(t, int64(9001), *depObj.Spec.Template.Spec.SecurityContext.RunAsUser)
				} else {
					require.Equal(t, int64(2), *depObj.Spec.Template.Spec.SecurityContext.FSGroup)
					require.Nil(t, depObj.Spec.Template.Spec.SecurityContext.RunAsUser)
				}
			}
		}
	}
}
