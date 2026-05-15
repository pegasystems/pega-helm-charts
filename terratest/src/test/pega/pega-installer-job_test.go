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
	action         string
}

var volDefaultMode int32 = 420
var customArtifactorySecret = "artifactory_secret"
var volDefaultModePointer = &volDefaultMode

func TestPegaInstallerJob(t *testing.T) {
    runPegaInstallerJobTest(t, "")
}

func TestPegaInstallerJobWithExternalPegaRESTSecret(t *testing.T) {
    runPegaInstallerJobTest(t, "my-pega-rest-external-secret")
}

func runPegaInstallerJobTest(t *testing.T, externalSecretName string) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}
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
							"global.deployment.name": depName,
							"global.provider":        vendor,
							"global.actions.execute": operation,
							"global.customArtifactory.authentication.external_secret_name": customArtifactorySecret,
							"installer.imagePullPolicy":                                    pullPolicy,
							"installer.upgrade.upgradeType":                                "zero-downtime",
							"installer.upgrade.pegaRESTUsername": "username",
                            "installer.upgrade.pegaRESTPassword": "password",
                            "installer.upgrade.pega_rest_external_secret_name": externalSecretName,
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
									expectedJob = pegaDbJob{"pega-pre-upgrade", []string{}, "pega-upgrade-environment-config", "pega-installer", "pre-upgrade"}
								} else if index == 2 {
									expectedJob = pegaDbJob{"pega-zdt-upgrade", []string{"wait-for-pre-dbupgrade"}, "pega-upgrade-environment-config", "pega-installer", "upgrade"}
								} else if index == 3 {
									expectedJob = pegaDbJob{"pega-post-upgrade", []string{"wait-for-pegaupgrade", "wait-for-rolling-updates"}, "pega-upgrade-environment-config", "pega-installer", "post-upgrade"}
								}

								assertJob(t, jobInfo, expectedJob, options, pullPolicy, externalSecretName)
							}

						}
					} else {
						if operation == "install" || operation == "install-deploy" {
							assertJob(t, yamlSplit[1], pegaDbJob{"pega-db-install", []string{"jdbc-lib-downloader"}, "pega-install-environment-config", "pega-installer", "install"}, options, pullPolicy, externalSecretName)
						} else {
							assertJob(t, yamlSplit[1], pegaDbJob{"pega-pre-upgrade", []string{"jdbc-lib-downloader"}, "pega-upgrade-environment-config", "pega-installer", "pre-upgrade"}, options, pullPolicy, externalSecretName)
						}
					}

				}

			}
		}
	}
}


func TestPegaInstallerJobCustomInitContainers(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
                var options = &helm.Options{
                    ValuesFiles: []string{"data/values_init_conts.yaml"},
                    SetValues: map[string]string{
                        "global.deployment.name": depName,
                        "global.provider":        vendor,
                        "global.actions.execute": operation,
                        "installer.upgrade.upgradeType": "zero-downtime",
                        "installer.upgrade.pegaRESTUsername": "username",
                        "installer.upgrade.pegaRESTPassword": "password",
                        "installer.upgrade.pega_rest_external_secret_name": "",
                    },
                }

                yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
                yamlSplit := strings.Split(yamlContent, "---")

                count := len(yamlSplit)
                for i := 1; i < count; i++ {
                    assertJobCustomInitContainers(t, yamlSplit[i], options)
                }
            }
        }
    }
}

func assertJobCustomInitContainers(t *testing.T, jobYaml string, options *helm.Options) {
    var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, jobYaml, &jobObj)

    if (len(jobObj.Name)>0) {  //filter out chunks of yaml containing comments
        jobSpec := jobObj.Spec.Template.Spec

        actualInitContainers := jobSpec.InitContainers
        count := len(actualInitContainers)
        var foundCustomInitContainer1 bool = false
        var foundCustomInitContainer2 bool = false
        for i := 0; i < count; i++ {
            if (actualInitContainers[i].Name=="pre-deployment-action1") {
                foundCustomInitContainer1 = true
            }
            if (actualInitContainers[i].Name=="pre-deployment-action2") {
                foundCustomInitContainer2 = true
            }
        }
        require.Equal(t, true, foundCustomInitContainer1)
        require.Equal(t, true, foundCustomInitContainer2)
    }
}



func assertJob(t *testing.T, jobYaml string, expectedJob pegaDbJob, options *helm.Options, pullPolicy string, externalSecretName string) {
	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, jobYaml, &jobObj)

	jobSpec := jobObj.Spec.Template.Spec
	jobContainers := jobObj.Spec.Template.Spec.Containers
	var volumes = jobSpec.Volumes
	var volumeMounts = jobSpec.Containers[0].VolumeMounts

	var containerPort int32 = 8080

	require.Empty(t, jobSpec.Affinity)
	require.Empty(t, jobSpec.Tolerations)
	var pegaInstallerCredentialVolume = findNamedVolume(volumes, "pega-installer-credentials-volume")
	require.NotNil(t, pegaInstallerCredentialVolume)
	require.Equal(t, pegaInstallerCredentialVolume.VolumeSource.Projected.Sources[0].Secret.Name, getObjName(options, "-db-secret"))

	if externalSecretName == "" {
	    require.Equal(t, pegaInstallerCredentialVolume.VolumeSource.Projected.Sources[1].Secret.Name, getObjName(options, "-upgrade-rest-secret"))
	} else {
	    require.Equal(t, pegaInstallerCredentialVolume.VolumeSource.Projected.Sources[1].Secret.Name, externalSecretName)
	}

	require.Equal(t, pegaInstallerCredentialVolume.VolumeSource.Projected.DefaultMode, volDefaultModePointer)

	var pegaVolumeInstaller = findNamedVolume(volumes, "pega-volume-installer")
	require.NotNil(t, pegaVolumeInstaller)
	if pegaVolumeInstaller.VolumeSource.ConfigMap.LocalObjectReference.Name == "pega-install-config" {
		require.Equal(t, pegaVolumeInstaller.VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-install-config")
	}
	if pegaVolumeInstaller.VolumeSource.ConfigMap.LocalObjectReference.Name == "pega-upgrade-config" {
		require.Equal(t, pegaVolumeInstaller.VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-upgrade-config")
	}
	require.Equal(t, pegaVolumeInstaller.VolumeSource.ConfigMap.DefaultMode, volDefaultModePointer)

	require.Equal(t, string(jobContainers[0].ImagePullPolicy), pullPolicy)

	require.Equal(t, jobContainers[0].Name, expectedJob.containerName)
	require.Equal(t, "YOUR_INSTALLER_IMAGE:TAG", jobContainers[0].Image)
	require.Equal(t, jobContainers[0].Ports[0].ContainerPort, containerPort)

	var pegaVolumeInstallerMount = findNamedVolumeMount(volumeMounts, "pega-volume-installer")
    require.NotNil(t, pegaVolumeInstallerMount)
	require.Equal(t, pegaVolumeInstallerMount.MountPath, "/opt/pega/config")

    var pegaInstallerCredentialVolumeMount = findNamedVolumeMount(volumeMounts, "pega-installer-credentials-volume")
    require.NotNil(t, pegaInstallerCredentialVolumeMount)
	require.Equal(t, pegaInstallerCredentialVolumeMount.MountPath, "/opt/pega/secrets")

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
	VerifyInitContainerData(t, actualInitContainers, options, "job")
}

func TestPegaInstallerJobWithNodeSelector(t *testing.T) {
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.deployment.name":         "install-ns",
			"global.provider":                "k8s",
			"global.actions.execute":         "install",
			"installer.imagePullPolicy":      "Always",
			"installer.upgrade.upgradeType":  "zero-downtime",
			"installer.nodeSelector.mylabel": "myvalue",
		},
	}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

	require.Equal(t, "myvalue", jobObj.Spec.Template.Spec.NodeSelector["mylabel"])

}

func TestPegaInstallerJobWithOverriddenImageAndResources(t *testing.T) {
	var options = &helm.Options{
	    ValuesFiles: []string{"data/pega-installer-override-resources_values.yaml"},
		SetValues: map[string]string{
			"global.deployment.name":         "install-ns",
			"global.provider":                "k8s",
			"global.actions.execute":         "install",
			"installer.imagePullPolicy":      "Always",
			"installer.upgrade.upgradeType":  "zero-downtime",
		},
	}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

    require.Equal(t, "MY_INSTALLER_IMAGE:TAG", jobObj.Spec.Template.Spec.Containers[0].Image)

    require.Equal(t, "200m", jobObj.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
    require.Equal(t, "1Gi", jobObj.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())

    require.Equal(t, "200m", jobObj.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
    require.Equal(t, "1Gi", jobObj.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

func TestPegaInstallerJobWithNonOverriddenImageAndResources(t *testing.T) {
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.deployment.name":         "install-ns",
			"global.provider":                "k8s",
			"global.actions.execute":         "install",
			"installer.imagePullPolicy":      "Always",
			"installer.upgrade.upgradeType":  "zero-downtime",
		},
	}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

    require.Equal(t, "YOUR_INSTALLER_IMAGE:TAG", jobObj.Spec.Template.Spec.Containers[0].Image)

    require.Equal(t, "1", jobObj.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
    require.Equal(t, "12Gi", jobObj.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())

    require.Equal(t, "2", jobObj.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
    require.Equal(t, "12Gi", jobObj.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

func TestPegaInstallerJobWithAffinity(t *testing.T) {

	var affinityBasePath = "installer.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0]."
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.deployment.name":        "install-ns",
					"global.provider":               vendor,
					"global.actions.execute":        operation,
					"installer.imagePullPolicy":     "Always",
					"installer.upgrade.upgradeType": "zero-downtime",
					affinityBasePath + "key":        "kubernetes.io/os",
					affinityBasePath + "operator":   "In",
					affinityBasePath + "values[0]":  "linux",
				},
			}

			helmChartPath, err := filepath.Abs(PegaHelmChartPath)
			require.NoError(t, err)

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")

			var jobObj k8sbatch.Job
			UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

			jobAffinity := jobObj.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution

			require.Equal(t, "kubernetes.io/os", jobAffinity.NodeSelectorTerms[0].MatchExpressions[0].Key)
			require.Equal(t, "In", string(jobAffinity.NodeSelectorTerms[0].MatchExpressions[0].Operator))
			require.Equal(t, "linux", jobAffinity.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
		}
	}
}

func TestPegaInstallerJobWithTolerations(t *testing.T) {

	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.deployment.name":            "install-ns",
					"global.provider":                   vendor,
					"global.actions.execute":            operation,
					"installer.imagePullPolicy":         "Always",
					"installer.upgrade.upgradeType":     "zero-downtime",
					"installer.tolerations[0].key":      "key1",
					"installer.tolerations[0].operator": "Equal",
					"installer.tolerations[0].value":    "value1",
					"installer.tolerations[0].effect":   "NoSchedule",
				},
			}

			helmChartPath, err := filepath.Abs(PegaHelmChartPath)
			require.NoError(t, err)

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")

			var jobObj k8sbatch.Job
			UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

			jobTolerations := jobObj.Spec.Template.Spec.Tolerations

			require.Equal(t, "key1", jobTolerations[0].Key)
			require.Equal(t, "Equal", string(jobTolerations[0].Operator))
			require.Equal(t, "value1", jobTolerations[0].Value)
			require.Equal(t, "NoSchedule", string(jobTolerations[0].Effect))
		}
	}
}

func TestPegaInstallerJobResourcesWithEphemeralStorage(t *testing.T) {
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.deployment.name":                        "pega",
			"global.provider":                               "k8s",
			"global.actions.execute":                        "install",
			"installer.resources.requests.ephemeralStorage": "10Gi",
			"installer.resources.limits.ephemeralStorage":   "11Gi",
			"installer.imagePullPolicy":                     "Always",
		},
	}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

	require.Equal(t, "11Gi", jobObj.Spec.Template.Spec.Containers[0].Resources.Limits.StorageEphemeral().String())
	require.Equal(t, "10Gi", jobObj.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().String())
}

func TestPegaInstallerJobResourcesWithNoEphemeralStorage(t *testing.T) {
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.deployment.name":    "pega",
			"global.provider":           "k8s",
			"global.actions.execute":    "install",
			"installer.imagePullPolicy": "Always",
		},
	}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
	yamlSplit := strings.Split(yamlContent, "---")

	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, yamlSplit[1], &jobObj)

	require.Equal(t, "0", jobObj.Spec.Template.Spec.Containers[0].Resources.Limits.StorageEphemeral().String())
	require.Equal(t, "0", jobObj.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().String())

}


func TestPegaInstallerJobNoICDownload(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}

    helmChartPath, err := filepath.Abs(PegaHelmChartPath)
    require.NoError(t, err)

    testsPath,err := filepath.Abs(PegaHelmChartTestsPath)
    require.NoError(t, err)

	for _, vendor := range supportedVendors {
    	for _, operation := range supportedOperations {
            var options = &helm.Options{
                SetValues: map[string]string{
                    "global.deployment.name":    "pega",
                    "global.provider":           vendor,
                    "global.actions.execute":    operation,
                    "installer.imagePullPolicy": "Always",
                    "installer.upgrade.upgradeType":  "zero-downtime",
                },
            }

            yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"}, "--values", testsPath + "/data/values_multidriver.yaml")
            yamlSplit := strings.Split(yamlContent, "---")

            for _, yaml := range yamlSplit {
                trimmedYaml := strings.TrimSpace(yaml)
                if (len(trimmedYaml) > 0 && !strings.HasPrefix(trimmedYaml, "#")) { //filter out empty chunks of yaml and comments
                    assertJobNoICDownloader(t, yaml, options)
                }
            }
        }
    }
}

func assertJobNoICDownloader(t *testing.T, yaml string, options *helm.Options) {
    var job k8sbatch.Job
    UnmarshalK8SYaml(t, yaml, &job)

    var jobObj = job.Spec.Template.Spec
    var initContainers = jobObj.InitContainers
    var volumes = jobObj.Volumes
    var volumeMounts = jobObj.Containers[0].VolumeMounts

    var ic = findNamedInitContainer(initContainers, "jdbc-lib-downloader")
    require.Equal(t, (*k8score.Container)(nil), ic)

    var jdbcLibVolume = findNamedVolume(volumes, "jdbc-lib-volume")
    require.Equal(t, (*k8score.Volume)(nil), jdbcLibVolume)

    var scriptVolume = findNamedVolume(volumes, "download-script-volume")
    require.Equal(t, (*k8score.Volume)(nil), scriptVolume)

    var jdbcLibVolumeMount = findNamedVolumeMount(volumeMounts, "jdbc-lib-volume")
    require.Equal(t, (*k8score.VolumeMount)(nil), jdbcLibVolumeMount)
}

func TestPegaInstallerJobWithICDownload(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}

    helmChartPath, err := filepath.Abs(PegaHelmChartPath)
    require.NoError(t, err)

    testsPath,err := filepath.Abs(PegaHelmChartTestsPath)
    require.NoError(t, err)

	for _, vendor := range supportedVendors {
    	for _, operation := range supportedOperations {
            var options = &helm.Options{
                SetValues: map[string]string{
                    "global.deployment.name":    "pega",
                    "global.provider":           vendor,
                    "global.actions.execute":    operation,
                    "installer.imagePullPolicy": "Always",
                    "installer.upgrade.upgradeType":  "zero-downtime",
                    "global.downloadContainer.image": "IC_DOWNLOAD_CONTAINER:1.0",
                },
            }

            yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"}, "--values", testsPath + "/data/values_multidriver.yaml")
            yamlSplit := strings.Split(yamlContent, "---")

            for _, yaml := range yamlSplit {
                trimmedYaml := strings.TrimSpace(yaml)
                if (len(trimmedYaml) > 0 && !isOnlyYamlComments(trimmedYaml)) { //filter out empty chunks of yaml and comments
                    assertJobICDownloadComponents(t, yaml, options, false, "10Mi")
                }
            }
        }
    }
}

func TestPegaInstallerJobWithICDownloadWithCert(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}

    helmChartPath, err := filepath.Abs(PegaHelmChartPath)
    require.NoError(t, err)

    testsPath,err := filepath.Abs(PegaHelmChartTestsPath)
    require.NoError(t, err)

	for _, vendor := range supportedVendors {
    	for _, operation := range supportedOperations {
            var options = &helm.Options{
                SetValues: map[string]string{
                    "global.deployment.name":    "pega",
                    "global.provider":           vendor,
                    "global.actions.execute":    operation,
                    "installer.imagePullPolicy": "Always",
                    "installer.upgrade.upgradeType":  "zero-downtime",
                    "global.downloadContainer.image": "IC_DOWNLOAD_CONTAINER:1.0",
                    "global.downloadContainer.sharedVolumeSize": "50Mi",
                },
            }

            yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"}, "--values", testsPath + "/data/values_multidriver-with-cert.yaml")
            yamlSplit := strings.Split(yamlContent, "---")

            for _, yaml := range yamlSplit {
                trimmedYaml := strings.TrimSpace(yaml)
                if (len(trimmedYaml) > 0 && !isOnlyYamlComments(trimmedYaml)) { //filter out empty chunks of yaml and comments
                    assertJobICDownloadComponents(t, yaml, options, true, "50Mi")
                }
            }
        }
    }
}

func isOnlyYamlComments(s string) bool {
    lines := strings.Split(s, "\n")
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if len(trimmed) > 0 && !strings.HasPrefix(trimmed, "#") {
            return false
        }
    }
    return true
}

func assertJobICDownloadComponents(t *testing.T, yaml string, options *helm.Options,  shouldHaveCert bool, volSize string) {
    var job k8sbatch.Job
	UnmarshalK8SYaml(t, yaml, &job)

	var jobSpec = job.Spec.Template.Spec
	var actualInitContainers = jobSpec.InitContainers
	var volumes = jobSpec.Volumes
	var volumeMounts = jobSpec.Containers[0].VolumeMounts

	assertDownloaderIC(t, findNamedInitContainer(actualInitContainers, "jdbc-lib-downloader"), "http://driverhost/drivers/driver1.jar,http://driverhost/drivers/driver2.jar", shouldHaveCert)

    var jdbcLibVolume = findNamedVolume(volumes, "jdbc-lib-volume")
    require.NotNil(t, jdbcLibVolume)
    require.Equal(t, volSize, jdbcLibVolume.VolumeSource.EmptyDir.SizeLimit.String())

    var scriptVolume = findNamedVolume(volumes, "download-script-volume")
    require.NotNil(t, scriptVolume)
    require.Equal(t, int32(0555), *scriptVolume.ConfigMap.DefaultMode)

    var jdbcLibVolumeMount = findNamedVolumeMount(volumeMounts, "jdbc-lib-volume")
    require.NotNil(t, jdbcLibVolumeMount)
    require.Equal(t, "/opt/pega/lib", jdbcLibVolumeMount.MountPath)
}