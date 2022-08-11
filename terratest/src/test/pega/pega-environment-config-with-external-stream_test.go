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

func TestPegaExternalStreamEnvironmentConfig(t *testing.T){
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
					"stream.enabled": "true",
					"stream.bootstrapServer": "localhost:9092",
					"stream.securityProtocol": "PLAINTEXT",
					"stream.saslMechanism": "PLAIN",
					"stream.streamNamePattern": "pega-{stream.name}",
					"stream.replicationFactor": "1",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalStreamEnvironmentConfig(t,yamlContent, options)

		}
	}


}

func VerifyPegaWithExternalStreamEnvironmentConfig(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &envConfigMap)
			envConfigData := envConfigMap.Data
			require.Equal(t, envConfigData["EXTERNAL_STREAM"], "true")
			require.Equal(t, envConfigData["STREAM_BOOTSTRAP_SERVERS"], "localhost:9092")
			require.Equal(t, envConfigData["STREAM_SECURITY_PROTOCOL"], "PLAINTEXT")
			require.Equal(t, envConfigData["STREAM_SASL_MECHANISM"], "PLAIN")
			require.Equal(t, envConfigData["STREAM_TRUSTSTORE_TYPE"], "")
			require.Equal(t, envConfigData["STREAM_KEYSTORE_TYPE"], "")
			require.Equal(t, envConfigData["STREAM_NAME_PATTERN"], "pega-{stream.name}")
			require.Equal(t, envConfigData["STREAM_REPLICATION_FACTOR"], "1")
		}
	}
}
