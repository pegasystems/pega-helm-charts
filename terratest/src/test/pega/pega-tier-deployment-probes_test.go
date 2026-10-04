package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func renderWebTierProbes(t *testing.T, setValues map[string]string) appsv1.Deployment {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	values := map[string]string{
		"global.provider":        "k8s",
		"global.actions.execute": "deploy",
		"global.deployment.name": "pega",
		"global.tier[0].name":    "web",
	}
	for k, v := range setValues {
		values[k] = v
	}

	options := &helm.Options{SetValues: values}
	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	var depObj appsv1.Deployment
	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	return depObj
}

func TestPegaTierDeploymentProbesDefaultProfile(t *testing.T) {
	depObj := renderWebTierProbes(t, nil)
	container := depObj.Spec.Template.Spec.Containers[0]
	expectedPath := "/prweb/PRRestService/monitor/pingService/ping"

	require.Equal(t, expectedPath, container.LivenessProbe.HTTPGet.Path)
	require.Equal(t, expectedPath, container.ReadinessProbe.HTTPGet.Path)
	require.NotNil(t, container.StartupProbe)
	require.Equal(t, expectedPath, container.StartupProbe.HTTPGet.Path)
	require.Equal(t, intstr.FromInt(8080), container.ReadinessProbe.HTTPGet.Port)
}

func TestPegaTierDeploymentProbesStandardProfile(t *testing.T) {
	depObj := renderWebTierProbes(t, map[string]string{"global.probes.profile": "standard"})
	container := depObj.Spec.Template.Spec.Containers[0]
	expectedPath := "/prweb/PRRestService/monitor/pingService/ping"

	require.Equal(t, expectedPath, container.LivenessProbe.HTTPGet.Path)
	require.Equal(t, expectedPath, container.ReadinessProbe.HTTPGet.Path)
	require.Equal(t, expectedPath, container.StartupProbe.HTTPGet.Path)
	require.Equal(t, intstr.FromInt(8080), container.ReadinessProbe.HTTPGet.Port)
}

func TestPegaTierDeploymentProbesEnhancedProfile(t *testing.T) {
	depObj := renderWebTierProbes(t, map[string]string{"global.probes.profile": "enhanced"})
	container := depObj.Spec.Template.Spec.Containers[0]

	require.Equal(t, "/prweb/PRRestService/monitor/pingService/liveness", container.LivenessProbe.HTTPGet.Path)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/readiness", container.ReadinessProbe.HTTPGet.Path)
	require.NotNil(t, container.StartupProbe)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/startup", container.StartupProbe.HTTPGet.Path)

	// readiness defaults to 8081 with the enhanced profile
	require.Equal(t, intstr.FromInt(8081), container.ReadinessProbe.HTTPGet.Port)
}

func TestPegaTierDeploymentProbesEnhancedProfileReadinessPortOverride(t *testing.T) {
	depObj := renderWebTierProbes(t, map[string]string{
		"global.probes.profile":              "enhanced",
		"global.tier[0].readinessProbe.port": "9090",
	})
	container := depObj.Spec.Template.Spec.Containers[0]

	require.Equal(t, intstr.FromInt(9090), container.ReadinessProbe.HTTPGet.Port)
}

func TestPegaTierDeploymentProbesInvalidProfile(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
			"global.deployment.name": "pega",
			"global.tier[0].name":    "web",
			"global.probes.profile":  "bogus",
		},
	}

	_, err = helm.RenderTemplateE(t, options, helmChartPath, "pega", []string{"templates/pega-tier-deployment.yaml"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "global.probes.profile must be either 'standard' or 'enhanced'")
}

