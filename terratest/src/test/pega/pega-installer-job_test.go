package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sbatch "k8s.io/api/batch/v1"
	k8score "k8s.io/api/core/v1"
)

type pegaDbJob struct {
	name           string
	initContainers []string
	configMapName  string
	containerName  string
}

var volDefaultMode int32 = 420
var volDefaultModePointer = &volDefaultMode

func TestPegaInstallerJob(t *testing.T) {

	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	var imagePullPolicy = []string{"", "IfNotPresent", "Always"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				for _, pullPolicy := range imagePullPolicy {
					var options = &helm.Options{
						SetValues: map[string]string{
							"global.deployment.name":        depName,
							"global.provider":               vendor,
							"global.actions.execute":        operation,
							"installer.imagePullPolicy":     pullPolicy,
							"installer.upgrade.upgradeType": "zero-downtime",
						},
					}
					yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
					yamlSplit := strings.Split(yamlContent, "---")

					// If there are three slices, it means that it is a pega-upgrade-deploy job
					if len(yamlSplit) == 4 {
						var expectedJob pegaDbJob
						for index, jobInfo := range yamlSplit {
							if index >= 1 && index <= 3 {
								if index == 1 {
									expectedJob = pegaDbJob{"pega-pre-upgrade", []string{}, "pega-upgrade-environment-config", "pega-installer"}
								} else if index == 2 {
									expectedJob = pegaDbJob{"pega-zdt-upgrade", []string{"wait-for-pre-dbupgrade"}, "pega-upgrade-environment-config", "pega-installer"}
								} else if index == 3 {
									expectedJob = pegaDbJob{"pega-post-upgrade", []string{"wait-for-pegaupgrade", "wait-for-rolling-updates"}, "pega-upgrade-environment-config", "pega-installer"}
								}

								assertJob(t, jobInfo, expectedJob, options, pullPolicy)
							}

						}
					} else {
						if operation == "install" || operation == "install-deploy" {
							assertJob(t, yamlSplit[1], pegaDbJob{"pega-db-install", []string{}, "pega-install-environment-config", "pega-installer"}, options, pullPolicy)
						} else {
							assertJob(t, yamlSplit[1], pegaDbJob{"pega-pre-upgrade", []string{}, "pega-upgrade-environment-config", "pega-installer"}, options, pullPolicy)
						}
					}

				}

			}
		}
	}
}

func assertJob(t *testing.T, jobYaml string, expectedJob pegaDbJob, options *helm.Options, pullPolicy string) {
	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, jobYaml, &jobObj)

	jobSpec := jobObj.Spec.Template.Spec
	jobContainers := jobObj.Spec.Template.Spec.Containers

	var containerPort int32 = 8080

	require.Equal(t, jobSpec.Volumes[0].Name, "pega-volume-credentials")
	require.Equal(t, jobSpec.Volumes[0].VolumeSource.Projected.Sources[0].Secret.Name, getObjName(options, "-credentials-secret"))
	require.Equal(t, jobSpec.Volumes[0].VolumeSource.Projected.DefaultMode, volDefaultModePointer)
	require.Equal(t, jobSpec.Volumes[1].Name, "pega-volume-installer")
	if jobSpec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name == "pega-install-config" {
		require.Equal(t, jobSpec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-install-config")
	}
	if jobSpec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name == "pega-upgrade-config" {
		require.Equal(t, jobSpec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-upgrade-config")
	}
	require.Equal(t, jobSpec.Volumes[1].VolumeSource.ConfigMap.DefaultMode, volDefaultModePointer)

	require.Equal(t, string(jobContainers[0].ImagePullPolicy), pullPolicy)

	require.Equal(t, jobContainers[0].Name, expectedJob.containerName)
	require.Equal(t, "YOUR_INSTALLER_IMAGE:TAG", jobContainers[0].Image)
	require.Equal(t, jobContainers[0].Ports[0].ContainerPort, containerPort)
	require.Equal(t, jobContainers[0].VolumeMounts[0].Name, "pega-volume-installer")
	require.Equal(t, jobContainers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, jobContainers[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, jobContainers[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")
	require.Equal(t, jobContainers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, expectedJob.configMapName)

	require.Equal(t, jobSpec.ImagePullSecrets[0].Name, getObjName(options, "-registry-secret"))

	require.Equal(t, jobSpec.RestartPolicy, k8score.RestartPolicy("Never"))

	actualInitContainers := jobSpec.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	require.Equal(t, expectedJob.initContainers, actualInitContainerNames)
	VerifyInitContainerData(t, actualInitContainers, options)
}
