package test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"bytes"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	k8sbatch "k8s.io/api/batch/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func TestInstallActionShouldNotRenderDeployments(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to install
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "upgrade",
		},
	}

	// with action as 'install' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
        "templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-installer-role.yaml",
		"templates/pega-environment-config.yaml",
		"charts/installer/templates/pega-installer-status-rolebinding.yaml",
		"charts/pegasearch/templates/pega-search-deployment.yaml",
		"charts/pegasearch/templates/pega-search-service.yaml",
		"charts/pegasearch/templates/pega-search-transport-service.yaml",
		"charts/installer/templates/pega-install-environment-config.yaml",
		"templates/pega-tier-config.yaml",
		"templates/pega-tier-deployment.yaml",
		"templates/pega-tier-hpa.yaml",
		"templates/pega-tier-ingress.yaml",
		"templates/pega-tier-service.yaml",

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

	// pega-registry-secret.yaml
	registrySecret := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})
	
	var registrySecretObj k8score.Secret
	helm.UnmarshalK8SYaml(t, registrySecret, &registrySecretObj)
	reqgistrySecretData := registrySecretObj.Data

	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
	//Base64 decode of below value is --> YOUR_DOCKER_REGISTRY_USERNAME:YOUR_DOCKER_REGISTRY_PASSWORD
	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "WU9VUl9ET0NLRVJfUkVHSVNUUllfVVNFUk5BTUU6WU9VUl9ET0NLRVJfUkVHSVNUUllfUEFTU1dPUkQ=")

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
	require.Equal(t, upgradeEnvConfigData["BYPASS_UDF_GENERATION"], "false")
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
	
	// pega-installer-config.yaml
    installerConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
	var installConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installerConfig, &installConfigMap)

	/*installConfigData := installConfigMap.Data

	compareConfigMapData(t, []byte(installConfigData["prconfig.xml.tmpl"]), "expectedPrconfig.xml")
	compareConfigMapData(t, []byte(installConfigData["setupDatabase.properties.tmpl"]), "expectedsetupDatabase.properties")
	compareConfigMapData(t, []byte(installConfigData["prbootstrap.properties.tmpl"]), "expectedPRbootstrap.properties")
	compareConfigMapData(t, []byte(installConfigData["prlog4j2.xml"]), "expectedPRlog4j2.xml")
	*/
	// pega-installer-job.yaml
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

	require.Equal(t, installerJobConatiners[0].Name, "pega-db-upgrade")
	require.Equal(t, installerJobConatiners[0].Image, "YOUR_INSTALLER_IMAGE:TAG")
	require.Equal(t, installerJobConatiners[0].Ports[0].ContainerPort, containerPort)
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].Name, "pega-volume-installer")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, installerJobConatiners[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")
	require.Equal(t, installerJobConatiners[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-upgrade-environment-config")
	
	require.Equal(t, installerJobSpec.ImagePullSecrets[0].Name, "pega-registry-secret")
	
	require.Equal(t, installerJobSpec.RestartPolicy, k8score.RestartPolicy("Never")) 



/* 	//Cassandra
		// pega-search-service.yaml
	searchService := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/cassandra/templates/service.yaml"})
	var searchServiceObj k8score.Service

	helm.UnmarshalK8SYaml(t, searchService, &searchServiceObj)
	var servicePort int32 = 80
	
	println(searchServiceObj.Spec.Ports)


	require.Equal(t, searchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, searchServiceObj.Spec.Selector["app"], "cassandra")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Name, "http")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Port, servicePort)
	require.Equal(t, searchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9200)) */
}

// util function for comparing
func compareConfigMapData(t *testing.T, actualFile []byte, expectedFileName string) {	
    expectedPrconfig, err := ioutil.ReadFile(expectedFileName)
	require.Empty(t, err)
	equal := bytes.Equal(expectedPrconfig, actualFile)
	require.Equal(t, true, equal)
}
