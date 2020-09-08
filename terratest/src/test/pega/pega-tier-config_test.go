package pega


import(
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"testing"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"fmt"
	"strings"

)



func TestPegaTierConfig(t *testing.T){
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

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
			VerifyTierConfg(t,yamlContent, options)

		}
	}


}

// VerifyTierConfg - Performs the tier specific configuration assetions with the values as provided in default values.yaml
func VerifyTierConfg(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(yamlContent, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			pegaConfigMapData := pegaConfigMap.Data
			compareConfigMapData(t, pegaConfigMapData["prconfig.xml"], "data/expectedInstallDeployPrconfig.xml")
			compareConfigMapData(t, pegaConfigMapData["context.xml.tmpl"], "data/expectedInstallDeployContext.xml")
			compareConfigMapData(t, pegaConfigMapData["prlog4j2.xml"], "data/expectedInstallDeployPRlog4j2.xml")
		}
	}
}