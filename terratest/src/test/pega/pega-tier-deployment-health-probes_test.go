package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

// TestPegaDeploymentLegacyHealthProbes verifies that when newHealthProbes is not enabled,
// probes use the legacy /ping endpoint with legacy default values.
func TestPegaDeploymentLegacyHealthProbes(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":        "eks",
			"global.actions.execute": "deploy",
			"global.deployment.name": "pega",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	verifyLegacyProbes(t, &pod)
}

// TestPegaDeploymentNewHealthProbesEnabled verifies that when newHealthProbes.enabled=true
// and pegaVersion >= 26.1.1, probes use the new separated endpoints with new defaults.
func TestPegaDeploymentNewHealthProbesEnabled(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "true",
			"global.pegaVersion":           "26.2.0",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	verifyNewProbes(t, &pod)

	// Verify no fallback warning annotation
	require.Empty(t, depObj.ObjectMeta.Annotations["pega.io/health-probes-warning"])
}

// TestPegaDeploymentNewHealthProbesWithExactVersion verifies behavior with exact threshold version 26.1.1
func TestPegaDeploymentNewHealthProbesWithExactVersion(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "true",
			"global.pegaVersion":           "26.1.1",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	// Exact version 26.1.1 should activate new probes (>= check)
	verifyNewProbes(t, &pod)
}

// TestPegaDeploymentNewHealthProbesFallbackOlderVersion verifies that when newHealthProbes.enabled=true
// but pegaVersion < 26.1.1, probes fall back to legacy /ping and a warning annotation is added.
func TestPegaDeploymentNewHealthProbesFallbackOlderVersion(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "true",
			"global.pegaVersion":           "25.1.0",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	// Should fall back to legacy probes
	verifyLegacyProbes(t, &pod)

	// Should have fallback warning annotation
	require.Contains(t, depObj.ObjectMeta.Annotations["pega.io/health-probes-warning"], "Falling back to legacy /ping probes")
	require.Contains(t, depObj.ObjectMeta.Annotations["pega.io/health-probes-warning"], "25.1.0")
}

// TestPegaDeploymentNewHealthProbesFallbackNoVersion verifies fallback when pegaVersion is not set at all.
func TestPegaDeploymentNewHealthProbesFallbackNoVersion(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "true",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	// Should fall back to legacy probes when no version is set
	verifyLegacyProbes(t, &pod)

	// Should have fallback warning annotation
	require.Contains(t, depObj.ObjectMeta.Annotations["pega.io/health-probes-warning"], "Falling back to legacy /ping probes")
}

// TestPegaDeploymentNewHealthProbesVersion26_1_0 verifies that version 26.1.0 (just below threshold)
// falls back to legacy probes.
func TestPegaDeploymentNewHealthProbesVersion26_1_0(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "true",
			"global.pegaVersion":           "26.1.0",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	// 26.1.0 is below 26.1.1, should use legacy
	verifyLegacyProbes(t, &pod)

	// Should have fallback warning annotation
	require.Contains(t, depObj.ObjectMeta.Annotations["pega.io/health-probes-warning"], "Falling back to legacy /ping probes")
}

// TestPegaDeploymentNewHealthProbesAppliedToAllTiers verifies that new probes apply
// to all tiers (web and batch) when enabled globally.
func TestPegaDeploymentNewHealthProbesAppliedToAllTiers(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var webDep appsv1.Deployment
	var batchDep appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "true",
			"global.pegaVersion":           "26.2.0",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	// Web tier
	UnmarshalK8SYaml(t, yamlSplit[1], &webDep)
	verifyNewProbes(t, &webDep.Spec.Template.Spec)

	// Batch tier
	UnmarshalK8SYaml(t, yamlSplit[2], &batchDep)
	verifyNewProbes(t, &batchDep.Spec.Template.Spec)
}

// TestPegaDeploymentNewHealthProbesWithCustomOverrides verifies that user-specified
// probe values take precedence over new probe defaults.
func TestPegaDeploymentNewHealthProbesWithCustomOverrides(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":                         "eks",
			"global.actions.execute":                  "deploy",
			"global.deployment.name":                  "pega",
			"global.newHealthProbes.enabled":            "true",
			"global.pegaVersion":                      "26.2.0",
			"global.tier[0].name":                     "web",
			"global.tier[0].livenessProbe.timeoutSeconds":    "15",
			"global.tier[0].livenessProbe.periodSeconds":     "25",
			"global.tier[0].livenessProbe.failureThreshold":  "5",
			"global.tier[0].readinessProbe.timeoutSeconds":   "8",
			"global.tier[0].readinessProbe.periodSeconds":    "15",
			"global.tier[0].startupProbe.timeoutSeconds":     "12",
			"global.tier[0].startupProbe.failureThreshold":   "45",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	// Paths should still be the new separated endpoints
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/liveness", pod.Containers[0].LivenessProbe.HTTPGet.Path)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/readiness", pod.Containers[0].ReadinessProbe.HTTPGet.Path)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/startup", pod.Containers[0].StartupProbe.HTTPGet.Path)

	// Custom values should override the new defaults
	require.Equal(t, int32(15), pod.Containers[0].LivenessProbe.TimeoutSeconds)
	require.Equal(t, int32(25), pod.Containers[0].LivenessProbe.PeriodSeconds)
	require.Equal(t, int32(5), pod.Containers[0].LivenessProbe.FailureThreshold)
	require.Equal(t, int32(8), pod.Containers[0].ReadinessProbe.TimeoutSeconds)
	require.Equal(t, int32(15), pod.Containers[0].ReadinessProbe.PeriodSeconds)
	require.Equal(t, int32(12), pod.Containers[0].StartupProbe.TimeoutSeconds)
	require.Equal(t, int32(45), pod.Containers[0].StartupProbe.FailureThreshold)
}

// TestPegaDeploymentNewHealthProbesDisabledExplicitly verifies that when explicitly set to false,
// legacy probes are used even with a compatible version.
func TestPegaDeploymentNewHealthProbesDisabledExplicitly(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	options := &helm.Options{
		SetValues: map[string]string{
			"global.provider":              "eks",
			"global.actions.execute":       "deploy",
			"global.deployment.name":       "pega",
			"global.newHealthProbes.enabled": "false",
			"global.pegaVersion":           "26.2.0",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
	pod := depObj.Spec.Template.Spec

	// Should use legacy probes
	verifyLegacyProbes(t, &pod)

	// No fallback warning since it was explicitly disabled
	require.Empty(t, depObj.ObjectMeta.Annotations["pega.io/health-probes-warning"])
}

// TestPegaDeploymentNewHealthProbesFutureVersion verifies new probes work with future versions.
func TestPegaDeploymentNewHealthProbesFutureVersion(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	versions := []string{"27.0.0", "26.5.3", "30.1.0"}

	for _, version := range versions {
		options := &helm.Options{
			SetValues: map[string]string{
				"global.provider":              "eks",
				"global.actions.execute":       "deploy",
				"global.deployment.name":       "pega",
				"global.newHealthProbes.enabled": "true",
				"global.pegaVersion":           version,
			},
		}

		yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
		yamlSplit := strings.Split(yamlContent, "---")

		UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
		pod := depObj.Spec.Template.Spec

		verifyNewProbes(t, &pod)
	}
}

// verifyLegacyProbes asserts that all probes use the legacy /ping endpoint with legacy defaults.
func verifyLegacyProbes(t *testing.T, pod *k8score.PodSpec) {
	t.Helper()

	// Liveness probe - legacy
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", pod.Containers[0].LivenessProbe.HTTPGet.Path)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].LivenessProbe.HTTPGet.Scheme)
	require.Equal(t, intstr.FromInt(8081), pod.Containers[0].LivenessProbe.HTTPGet.Port)
	require.Equal(t, int32(0), pod.Containers[0].LivenessProbe.InitialDelaySeconds)
	require.Equal(t, int32(20), pod.Containers[0].LivenessProbe.TimeoutSeconds)
	require.Equal(t, int32(30), pod.Containers[0].LivenessProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].LivenessProbe.SuccessThreshold)
	require.Equal(t, int32(3), pod.Containers[0].LivenessProbe.FailureThreshold)

	// Readiness probe - legacy
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", pod.Containers[0].ReadinessProbe.HTTPGet.Path)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].ReadinessProbe.HTTPGet.Scheme)
	require.Equal(t, intstr.FromInt(8080), pod.Containers[0].ReadinessProbe.HTTPGet.Port)
	require.Equal(t, int32(0), pod.Containers[0].ReadinessProbe.InitialDelaySeconds)
	require.Equal(t, int32(10), pod.Containers[0].ReadinessProbe.TimeoutSeconds)
	require.Equal(t, int32(10), pod.Containers[0].ReadinessProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].ReadinessProbe.SuccessThreshold)
	require.Equal(t, int32(3), pod.Containers[0].ReadinessProbe.FailureThreshold)

	// Startup probe - legacy
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", pod.Containers[0].StartupProbe.HTTPGet.Path)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].StartupProbe.HTTPGet.Scheme)
	require.Equal(t, intstr.FromInt(8080), pod.Containers[0].StartupProbe.HTTPGet.Port)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.InitialDelaySeconds)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.TimeoutSeconds)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].StartupProbe.SuccessThreshold)
	require.Equal(t, int32(30), pod.Containers[0].StartupProbe.FailureThreshold)
}

// verifyNewProbes asserts that all probes use the new separated endpoints with new defaults.
func verifyNewProbes(t *testing.T, pod *k8score.PodSpec) {
	t.Helper()

	// Liveness probe - new separated endpoint
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/liveness", pod.Containers[0].LivenessProbe.HTTPGet.Path)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].LivenessProbe.HTTPGet.Scheme)
	require.Equal(t, intstr.FromInt(8081), pod.Containers[0].LivenessProbe.HTTPGet.Port)
	require.Equal(t, int32(0), pod.Containers[0].LivenessProbe.InitialDelaySeconds)
	require.Equal(t, int32(5), pod.Containers[0].LivenessProbe.TimeoutSeconds)
	require.Equal(t, int32(20), pod.Containers[0].LivenessProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].LivenessProbe.SuccessThreshold)
	require.Equal(t, int32(3), pod.Containers[0].LivenessProbe.FailureThreshold)

	// Readiness probe - new separated endpoint
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/readiness", pod.Containers[0].ReadinessProbe.HTTPGet.Path)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].ReadinessProbe.HTTPGet.Scheme)
	require.Equal(t, intstr.FromInt(8080), pod.Containers[0].ReadinessProbe.HTTPGet.Port)
	require.Equal(t, int32(0), pod.Containers[0].ReadinessProbe.InitialDelaySeconds)
	require.Equal(t, int32(5), pod.Containers[0].ReadinessProbe.TimeoutSeconds)
	require.Equal(t, int32(10), pod.Containers[0].ReadinessProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].ReadinessProbe.SuccessThreshold)
	require.Equal(t, int32(3), pod.Containers[0].ReadinessProbe.FailureThreshold)

	// Startup probe - new separated endpoint
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/startup", pod.Containers[0].StartupProbe.HTTPGet.Path)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].StartupProbe.HTTPGet.Scheme)
	require.Equal(t, intstr.FromInt(8080), pod.Containers[0].StartupProbe.HTTPGet.Port)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.InitialDelaySeconds)
	require.Equal(t, int32(5), pod.Containers[0].StartupProbe.TimeoutSeconds)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].StartupProbe.SuccessThreshold)
	require.Equal(t, int32(60), pod.Containers[0].StartupProbe.FailureThreshold)
}
