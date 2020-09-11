package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/ingress-gce/pkg/apis/backendconfig/v1beta1"
	"path/filepath"
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
			verifyBackendConfig(t, yamlContent)			

		}

}

func verifyBackendConfig(t *testing.T, yamlContent string) {
	var backendConfig v1beta1.BackendConfig
	UnmarshalK8SYaml(t, yamlContent, &backendConfig)

	require.Equal(t, "pega-backend-config", backendConfig.Name)
	require.Equal(t, 40, int(*backendConfig.Spec.TimeoutSec))
	require.Equal(t, 60, int(backendConfig.Spec.ConnectionDraining.DrainingTimeoutSec))
	require.Equal(t, "GENERATED_COOKIE", backendConfig.Spec.SessionAffinity.AffinityType)
	require.Equal(t, 3720, int(*backendConfig.Spec.SessionAffinity.AffinityCookieTtlSec))
}
