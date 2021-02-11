package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/ingress-gce/pkg/apis/backendconfig/v1beta1"
	"path/filepath"
	"strings"
	"testing"
)


func TestPegaGKEBackendConfig(t *testing.T) {
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _,operation := range supportedOperations{

			var options = &helm.Options{			
				SetValues: map[string]string{
					"global.provider":        "gke",
					"global.actions.execute": operation,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-gke-backend-config.yaml"})
			verifyBackendConfigs(t, yamlContent)
		}
}

func verifyBackendConfigs(t *testing.T, yamlContent string) {
	var backendConfig v1beta1.BackendConfig
	backendConfigSlice := strings.Split(yamlContent, "---")
	for index, backendConfigStr := range backendConfigSlice {
		if index >= 1 && index <= 2 {
			UnmarshalK8SYaml(t, backendConfigStr, &backendConfig)
			if index == 1 {
				verifyBackendConfig(t, &backendConfig, "pega-web", 8080)
			} else {
				// web and stream health check will happen on 8080 port
				verifyBackendConfig(t, &backendConfig, "pega-stream", 8080)
			}
		}
	}
}

func verifyBackendConfig(t *testing.T, backendConfig *v1beta1.BackendConfig, name string, port int) {
	require.Equal(t, name, backendConfig.Name)
	require.Equal(t, 40, int(*backendConfig.Spec.TimeoutSec))
	require.Equal(t, 60, int(backendConfig.Spec.ConnectionDraining.DrainingTimeoutSec))
	require.Equal(t, "GENERATED_COOKIE", backendConfig.Spec.SessionAffinity.AffinityType)
	require.Equal(t, 3720, int(*backendConfig.Spec.SessionAffinity.AffinityCookieTtlSec))

	require.Equal(t, 5, int(*backendConfig.Spec.HealthCheck.CheckIntervalSec))
	require.Equal(t, 1, int(*backendConfig.Spec.HealthCheck.HealthyThreshold))
	require.Equal(t, port, int(*backendConfig.Spec.HealthCheck.Port))
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", *backendConfig.Spec.HealthCheck.RequestPath)
	require.Equal(t, 5, int(*backendConfig.Spec.HealthCheck.TimeoutSec))
	require.Equal(t, "HTTP", *backendConfig.Spec.HealthCheck.Type)
	require.Equal(t, 2, int(*backendConfig.Spec.HealthCheck.UnhealthyThreshold))
}
