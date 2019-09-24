package test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sbatch "k8s.io/api/batch/v1"
	k8score "k8s.io/api/core/v1"
)

type pegaJob struct {
	name           string
	initContainers []string
	configMapName  string
}

/* func VerifyPegaStandardTierDeployment(t *testing.T, helmChartPath string, options *helm.Options, initContainers []string) {

	// Deployment objects
	VerifyPegaJobs(t, helmChartPath, options, initContainers)
} */

/* func VerifyPegaJobs(t *testing.T, helmChartPath string, options *helm.Options, expectedJob pegaJob) {
	deployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	var deploymentObj appsv1.Deployment
	deploymentSlice := strings.Split(deployment, "---")
	for index, deploymentInfo := range deploymentSlice {
		if index >= 1 && index <= 3 {
			helm.UnmarshalK8SYaml(t, deploymentInfo, &deploymentObj)
			if index == 1 {
				VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-web", initContainers, "WebUser"})
			} else if index == 2 {
				VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-batch", initContainers, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX"})
			}
		}
	}
} */

//To verify Peg ajob
func VerifyPegaJob(t *testing.T, pegaHelmChartPath string, options *helm.Options, expectedJob pegaJob) {
	//t.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerJob := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	var installerJobObj k8sbatch.Job
	helm.UnmarshalK8SYaml(t, installerJob, &installerJobObj)
	installerJobSpec := installerJobObj.Spec.Template.Spec
	installerJobConatiners := installerJobObj.Spec.Template.Spec.Containers

	var containerPort int32 = 8080

	require.Equal(t, installerJobSpec.Volumes[0].Name, "pega-volume-credentials")
	require.Equal(t, installerJobSpec.Volumes[0].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, installerJobSpec.Volumes[0].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, installerJobSpec.Volumes[1].Name, "pega-volume-installer")
	require.Equal(t, installerJobSpec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-installer-config")
	require.Equal(t, installerJobSpec.Volumes[1].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)

	require.Equal(t, installerJobConatiners[0].Name, expectedJob.name)
	require.Equal(t, "YOUR_INSTALLER_IMAGE:TAG", installerJobConatiners[0].Image)
	require.Equal(t, installerJobConatiners[0].Ports[0].ContainerPort, containerPort)
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].Name, "pega-volume-installer")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")
	require.Equal(t, installerJobConatiners[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, expectedJob.configMapName)

	require.Equal(t, installerJobSpec.ImagePullSecrets[0].Name, "pega-registry-secret")

	require.Equal(t, installerJobSpec.RestartPolicy, k8score.RestartPolicy("Never"))

	actualInitContainers := installerJobSpec.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	require.Equal(t, expectedJob.initContainers, actualInitContainerNames)
	VerifyInstallerInitContinerData(t, actualInitContainers)

}

func VerifyInstallerInitContinerData(t *testing.T, containers []k8score.Container) {

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
		} else {
			fmt.Println("in last else", name)
			t.Fail()
		}
	}
}
