package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
	"fmt"
)



func TestPegaEnvironmentConfig(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{			
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyEnvironmentConfig(t,yamlContent, options)

		}
	}


}

func VerifyEnvironmentConfig(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	envConfigData := envConfigMap.Data
	require.Equal(t, envConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, envConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, envConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, envConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	if options.SetValues["global.actions.execute"] == "upgrade-deploy" {
		require.Equal(t, envConfigData["RULES_SCHEMA"], "")
	} else {
		require.Equal(t, envConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	}
	require.Equal(t, envConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, envConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, envConfigData["JDBC_CONNECTION_PROPERTIES"], "")
	require.Equal(t, envConfigData["PEGA_SEARCH_URL"], "http://pega-search")
	require.Equal(t, envConfigData["CASSANDRA_CLUSTER"], "true")
	require.Equal(t, envConfigData["CASSANDRA_NODES"], "pega-cassandra")
	require.Equal(t, envConfigData["CASSANDRA_PORT"], "9042")
}