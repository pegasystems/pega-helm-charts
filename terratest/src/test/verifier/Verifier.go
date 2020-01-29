package verifier

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	. "test/deployment"
	"testing"
)

var terminationGracePeriodSeconds int64 = 300
var terminationGracePeriodSecondsPtr = &terminationGracePeriodSeconds
var volumeDefaultMode int32 = 420
var volumeDefaultModePtr = &volumeDefaultMode

type Verifier interface {
	Verify()
}

type VerifierImpl struct {
	Verifier
	PegaDeployment
	k8sInformationExtractor K8sInformationExtractor
	t                       *testing.T
	_helmOptions            *helm.Options
	_nodeType               string
	_initContainers         []string
	_passivationTimeout     string
}

func NewDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string) *VerifierImpl {
	return &VerifierImpl{
		t:               t,
		_helmOptions:    helmOptions,
		_initContainers: initContainers,
	}
}

func (v *VerifierImpl) Verify() {
	v.verifyPod(v.k8sInformationExtractor.GetPod())
	//v.verifyVolume(v.getVolume())
	//require.Equal(v.t, v.k8sDeployment.Spec.Replicas, replicasPtr)
	//require.Equal(v.t, v.k8sDeployment.Spec.ProgressDeadlineSeconds, ProgressDeadlineSecondsPtr)
	//require.Equal(v.t, v.deploymentName, v.k8sDeployment.Spec.Selector.MatchLabels["app"])
	//require.Equal(v.t, v.k8sDeployment.Spec.Strategy.RollingUpdate.MaxSurge, rollingUpdatePtr)
	//require.Equal(v.t, v.k8sDeployment.Spec.Strategy.RollingUpdate.MaxUnavailable, rollingUpdatePtr)
	//require.Equal(v.t, v.k8sDeployment.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))
	//require.Equal(v.t, v.deploymentName, v.k8sDeployment.Spec.Template.Labels["app"])
	//require.NotEmpty(v.t, v.k8sDeployment.Spec.Template.Annotations["config-check"])
	//v.verifyDeployment()
}

func (v *VerifierImpl) verifyPod(pod *k8score.PodSpec) {
	require.Equal(v.t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(v.t, v.k8sInformationExtractor.GetDeploymentMetadata().Name, pod.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	require.Equal(v.t, pod.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(v.t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(v.t, pod.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(v.t, pod.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	actualInitContainers := pod.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	require.Equal(v.t, v._initContainers, actualInitContainerNames)
	v.verifyInitContinerData(v.t, actualInitContainers, v._helmOptions)
	require.Equal(v.t, pod.Containers[0].Name, "pega-web-tomcat")
	require.Equal(v.t, pod.Containers[0].Image, "pegasystems/pega")
	require.Equal(v.t, pod.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(v.t, pod.Containers[0].Ports[0].ContainerPort, int32(8080))
	var envIndex int32 = 0
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Name, "NODE_TYPE")
	require.Equal(v.t, v._nodeType, pod.Containers[0].Env[envIndex].Value)
	if v.k8sInformationExtractor.GetDeploymentMetadata().Name == "pega-web" || v.k8sInformationExtractor.GetDeploymentMetadata().Name == "pega-stream" {
		envIndex++
		require.Equal(v.t, pod.Containers[0].Env[envIndex].Name, "REQUESTOR_PASSIVATION_TIMEOUT")
		require.Equal(v.t, v._passivationTimeout, pod.Containers[0].Env[envIndex].Value)
	}
	envIndex++
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Name, "JAVA_OPTS")
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Value, "")
	envIndex++
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Name, "INITIAL_HEAP")
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Value, "4096m")
	envIndex++
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Name, "MAX_HEAP")
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Value, "7168m")
	require.Equal(v.t, pod.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")

	require.Equal(v.t, "2", pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(v.t, "8Gi", pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(v.t, "200m", pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(v.t, "6Gi", pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(v.t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(v.t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")

	require.Equal(v.t, pod.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(v.t, pod.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(v.t, pod.Containers[0].LivenessProbe.PeriodSeconds, int32(20))
	require.Equal(v.t, pod.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(v.t, pod.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(v.t, pod.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(v.t, pod.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(v.t, pod.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(v.t, pod.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.PeriodSeconds, int32(20))
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(v.t, pod.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(v.t, pod.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(v.t, pod.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(v.t, pod.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)
	require.Equal(v.t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(v.t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(v.t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(v.t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(v.t, pod.Volumes[1].Secret.SecretName, "pega-credentials-secret")
}
func (v *VerifierImpl) verifyInitContinerData(t *testing.T, containers []k8score.Container, options *helm.Options) {

	if len(containers) == 0 {
		println("no init containers")
	}

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
		} else if name == "wait-for-cassandra" {
			require.Equal(t, "cassandra:3.11.3", container.Image)
			require.Equal(t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" release-name-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-pegaupgrade" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-db-upgrade"}, container.Args)
			aksSpecificUpgraderDeployEnvs(t, options, container)
		} else if name == "wait-for-pre-dbupgrade" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-pre-upgrade"}, container.Args)
		} else if name == "wait-for-rolling-updates" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"sh", "-c", " kubectl rollout status deployment/pega-web --namespace default && kubectl rollout status deployment/pega-batch --namespace default && kubectl rollout status statefulset/pega-stream --namespace default"}, container.Command)
		} else {
			fmt.Println("invalid init containers found.. please check the list", name)
			t.Fail()
		}
	}
}

//aksSpecificUpgraderDeployEnvs - Test aks specific upgrade-deploy environmnet variables in case of upgrade-deploy
func aksSpecificUpgraderDeployEnvs(t *testing.T, options *helm.Options, container k8score.Container) {
	if options.SetValues["global.provider"] == "aks" && options.SetValues["global.actions.execute"] == "upgrade-deploy" {
		require.Equal(t, container.Env[0].Name, "KUBERNETES_SERVICE_HOST")
		require.Equal(t, container.Env[0].Value, "API_SERVICE_ADDRESS")
		require.Equal(t, container.Env[1].Name, "KUBERNETES_SERVICE_PORT_HTTPS")
		require.Equal(t, container.Env[1].Value, "SERVICE_PORT_HTTPS")
		require.Equal(t, container.Env[2].Name, "KUBERNETES_SERVICE_PORT")
		require.Equal(t, container.Env[2].Value, "SERVICE_PORT_HTTPS")
	}
}
