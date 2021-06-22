package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)


func TestPegaCredentialsSecret(t *testing.T){
    var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"install","install-deploy", "upgrade", "upgrade-deploy"}
    var deploymentNames = []string{"pega","myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

            for _, depName := range deploymentNames {
                fmt.Println(vendor + "-" + operation)

                var options = &helm.Options{
                    SetValues: map[string]string{
                        "global.deployment.name": depName,
                        "global.provider":        vendor,
                        "global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
                    },
                }

                yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
                VerifyCredentialsSecret(t, yamlContent, options)
			}
		}
	}

}

func VerifyCredentialsSecret(t *testing.T, yamlContent string, options *helm.Options) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)

	require.Equal(t, secretobj.ObjectMeta.Name, getObjName(options, "-credentials-secret"))
	secretData := secretobj.Data
	require.Equal(t, "YOUR_JDBC_USERNAME", string(secretData["DB_USERNAME"]))
	require.Equal(t, "YOUR_JDBC_PASSWORD", string(secretData["DB_PASSWORD"]))
	require.Equal(t, "", string(secretData["CASSANDRA_TRUSTSTORE_PASSWORD"]))
	require.Equal(t, "", string(secretData["CASSANDRA_KEYSTORE_PASSWORD"]))
}
