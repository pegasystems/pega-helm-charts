package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

var initContainers = []string{"wait-for-pegasearch", "wait-for-cassandra"}

func TestPegaTierDeployment(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"global.deployment.name":        depName,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.storageClassName": "storage-class",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				assertWeb(t, yamlSplit[1], options)
				assertBatch(t, yamlSplit[2], options)
				assertStream(t, yamlSplit[3], options)
				assertStreamWithSorageClass(t, yamlSplit[3], options)

			}
		}
	}
}

func assertStreamWithSorageClass(t *testing.T, streamYaml string, options *helm.Options) {
	var statefulsetObj appsv1beta2.StatefulSet
	UnmarshalK8SYaml(t, streamYaml, &statefulsetObj)
	require.Equal(t, statefulsetObj.ObjectMeta.Name, getObjName(options, "-stream"))
	storageClassName := "storage-class"
	require.Equal(t, &storageClassName, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.StorageClassName)
}

func TestPegaTierDeploymentWithFSGroup(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks"}
	customFsGroups := map[string]int64{
		"1000": 1000,
		"2000": 2000,
		"3000": 3000,
	}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var depObj appsv1.Deployment

	for _, vendor := range supportedVendors {
		for key, value := range customFsGroups {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                        vendor,
					"global.actions.execute":                 "deploy",
					"global.deployment.name":                 "pega",
					"installer.upgrade.upgradeType":          "zero-downtime",
					"global.tier[0].name":                    "web",
					"global.tier[1].name":                    "batch",
					"global.tier[2].name":                    "stream",
					"global.tier[0].securityContext.fsGroup": key, // web tier
					"global.tier[1].securityContext.fsGroup": key, // batch tier
					"global.tier[2].securityContext.fsGroup": key, // stream tier
				},
			}
			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")

			UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
			require.Equal(t, value, *depObj.Spec.Template.Spec.SecurityContext.FSGroup)
		}
	}
}

func assertStream(t *testing.T, streamYaml string, options *helm.Options) {
	var statefulsetObj appsv1beta2.StatefulSet
	UnmarshalK8SYaml(t, streamYaml, &statefulsetObj)
	require.Equal(t, statefulsetObj.ObjectMeta.Name, getObjName(options, "-stream"))
	VerifyPegaStatefulSet(t, &statefulsetObj, pegaDeployment{getObjName(options, "-stream"), initContainers, "Stream", "900"}, options)
}

func assertBatch(t *testing.T, batchYaml string, options *helm.Options) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, batchYaml, &deploymentObj)
	require.Equal(t, deploymentObj.ObjectMeta.Name, getObjName(options, "-batch"))
	VerifyPegaDeployment(t, &deploymentObj,
		pegaDeployment{getObjName(options, "-batch"), initContainers, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX", ""},
		options)

}

func assertWeb(t *testing.T, webYaml string, options *helm.Options) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, webYaml, &deploymentObj)
	require.Equal(t, deploymentObj.ObjectMeta.Name, getObjName(options, "-web"))
	VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{getObjName(options, "-web"), initContainers, "WebUser", "900"}, options)
}

// VerifyPegaStatefulSet - Performs specific Pega statefulset assertions with the values as provided in default values.yaml
func VerifyPegaStatefulSet(t *testing.T, statefulsetObj *appsv1beta2.StatefulSet, expectedStatefulset pegaDeployment, options *helm.Options) {
	require.Equal(t, getObjName(options, "-stream"), statefulsetObj.Spec.VolumeClaimTemplates[0].Name)
	require.Equal(t, k8score.PersistentVolumeAccessMode("ReadWriteOnce"), statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0])
	require.Equal(t, getObjName(options, "-stream"), statefulsetObj.Spec.ServiceName)
	statefulsetSpec := statefulsetObj.Spec.Template.Spec
	require.Equal(t, getObjName(options, "-stream"), statefulsetSpec.Containers[0].VolumeMounts[1].Name)
	require.Equal(t, "/opt/pega/kafkadata", statefulsetSpec.Containers[0].VolumeMounts[1].MountPath)
	require.Equal(t, "pega-volume-credentials", statefulsetSpec.Containers[0].VolumeMounts[2].Name)
	require.Equal(t, "/opt/pega/secrets", statefulsetSpec.Containers[0].VolumeMounts[2].MountPath)
	VerifyDeployment(t, &statefulsetSpec, expectedStatefulset, options)
}

// VerifyPegaDeployment - Performs specific Pega deployment assertions with the values as provided in default values.yaml
func VerifyPegaDeployment(t *testing.T, deploymentObj *appsv1.Deployment, expectedDeployment pegaDeployment, options *helm.Options) {
	require.Equal(t, int32(1), *deploymentObj.Spec.Replicas)
	require.Equal(t, int32(2147483647), *deploymentObj.Spec.ProgressDeadlineSeconds)
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Selector.MatchLabels["app"])
	require.Equal(t, intstr.FromInt(1), *deploymentObj.Spec.Strategy.RollingUpdate.MaxSurge)
	require.Equal(t, intstr.FromInt(0), *deploymentObj.Spec.Strategy.RollingUpdate.MaxUnavailable)
	require.Equal(t, appsv1.DeploymentStrategyType("RollingUpdate"), deploymentObj.Spec.Strategy.Type)
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Template.Labels["app"])
	require.NotEmpty(t, deploymentObj.Spec.Template.Annotations["config-check"])
	require.NotEmpty(t, deploymentObj.Spec.Template.Annotations["config-tier-check"])
	deploymentSpec := deploymentObj.Spec.Template.Spec
	VerifyDeployment(t, &deploymentSpec, expectedDeployment, options)
}

// VerifyDeployment - Performs common pega deployment/statefulset assertions with the values as provided in default values.yaml
func VerifyDeployment(t *testing.T, pod *k8score.PodSpec, expectedSpec pegaDeployment, options *helm.Options) {
	require.Equal(t, "pega-volume-config", pod.Volumes[0].Name)
	require.Equal(t, expectedSpec.name, pod.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	require.Equal(t, volumeDefaultModePtr, pod.Volumes[0].VolumeSource.ConfigMap.DefaultMode)
	require.Equal(t, "pega-volume-credentials", pod.Volumes[1].Name)
	require.Equal(t, getObjName(options, "-credentials-secret"), pod.Volumes[1].VolumeSource.Projected.Sources[0].Secret.Name)
	require.Equal(t, volumeDefaultModePtr, pod.Volumes[1].VolumeSource.Projected.DefaultMode)

	actualInitContainers := pod.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	//require.Equal(t, expectedSpec.initContainers, actualInitContainerNames) NEED TO CHANGE FOR "install-deploy"
	VerifyInitContainerData(t, actualInitContainers, options)
	require.Equal(t, "pega-web-tomcat", pod.Containers[0].Name)
	require.Equal(t, "pegasystems/pega", pod.Containers[0].Image)
	require.Equal(t, "pega-web-port", pod.Containers[0].Ports[0].Name)
	require.Equal(t, int32(8080), pod.Containers[0].Ports[0].ContainerPort)
	var envIndex int32 = 0
	require.Equal(t, "NODE_TYPE", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, expectedSpec.nodeType, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "PEGA_APP_CONTEXT_PATH", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "prweb", pod.Containers[0].Env[envIndex].Value)
	if expectedSpec.name == getObjName(options, "-web") || expectedSpec.name == getObjName(options, "-stream") {
		envIndex++
		require.Equal(t, "REQUESTOR_PASSIVATION_TIMEOUT", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, expectedSpec.passivationTimeout, pod.Containers[0].Env[envIndex].Value)
	}
	if options.SetValues["constellation.enabled"] == "true" {
		envIndex++
		require.Equal(t, "COSMOS_SETTINGS", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "Pega-UIEngine/cosmosservicesURI=/c11n", pod.Containers[0].Env[envIndex].Value)
	}
	envIndex++
	require.Equal(t, "JAVA_OPTS", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "", pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "CATALINA_OPTS", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "", pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "INITIAL_HEAP", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "4096m", pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "MAX_HEAP", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "8192m", pod.Containers[0].Env[envIndex].Value)
	require.Equal(t, getObjName(options, "-environment-config"), pod.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name)
	require.Equal(t, "4", pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "12Gi", pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, "3", pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "12Gi", pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, "pega-volume-config", pod.Containers[0].VolumeMounts[0].Name)
	require.Equal(t, "/opt/pega/config", pod.Containers[0].VolumeMounts[0].MountPath)

	//If these tests start failing, helm version in use is compiled against K8s version < 1.18
	//https://helm.sh/docs/topics/version_skew/#supported-version-skew
	require.Equal(t, int32(0), pod.Containers[0].LivenessProbe.InitialDelaySeconds)
	require.Equal(t, int32(20), pod.Containers[0].LivenessProbe.TimeoutSeconds)
	require.Equal(t, int32(30), pod.Containers[0].LivenessProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].LivenessProbe.SuccessThreshold)
	require.Equal(t, int32(3), pod.Containers[0].LivenessProbe.FailureThreshold)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", pod.Containers[0].LivenessProbe.HTTPGet.Path)
	require.Equal(t, intstr.FromInt(8081), pod.Containers[0].LivenessProbe.HTTPGet.Port)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].LivenessProbe.HTTPGet.Scheme)

	require.Equal(t, int32(0), pod.Containers[0].ReadinessProbe.InitialDelaySeconds)
	require.Equal(t, int32(10), pod.Containers[0].ReadinessProbe.TimeoutSeconds)
	require.Equal(t, int32(10), pod.Containers[0].ReadinessProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].ReadinessProbe.SuccessThreshold)
	require.Equal(t, int32(3), pod.Containers[0].ReadinessProbe.FailureThreshold)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", pod.Containers[0].ReadinessProbe.HTTPGet.Path)
	require.Equal(t, intstr.FromInt(8080), pod.Containers[0].ReadinessProbe.HTTPGet.Port)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].ReadinessProbe.HTTPGet.Scheme)

	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.InitialDelaySeconds)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.TimeoutSeconds)
	require.Equal(t, int32(10), pod.Containers[0].StartupProbe.PeriodSeconds)
	require.Equal(t, int32(1), pod.Containers[0].StartupProbe.SuccessThreshold)
	require.Equal(t, int32(30), pod.Containers[0].StartupProbe.FailureThreshold)
	require.Equal(t, "/prweb/PRRestService/monitor/pingService/ping", pod.Containers[0].StartupProbe.HTTPGet.Path)
	require.Equal(t, intstr.FromInt(8080), pod.Containers[0].StartupProbe.HTTPGet.Port)
	require.Equal(t, k8score.URIScheme("HTTP"), pod.Containers[0].StartupProbe.HTTPGet.Scheme)

	require.Equal(t, getObjName(options, "-registry-secret"), pod.ImagePullSecrets[0].Name)
	require.Equal(t, k8score.RestartPolicy("Always"), pod.RestartPolicy)
	require.Equal(t, int64(300), *pod.TerminationGracePeriodSeconds)
	require.Equal(t, "pega-volume-config", pod.Containers[0].VolumeMounts[0].Name)
	require.Equal(t, "/opt/pega/config", pod.Containers[0].VolumeMounts[0].MountPath)
	require.Equal(t, "pega-volume-config", pod.Volumes[0].Name)
	require.Equal(t, "pega-volume-credentials", pod.Volumes[1].Name)
	require.Equal(t, getObjName(options, "-credentials-secret"), pod.Volumes[1].Projected.Sources[0].Secret.Name)

}

type pegaDeployment struct {
	name               string
	initContainers     []string
	nodeType           string
	passivationTimeout string
}
