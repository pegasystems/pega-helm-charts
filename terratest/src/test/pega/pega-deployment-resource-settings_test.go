package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestDefaultsLatestResourceConfiguration(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		//ValuesFiles: []string{"data/values_new_resource_configuration.yaml"},
		SetValues: map[string]string{
			"global.deployment.name": "resource-settings-1",
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		},
	}

	var depObj appsv1.Deployment

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")
	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

	resources := &depObj.Spec.Template.Spec.Containers[0].Resources

	require.Equal(t, "4", resources.Limits.Cpu().String())
	require.Equal(t, "12Gi", resources.Limits.Memory().String())
	require.Equal(t, "0", resources.Limits.StorageEphemeral().String())

	require.Equal(t, "3", resources.Requests.Cpu().String())
	require.Equal(t, "12Gi", resources.Requests.Memory().String())
	require.Equal(t, "0", resources.Requests.StorageEphemeral().String())

}

func TestLegacyResouceConfiguration(t *testing.T) {

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		ValuesFiles: []string{"data/values_legacy_resource_configuration.yaml"},
		SetValues: map[string]string{
			"global.deployment.name": "resource-settings-1",
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		},
	}

	var depObj appsv1.Deployment

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")
	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

	resources := &depObj.Spec.Template.Spec.Containers[0].Resources

	require.Equal(t, "3", resources.Limits.Cpu().String())
	require.Equal(t, "6Gi", resources.Limits.Memory().String())
	require.Equal(t, "25Gi", resources.Limits.StorageEphemeral().String())
	require.Equal(t, "2", resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", resources.Requests.Memory().String())
	require.Equal(t, "25Gi", resources.Requests.StorageEphemeral().String())
}

func TestCPULimitNull(t *testing.T) {

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		SetValues: map[string]string{
			"global.deployment.name":                   "resource-settings-1",
			"global.provider":                          "k8s",
			"global.actions.execute":                   "deploy",
			"global.tier[0].name":                      "web",
			"global.tier[0].resources.requests.cpu":    "2",
			"global.tier[0].resources.requests.memory": "6Gi",
			"global.tier[0].resources.limits.memory":   "6Gi",
		},
	}

	var depObj appsv1.Deployment

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")
	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

	resources := &depObj.Spec.Template.Spec.Containers[0].Resources

	require.Equal(t, "0", resources.Limits.Cpu().String())
	require.Equal(t, "6Gi", resources.Limits.Memory().String())
	require.Equal(t, "2", resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", resources.Requests.Memory().String())
}

func TestResourceDefaultsWithLegacyValues(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		SetValues: map[string]string{
			"global.deployment.name": "resource-settings-1",
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
			// default values.yaml configuration is replaced with below for tier[0]. Hence neither resourcse nor deprecated
			// resource configuration is set. In the deployment values are set with defaults from code
			"global.tier[0].name": "web",
		},
	}

	var depObj appsv1.Deployment

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")
	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)

	resources := &depObj.Spec.Template.Spec.Containers[0].Resources

	require.Equal(t, "4", resources.Limits.Cpu().String())
	require.Equal(t, "12Gi", resources.Limits.Memory().String())
	require.Equal(t, "0", resources.Limits.StorageEphemeral().String())

	require.Equal(t, "3", resources.Requests.Cpu().String())
	require.Equal(t, "12Gi", resources.Requests.Memory().String())
	require.Equal(t, "0", resources.Requests.StorageEphemeral().String())
}
