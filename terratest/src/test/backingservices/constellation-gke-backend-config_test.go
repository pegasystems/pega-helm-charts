package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/ingress-gce/pkg/apis/backendconfig/v1"
)

func TestConstellationGKEBackendConfig(t *testing.T) {

	var deploymentName string = "constellation-gke"

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":         "true",
			"constellation.deployment.name": deploymentName,
			"constellation.provider":        "gke",
			"constellation.ingress.enabled": "true",
		}),
	)

	var cllnBackendConfig v1.BackendConfig
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: deploymentName,
		Kind: "BackendConfig",
	}, &cllnBackendConfig)

	require.Equal(t, deploymentName, cllnBackendConfig.Name)
	require.Equal(t, 40, int(*cllnBackendConfig.Spec.TimeoutSec))
	require.Equal(t, 5, int(*cllnBackendConfig.Spec.HealthCheck.CheckIntervalSec))
	require.Equal(t, 1, int(*cllnBackendConfig.Spec.HealthCheck.HealthyThreshold))
	require.Equal(t, 3000, int(*cllnBackendConfig.Spec.HealthCheck.Port))
	require.Equal(t, "/c11n/v860/ping", *cllnBackendConfig.Spec.HealthCheck.RequestPath)
	require.Equal(t, 5, int(*cllnBackendConfig.Spec.HealthCheck.TimeoutSec))
	require.Equal(t, "HTTP", *cllnBackendConfig.Spec.HealthCheck.Type)
	require.Equal(t, 2, int(*cllnBackendConfig.Spec.HealthCheck.UnhealthyThreshold))
}

func TestConstellationMessagingGKEBackendConfig(t *testing.T) {

	var deploymentName string = "constellation-messaging-gke"

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation-messaging.enabled":         "true",
			"constellation-messaging.deployment.name": deploymentName,
			"constellation-messaging.provider":        "gke",
			"constellation-messaging.ingress.enabled": "true",
		}),
	)

	var cllnBackendConfig v1.BackendConfig
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: deploymentName,
		Kind: "BackendConfig",
	}, &cllnBackendConfig)

	require.Equal(t, deploymentName, cllnBackendConfig.Name)
	require.Equal(t, 40, int(*cllnBackendConfig.Spec.TimeoutSec))
	require.Equal(t, 5, int(*cllnBackendConfig.Spec.HealthCheck.CheckIntervalSec))
	require.Equal(t, 1, int(*cllnBackendConfig.Spec.HealthCheck.HealthyThreshold))
	require.Equal(t, 3000, int(*cllnBackendConfig.Spec.HealthCheck.Port))
	require.Equal(t, "/c11n-messaging/ping", *cllnBackendConfig.Spec.HealthCheck.RequestPath)
	require.Equal(t, 5, int(*cllnBackendConfig.Spec.HealthCheck.TimeoutSec))
	require.Equal(t, "HTTP", *cllnBackendConfig.Spec.HealthCheck.Type)
	require.Equal(t, 2, int(*cllnBackendConfig.Spec.HealthCheck.UnhealthyThreshold))
}
