package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierDeploymentWithTolerations(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                        vendor,
						"global.actions.execute":                 operation,
						"global.deployment.name":                 deploymentName,
						"global.tier[0].name":                    "web",
						"global.tier[0].tolerations[0].key":      "availability-zone",
						"global.tier[0].tolerations[0].value":    "us-east-1",
						"global.tier[0].tolerations[0].operator": "Equal",
						"global.tier[0].tolerations[0].effect":   "NotSchedule",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.Tolerations
				require.Equal(t, "availability-zone", constraints[0].Key)
				require.Equal(t, "us-east-1", constraints[0].Value)
				require.Equal(t, "Equal", string(constraints[0].Operator))
				require.Equal(t, "NotSchedule", string(constraints[0].Effect))
			}
		}
	}
}

func TestPegaTierDeploymentWithMultipleTolerations(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                        vendor,
						"global.actions.execute":                 operation,
						"global.deployment.name":                 deploymentName,
						"global.tier[0].name":                    "web",
						"global.tier[0].tolerations[0].key":      "availability-zone",
						"global.tier[0].tolerations[0].value":    "us-east-1",
						"global.tier[0].tolerations[0].operator": "Equal",
						"global.tier[0].tolerations[0].effect":   "NotSchedule",
						"global.tier[0].tolerations[1].key":      "availability-zone",
						"global.tier[0].tolerations[1].operator": "Exists",
						"global.tier[0].tolerations[1].effect":   "NoExecute",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.Tolerations
				require.Equal(t, "availability-zone", constraints[0].Key)
				require.Equal(t, "us-east-1", constraints[0].Value)
				require.Equal(t, "Equal", string(constraints[0].Operator))
				require.Equal(t, "NotSchedule", string(constraints[0].Effect))
				require.Equal(t, "availability-zone", constraints[1].Key)
				require.Equal(t, "Exists", string(constraints[1].Operator))
				require.Equal(t, "NoExecute", string(constraints[1].Effect))
			}
		}
	}
}

func TestPegaTierDeploymentWithoutTolerations(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"global.deployment.name": deploymentName,
						"global.tier[0].name":    "web",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.Tolerations
				require.Empty(t, constraints)
			}
		}
	}
}
