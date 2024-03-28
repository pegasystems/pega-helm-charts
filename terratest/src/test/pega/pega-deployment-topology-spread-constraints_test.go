package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierDeploymentWithMultiTopologySpreadConstraints(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                                                           vendor,
						"global.actions.execute":                                                    operation,
						"global.deployment.name":                                                    deploymentName,
						"installer.upgrade.upgradeType":                                             "zero-downtime",
						"global.tier[0].name":                                                       "web",
						"global.tier[0].topologySpreadConstraints[0].maxSkew":                       "1",
						"global.tier[0].topologySpreadConstraints[0].topologyKey":                   "zone",
						"global.tier[0].topologySpreadConstraints[0].whenUnsatisfiable":             "DoNotSchedule",
						"global.tier[0].topologySpreadConstraints[0].labelSelector.matchLabels.key": "web-pod",
						"global.tier[0].topologySpreadConstraints[1].maxSkew":                       "2",
						"global.tier[0].topologySpreadConstraints[1].topologyKey":                   "node",
						"global.tier[0].topologySpreadConstraints[1].whenUnsatisfiable":             "ScheduleAnyway",
						"global.tier[0].topologySpreadConstraints[1].labelSelector.matchLabels.key": "web-pod2",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.TopologySpreadConstraints
				require.Equal(t, "zone", constraints[0].TopologyKey)
				require.Equal(t, 1, int(constraints[0].MaxSkew))
				require.Equal(t, "DoNotSchedule", string(constraints[0].WhenUnsatisfiable))
				require.Equal(t, "web-pod", constraints[0].LabelSelector.MatchLabels["key"])
				require.Equal(t, "node", constraints[1].TopologyKey)
				require.Equal(t, 2, int(constraints[1].MaxSkew))
				require.Equal(t, "ScheduleAnyway", string(constraints[1].WhenUnsatisfiable))
				require.Equal(t, "web-pod2", constraints[1].LabelSelector.MatchLabels["key"])
			}
		}
	}
}

func TestPegaTierDeploymentWithSingleTopologySpreadConstraints(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                                                           vendor,
						"global.actions.execute":                                                    operation,
						"global.deployment.name":                                                    deploymentName,
						"installer.upgrade.upgradeType":                                             "zero-downtime",
						"global.tier[0].name":                                                       "web",
						"global.tier[0].topologySpreadConstraints[0].maxSkew":                       "2",
						"global.tier[0].topologySpreadConstraints[0].topologyKey":                   "zoneName",
						"global.tier[0].topologySpreadConstraints[0].whenUnsatisfiable":             "ScheduleAnyway",
						"global.tier[0].topologySpreadConstraints[0].labelSelector.matchLabels.app": "web-pod",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.TopologySpreadConstraints
				require.Equal(t, "zoneName", constraints[0].TopologyKey)
				require.Equal(t, 2, int(constraints[0].MaxSkew))
				require.Equal(t, "ScheduleAnyway", string(constraints[0].WhenUnsatisfiable))
				require.Equal(t, "web-pod", constraints[0].LabelSelector.MatchLabels["app"])
			}
		}
	}
}

func TestPegaTierDeploymentWithoutTopologySpreadConstraints(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"global.deployment.name":        deploymentName,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.tier[0].name":           "web",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.TopologySpreadConstraints
				require.Empty(t, constraints)
			}
		}
	}
}

func TestPegaSearchDeploymentWithTopologySpreadConstraints(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                                           vendor,
						"global.actions.execute":                                    operation,
						"global.deployment.name":                                    deploymentName,
						"installer.upgrade.upgradeType":                             "zero-downtime",
						"pegasearch.replicas":                                       "2",
						"pegasearch.topologySpreadConstraints[0].maxSkew":           "2",
						"pegasearch.topologySpreadConstraints[0].topologyKey":       "az-name",
						"pegasearch.topologySpreadConstraints[0].whenUnsatisfiable": "ScheduleAnyway",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.TopologySpreadConstraints
				require.Equal(t, "az-name", constraints[0].TopologyKey)
				require.Equal(t, 2, int(constraints[0].MaxSkew))
				require.Equal(t, "ScheduleAnyway", string(constraints[0].WhenUnsatisfiable))
			}
		}
	}
}

func TestPegaSearchDeploymentWithoutTopologySpreadConstraints(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var depObj appsv1.Deployment
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"global.deployment.name":        deploymentName,
						"installer.upgrade.upgradeType": "zero-downtime",
						"pegasearch.replicas":           "2",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
				constraints := depObj.Spec.Template.Spec.TopologySpreadConstraints
				require.Empty(t, constraints)
			}
		}
	}
}
