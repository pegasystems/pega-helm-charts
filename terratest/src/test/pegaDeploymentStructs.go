package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"strings"
	"testing"
	"testutility"
)

type pegaStandardDeploymentTest struct {
	provider        string
	action          string
	verifiers       []testutility.VerifierImpl
	t               *testing.T
	_helmOptions    *helm.Options
	_initContainers []string
	_helmChartPath  string
}

func NewPegaStandardDeploymentTest(provider string, action string, initContainers []string, t *testing.T) *pegaStandardDeploymentTest {
	d := &pegaStandardDeploymentTest{
		provider:        provider,
		action:          action,
		_initContainers: initContainers,
		t:               t,
	}
	d.t.Parallel()

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	d._helmChartPath = helmChartPath
	require.NoError(d.t, err)

	d._helmOptions = &helm.Options{
		SetValues: map[string]string{
			"global.provider":        d.provider,
			"global.actions.execute": d.action,
		},
	}
	return d
}

type DeploymentVerifier struct {
	testutility.VerifierImpl
	k8sDeployment *v1.Deployment
}

func (d *pegaStandardDeploymentTest) Run() {
	d.splitAndVerifyPegaDeployments()
}

func (d *pegaStandardDeploymentTest) splitAndVerifyPegaDeployments() {
	deployment := helm.RenderTemplate(d.t, d._helmOptions, d._helmChartPath, []string{DeploymentTierPath})
	deploymentSlice := strings.Split(deployment, "---")

	var processedDeployments int
	for _, deploymentInfo := range deploymentSlice {
		if len(deploymentInfo) == 0 {
			continue
		}
		var deployMeta testutility.DeploymentMetadata
		helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployMeta)
		switch deployMeta.Name {
		case "pega-web":
			var deployment v1.Deployment
			helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployment)
			NewPegaWebDeployVerifier(d.t, d._helmOptions, d._initContainers, &deployment).verify()
			processedDeployments++
			break
		case "pega-batch":
			var deployment v1.Deployment
			helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployment)
			NewPegaBatchDeployVerifier(d.t, d._helmOptions, d._initContainers, &deployment).verify()
			processedDeployments++
			break
		case "pega-stream":
			var deployment v1beta2.StatefulSet
			helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployment)
			NewPegaStreamDeployVerifier(d.t, d._helmOptions, d._initContainers, &deployment).verify()
			processedDeployments++
			break
		}
	}
	require.Equal(d.t, 3, processedDeployments, "Should process all deployments")
}

func NewDeploymentVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *v1.Deployment) *DeploymentVerifier {
	verifierImpl := *NewDeployVerifier(t, helmOptions, initContainers)
	verifierImpl.k8sInformationExtractor = testutility._k8sDeploymentExtractor{
		_k8sDeployment: deployment,
	}
	return &DeploymentVerifier{
		k8sDeployment: deployment,
		VerifierImpl:  verifierImpl,
	}
}

func NewDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string) *testutility.VerifierImpl {
	return &testutility.VerifierImpl{
		t:               t,
		_helmOptions:    helmOptions,
		_initContainers: initContainers,
	}
}

func NewPegaWebDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *v1.Deployment) *PegaWebDeployVerifier {
	v := &PegaWebDeployVerifier{
		DeploymentVerifier: *NewDeploymentVerifier(t, helmOptions, initContainers, deployment),
	}
	v._nodeType = "WebUser"
	v._passivationTimeout = "900"
	return v
}

type PegaBatchDeployVerifier struct {
	DeploymentVerifier
}

type PegaStreamDeployVerifier struct {
	StatefulSetVerifier
}

func NewPegaStreamDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *v1beta2.StatefulSet) *PegaStreamDeployVerifier {
	v := &PegaStreamDeployVerifier{
		StatefulSetVerifier: StatefulSetVerifier{
			VerifierImpl:  *NewDeployVerifier(t, helmOptions, initContainers),
			k8sDeployment: deployment,
		},
	}
	v._nodeType = "Stream"
	v._passivationTimeout = "900"
	return v
}

type PegaWebDeployVerifier struct {
	DeploymentVerifier
}

func (v *PegaWebDeployVerifier) getPod() *k8score.PodSpec {
	return &v.k8sDeployment.Spec.Template.Spec
}

func NewPegaBatchDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *v1.Deployment) *PegaBatchDeployVerifier {
	v := &PegaBatchDeployVerifier{
		DeploymentVerifier: *NewDeploymentVerifier(t, helmOptions, initContainers, deployment),
	}
	v._nodeType = "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX"
	return v
}

type StatefulSetVerifier struct {
	testutility.VerifierImpl
	k8sDeployment *v1beta2.StatefulSet
}

func (p PegaStreamDeployVerifier) verify() {
	require.Equal(p.t, p.k8sDeployment.Spec.VolumeClaimTemplates[0].Name, "pega-stream")
	require.Equal(p.t, p.k8sDeployment.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(p.t, p.k8sDeployment.Spec.ServiceName, "pega-stream")
	statefulsetSpec := p.k8sDeployment.Spec.Template.Spec
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[1].Name, "pega-stream")
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/streamvol")
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[2].Name, "pega-volume-credentials")
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[2].MountPath, "/opt/pega/secrets")
}

func (v *testutility.VerifierImpl) verify() {
	v.verifyPod(v.k8sInformationExtractor.getPod())
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

func (v *testutility.VerifierImpl) verifyPod(pod *k8score.PodSpec) {
	require.Equal(v.t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(v.t, v.k8sInformationExtractor.getDeploymentMetadata().Name, pod.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
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
	VerifyInitContinerData(v.t, actualInitContainers, v._helmOptions)
	require.Equal(v.t, pod.Containers[0].Name, "pega-web-tomcat")
	require.Equal(v.t, pod.Containers[0].Image, "pegasystems/pega")
	require.Equal(v.t, pod.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(v.t, pod.Containers[0].Ports[0].ContainerPort, int32(8080))
	var envIndex int32 = 0
	require.Equal(v.t, pod.Containers[0].Env[envIndex].Name, "NODE_TYPE")
	require.Equal(v.t, v._nodeType, pod.Containers[0].Env[envIndex].Value)
	if v.k8sInformationExtractor.getDeploymentMetadata().Name == "pega-web" || v.k8sInformationExtractor.getDeploymentMetadata().Name == "pega-stream" {
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

func (e testutility._k8sDeploymentExtractor) getPod() *k8score.PodSpec {
	return &e._k8sDeployment.Spec.Template.Spec
}

func (e testutility._k8sDeploymentExtractor) getDeploymentMetadata() testutility.DeploymentMetadata {
	return testutility.DeploymentMetadata{
		ObjectMeta: e._k8sDeployment.ObjectMeta,
		TypeMeta:   e._k8sDeployment.TypeMeta,
	}
}
