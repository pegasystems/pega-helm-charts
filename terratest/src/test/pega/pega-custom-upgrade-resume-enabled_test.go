package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sbatch "k8s.io/api/batch/v1"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaCustomUpgradeJobResumeEnabled(t *testing.T) {

	var supportedVendors = []string{"eks"}
	var supportedOperations = []string{"upgrade", "upgrade-deploy"}
	var deploymentNames = []string{"pega"}
	var upgradeType = []string{"custom"}
	var upgradeSteps = []string{"rules_upgrade"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				for _, upType := range upgradeType {
					for _, upSteps := range upgradeSteps {
						var options = &helm.Options{
							ValuesFiles: []string{"data/values_with_automatic_resume_enabled.yaml"},
							SetValues: map[string]string{
								"global.deployment.name":         depName,
								"global.provider":                vendor,
								"global.actions.execute":         operation,
								"installer.upgrade.upgradeType":  upType,
								"installer.upgrade.upgradeSteps": upSteps,
							},
						}
						yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
						yamlSplit := strings.Split(yamlContent, "---")
						assertUpgradeJobJobResumeEnabled(t, yamlSplit[1], pegaDbJob{"pega-db-custom-upgrade", []string{}, "pega-upgrade-environment-config", "pega-installer", "upgrade"}, options)

					}
				}
			}
		}
	}
}

func assertUpgradeJobJobResumeEnabled(t *testing.T, jobYaml string, expectedJob pegaDbJob, options *helm.Options) {
	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, jobYaml, &jobObj)

	jobSpec := jobObj.Spec.Template.Spec
	jobContainers := jobObj.Spec.Template.Spec.Containers

	var containerPort int32 = 8080

	require.Equal(t, jobSpec.Volumes[0].Name, "pega-installer-mount-volume")
	require.Equal(t, jobSpec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName, "test-claim")

	if jobContainers[0].Name == "pega-db-upgrade-rules-migration" || jobContainers[0].Name == "pega-db-upgrade-rules-upgrade" || jobContainers[0].Name == "pega-db-upgrade-data-upgrade" {
		require.Equal(t, jobContainers[0].Name, "pega-installer")
	}

	require.Equal(t, "YOUR_INSTALLER_IMAGE:TAG", jobContainers[0].Image)
	require.Equal(t, jobContainers[0].Ports[0].ContainerPort, containerPort)
	require.Equal(t, jobContainers[0].VolumeMounts[0].Name, "pega-installer-mount-volume")
	require.Equal(t, jobContainers[0].VolumeMounts[0].MountPath, "/opt/pega/mount/installer")
	require.Equal(t, jobContainers[0].Env[0].Name, "ACTION")
	require.Equal(t, jobContainers[0].Env[0].Value, expectedJob.action)
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
