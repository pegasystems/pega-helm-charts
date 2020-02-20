package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/ingress-gce/pkg/apis/backendconfig/v1beta1"
	"path/filepath"
	"testing"
)

func TestPegaGkeStandardTierDeployment(t *testing.T) {
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.provider":        "gke",
			"global.actions.execute": "deploy",
		},
	}

	t.Parallel()
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, options, []string{"wait-for-pegasearch", "wait-for-cassandra"})
	verifyBackendConfig(t, helmChartPath, options)
}

func verifyBackendConfig(t *testing.T, helmChartPath string, options *helm.Options) {
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-gke-backend-config.yaml"})

	var backendConfig v1beta1.BackendConfig
	helm.UnmarshalK8SYaml(t, output, &backendConfig)

	require.Equal(t, "pega-backend-config", backendConfig.Name)
	require.Equal(t, 40, int(*backendConfig.Spec.TimeoutSec))
	require.Equal(t, 60, int(backendConfig.Spec.ConnectionDraining.DrainingTimeoutSec))
	require.Equal(t, "GENERATED_COOKIE", backendConfig.Spec.SessionAffinity.AffinityType)
	require.Equal(t, 3720, int(*backendConfig.Spec.SessionAffinity.AffinityCookieTtlSec))
}
