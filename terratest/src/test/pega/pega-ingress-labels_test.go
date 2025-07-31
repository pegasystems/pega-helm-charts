package pega

import (
	"testing"
	"bytes"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"path/filepath"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestPegaIngressLabels(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err, "Failed to resolve Helm chart path")

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, deploymentName := range deploymentNames {
				// Set the Helm values to configure the labels
				options := &helm.Options{
					SetValues: map[string]string{
						"global.provider":                        vendor,
						"global.actions.execute":                 operation,
						"global.deployment.name":                 deploymentName,
						"global.tier[0].name":                    "web",
						"global.tier[0].ingress.enabled":         "true",
						"global.tier[0].ingress.labels.test-label": "test-value",
						"global.tier[0].ingress.domain":          "pega.local",
						"global.tier[0].service.port":            "80",
					},
				}

				// Render the Kubernetes manifests using Helm
				output := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})

				// Create a YAML decoder from the output
				decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(output)), 4096)
				for {
					var resource unstructured.Unstructured
					err := decoder.Decode(&resource)
					if err != nil {
						// Break on EOF
						break
					}

					if vendor == "openshift"{
						if resource.GetKind() == "Route" {
							labels := resource.GetLabels()

							// Perform assertions on the labels
							assert.Contains(t, labels, "test-label", "Route object is missing the label test-label")
							assert.Equal(t, "test-value", labels["test-label"], "Expected Route to have label test-label=test-value")
						}
					} else{
						if resource.GetKind() == "Ingress" {
							labels := resource.GetLabels()

							// Perform assertions on the labels
							assert.Contains(t, labels, "test-label", "Ingress object is missing the label test-label")
							assert.Equal(t, "test-value", labels["test-label"], "Expected Ingress to have label test-label=test-value")
						}
					}
				}
			}
		}
	}
}


func TestPegaEKSIngressHealthcheck(t *testing.T) {
	var deploymentNames = []string{"pega", "myapp-dev"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err, "Failed to resolve Helm chart path")

    for _, deploymentName := range deploymentNames {
        // Set the Helm values to configure the labels
        options := &helm.Options{
            SetValues: map[string]string{
                "global.provider":                        "eks",
                "global.actions.execute":                 "deploy",
                "global.deployment.name":                 deploymentName,
                "global.tier[0].name":                    "web",
                "global.tier[0].ingress.enabled":         "true",
                "global.tier[0].ingress.domain":          "pega.local",
                "global.tier[0].service.port":            "80",
            },
        }

        // Render the Kubernetes manifests using Helm
        output := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})

        // Create a YAML decoder from the output
        decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(output)), 4096)
        for {
            var resource unstructured.Unstructured
            err := decoder.Decode(&resource)
            if err != nil {
                // Break on EOF
                break
            }

            assert.Equal(t, resource.GetKind(), "Ingress", "Should be an ingress")
            annotations := resource.GetAnnotations()
            assert.Contains(t, annotations, "alb.ingress.kubernetes.io/healthcheck-path", "Expected ingress to specify healthcheck path annotation.")
            assert.Equal(t, annotations["alb.ingress.kubernetes.io/healthcheck-path"], "/healthz.html", "Annotation should have the correct value.")
        }
    }
}
