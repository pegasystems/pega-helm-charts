package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
	"fmt"
)



func TestPegaCredentialsSecret(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"install","install-deploy", "upgrade", "upgrade-deploy"}

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

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
			VerifyCredentialsSecret(t,yamlContent)
		}
	}


}

func VerifyCredentialsSecret(t *testing.T, yamlContent string) {
	
	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["DB_USERNAME"]), "YOUR_JDBC_USERNAME")
	require.Equal(t, string(secretData["DB_PASSWORD"]), "YOUR_JDBC_PASSWORD")
}