package test

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

const PegaHelmChartPath = "../../../charts/pega"

// set action execute to install
var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "deploy",
		"global.provider":        "openshift",
	},
}

// TestPegaStandardTierDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestOpenshiftPegaTierDeployment(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, options, []string{"wait-for-pegasearch", "wait-for-cassandra"})
}
