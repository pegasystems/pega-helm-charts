package test

import (
	"path/filepath"
	"testing"

	//"fmt"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	k8srbac "k8s.io/api/rbac/v1"

	//k8sresource "k8s.io/apimachinery/pkg/api/resource"

	"github.com/gruntwork-io/terratest/modules/helm"
)

const pegaHelmChartPath = "../../../charts/pega"

var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "install-deploy",
	},
}

func TestInstallDeployActionSkippedTemplates(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	// with action as 'install-deploy' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-upgrade-environment-config.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	// assert that above templates are not rendered
	require.Empty(t, deployment)
}

func TestInstallDeployActionInstallerRole(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	deployRole := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-role.yaml"})
	var deployRoleObj k8srbac.Role
	helm.UnmarshalK8SYaml(t, deployRole, &deployRoleObj)
	require.Equal(t, deployRoleObj.Rules[0].APIGroups, []string{"", "batch", "extensions", "apps"})
	require.Equal(t, deployRoleObj.Rules[0].Resources, []string{"jobs", "deployments", "statefulsets"})
	require.Equal(t, deployRoleObj.Rules[0].Verbs, []string{"get", "watch", "list"})

}

func TestInstallDeployActionInstallerRoleBinding(t *testing.T) {
	t.Parallel()

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

func TestInstallDeployActionInstallerJob(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaJob(t, helmChartPath, options, pegaJob{"pega-db-install", []string{}, "pega-install-environment-config"})

}

func TestInstallDeployActionInstallerConfig(t *testing.T) {
	t.Parallel()
	t.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
	var installConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installerConfig, &installConfigMap)

	installConfigData := installConfigMap.Data
	compareConfigMapData(t, []byte(installConfigData["prconfig.xml.tmpl"]), "expectedPrconfig.xml")
	compareConfigMapData(t, []byte(installConfigData["setupDatabase.properties.tmpl"]), "expectedsetupDatabase.properties")
	compareConfigMapData(t, []byte(installConfigData["prbootstrap.properties.tmpl"]), "expectedPRbootstrap.properties")
	compareConfigMapData(t, []byte(installConfigData["prlog4j2.xml"]), "expectedPRlog4j2.xml")

}

func TestInstallDeployActionInstallerEnvConfig(t *testing.T) {
	t.Parallel()
	//.Skip()
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
	require.Equal(t, installEnvConfigData["BYPASS_UDF_GENERATION"], "false")
	require.Equal(t, installEnvConfigData["BYPASS_TRUNCATE_UPDATESCACHE"], "false")
	require.Equal(t, installEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, installEnvConfigData["MAX_IDLE"], "5")
	require.Equal(t, installEnvConfigData["MAX_WAIT"], "-1")
	require.Equal(t, installEnvConfigData["MAX_ACTIVE"], "10")
	require.Equal(t, installEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, installEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, installEnvConfigData["ACTION"], "install-deploy")

}

func TestInstallDeployActionStandardDeployment(t *testing.T) {
	t.Parallel()
	//t.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, options, []string{"wait-for-pegainstall", "wait-for-pegasearch", "wait-for-cassandra"})

}
