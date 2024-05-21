package backingservices

import (
	"github.com/stretchr/testify/require"
	"k8s.io/ingress-gce/pkg/apis/backendconfig/v1"
	"testing"
)

func TestConstellationGKEBackendConfig(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled": 			"true",
			"constellation.deployment.name":    "clln-gke",
			"constellation.provider":      		"gke",
		}),
	)

	var cllnBackendConfig v1.BackendConfig
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "clln-gke",
		Kind: "BackendConfig",
	}, &cllnBackendConfig)

	require.Equal(t, 40, int64(*cllnBackendConfig.Spec.TimeoutSec))

}

/*
func VerifySRSDeployment(t *testing.T, deploymentObj appsv1.Deployment, expectedDeployment srsDeployment) {
	require.Equal(t, expectedDeployment.replicaCount, *deploymentObj.Spec.Replicas )
	require.Equal(t, expectedDeployment.appName, deploymentObj.Spec.Selector.MatchLabels["app.kubernetes.io/name"])
	if expectedDeployment.internetEgress {
		require.Equal(t, "true", deploymentObj.Spec.Selector.MatchLabels["networking/allow-internet-egress"])
	}
	require.Equal(t, expectedDeployment.appName, deploymentObj.Spec.Template.Labels["app.kubernetes.io/name"])
	deploymentSpec := deploymentObj.Spec.Template.Spec
	VerifyDeployment(t, &deploymentSpec, expectedDeployment)
}

func verifyBackendConfigs(t *testing.T, yamlContent string, options *helm.Options) {
	var backendConfig v1.BackendConfig
	backendConfigSlice := strings.Split(yamlContent, "---")
	for index, backendConfigStr := range backendConfigSlice {
		if index >= 1 && index <= 2 {
			UnmarshalK8SYaml(t, backendConfigStr, &backendConfig)
			if index == 1 {
				verifyBackendConfig(t, &backendConfig, getObjName(options, "-web"), 8080)
			} else {
				// web and stream health check will happen on 8080 port
				verifyBackendConfig(t, &backendConfig, getObjName(options, "-stream"), 8080)
			}
		}
	}
}

func verifyBackendConfig(t *testing.T, backendConfig *v1.BackendConfig, name string, port int) {
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

func TestPegaGKECustomBackendConfig(t *testing.T) {
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	//testsPath, err := filepath.Abs(PegaHelmChartTestsPath)
	require.NoError(t, err)

	for _, operation := range supportedOperations {

		for _, depName := range deploymentNames {

			var options = &helm.Options{
				ValuesFiles: []string{"data/values_backend_config_gke.yaml"},
				SetValues: map[string]string{
					"global.deployment.name":        depName,
					"global.provider":               "gke",
					"global.actions.execute":        operation,
					"global.jdbc.url":               "true",
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-gke-backend-config.yaml"})

			verifyCustomBackendConfigs(t, yamlContent, options)
		}
	}
}

func verifyCustomBackendConfigs(t *testing.T, yamlContent string, options *helm.Options) {
	var backendConfig v1.BackendConfig
	backendConfigSlice := strings.Split(yamlContent, "---")
	for index, backendConfigStr := range backendConfigSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, backendConfigStr, &backendConfig)
			verifyCustomTimeouts(t, &backendConfig)
		}
	}
}

func verifyCustomTimeouts(t *testing.T, backendConfig *v1.BackendConfig) {
	require.Equal(t, 100, int(backendConfig.Spec.ConnectionDraining.DrainingTimeoutSec))
}
*/