package deployment

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

var terminationGracePeriodSeconds int64 = 300
var terminationGracePeriodSecondsPtr = &terminationGracePeriodSeconds
var volumeDefaultMode int32 = 420
var volumeDefaultModePtr = &volumeDefaultMode

type Verifier struct {
	o                       PegaDeploymentTestOptions
	k8sInformationExtractor K8sInformationExtractor
	t                       *testing.T
	podVerifier             PodVerifierImpl
	_helmOptions            *helm.Options
	_nodeType               string
	_initContainers         []string
	_passivationTimeout     string
}

func NewDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string) *Verifier {

	return &Verifier{
		t:               t,
		podVerifier:     PodVerifierImpl{t: t},
		_helmOptions:    helmOptions,
		_initContainers: initContainers,
	}
}

func (v Verifier) Verify() {
	v.podVerifier.verifyPod(v)
}

func (v *Verifier) getPod() *k8score.PodSpec {
	return v.k8sInformationExtractor.GetPod()
}

func (v *Verifier) getDeploymentMetadata() DeploymentMetadata {
	return v.k8sInformationExtractor.GetDeploymentMetadata()
}

type PodValidator interface {
	verifyPod(v Verifier)
	verifyPodEnvironment(t *testing.T, pod *k8score.PodSpec, v Verifier)
	specificContainerEnvironmentValidation(t *testing.T, pod *k8score.PodSpec, options *helm.Options)
}
type PodVerifierImpl struct {
	PodValidator
}

func (p *PodVerifierImpl) verifyPod(v Verifier) {
	t := v.t
	pod := v.getPod()
	require.Equal(t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, v.getDeploymentMetadata().Name, pod.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
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
	p.verifyInitContinerData(pod, v._helmOptions)
	p.verifyPodEnvironment(pod, v)

	require.Equal(t, v._initContainers, actualInitContainerNames)
	require.Equal(t, pod.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, pod.Containers[0].Image, "pegasystems/pega")
	require.Equal(t, pod.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, pod.Containers[0].Ports[0].ContainerPort, int32(8080))

	require.Equal(t, "2", pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "8Gi", pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, "200m", pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")

	require.Equal(t, pod.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, pod.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, pod.Containers[0].LivenessProbe.PeriodSeconds, int32(20))
	require.Equal(t, pod.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, pod.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, pod.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, pod.Containers[0].ReadinessProbe.PeriodSeconds, int32(20))
	require.Equal(t, pod.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, pod.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, pod.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, pod.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)
	require.Equal(t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, pod.Volumes[1].Secret.SecretName, "pega-credentials-secret")
}

func (p *PodVerifierImpl) verifyPodEnvironment(pod *k8score.PodSpec, v Verifier) {
	var envIndex int32 = 0
	env := pod.Containers[0].Env
	require.Contains(v.t, env, k8score.EnvVar{
		Name:  "NODE_TYPE",
		Value: v._nodeType,
	})
	require.Equal(v.t, env[envIndex].Name, "NODE_TYPE")
	require.Equal(v.t, v._nodeType, env[envIndex].Value)
	if v.k8sInformationExtractor.GetDeploymentMetadata().Name == "pega-web" || v.k8sInformationExtractor.GetDeploymentMetadata().Name == "pega-stream" {
		envIndex++
		require.Equal(v.t, env[envIndex].Name, "REQUESTOR_PASSIVATION_TIMEOUT")
		require.Equal(v.t, v._passivationTimeout, env[envIndex].Value)
	}
	envIndex++
	require.Equal(v.t, env[envIndex].Name, "JAVA_OPTS")
	require.Equal(v.t, env[envIndex].Value, "")
	envIndex++
	require.Equal(v.t, env[envIndex].Name, "INITIAL_HEAP")
	require.Equal(v.t, env[envIndex].Value, "4096m")
	envIndex++
	require.Equal(v.t, env[envIndex].Name, "MAX_HEAP")
	require.Equal(v.t, env[envIndex].Value, "7168m")
	require.Equal(v.t, pod.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")
}

func (p *PodVerifierImpl) verifyInitContinerData(pod *k8score.PodSpec, options *helm.Options) {
	containers := pod.InitContainers
	if len(containers) == 0 {
		println("no init containers")
	}

	for i := 0; i < len(containers); i++ {
		container := containers[i]
		name := container.Name
		if name == "wait-for-pegainstall" {
			require.Equal(p.t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(p.t, []string{"job", "pega-db-install"}, container.Args)
		} else if name == "wait-for-pegasearch" {
			require.Equal(p.t, "busybox:1.31.0", container.Image)
			require.Equal(p.t, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-cassandra" {
			require.Equal(p.t, "cassandra:3.11.3", container.Image)
			require.Equal(p.t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" release-name-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-cassandra" {
			require.Equal(p.t, "cassandra:3.11.3", container.Image)
			require.Equal(p.t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" release-name-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-pegaupgrade" {
			require.Equal(p.t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(p.t, []string{"job", "pega-db-upgrade"}, container.Args)
			aksSpecificUpgraderDeployEnvs(p.t, options, container)
			p.specificContainerEnvironmentValidation(pod, options)
		} else if name == "wait-for-pre-dbupgrade" {
			require.Equal(p.t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(p.t, []string{"job", "pega-pre-upgrade"}, container.Args)
		} else if name == "wait-for-rolling-updates" {
			require.Equal(p.t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(p.t, []string{"sh", "-c", " kubectl rollout status deployment/pega-web --namespace default && kubectl rollout status deployment/pega-batch --namespace default && kubectl rollout status statefulset/pega-stream --namespace default"}, container.Command)
		} else {
			fmt.Println("invalid init containers found.. please check the list", name)
			p.t.Fail()
		}
	}
}

func (p *PodVerifierImpl) specificContainerEnvironmentValidation(pod *k8score.PodSpec, options *helm.Options) {

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
