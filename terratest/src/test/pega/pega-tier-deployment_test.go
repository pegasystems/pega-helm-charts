package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"strings"
	"testing"
	"fmt"
)


var initContainers = []string{"wait-for-pegasearch", "wait-for-cassandra"}

func TestPegaTierDeployment(t *testing.T){
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

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")			
			assertWeb(t,yamlSplit[1],options)
			assertBatch(t,yamlSplit[2],options)
			assertStream(t,yamlSplit[3],options)

		}
	}
}

func assertStream(t *testing.T, streamYaml string, options *helm.Options){
	var statefulsetObj appsv1beta2.StatefulSet
	UnmarshalK8SYaml(t,streamYaml,&statefulsetObj)
	VerifyPegaStatefulSet(t, &statefulsetObj, pegaDeployment{"pega-stream", initContainers, "Stream", "900"},options)
}



func assertBatch(t *testing.T, batchYaml string, options *helm.Options){
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t,batchYaml,&deploymentObj)
	VerifyPegaDeployment(t, &deploymentObj,
		pegaDeployment{"pega-batch", initContainers, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX", ""},
		options)
	
}

func assertWeb(t *testing.T, webYaml string, options *helm.Options){
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t,webYaml,&deploymentObj)
	VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-web", initContainers, "WebUser", "900"}, options)
	
	
}

// VerifyPegaStatefulSet - Performs specific Pega statefulset assertions with the values as provided in default values.yaml
func VerifyPegaStatefulSet(t *testing.T, statefulsetObj *appsv1beta2.StatefulSet, expectedStatefulset pegaDeployment, options *helm.Options) {
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Name, "pega-stream")
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(t, statefulsetObj.Spec.ServiceName, "pega-stream")
	statefulsetSpec := statefulsetObj.Spec.Template.Spec
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].Name, "pega-stream")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/streamvol")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[2].Name, "pega-volume-credentials")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[2].MountPath, "/opt/pega/secrets")
	VerifyDeployment(t, &statefulsetSpec, expectedStatefulset, options)
}


// VerifyPegaDeployment - Performs specific Pega deployment assertions with the values as provided in default values.yaml
func VerifyPegaDeployment(t *testing.T, deploymentObj *appsv1.Deployment, expectedDeployment pegaDeployment, options *helm.Options) {
	require.Equal(t, *deploymentObj.Spec.Replicas, int32(1))
	require.Equal(t, *deploymentObj.Spec.ProgressDeadlineSeconds, int32(2147483647))
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Selector.MatchLabels["app"])
	require.Equal(t, *deploymentObj.Spec.Strategy.RollingUpdate.MaxSurge, intstr.FromString("25%"))
	require.Equal(t, *deploymentObj.Spec.Strategy.RollingUpdate.MaxUnavailable, intstr.FromString("25%"))
	require.Equal(t, deploymentObj.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Template.Labels["app"])
	require.NotEmpty(t, deploymentObj.Spec.Template.Annotations["config-check"])
	deploymentSpec := deploymentObj.Spec.Template.Spec
	VerifyDeployment(t, &deploymentSpec, expectedDeployment, options)
}

// VerifyDeployment - Performs common pega deployment/statefulset assertions with the values as provided in default values.yaml
func VerifyDeployment(t *testing.T, pod *k8score.PodSpec, expectedSpec pegaDeployment, options *helm.Options) {
	require.Equal(t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, expectedSpec.name, pod.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	require.Equal(t, pod.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, pod.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, pod.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	actualInitContainers := pod.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	//require.Equal(t, expectedSpec.initContainers, actualInitContainerNames) NEED TO CHANGE FOR "install-deploy"
	VerifyInitContinerData(t, actualInitContainers, options)
	require.Equal(t, pod.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, pod.Containers[0].Image, "pegasystems/pega")
	require.Equal(t, pod.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, pod.Containers[0].Ports[0].ContainerPort, int32(8080))
	var envIndex int32 = 0
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "NODE_TYPE")
	require.Equal(t, expectedSpec.nodeType, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "PEGA_APP_CONTEXT_PATH")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "prweb")
	if expectedSpec.name == "pega-web" || expectedSpec.name == "pega-stream" {
		envIndex++
		require.Equal(t, pod.Containers[0].Env[envIndex].Name, "REQUESTOR_PASSIVATION_TIMEOUT")
		require.Equal(t, expectedSpec.passivationTimeout, pod.Containers[0].Env[envIndex].Value)
	}
	if options.SetValues["constellation.enabled"] == "true" {
		envIndex++
		require.Equal(t, pod.Containers[0].Env[envIndex].Name, "COSMOS_SETTINGS")
		require.Equal(t, "Pega-UIEngine/cosmosservicesURI=/c11n", pod.Containers[0].Env[envIndex].Value)
	}
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "JAVA_OPTS")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "")
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "INITIAL_HEAP")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "4096m")
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "MAX_HEAP")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "7168m")
	require.Equal(t, pod.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")
	require.Equal(t, "4", pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "8Gi", pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, "2", pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")

	require.Equal(t, pod.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, pod.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, pod.Containers[0].LivenessProbe.PeriodSeconds, int32(30))
	require.Equal(t, pod.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, pod.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, pod.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, pod.Containers[0].ReadinessProbe.PeriodSeconds, int32(30))
	require.Equal(t, pod.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, pod.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, pod.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, *pod.TerminationGracePeriodSeconds, int64(300))
	require.Equal(t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, pod.Volumes[1].Secret.SecretName, "pega-credentials-secret")

}

type pegaDeployment struct {
	name               string
	initContainers     []string
	nodeType           string
	passivationTimeout string
}
