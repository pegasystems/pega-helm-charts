package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/ingress-gce/pkg/apis/backendconfig/v1beta1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaGKECustomBackendConfig(t *testing.T) {
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, operation := range supportedOperations {

		for _, depName := range deploymentNames {

			var options = &helm.Options{
				//ValuesFiles: []string{"data/values_gke_backend_config.yaml"},
				SetValues: map[string]string{
					"global.deployment.name":        depName,
					"global.provider":               "gke",
					"global.actions.execute":        operation,
					"installer.upgrade.upgradeType": "zero-downtime",
					"global.tier[0].backendConfig.connectionDraining.drainingTimeoutSec": "80",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-gke-backend-config.yaml"})
			verifyCustomBackendConfigs(t, yamlContent, options)
		}
	}
}

func verifyCustomBackendConfigs(t *testing.T, yamlContent string, options *helm.Options) {
	var backendConfig v1beta1.BackendConfig
	backendConfigSlice := strings.Split(yamlContent, "---")
	for index, backendConfigStr := range backendConfigSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, backendConfigStr, &backendConfig)
			verifyCustomTimeouts(t, &backendConfig, getObjName(options, "-web"), 8080)
		}
	}
}

func verifyCustomTimeouts(t *testing.T, backendConfig *v1beta1.BackendConfig, name string, port int) {
	require.Equal(t, 60, int(backendConfig.Spec.ConnectionDraining.DrainingTimeoutSec))
}
