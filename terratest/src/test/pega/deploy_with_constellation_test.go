package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

// set action execute to install
var constellation_options = &helm.Options{
	SetValues: map[string]string{
		"global.provider":        "k8s",
		"global.actions.execute": "deploy",
		"constellation.enabled":  "true",
	},
}

// TestPegaStandardTierDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestPegaStandardTierDeploymentWithConstellation(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, constellation_options, []string{"wait-for-pegasearch", "wait-for-cassandra"})
}
