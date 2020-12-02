package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaHazelcastEnvironmentConfigForClient(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"hazelcast.enabled": "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaHazelcastEnvironmentConfigForClient(t,yamlContent, options)

		}
	}


}

func VerifyPegaHazelcastEnvironmentConfigForClient(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &envConfigMap)
			envConfigData := envConfigMap.Data
			require.Equal(t, envConfigData["HZ_DISCOVERY_K8S"], "true")
			require.Equal(t, envConfigData["HZ_CLIENT_MODE"], "true")
			require.Equal(t, envConfigData["HZ_CLUSTER_NAME"], "PRPC")
			require.Equal(t, envConfigData["HZ_SERVER_HOSTNAME"], "pega-hazelcast-service.default.svc.cluster.local")
		}
	}
}
