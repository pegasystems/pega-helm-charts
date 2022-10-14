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
					"stream.trustStore": "truststore.jks",
					"stream.keyStore": "keystore.jks",
					"stream.streamNamePattern": "pega-{stream.name}",
					"stream.replicationFactor": "1",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalStreamEnvironmentConfig(t,yamlContent, "/opt/pega/certs/truststore.jks", "/opt/pega/certs/keystore.jks", "","", options)

		}
	}
}

func VerifyPegaWithExternalStreamEnvironmentConfig(t *testing.T, yamlContent string, truststore string, keystore string, jksCertType string, keyCertType string, options *helm.Options) {

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
            require.Equal(t, envConfigData["STREAM_TRUSTSTORE"], truststore)
            require.Equal(t, envConfigData["STREAM_KEYSTORE"], keystore)
			require.Equal(t, envConfigData["STREAM_TRUSTSTORE_TYPE"], jksCertType)
			require.Equal(t, envConfigData["STREAM_KEYSTORE_TYPE"], keyCertType)
			require.Equal(t, envConfigData["STREAM_NAME_PATTERN"], "pega-{stream.name}")
			require.Equal(t, envConfigData["STREAM_REPLICATION_FACTOR"], "1")
		}
	}
}

func TestPegaExternalStreamEnvironmentConfigWithoutSSL(t *testing.T){
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
			VerifyPegaWithExternalStreamEnvironmentConfig(t,yamlContent, "", "", "","", options)

		}
	}
}

func TestPegaExternalStreamEnvironmentConfigWithPEM(t *testing.T){
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
					"stream.trustStore": "truststore.pem",
					"stream.trustStoreType": "PEM",
					"stream.keyStore": "keystore.pem",
					"stream.keyStoreType": "PEM",
					"stream.streamNamePattern": "pega-{stream.name}",
					"stream.replicationFactor": "1",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalStreamEnvironmentConfig(t,yamlContent, "/opt/pega/certs/truststore.pem", "/opt/pega/certs/keystore.pem", "PEM", "PEM", options)

		}
	}
}
