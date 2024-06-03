package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/ingress-gce/pkg/apis/backendconfig/v1"
)

func TestConstellationGKEBackendConfig(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":         "true",
			"constellation.deployment.name": "constellation-gke",
			"constellation.provider":        "gke",
			"constellation.ingress.enabled": "true",
		}),
	)

	var cllnBackendConfig v1.BackendConfig
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "constellation-gke",
		Kind: "BackendConfig",
	}, &cllnBackendConfig)

	require.Equal(t, "constellation-gke", cllnBackendConfig.Name)
	require.Equal(t, 40, int(*cllnBackendConfig.Spec.TimeoutSec))
	require.Equal(t, 5, int(*cllnBackendConfig.Spec.HealthCheck.CheckIntervalSec))
	require.Equal(t, 1, int(*cllnBackendConfig.Spec.HealthCheck.HealthyThreshold))
	require.Equal(t, 3000, int(*cllnBackendConfig.Spec.HealthCheck.Port))
	require.Equal(t, "/c11n/v860/ping", *cllnBackendConfig.Spec.HealthCheck.RequestPath)
	require.Equal(t, 5, int(*cllnBackendConfig.Spec.HealthCheck.TimeoutSec))
	require.Equal(t, "HTTP", *cllnBackendConfig.Spec.HealthCheck.Type)
	require.Equal(t, 2, int(*cllnBackendConfig.Spec.HealthCheck.UnhealthyThreshold))
}
