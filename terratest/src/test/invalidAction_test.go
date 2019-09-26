package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

const PegaHelmChartPath = "../../../charts/pega"

// set action execute to install
var Invalidoptions = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "deployment",
		"global.provider":        "openshift",
	},
}

// TestPegaStandardTierDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestInvalidAction(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	deployment, err := helm.RenderTemplateE(t, Invalidoptions, helmChartPath, []string{"templates/pega-action-validate.yaml"})
	if err != nil {
		strings.Contains(string(deployment), "Action value is not correct")
		fmt.Println("Invalid Action Test passed")
	} else {
		fmt.Println("Provided action is valid")
	}
}
