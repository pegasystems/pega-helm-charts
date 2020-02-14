package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

// TestInvalidAction - Tests in valid action correctly rendering error
func TestInvalidAction(t *testing.T) {
	t.Parallel()

	// set action execute to install
	var Invalidoptions = &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "invalid-action",
			"global.provider":        "openshift",
		},
	}
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	deployment, err := helm.RenderTemplateE(t, Invalidoptions, helmChartPath, []string{"templates/pega-action-validate.yaml"})

	require.Error(t, err)
	require.Contains(t, string(deployment), "Action value is not correct")

}

// TestValidAction - Tests valid action
func TestValidAction(t *testing.T) {
	t.Parallel()
	// set action execute to install
	var Invalidoptions = &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
			"global.provider":        "openshift",
		},
	}
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	deployment, err := helm.RenderTemplateE(t, Invalidoptions, helmChartPath, []string{"templates/pega-action-validate.yaml"})
	require.NoError(t, err)
	require.NotContains(t, string(deployment), "Action value is not correct")
}
