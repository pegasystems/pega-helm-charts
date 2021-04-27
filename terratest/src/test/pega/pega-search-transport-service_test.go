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



func TestPegaSearchTransportService(t *testing.T){
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}
	var supportedVendors = []string{"k8s", "openshift", "eks","gke","aks","pks"}
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

                yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-transport-service.yaml"})
                VerifySearchTransportService(t, yamlContent, options)
            }
		}
	}


}
// VerifySearchTransportService - Performs the search transport service assertions deployed with the values as provided in default values.yaml
func VerifySearchTransportService(t *testing.T, yamlContent string, options *helm.Options) {
	var transportSearchServiceObj k8score.Service
    UnmarshalK8SYaml(t, yamlContent, &transportSearchServiceObj)
    require.Equal(t, transportSearchServiceObj.ObjectMeta.Name, getObjName(options, "-search-transport"))
	require.Equal(t, transportSearchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, transportSearchServiceObj.Spec.Selector["app"], getObjName(options, "-search"))
	require.Equal(t, transportSearchServiceObj.Spec.ClusterIP, "None")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Name, "transport")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Port, int32(80))
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9300))
}