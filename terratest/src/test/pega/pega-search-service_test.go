package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"testing"
	"fmt"
)



func TestPegaSearchService(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}
    var deploymentNames = []string{"pega","myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

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

                yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-service.yaml"})
                VerifySearchService(t, yamlContent, options)
            }
		}
	}


}

func VerifySearchService(t *testing.T, yamlContent string, options *helm.Options) {
	var searchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, yamlContent, &searchServiceObj)
	require.Equal(t, searchServiceObj.ObjectMeta.Name, getObjName(options, "-search"))
	require.Equal(t, searchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, searchServiceObj.Spec.Selector["app"], getObjName(options, "-search"))
	require.Equal(t, searchServiceObj.Spec.Ports[0].Name, "http")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Port, int32(80))
	require.Equal(t, searchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9200))
}