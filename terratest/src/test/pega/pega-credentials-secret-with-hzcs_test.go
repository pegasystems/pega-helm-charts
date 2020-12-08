package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
	"fmt"
)

var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
const HZ_CS_AUTH_USERNAME = "HZClusterUser"
const HZ_CS_AUTH_PASSWORD = "HZclusterPassword"

func TestPegaCredentialsSecretWhenHazelcastIsEnabled(t *testing.T){
	var supportedOperations =  []string{"deploy","install-deploy", "upgrade-deploy"}

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
					"hazelcast.username": HZ_CS_AUTH_USERNAME,
					"hazelcast.password": HZ_CS_AUTH_PASSWORD,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
			VerifyCredentialsSecretWhenHazelcastIsEnabled(t,yamlContent)
		}
	}


}

func VerifyCredentialsSecretWhenHazelcastIsEnabled(t *testing.T, yamlContent string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["DB_USERNAME"]), "YOUR_JDBC_USERNAME")
	require.Equal(t, string(secretData["DB_PASSWORD"]), "YOUR_JDBC_PASSWORD")
	require.Equal(t, string(secretData["HZ_CS_AUTH_USERNAME"]), HZ_CS_AUTH_USERNAME)
    require.Equal(t, string(secretData["HZ_CS_AUTH_PASSWORD"]), HZ_CS_AUTH_PASSWORD)
}