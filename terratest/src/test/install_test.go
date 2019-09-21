package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	require.Equal(t, installEnvConfigData["STATIC_ASSEMBLER"], "<nil>")
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

	// pega-registry-secret.yaml
	registrySecret := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})

	var registrySecretObj k8score.Secret
	helm.UnmarshalK8SYaml(t, registrySecret, &registrySecretObj)
	reqgistrySecretData := registrySecretObj.Data

	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
}

func TestInstallActionCredentailsSecret(t *testing.T) {
	//t.Skip("skipping this")
	t.Parallel()
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	secretOutput := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
	var secretobj k8score.Secret
	helm.UnmarshalK8SYaml(t, secretOutput, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["DB_USERNAME"]), "YOUR_JDBC_USERNAME")
	require.Equal(t, string(secretData["DB_PASSWORD"]), "YOUR_JDBC_PASSWORD")
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

/*func TestInstallActionShouldNotRenderDeployments(t *testing.T) {
	t.Skip("skipping this")
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to install
	options := &helm.Options{
		SetValues: map[string]string{
			"actions.execute": "install",
		},
	}

	// with action as 'install' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"charts/templates/pega-action-validate.yaml",
		"templates/pega-batch-config.yaml",
		"templates/pega-batch-deployment.yaml",
		"templates/pega-batch-hpa.yaml",
		"templates/pega-deploy-role.yaml",
		"templates/pega-environment-config.yaml",
		"templates/pega-installer-status-rolebinding.yaml",
		"templates/pega-search-deployment.yaml",
		"templates/pega-search-service.yaml",
		"templates/pega-search-transport-service.yaml",
		"templates/pega-stream-config.yaml",
		"templates/pega-stream-deployment.yaml",
		"templates/pega-stream-ingress.yaml",
		"templates/pega-stream-service.yaml",
		"templates/pega-upgrade-environment-config.yaml",
		"templates/pega-web-config.yaml",
		"templates/pega-web-deployment.yaml",
		"templates/pega-web-hpa.yaml",
		"templates/pega-web-ingress.yaml",
		"templates/pega-web-service.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	// assert that above templates are not rendered
	require.Empty(t, deployment)

	// pega-credentials-secret.yaml
	secretOutput := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
	var secretobj k8score.Secret
	helm.UnmarshalK8SYaml(t, secretOutput, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["DB_USERNAME"]), "YOUR_JDBC_USERNAME")
	require.Equal(t, string(secretData["DB_PASSWORD"]), "YOUR_JDBC_PASSWORD")

	// pega-install-environment-config.yaml
	installEnvConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-install-environment-config.yaml"})
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
	require.Equal(t, installEnvConfigData["ADMIN_PASSWORD"], "")
	require.Equal(t, installEnvConfigData["STATIC_ASSEMBLER"], "<nil>")
	require.Equal(t, installEnvConfigData["BYPASS_UDF_GENERATION"], "false")
	require.Equal(t, installEnvConfigData["BYPASS_TRUNCATE_UPDATESCACHE"], "false")
	require.Equal(t, installEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, installEnvConfigData["MAX_IDLE"], "5")
	require.Equal(t, installEnvConfigData["MAX_WAIT"], "-1")
	require.Equal(t, installEnvConfigData["MAX_ACTIVE"], "10")
	require.Equal(t, installEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, installEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, installEnvConfigData["ACTION"], "install")

	// pega-registry-secret.yaml
	registrySecret := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})

	var registrySecretObj k8score.Secret
	helm.UnmarshalK8SYaml(t, registrySecret, &registrySecretObj)
	reqgistrySecretData := registrySecretObj.Data

	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")

	// pega-installer-config.yaml
	installerConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-installer-config.yaml"})
	var installConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installerConfig, &installConfigMap)

	installConfigData := installConfigMap.Data
	compareConfigMapData(t, []byte(installConfigData["prconfig.xml.tmpl"]), "expectedPrconfig.xml")
	compareConfigMapData(t, []byte(installConfigData["setupDatabase.properties.tmpl"]), "expectedsetupDatabase.properties")
	compareConfigMapData(t, []byte(installConfigData["prbootstrap.properties.tmpl"]), "expectedPRbootstrap.properties")
	compareConfigMapData(t, []byte(installConfigData["prlog4j2.xml"]), "expectedPRlog4j2.xml")

	// pega-installer-job.yaml
	installerJob := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-installer-job.yaml"})
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
	require.Equal(t, installerJobConatiners[0].Image, "YOUR_PEGA_INSTALLER_IMAGE:TAG")
	require.Equal(t, installerJobConatiners[0].Ports[0].ContainerPort, containerPort)
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].Name, "pega-volume-installer")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")
	require.Equal(t, installerJobConatiners[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-install-environment-config")

	require.Equal(t, installerJobSpec.ImagePullSecrets[0].Name, "pega-registry-secret")

	require.Equal(t, installerJobSpec.RestartPolicy, k8score.RestartPolicy("Never"))
}*/

// util function for comparing
func compareConfigMapData(t *testing.T, actualFile []byte, expectedFileName string) {
	fmt.Println("expectedFileName")
	fmt.Print(expectedFileName)
	expectedPrconfig, err := ioutil.ReadFile(expectedFileName)
	require.Empty(t, err)

	equal := bytes.Equal(expectedPrconfig, actualFile)
	require.Equal(t, true, equal)
}
