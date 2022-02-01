package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"strings"
	"testing"
)


func TestHazelcastDeployment(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"hazelcast.enabled":  "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/pega-hz-deployment.yaml"})
			VerifyHazelcastDeployment(t,yamlContent)

		}
	}
}

func VerifyHazelcastDeployment(t *testing.T, yamlContent string) {
	var statefulsetObj appsv1beta2.StatefulSet
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &statefulsetObj)
			require.Equal(t, *statefulsetObj.Spec.Replicas, int32(3))
			require.Equal(t, statefulsetObj.Spec.ServiceName, "pega-hazelcast-service")
			statefulsetSpec := statefulsetObj.Spec.Template.Spec
			require.Equal(t, "/hazelcast/health/ready", statefulsetSpec.Containers[0].LivenessProbe.HTTPGet.Path)
			require.Equal(t, intstr.FromInt(5701), statefulsetSpec.Containers[0].LivenessProbe.HTTPGet.Port)
			require.Equal(t, "/hazelcast/health/ready", statefulsetSpec.Containers[0].ReadinessProbe.HTTPGet.Path)
			require.Equal(t, intstr.FromInt(5701), statefulsetSpec.Containers[0].ReadinessProbe.HTTPGet.Port)
			require.Equal(t, "1", statefulsetSpec.Containers[0].Resources.Requests.Cpu().String())
			require.Equal(t, "1Gi", statefulsetSpec.Containers[0].Resources.Requests.Memory().String())
			require.Equal(t, statefulsetSpec.Volumes[0].Name, "logs")
            require.Equal(t, statefulsetSpec.Volumes[1].Name, "pega-volume-credentials")
            require.Equal(t, statefulsetSpec.Volumes[1].Secret.SecretName, "pega-credentials-secret")
            require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].Name, "logs")
            require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/hazelcast/logs")
            require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].Name, "pega-volume-credentials")
            require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/hazelcast/secrets")
		}
	}
}