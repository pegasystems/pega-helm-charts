package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
)

func TestPegaSearchDeployment(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.storageClassName":       "storage-class",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
				VerifySearchDeployment(t, yamlContent, options)
			}
		}
	}

}

func VerifySearchDeployment(t *testing.T, yamlContent string, options *helm.Options) {
	var statefulsetObj appsv1beta2.StatefulSet
	storageClassName := "storage-class"
	UnmarshalK8SYaml(t, yamlContent, &statefulsetObj)
	require.Equal(t, statefulsetObj.ObjectMeta.Name, getObjName(options, "-search"))
	require.Equal(t, *statefulsetObj.Spec.Replicas, int32(1))
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Name, "esstorage")
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.StorageClassName, &storageClassName)
	require.Equal(t, statefulsetObj.Spec.ServiceName, getObjName(options, "-search"))
	statefulsetSpec := statefulsetObj.Spec.Template.Spec
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].Name, "esstorage")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].MountPath, "/usr/share/elasticsearch/data")
	statefulsetAffinity := statefulsetObj.Spec.Template.Spec.Affinity
	require.Empty(t, statefulsetAffinity)
}

func TestPegaSearchDeploymentWithPodAffinity(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var affintiyBasePath = "pegasearch.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0]."

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.storageClassName":       "storage-class",
						affintiyBasePath + "key":        "kubernetes.io/os",
						affintiyBasePath + "operator":   "In",
						affintiyBasePath + "values[0]":  "linux",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
				VerifySearchDeploymentWithAffinity(t, yamlContent, options)
			}
		}
	}
}

func VerifySearchDeploymentWithAffinity(t *testing.T, yamlContent string, options *helm.Options) {
	var statefulsetObj appsv1beta2.StatefulSet
	storageClassName := "storage-class"
	UnmarshalK8SYaml(t, yamlContent, &statefulsetObj)
	require.Equal(t, statefulsetObj.ObjectMeta.Name, getObjName(options, "-search"))
	require.Equal(t, *statefulsetObj.Spec.Replicas, int32(1))
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Name, "esstorage")
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.StorageClassName, &storageClassName)
	require.Equal(t, statefulsetObj.Spec.ServiceName, getObjName(options, "-search"))
	statefulsetAffinity := statefulsetObj.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	require.Equal(t, "kubernetes.io/os", statefulsetAffinity.NodeSelectorTerms[0].MatchExpressions[0].Key)
	require.Equal(t, "In", string(statefulsetAffinity.NodeSelectorTerms[0].MatchExpressions[0].Operator))
	require.Equal(t, "linux", statefulsetAffinity.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
}
