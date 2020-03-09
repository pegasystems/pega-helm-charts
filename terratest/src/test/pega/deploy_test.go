package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPegaStandardTierDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestPegaStandardTierDeployment2(t *testing.T) {
	// set action execute to install
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		},
	}

	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, options, []string{"wait-for-pegasearch", "wait-for-cassandra"})
}
