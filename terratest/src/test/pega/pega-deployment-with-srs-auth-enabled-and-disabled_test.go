package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithSRSDisabled(t *testing.T) {
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
			deployments := strings.Split(deploymentYaml, "---")
			for _, deployment := range deployments {
				assertNoSRSAuthSettings(t, deployment)
			}
		}
	}
}

func TestPegaDeploymentWithSRSAuthDisabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                  vendor,
					"global.actions.execute":           operation,
					"pegasearch.externalSearchService": "true",
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			deployments := strings.Split(deploymentYaml, "---")
			for _, deployment := range deployments {
				assertNoSRSAuthSettings(t, deployment)
			}
		}
	}
}

func TestPegaDeploymentWithSRSAuthEnabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                  vendor,
					"global.actions.execute":           operation,
					"pegasearch.externalSearchService": "true",
					"pegasearch.srsAuth.enabled":       "true",
					"pegasearch.srsAuth.privateKey":    SRSAuthPrivateKeyExample,
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			deployments := strings.Split(deploymentYaml, "---")
			for _, deployment := range deployments {
				assertHasSRSAuthSettings(t, deployment)
			}
		}
	}
}

func assertNoSRSAuthSettings(t *testing.T, pegaTierDeployment string) {
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, pegaTierDeployment, &deployment)
	for _, container := range deployment.Spec.Template.Spec.Containers {
		for _, envVar := range container.Env {
			if "SERV_AUTH_PRIVATE_KEY" == envVar.Name {
				require.Fail(t, "container '"+container.Name+"' should not have 'SERV_AUTH_PRIVATE_KEY' environment variable")
			}
		}
	}
}

func assertHasSRSAuthSettings(t *testing.T, pegaTierDeployment string) {
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, pegaTierDeployment, &deployment)
	for _, container := range deployment.Spec.Template.Spec.Containers {
		hasPrivateKey := false
		for _, envVar := range container.Env {
			if "SERV_AUTH_PRIVATE_KEY" == envVar.Name {
				require.Equal(t, "pega-srs-auth-secret", envVar.ValueFrom.SecretKeyRef.Name)
				require.Equal(t, "privateKey", envVar.ValueFrom.SecretKeyRef.Key)
				hasPrivateKey = true
			}
		}
		require.True(t, hasPrivateKey, "container '"+container.Name+"' should have 'SERV_AUTH_PRIVATE_KEY' environment variable")
	}
}
