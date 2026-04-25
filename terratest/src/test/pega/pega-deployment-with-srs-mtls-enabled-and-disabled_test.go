package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithSRSandMTLSDisabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks", "openshift"}
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
				assertNoSRSmTLSSettings(t, deployment)
			}
		}
	}
}

func TestPegaDeploymentWithSRSEnabledAndmTLSDisabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks", "openshift"}
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
				assertNoSRSmTLSSettings(t, deployment)
			}
		}
	}
}

func TestPegaDeploymentWithSRSmTLSEnabledAndAllParamsGiven(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks", "openshift"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var supportedExternalSecrets = []string{"", "test-external-secret"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, externalSecret := range supportedExternalSecrets {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                         vendor,
						"global.actions.execute":                  operation,
						"pegasearch.externalSearchService":        "true",
						"pegasearch.srsMTLS.enabled":              "true",
						"pegasearch.srsMTLS.trustStore":           "trustStore.jsk",
						"pegasearch.srsMTLS.keyStore":             "keyStore.jks",
						"pegasearch.srsAuth.external_secret_name": externalSecret,
					},
				}
				deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				deployments := strings.Split(deploymentYaml, "---")
				for _, deployment := range deployments {
					assertHasSRSmTLSSettings(t, deployment, externalSecret)
				}
			}
		}
	}
}

func TestPegaDeploymentWithSRSmTLSVaultEnabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks", "openshift"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                         vendor,
					"global.actions.execute":                  operation,
					"pegasearch.externalSearchService":        "true",
					"pegasearch.srsMTLS.enabled":              "true",
					"pegasearch.srsMTLS.vault.enabled":        "true",
					"pegasearch.srsMTLS.vault.url":            "https://vault.example.com",
					"pegasearch.srsMTLS.vault.role":           "pega-role",
					"pegasearch.srsMTLS.vault.secretPath":     "secret/data/pega/mtls",
					"pegasearch.srsMTLS.vault.tokenSecret":    "vault-token-secret",
					"pegasearch.srsMTLS.vault.tokenSecretKey": "token",
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			deployments := strings.Split(deploymentYaml, "---")
			for _, deployment := range deployments {
				if len(strings.TrimSpace(deployment)) == 0 {
					continue
				}
				assertHasVaultTokenVolume(t, deployment)
				assertNoSRSPasswordSecrets(t, deployment)
			}
		}
	}
}

func assertNoSRSmTLSSettings(t *testing.T, pegaTierDeployment string) {
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, pegaTierDeployment, &deployment)
	for _, container := range deployment.Spec.Template.Spec.Containers {
		for _, envVar := range container.Env {
			if "SRS_TRUSTSTORE_PASSWORD" == envVar.Name {
				require.Fail(t, "container '"+container.Name+"' should not have 'SRS_TRUSTSTORE_PASSWORD' environment variable")
			}
			if "SRS_KEYSTORE_PASSWORD" == envVar.Name {
				require.Fail(t, "container '"+container.Name+"' should not have 'SRS_KEYSTORE_PASSWORD' environment variable")
			}
		}
	}
}

func assertHasSRSmTLSSettings(t *testing.T, pegaTierDeployment string /* authType string, */, externalSecret string) {
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, pegaTierDeployment, &deployment)
	for _, container := range deployment.Spec.Template.Spec.Containers {
		for _, envVar := range container.Env {
			if "SRS_TRUSTSTORE_PASSWORD" == envVar.Name {
				if externalSecret == "" {
					require.Equal(t, "pega-srs-mtls-secret", envVar.ValueFrom.SecretKeyRef.Name)
					require.Equal(t, "SRS_TRUSTSTORE_PASSWORD", envVar.ValueFrom.SecretKeyRef.Key)
				} else {
					require.Equal(t, externalSecret, envVar.ValueFrom.SecretKeyRef.Name)
					require.Equal(t, "SRS_TRUSTSTORE_PASSWORD", envVar.ValueFrom.SecretKeyRef.Key)
				}
			}
			if "SRS_KEYSTORE_PASSWORD" == envVar.Name {
				if externalSecret == "" {
					require.Equal(t, "pega-srs-mtls-secret", envVar.ValueFrom.SecretKeyRef.Name)
					require.Equal(t, "SRS_KEYSTORE_PASSWORD", envVar.ValueFrom.SecretKeyRef.Key)
				} else {
					require.Equal(t, externalSecret, envVar.ValueFrom.SecretKeyRef.Name)
					require.Equal(t, "SRS_KEYSTORE_PASSWORD", envVar.ValueFrom.SecretKeyRef.Key)
				}
			}
		}
	}
}

func assertHasVaultTokenVolume(t *testing.T, pegaTierDeployment string) {
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, pegaTierDeployment, &deployment)

	// Verify vault token volume exists
	foundVolume := false
	for _, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.Name == "pega-vault-token-volume" {
			foundVolume = true
			require.NotNil(t, volume.Secret, "vault token volume should be a secret volume")
			require.Equal(t, "vault-token-secret", volume.Secret.SecretName)
			require.Equal(t, 1, len(volume.Secret.Items))
			require.Equal(t, "token", volume.Secret.Items[0].Key)
			require.Equal(t, "token", volume.Secret.Items[0].Path)
			break
		}
	}
	require.True(t, foundVolume, "deployment should have 'pega-vault-token-volume' volume")

	// Verify vault token volume mount exists
	for _, container := range deployment.Spec.Template.Spec.Containers {
		foundMount := false
		for _, mount := range container.VolumeMounts {
			if mount.Name == "pega-vault-token-volume" {
				foundMount = true
				require.Equal(t, "/opt/pega/secrets/vault", mount.MountPath)
				require.True(t, mount.ReadOnly)
				break
			}
		}
		require.True(t, foundMount, "container '"+container.Name+"' should have vault token volume mount")
	}
}

func assertNoSRSPasswordSecrets(t *testing.T, pegaTierDeployment string) {
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, pegaTierDeployment, &deployment)
	for _, container := range deployment.Spec.Template.Spec.Containers {
		for _, envVar := range container.Env {
			if "SRS_TRUSTSTORE_PASSWORD" == envVar.Name {
				require.Fail(t, "container '"+container.Name+"' should not have 'SRS_TRUSTSTORE_PASSWORD' when vault is enabled")
			}
			if "SRS_KEYSTORE_PASSWORD" == envVar.Name {
				require.Fail(t, "container '"+container.Name+"' should not have 'SRS_KEYSTORE_PASSWORD' when vault is enabled")
			}
		}
	}
}
