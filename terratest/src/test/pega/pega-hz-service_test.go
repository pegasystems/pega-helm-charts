package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"strings"
	"testing"
	"fmt"
)



func TestHazelcastService(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"hazelcast.enabled":  "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/pega-hz-service.yaml"})
			VerifyHazelcastService(t,yamlContent,options)
		}
	}
}

func VerifyHazelcastService(t *testing.T, yamlContent string, options *helm.Options) {
	var hazelcastServiceObj k8score.Service
	serviceSlice := strings.Split(yamlContent, "---")
	for index, serviceInfo := range serviceSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, serviceInfo, &hazelcastServiceObj)
			require.Equal(t, "pega-hazelcast-service", hazelcastServiceObj.Name)
			require.Equal(t, "Hazelcast", hazelcastServiceObj.Spec.Selector["component"])
			require.Equal(t, "pega-hazelcast", hazelcastServiceObj.Spec.Selector["app"])
			require.Equal(t, "tcp-hzport", hazelcastServiceObj.Spec.Ports[0].Name)
			require.Equal(t, int32(5701), hazelcastServiceObj.Spec.Ports[0].Port)
			require.Equal(t, intstr.FromInt(5701), hazelcastServiceObj.Spec.Ports[0].TargetPort)
		}
	}
}