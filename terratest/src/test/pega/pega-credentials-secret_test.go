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

func TestPegaCredentialsSecretWithArtifactoryBasicAuth(t *testing.T){
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
						"global.customArtifactory.authentication.basic.username": "username",
						"global.customArtifactory.authentication.basic.password": "pwd",
						"global.customArtifactory.authentication.apiKey.headerName": "",
						"global.customArtifactory.authentication.apiKey.value": "",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
				VerifyCredentialsSecretArtifactoryBasicAuth(t, yamlContent, options, true)
			}
		}
	}

}

func TestPegaCredentialsSecretWithNoArtifactoryBasicAuth(t *testing.T){
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
						"global.customArtifactory.authentication.basic.username": "",
						"global.customArtifactory.authentication.basic.password": "pwd",
						"global.customArtifactory.authentication.apiKey.headerName": "",
						"global.customArtifactory.authentication.apiKey.value": "",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
				VerifyCredentialsSecretArtifactoryBasicAuth(t, yamlContent, options, false)
			}
		}
	}

}

func TestPegaCredentialsSecretWithArtifactoryApiKeyAuth(t *testing.T){
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
						"global.customArtifactory.authentication.apiKey.headerName": "headerName",
						"global.customArtifactory.authentication.apiKey.value": "apiKey-value",
						"global.customArtifactory.authentication.basic.username": "",
						"global.customArtifactory.authentication.basic.password": "",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
				VerifyCredentialsSecretArtifactoryApiKeyAuth(t, yamlContent, options, true)
			}
		}
	}

}

func TestPegaCredentialsSecretWithNoArtifactoryApiKeyAuth(t *testing.T){
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
						"global.customArtifactory.authentication.apiKey.headerName": "headerName",
						"global.customArtifactory.authentication.apiKey.value": "",
						"global.customArtifactory.authentication.basic.username": "",
						"global.customArtifactory.authentication.basic.password": "",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
				VerifyCredentialsSecretArtifactoryApiKeyAuth(t, yamlContent, options, false)
			}
		}
	}

}

func VerifyCredentialsSecretArtifactoryBasicAuth(t *testing.T, yamlContent string, options *helm.Options, expectBasicAuthData bool) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)

	require.Equal(t, secretobj.ObjectMeta.Name, getObjName(options, "-credentials-secret"))
	secretData := secretobj.Data
	if expectBasicAuthData == true {
		require.Equal(t, "username", string(secretData["CUSTOM_ARTIFACTORY_USERNAME"]))
		require.Equal(t, "pwd", string(secretData["CUSTOM_ARTIFACTORY_PASSWORD"]))
	} else {
		require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_USERNAME"]))
		require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_PASSWORD"]))
	}

	require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_APIKEY_HEADER"]))
	require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_APIKEY"]))
}

func VerifyCredentialsSecretArtifactoryApiKeyAuth(t *testing.T, yamlContent string, options *helm.Options, expectApiKeyAuthData bool) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)

	require.Equal(t, secretobj.ObjectMeta.Name, getObjName(options, "-credentials-secret"))
	secretData := secretobj.Data
	if expectApiKeyAuthData == true {
		require.Equal(t, "headerName", string(secretData["CUSTOM_ARTIFACTORY_APIKEY_HEADER"]))
		require.Equal(t, "apiKey-value", string(secretData["CUSTOM_ARTIFACTORY_APIKEY"]))
	} else {
		require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_APIKEY_HEADER"]))
		require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_APIKEY"]))
	}

	require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_USERNAME"]))
	require.Empty(t, string(secretData["CUSTOM_ARTIFACTORY_PASSWORD"]))
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
