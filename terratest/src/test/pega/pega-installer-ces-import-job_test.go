package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sbatch "k8s.io/api/batch/v1"
)

const cesImportJobTemplate = "charts/installer/templates/pega-installer-ces-import-job.yaml"
const cesImage = "pega-docker.downloads.pega.com/platform-services/ces:1.2.2"

func renderCesImportJob(t *testing.T, values map[string]string) (string, error) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	values["global.provider"] = "k8s"
	return RenderTemplateE(t, &helm.Options{SetValues: values}, helmChartPath, []string{cesImportJobTemplate})
}

func TestPegaInstallerCesImportJob(t *testing.T) {
	for _, operation := range []string{"deploy", "install-deploy", "upgrade-deploy"} {
		yamlContent, err := renderCesImportJob(t, map[string]string{
			"global.actions.execute":        operation,
			"installer.upgrade.upgradeType": "zero-downtime",
			"installer.image":               "platform/installer:26.1.0",
			"installer.ces.image":           cesImage,
		})
		require.NoError(t, err)

		var job k8sbatch.Job
		UnmarshalK8SYaml(t, yamlContent, &job)
		require.Regexp(t, "^pega-ces-import-[0-9a-f]{10}$", job.Name)
		require.Equal(t, int32(0), *job.Spec.BackoffLimit)

		pod := job.Spec.Template.Spec
		require.Equal(t, "wait-for-rolling-updates", pod.InitContainers[0].Name)
		require.Equal(t, "ces-files", pod.InitContainers[1].Name)
		require.Equal(t, cesImage, pod.InitContainers[1].Image)
		require.Equal(t, "/work", pod.InitContainers[1].VolumeMounts[0].MountPath)

		importer := pod.Containers[0]
		require.Equal(t, "platform/installer:26.1.0", importer.Image)
		require.Equal(t, []string{"sh", "/opt/pega/ces/import-ces.sh"}, importer.Command)
		require.Equal(t, "/opt/pega/ces", importer.VolumeMounts[0].MountPath)
		require.Equal(t, "/opt/pega/secrets", importer.VolumeMounts[1].MountPath)
		require.Equal(t, "http://pega-web:80/prweb/PRRestService", importer.Env[1].Value)
		require.Equal(t, "skip", importer.Env[2].Value)
		require.Equal(t, "pega-upgrade-rest-secret", pod.Volumes[1].Projected.Sources[0].Secret.Name)
	}
}

func TestPegaInstallerCesImportJobNotRendered(t *testing.T) {
	cases := []map[string]string{
		{"global.actions.execute": "deploy"},
		{"global.actions.execute": "install", "installer.ces.image": cesImage},
		{"global.actions.execute": "upgrade", "installer.ces.image": cesImage},
	}
	for _, values := range cases {
		_, err := renderCesImportJob(t, values)
		require.Error(t, err, "CES import Job must not render for %v", values)
	}
}

func TestPegaInstallerCesImportJobNameFollowsImage(t *testing.T) {
	names := map[string]bool{}
	for _, tag := range []string{"1.2.2", "1.2.2", "1.2.3"} {
		yamlContent, err := renderCesImportJob(t, map[string]string{
			"global.actions.execute": "deploy",
			"installer.ces.image":    "pega-docker.downloads.pega.com/platform-services/ces:" + tag,
		})
		require.NoError(t, err)
		var job k8sbatch.Job
		UnmarshalK8SYaml(t, yamlContent, &job)
		names[job.Name] = true
	}
	require.Len(t, names, 2)
}

func TestPegaInstallerCesImportJobWithExternalRESTSecretAndURL(t *testing.T) {
	yamlContent, err := renderCesImportJob(t, map[string]string{
		"global.actions.execute":                           "deploy",
		"installer.ces.image":                              cesImage,
		"installer.ces.importExistingInstances":            "override",
		"installer.upgrade.pega_rest_external_secret_name": "my-pega-rest-secret",
		"installer.upgrade.pegaRESTServerURL":              "https://pega.example.com/prweb/PRRestService",
	})
	require.NoError(t, err)
	var job k8sbatch.Job
	UnmarshalK8SYaml(t, yamlContent, &job)
	pod := job.Spec.Template.Spec
	require.Equal(t, "my-pega-rest-secret", pod.Volumes[1].Projected.Sources[0].Secret.Name)
	require.Equal(t, "https://pega.example.com/prweb/PRRestService", pod.Containers[0].Env[1].Value)
	require.Equal(t, "override", pod.Containers[0].Env[2].Value)
}

func TestPegaInstallerRoleWithCesImportOnDeploy(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	yamlContent := RenderTemplate(t, &helm.Options{SetValues: map[string]string{
		"global.provider":        "k8s",
		"global.actions.execute": "deploy",
		"installer.ces.image":    cesImage,
	}}, helmChartPath, []string{"charts/installer/templates/pega-installer-role.yaml"})
	assertInstallerRole(t, yamlContent)
}
