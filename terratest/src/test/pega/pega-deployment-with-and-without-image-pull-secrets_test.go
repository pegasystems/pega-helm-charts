package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithoutImagePullSecrets(t *testing.T) {

	var supportedVendors = []string{"k8s", "eks", "gke", "aks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(deploymentYaml, "---")
			assertWithoutImagePullSecrets(t, yamlSplit[1])
		}
	}
}

func TestPegaDeploymentWithImagePullSecrets(t *testing.T) {

	var supportedVendors = []string{"k8s", "eks", "gke", "aks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                    vendor,
					"global.actions.execute":             operation,
					"global.docker.imagePullSecretNames": "{\"secret1\",\"secret2\"}",
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(deploymentYaml, "---")
			assertWithImagePullSecrets(t, yamlSplit[1])
		}
	}
}

func TestHazelcastDeploymentWithoutImagePullSecrets(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"hazelcast.enabled":      "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/pega-hz-deployment.yaml"})
			assertWithoutImagePullSecrets(t, yamlContent)

		}
	}
}

func TestHazelcastDeploymentWithImagePullSecrets(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                    vendor,
					"global.actions.execute":             operation,
					"hazelcast.enabled":                  "true",
					"global.docker.imagePullSecretNames": "{\"secret1\",\"secret2\"}",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/pega-hz-deployment.yaml"})
			assertWithImagePullSecrets(t, yamlContent)

		}
	}
}

func TestPegaSearchDeploymentWithoutImagePullSecrets(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
				assertWithoutImagePullSecrets(t, yamlContent)
			}
		}
	}
}

func TestPegaSearchDeploymentWithImagePullSecrets(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":             depName,
						"global.provider":                    vendor,
						"global.actions.execute":             operation,
						"global.docker.imagePullSecretNames": "{\"secret1\",\"secret2\"}",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-deployment.yaml"})
				assertWithImagePullSecrets(t, yamlContent)
			}
		}
	}
}

func TestConstellationDeploymentWithoutImagePullSecrets(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"constellation.enabled":  "true",
				},
			}
			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/constellation/templates/clln-deployment.yaml"})
			assertWithoutImagePullSecrets(t, yamlContent)
		}
	}
}

func TestConstellationDeploymentWithImagePullSecrets(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                    vendor,
					"global.actions.execute":             operation,
					"constellation.enabled":              "true",
					"global.docker.imagePullSecretNames": "{\"secret1\",\"secret2\"}",
				},
			}
			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/constellation/templates/clln-deployment.yaml"})
			assertWithImagePullSecrets(t, yamlContent)
		}
	}
}

func TestPegaDeploymentWithoutRegistryBlock(t *testing.T) {

	var supportedVendors = []string{"k8s", "eks", "gke", "aks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"global.docker.registry": "",
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(deploymentYaml, "---")
			assertWithoutRegistryBlock(t, yamlSplit[1])
		}
	}
}

func TestPegaDeploymentWithoutRegistryBlockWithExternalSecrets(t *testing.T) {

	var supportedVendors = []string{"k8s", "eks", "gke", "aks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                    vendor,
					"global.actions.execute":             operation,
					"global.docker.registry":             "",
					"global.docker.imagePullSecretNames": "{\"secret1\",\"secret2\"}",
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(deploymentYaml, "---")
			assertWithoutRegistryBlockWithExternalSecrets(t, yamlSplit[1])
		}
	}
}

func assertWithoutImagePullSecrets(t *testing.T, webYaml string) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, webYaml, &deploymentObj)
	require.Equal(t, deploymentObj.Spec.Template.Spec.ImagePullSecrets[0].Name, "pega-registry-secret")
}

func assertWithImagePullSecrets(t *testing.T, webYaml string) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, webYaml, &deploymentObj)
	imagePullSecrets := deploymentObj.Spec.Template.Spec.ImagePullSecrets
	require.Equal(t, imagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, imagePullSecrets[1].Name, "secret1")
	require.Equal(t, imagePullSecrets[2].Name, "secret2")
	require.Equal(t, len(imagePullSecrets), 3)
}

func assertWithoutRegistryBlock(t *testing.T, webYaml string) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, webYaml, &deploymentObj)
	imagePullSecrets := deploymentObj.Spec.Template.Spec.ImagePullSecrets
	require.Equal(t, len(imagePullSecrets), 0)
}

func assertWithoutRegistryBlockWithExternalSecrets(t *testing.T, webYaml string) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, webYaml, &deploymentObj)
	imagePullSecrets := deploymentObj.Spec.Template.Spec.ImagePullSecrets
	require.Equal(t, len(imagePullSecrets), 2)
	require.Equal(t, imagePullSecrets[0].Name, "secret1")
	require.Equal(t, imagePullSecrets[1].Name, "secret2")
}
