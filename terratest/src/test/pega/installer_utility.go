package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sbatch "k8s.io/api/batch/v1"
	k8score "k8s.io/api/core/v1"
	k8srbac "k8s.io/api/rbac/v1"
)

type pegaJob struct {
	name           string
	initContainers []string
	configMapName  string
}

// ReturnJobSlices - returns string array of rendered yaml sepearted by delimiter as "---"
func ReturnJobSlices(t *testing.T, pegaHelmChartPath string, options *helm.Options) []string {
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerJob := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})

	installerSlice := strings.Split(installerJob, "---")
	return installerSlice
}

// VerifyPegaJob -  Tests installer jobs rendered with the values as provided in default values.yaml
func VerifyPegaJob(t *testing.T, options *helm.Options, installerJobObj *k8sbatch.Job, expectedJob pegaJob) {
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
	VerifyInitContinerData(t, actualInitContainers, options)
}

// VerifyUpgradeEnvConfig - Tests upgrade environment config rendered with the values as provided in default values.yaml
func VerifyUpgradeEnvConfig(t *testing.T, options *helm.Options, pegaHelmChartPath string) {
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)
	// pega-install-environment-config.yaml
	upgradeEnvConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-upgrade-environment-config.yaml"})
	var upgradeEnvConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, upgradeEnvConfig, &upgradeEnvConfigMap)

	upgradeEnvConfigData := upgradeEnvConfigMap.Data

	require.Equal(t, upgradeEnvConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, upgradeEnvConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, upgradeEnvConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, upgradeEnvConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, upgradeEnvConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, upgradeEnvConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, upgradeEnvConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, upgradeEnvConfigData["UPGRADE_TYPE"], "in-place")
	require.Equal(t, upgradeEnvConfigData["MULTITENANT_SYSTEM"], "false")
	require.Equal(t, upgradeEnvConfigData["BYPASS_UDF_GENERATION"], "true")
	require.Equal(t, upgradeEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, upgradeEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, upgradeEnvConfigData["TARGET_RULES_SCHEMA"], "")
	require.Equal(t, upgradeEnvConfigData["TARGET_ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, upgradeEnvConfigData["MIGRATION_DB_LOAD_COMMIT_RATE"], "100")
	require.Equal(t, upgradeEnvConfigData["UPDATE_EXISTING_APPLICATIONS"], "false")
	require.Equal(t, upgradeEnvConfigData["UPDATE_APPLICATIONS_SCHEMA"], "false")
	require.Equal(t, upgradeEnvConfigData["RUN_RULESET_CLEANUP"], "false")
	require.Equal(t, upgradeEnvConfigData["REBUILD_INDEXES"], "false")
	require.Equal(t, upgradeEnvConfigData["DISTRIBUTION_KIT_URL"], "")
}

// VerifyInstallEnvConfig - Tests Installer environment config rendered with the values as provided in default values.yaml
func VerifyInstallEnvConfig(t *testing.T, options *helm.Options, pegaHelmChartPath string) {

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)
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
	require.Equal(t, "ADMIN_PASSWORD", installEnvConfigData["ADMIN_PASSWORD"])
	require.Equal(t, "", installEnvConfigData["STATIC_ASSEMBLER"])
	require.Equal(t, installEnvConfigData["BYPASS_UDF_GENERATION"], "true")
	require.Equal(t, installEnvConfigData["BYPASS_TRUNCATE_UPDATESCACHE"], "false")
	require.Equal(t, installEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, installEnvConfigData["MAX_IDLE"], "5")
	require.Equal(t, installEnvConfigData["MAX_WAIT"], "-1")
	require.Equal(t, installEnvConfigData["MAX_ACTIVE"], "10")
	require.Equal(t, installEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, installEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, installEnvConfigData["DISTRIBUTION_KIT_URL"], "")
	require.Equal(t, installEnvConfigData["ACTION"], options.SetValues["global.actions.execute"])
	require.Equal(t, "", installEnvConfigData["DISTRIBUTION_KIT_URL"])

}

// VerifyInstallerRoleBinding - Tests Installer role binding rendered with the values as provided in default values.yaml
func VerifyInstallerRoleBinding(t *testing.T, options *helm.Options, pegaHelmChartPath string) {
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerRoleBinding := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-status-rolebinding.yaml"})
	var installerRoleBindingObj k8srbac.RoleBinding
	helm.UnmarshalK8SYaml(t, installerRoleBinding, &installerRoleBindingObj)
	require.Equal(t, installerRoleBindingObj.RoleRef.APIGroup, "rbac.authorization.k8s.io")
	require.Equal(t, installerRoleBindingObj.RoleRef.Kind, "Role")
	require.Equal(t, installerRoleBindingObj.RoleRef.Name, "jobs-reader")

	require.Equal(t, installerRoleBindingObj.Subjects[0].Kind, "ServiceAccount")
	require.Equal(t, installerRoleBindingObj.Subjects[0].Name, "default")
	require.Equal(t, installerRoleBindingObj.Subjects[0].Namespace, "default")
}

// VerifyInstallerRole - Tests Installer role rendered with the values as provided in default values.yaml
func VerifyInstallerRole(t *testing.T, options *helm.Options, pegaHelmChartPath string) {
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	deployRole := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-role.yaml"})
	var deployRoleObj k8srbac.Role
	helm.UnmarshalK8SYaml(t, deployRole, &deployRoleObj)
	require.Equal(t, deployRoleObj.Rules[0].APIGroups, []string{"apps"})
	require.Equal(t, deployRoleObj.Rules[0].Resources, []string{"deployments", "statefulsets"})
	require.Equal(t, deployRoleObj.Rules[0].Verbs, []string{"get", "watch", "list"})
	require.Equal(t, deployRoleObj.Rules[1].APIGroups, []string{"batch"})
	require.Equal(t, deployRoleObj.Rules[1].Resources, []string{"jobs"})
	require.Equal(t, deployRoleObj.Rules[1].Verbs, []string{"get", "watch", "list"})
	require.Equal(t, deployRoleObj.Rules[2].APIGroups, []string{"extensions"})
	require.Equal(t, deployRoleObj.Rules[2].Resources, []string{"deployments"})
	require.Equal(t, deployRoleObj.Rules[2].Verbs, []string{"get", "watch", "list"})
}

// VerifyInstallerConfigMaps - Tests Installer configuration rendered with the values as provided in default values.yaml
func VerifyInstallerConfigMaps(t *testing.T, options *helm.Options, pegaHelmChartPath string) {
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
	var installConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installerConfig, &installConfigMap)

	installConfigData := installConfigMap.Data

	compareConfigMapData(t, installConfigData["prconfig.xml.tmpl"], "data/expectedPrconfig.xml")
	compareConfigMapData(t, installConfigData["setupDatabase.properties.tmpl"], "data/expectedSetupdatabase.properties")
	compareConfigMapData(t, installConfigData["prbootstrap.properties.tmpl"], "data/expectedPRbootstrap.properties")
	compareConfigMapData(t, installConfigData["migrateSystem.properties.tmpl"], "data/expectedMigrateSystem.properties.tmpl")
	compareConfigMapData(t, installConfigData["prlog4j2.xml"], "data/expectedPRlog4j2.xml")
	compareConfigMapData(t, installConfigData["prpcUtils.properties.tmpl"], "data/expectedPRPCUtils.properties.tmpl")
}
