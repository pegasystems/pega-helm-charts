package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
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
}
