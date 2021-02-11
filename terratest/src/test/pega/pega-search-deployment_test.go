package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)



func TestPegaSearchDeployment(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			var options = &helm.Options{			
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
			VerifySearchDeployment(t,yamlContent)

		}
	}


}

func VerifySearchDeployment(t *testing.T, yamlContent string) {
	var statefulsetObj appsv1beta2.StatefulSet
	UnmarshalK8SYaml(t, yamlContent, &statefulsetObj)
	require.Equal(t, *statefulsetObj.Spec.Replicas, int32(1))
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Name, "esstorage")
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(t, statefulsetObj.Spec.ServiceName, "pega-search")
	statefulsetSpec := statefulsetObj.Spec.Template.Spec
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].Name, "esstorage")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].MountPath, "/usr/share/elasticsearch/data")
}