package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

var replicas int32 = 1
var replicasPtr = &replicas
var ProgressDeadlineSeconds int32 = 2147483647
var ProgressDeadlineSecondsPtr = &ProgressDeadlineSeconds
var rollingUpdate intstr.IntOrString = intstr.FromString("25%")
var rollingUpdatePtr = &rollingUpdate
var volumeDefaultMode int32 = 420
var volumeDefaultModePtr = &volumeDefaultMode
var terminationGracePeriodSeconds int64 = 300
var terminationGracePeriodSecondsPtr = &terminationGracePeriodSeconds

type pegaDeployment struct {
	name           string
	initContainers []string
	nodeType       string
}

func VerifyInitContinerData(t *testing.T, containers []k8score.Container) {

	for i := 0; i < len(containers); i++ {
		container := containers[i]
		name := container.Name
		if name == "wait-for-pegainstall" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-db-install"}, container.Args)
		} else if name == "wait-for-pegasearch" {
			require.Equal(t, "busybox:1.31.0", container.Image)
			require.Equal(t, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-cassandra" {
			require.Equal(t, "cassandra:3.11.3", container.Image)
			require.Equal(t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" release-name-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else {
			fmt.Println("in last else", name)
			t.Fail()
		}
	}
}

func VerifyPegaDeployment(t *testing.T, deploymentObj *appsv1.Deployment, expectedDeployment *pegaDeployment) {

	require.Equal(t, deploymentObj.Spec.Replicas, replicasPtr)
	require.Equal(t, deploymentObj.Spec.ProgressDeadlineSeconds, ProgressDeadlineSecondsPtr)
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Selector.MatchLabels["app"])
	require.Equal(t, deploymentObj.Spec.Strategy.RollingUpdate.MaxSurge, rollingUpdatePtr)
	require.Equal(t, deploymentObj.Spec.Strategy.RollingUpdate.MaxUnavailable, rollingUpdatePtr)
	require.Equal(t, deploymentObj.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))

	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Template.Labels["app"])
	require.NotEmpty(t, deploymentObj.Spec.Template.Annotations["config-check"])

	webDeploymentSpec := deploymentObj.Spec.Template.Spec

	require.Equal(t, webDeploymentSpec.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, expectedDeployment.name, webDeploymentSpec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	require.Equal(t, webDeploymentSpec.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, webDeploymentSpec.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, webDeploymentSpec.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, webDeploymentSpec.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	actualInitContainers := webDeploymentSpec.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	require.Equal(t, expectedDeployment.initContainers, actualInitContainerNames)
	VerifyInitContinerData(t, actualInitContainers)
	require.Equal(t, webDeploymentSpec.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, webDeploymentSpec.Containers[0].Image, "YOUR_PEGA_DEPLOY_IMAGE:TAG")
	require.Equal(t, webDeploymentSpec.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, webDeploymentSpec.Containers[0].Ports[0].ContainerPort, int32(8080))
	require.Equal(t, webDeploymentSpec.Containers[0].Env[0].Name, "NODE_TYPE")
	require.Equal(t, expectedDeployment.nodeType, webDeploymentSpec.Containers[0].Env[0].Value)
	require.Equal(t, webDeploymentSpec.Containers[0].Env[1].Name, "JAVA_OPTS")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[1].Value, "")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[2].Name, "INITIAL_HEAP")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[2].Value, "4096m")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[3].Name, "MAX_HEAP")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[3].Value, "7168m")
	require.Equal(t, webDeploymentSpec.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")

	require.Equal(t, "2", webDeploymentSpec.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "8Gi", webDeploymentSpec.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, "200m", webDeploymentSpec.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", webDeploymentSpec.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")

	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.PeriodSeconds, int32(10))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.PeriodSeconds, int32(10))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, webDeploymentSpec.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, webDeploymentSpec.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, webDeploymentSpec.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)
}

/*func VerifyPegaServices()

func VerifyPegaIngress()

func VerifyPegaStatefulset()

// Just verify what is exposed in the values yaml & k8s objects
func VerifyCassandraService()

// Just verify what is exposed in the values yaml & k8s objects
func VerifySearchService()

func VefifyEnvironmentConfig()

func VerifyTierConfig()

func VerifyHPA()
*/
