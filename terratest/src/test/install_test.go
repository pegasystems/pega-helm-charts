package test

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
	k8score "k8s.io/api/core/v1"
)

const PegaHelmChartPath = "../../../charts/pega"

// set action execute to install
var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "install",
		"cassandra.enabled":      "false",
	},
}

func TestInstallActionSkippedTemplates(t *testing.T) {
	//t.Skip("skipping this")
	t.Parallel()

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-action-validate.yaml",
		//		"templates/pega-credentials-secret.yaml",
		"templates/pega-environment-config.yaml",
		//		"templates/pega-registry-secret.yaml",
		"templates/pega-tier-config.yaml",
		"templates/pega-tier-deployment.yaml",
		"templates/pega-tier-hpa.yaml",
		"templates/pega-tier-ingress.yaml",
		"templates/pega-tier-service.yaml",
		"charts/installer/templates/pega-installer-role.yaml",
		"charts/installer/templates/pega-installer-status-rolebinding.yaml",
	})

	var emptyObjects appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &emptyObjects)

	// assert that above templates are not rendered
	require.Empty(t, emptyObjects)
}

func TestInstallActionInstallJob(t *testing.T) {
	//t.Skip("skipping this")
	t.Parallel()

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	installerJob := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	var installerJobObj k8sbatch.Job
	helm.UnmarshalK8SYaml(t, installerJob, &installerJobObj)
	installerJobSpec := installerJobObj.Spec.Template.Spec
	installerJobConatiners := installerJobObj.Spec.Template.Spec.Containers

	var containerPort int32 = 8080
	var volumeDefaultMode int32 = 420
	var volumeDefaultModePtr = &volumeDefaultMode

	require.Equal(t, installerJobSpec.Volumes[0].Name, "pega-volume-credentials")
	require.Equal(t, installerJobSpec.Volumes[0].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, installerJobSpec.Volumes[0].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, installerJobSpec.Volumes[1].Name, "pega-volume-installer")
	require.Equal(t, installerJobSpec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-installer-config")
	require.Equal(t, installerJobSpec.Volumes[1].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)

	require.Equal(t, installerJobConatiners[0].Name, "pega-db-install")
	require.Equal(t, installerJobConatiners[0].Image, "YOUR_INSTALLER_IMAGE:TAG")
	require.Equal(t, installerJobConatiners[0].Ports[0].ContainerPort, containerPort)
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].Name, "pega-volume-installer")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")
	require.Equal(t, installerJobConatiners[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-install-environment-config")

	require.Equal(t, installerJobSpec.ImagePullSecrets[0].Name, "pega-registry-secret")

	require.Equal(t, installerJobSpec.RestartPolicy, k8score.RestartPolicy("Never"))
}

func TestInstallActionInstallerEnvironmentConfig(t *testing.T) {
	//t.Skip("skipping this")
	t.Parallel()
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	// pega-install-environment-config.yaml
	installEnvConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-install-environment-config.yaml"})
	var installEnvConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installEnvConfig, &installEnvConfigMap)

	installEnvConfigData := installEnvConfigMap.Data
	require.Equal(t, installEnvConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, installEnvConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, installEnvConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, installEnvConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, installEnvConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, installEnvConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, installEnvConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, installEnvConfigData["SYSTEM_NAME"], "pega")
	require.Equal(t, installEnvConfigData["PRODUCTION_LEVEL"], "2")
	require.Equal(t, installEnvConfigData["MULTITENANT_SYSTEM"], "false")
	require.Equal(t, installEnvConfigData["ADMIN_PASSWORD"], "ADMIN_PASSWORD")
	require.Equal(t, "", installEnvConfigData["STATIC_ASSEMBLER"])
	require.Equal(t, installEnvConfigData["BYPASS_UDF_GENERATION"], "false")
	require.Equal(t, installEnvConfigData["BYPASS_TRUNCATE_UPDATESCACHE"], "false")
	require.Equal(t, installEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, installEnvConfigData["MAX_IDLE"], "5")
	require.Equal(t, installEnvConfigData["MAX_WAIT"], "-1")
	require.Equal(t, installEnvConfigData["MAX_ACTIVE"], "10")
	require.Equal(t, installEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, installEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, installEnvConfigData["ACTION"], "install")
}

func TestInstallActionRegistrySecret(t *testing.T) {
	//t.Skip("skipping this")
	t.Parallel()
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerfiyRegistrySecret(t, helmChartPath, options)
}

func TestInstallActionCredentailsSecret(t *testing.T) {
	//t.Skip("skipping this")
	t.Parallel()
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyCredentialsSecret(t, helmChartPath, options)
}

func TestInstallActionInstallerConfig(t *testing.T) {
	t.Skip("skipping this")
	t.Parallel()
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	// pega-installer-config.yaml
	installerConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
	var installConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installerConfig, &installConfigMap)

	installConfigData := installConfigMap.Data
	compareConfigMapData(t, []byte(installConfigData["prconfig.xml.tmpl"]), "expectedPrconfig.xml")
	compareConfigMapData(t, []byte(installConfigData["setupDatabase.properties.tmpl"]), "expectedsetupDatabase.properties")
	compareConfigMapData(t, []byte(installConfigData["prbootstrap.properties.tmpl"]), "expectedPRbootstrap.properties")
	compareConfigMapData(t, []byte(installConfigData["prlog4j2.xml"]), "expectedPRlog4j2.xml")
}
