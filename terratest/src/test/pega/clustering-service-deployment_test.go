package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestClusteringServiceDeployment(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                    vendor,
					"global.actions.execute":             operation,
					"hazelcast.clusteringServiceEnabled": "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/clustering-service-deployment.yaml"})
			VerifyClusteringServiceDeployment(t, yamlContent)

		}
	}
}

func VerifyClusteringServiceDeployment(t *testing.T, yamlContent string) {
	var statefulsetObj appsv1beta2.StatefulSet
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &statefulsetObj)
			require.Equal(t, *statefulsetObj.Spec.Replicas, int32(3))
			require.Equal(t, statefulsetObj.Spec.ServiceName, "clusteringservice-service")
			statefulsetSpec := statefulsetObj.Spec.Template.Spec
			require.Equal(t, "/hazelcast/health/ready", statefulsetSpec.Containers[0].LivenessProbe.HTTPGet.Path)
			require.Equal(t, intstr.FromInt(5701), statefulsetSpec.Containers[0].LivenessProbe.HTTPGet.Port)
			require.Equal(t, "/hazelcast/health/ready", statefulsetSpec.Containers[0].ReadinessProbe.HTTPGet.Path)
			require.Equal(t, intstr.FromInt(5701), statefulsetSpec.Containers[0].ReadinessProbe.HTTPGet.Port)
			require.Equal(t, "1", statefulsetSpec.Containers[0].Resources.Requests.Cpu().String())
			require.Equal(t, "1Gi", statefulsetSpec.Containers[0].Resources.Requests.Memory().String())
			require.Equal(t, statefulsetSpec.Volumes[0].Name, "logs")
			require.Equal(t, statefulsetSpec.Volumes[1].Name, "hazelcast-volume-credentials")
			require.Equal(t, statefulsetSpec.Volumes[1].Projected.Sources[0].Secret.Name, "pega-hz-secret")
			require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].Name, "logs")
			require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/hazelcast/logs")
			require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].Name, "hazelcast-volume-credentials")
			require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/hazelcast/secrets")
			statefulsetAffinity := statefulsetObj.Spec.Template.Spec.Affinity
			require.Empty(t, statefulsetAffinity)
		}
	}
}

func TestClusteringServiceDeploymentSecurityContext(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                     vendor,
					"global.actions.execute":              operation,
					"hazelcast.clusteringServiceEnabled":  "true",
					"hazelcast.securityContext.runAsUser": "1000",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/clustering-service-deployment.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")
			var statefulsetObj appsv1beta2.StatefulSet
			UnmarshalK8SYaml(t, yamlSplit[1], &statefulsetObj)

			require.Equal(t, int64(1000), *statefulsetObj.Spec.Template.Spec.SecurityContext.RunAsUser)

		}
	}
}

func TestClusteringServiceDeploymentWithAffinity(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var affintiyBasePath = "hazelcast.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0]."

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                    vendor,
					"global.actions.execute":             operation,
					"hazelcast.clusteringServiceEnabled": "true",
					affintiyBasePath + "key":             "kubernetes.io/os",
					affintiyBasePath + "operator":        "In",
					affintiyBasePath + "values[0]":       "linux",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/clustering-service-deployment.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")
			var statefulsetObj appsv1beta2.StatefulSet
			UnmarshalK8SYaml(t, yamlSplit[1], &statefulsetObj)

			statefulsetAffinity := statefulsetObj.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
			require.Equal(t, "kubernetes.io/os", statefulsetAffinity.NodeSelectorTerms[0].MatchExpressions[0].Key)
			require.Equal(t, "In", string(statefulsetAffinity.NodeSelectorTerms[0].MatchExpressions[0].Operator))
			require.Equal(t, "linux", statefulsetAffinity.NodeSelectorTerms[0].MatchExpressions[0].Values[0])

		}
	}
}
