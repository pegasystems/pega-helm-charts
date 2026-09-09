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

func TestPegaHostAliases(t *testing.T) {
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
					ValuesFiles: []string{"data/values_hostAliases.yaml"},
					SetValues: map[string]string{
						"global.provider":                        vendor,
						"global.actions.execute":                 operation,
						"global.deployment.name":                 deploymentName,
					},
				}

				// Render the Kubernetes manifests using Helm
				output := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})

				// Create a YAML decoder from the output
				decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(output)), 4096)
				for {
					var resource unstructured.Unstructured
					err := decoder.Decode(&resource)
					if err != nil {
						// Break on EOF
						break
					}

                    // Only check Deployment resources
                    if resource.GetKind() == "Deployment" {
                        // Extract and validate hostAliases
                        spec, found, err := unstructured.NestedMap(resource.Object, "spec", "template", "spec")
                        require.NoError(t, err, "Error extracting spec from resource")
                        if !found {
                            t.Errorf("spec.template.spec not found in Deployment %s", resource.GetName())
                            continue
                        }

                        hostAliases, found, err := unstructured.NestedSlice(spec, "hostAliases")
                        require.NoError(t, err, "Error extracting hostAliases from spec")
                        if !found {
                            t.Errorf("hostAliases not found in Deployment %s", resource.GetName())
                            continue
                        }

                        // Validate the content of hostAliases
                        assert.Equal(t, 1, len(hostAliases), "Expected exactly 1 hostAliases entry in Deployment %s", resource.GetName())
                        entry := hostAliases[0].(map[string]interface{})
                        assert.Equal(t, "127.0.0.1", entry["ip"], "Unexpected IP in hostAliases of Deployment %s", resource.GetName())
                        assert.Contains(t, entry["hostnames"], "test1.local", "Expected hostname test1.local in hostAliases of Deployment %s", resource.GetName())
                        assert.Contains(t, entry["hostnames"], "test2.local", "Expected hostname test2.local in hostAliases of Deployment %s", resource.GetName())
                    }
				}
			}
		}
	}
}
