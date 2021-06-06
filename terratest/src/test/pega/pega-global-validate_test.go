package pega

import (
	"fmt"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestPegaPRConfigGlobalConfig(t *testing.T) {
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy"}
	var deploymentNames = []string{"pega","myapp-dev"}

	var custom_global_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Prconfig\" value=\"prconfig.xml\" />\n</pegarules>";
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation + "-" +depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.configurations.prconfig" : custom_global_config,

					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				verifyPRConfigGlobalData(t,yamlContent, options)
			}
		}
	}
}

func TestPegaPRlog4j2GlobalConfig(t *testing.T) {
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy"}
	var deploymentNames = []string{"pega","myapp-dev"}

	var custom_global_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Prlog4j2\" value=\"prlog4j2.xml\" />\n</pegarules>";
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation + "-" +depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.configurations.prlog4j2" : custom_global_config,

					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				verifyPRlog4j2GlobalData(t,yamlContent, options)
			}
		}
	}
}
func TestPegaWebXMLGlobalConfig(t *testing.T) {
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy"}
	var deploymentNames = []string{"pega","myapp-dev"}

	var custom_global_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Web\" value=\"webXML.xml\" />\n</pegarules>";
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation + "-" +depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.configurations.webXML" : custom_global_config,

					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				verifyWebXMLGlobalData(t,yamlContent, options)
			}
		}
	}
}
func TestPegaServerXMLGlobalConfig(t *testing.T) {
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy"}
	var deploymentNames = []string{"pega","myapp-dev"}

	var custom_global_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Server\" value=\"serverXML.xml\" />\n</pegarules>";
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation + "-" +depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.configurations.serverXML" : custom_global_config,

					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				verifyserverXMLGlobalData(t,yamlContent, options)
			}
		}
	}
}
func TestPegacontextXMLGlobalConfig(t *testing.T) {
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}
	var deploymentNames = []string{"pega","myapp-dev"}

	var custom_global_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Context\" value=\"contextXML.xml\" />\n</pegarules>";
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation + "-" +depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						//	"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.configurations.contextXML" : custom_global_config,
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				verifyContextXMLGlobalData(t,yamlContent, options)
			}
		}
	}
}
func verifyPRConfigGlobalData(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	pegaConfigMapData := envConfigMap.Data

	compareConfigMapData(t, pegaConfigMapData["prconfig.xml"], "data/expectedInstallCustomPrconfig.xml")




}

func verifyPRlog4j2GlobalData(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	pegaConfigMapData := envConfigMap.Data

	compareConfigMapData(t, pegaConfigMapData["prlog4j2.xml"], "data/expectedInstallCustomPrlog4j2.xml")



}


func verifyWebXMLGlobalData(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	pegaConfigMapData := envConfigMap.Data

	compareConfigMapData(t, pegaConfigMapData["web.xml"], "data/expectedInstallCustomWeb.xml")



}


func verifyserverXMLGlobalData(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	pegaConfigMapData := envConfigMap.Data
	fmt.Print("==========")
	fmt.Print(pegaConfigMapData["serverXML.xml"])
	fmt.Print("++++++++++")

	compareConfigMapData(t, pegaConfigMapData["server.xml"], "data/expectedInstallCustomServer.xml")



}

func verifyContextXMLGlobalData(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	pegaConfigMapData := envConfigMap.Data

	compareConfigMapData(t, pegaConfigMapData["context.xml.tmpl"], "data/expectedInstallCustomContext.xml")




}

