package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaWebTierService(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					ValuesFiles: []string{"data/values_http_disabled.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
				VerifyPegaWebServices(t, yamlContent, options)
			}
		}
	}

}

// SplitAndVerifyPegaServices - Splits the services from the rendered template and asserts each service objects
func VerifyPegaWebServices(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaServiceObj k8score.Service
	serviceSlice := strings.Split(yamlContent, "---")
	for index, serviceInfo := range serviceSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, serviceInfo, &pegaServiceObj)

			if index == 1 {
				require.Equal(t, getObjName(options, "-web"), pegaServiceObj.ObjectMeta.Name)
				VerifyPegaWebService(t, &pegaServiceObj, pegaServices{getObjName(options, "-web"), int32(443), intstr.IntOrString{IntVal: 8443}})
			}
		}
	}
}

// VerifyPegaService - Performs Pega Service assertions with the values as provided in values_http_disabled.yaml
func VerifyPegaWebService(t *testing.T, serviceObj *k8score.Service, expectedService pegaServices) {
	require.Equal(t, serviceObj.Spec.Selector["app"], expectedService.Name)
	require.Equal(t, serviceObj.Spec.Ports[0].Port, expectedService.Port)
	require.Equal(t, serviceObj.Spec.Ports[0].TargetPort, expectedService.TargetPort)
}
