package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/api/policy/v1beta1"
)

// TestPegaTierPDBEnabled - verify that a PodDisruptionBudget is created when global.tier.pdb.enabled=true
func TestPegaTierPDBEnabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testsPath, err := filepath.Abs(PegaHelmChartTestsPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-pdb.yaml"}, "--values", testsPath+"/data/values_pdb_enabled.yaml")
				verifyPegaPDBs(t, yamlContent, options, []pdb{
					{
						name:         getObjName(options, "-web-pdb"),
						kind:         "PodDisruptionBudget",
						apiversion:   "policy/v1beta1",
						minAvailable: 1,
					},
					{
						name:         getObjName(options, "-batch-pdb"),
						kind:         "PodDisruptionBudget",
						apiversion:   "policy/v1beta1",
						minAvailable: 1,
					},
					{
						name:         getObjName(options, "-stream-pdb"),
						kind:         "PodDisruptionBudget",
						apiversion:   "policy/v1beta1",
						minAvailable: 1,
					},
				})
			}
		}
	}
}

// TestPegaTierPDBDisabled - verify that a PodDisruptionBudget is not created when global.tier.pdb.enabled=false
func TestPegaTierPDBDisabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testsPath, err := filepath.Abs(PegaHelmChartTestsPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name": depName,
						"global.provider":        vendor,
						"global.actions.execute": operation,
					},
				}

				_, err := RenderTemplateWithErr(t, options, helmChartPath, []string{"templates/pega-tier-pdb.yaml"}, "--values", testsPath+"/data/values_pdb_disabled.yaml")
				require.NotNil(t, err)
			}
		}
	}
}

// verifyPegaPDBs - Splits the PDB object from the rendered template and asserts each PDB object
func verifyPegaPDBs(t *testing.T, yamlContent string, options *helm.Options, expectedPdbs []pdb) {
	var pegaPdbObj v1beta1.PodDisruptionBudget
	pdbSlice := strings.SplitAfter(yamlContent, "---")
	pdbSlice = pdbSlice[1:]

	require.Equal(t, len(expectedPdbs), len(pdbSlice))

	for i, pdb := range pdbSlice {
		UnmarshalK8SYaml(t, pdb, &pegaPdbObj)
		verifyPegaPdb(t, &pegaPdbObj, expectedPdbs[i])
	}
}

// verifyPegaPdb - Performs Pega PDB assertions with the values as provided
func verifyPegaPdb(t *testing.T, pegaPdbObj *v1beta1.PodDisruptionBudget, expectedPdb pdb) {
	require.Equal(t, pegaPdbObj.TypeMeta.Kind, expectedPdb.kind)
	//if the below fails it means that the helm version used in testing is compiled against
	//kubernetes 1.21 or higher, and we should adjust this test to use the policy/v1 API version
	require.Equal(t, pegaPdbObj.TypeMeta.APIVersion, expectedPdb.apiversion)
	require.Equal(t, expectedPdb.minAvailable, pegaPdbObj.Spec.MinAvailable.IntVal)
}

type pdb struct {
	name         string
	kind         string
	apiversion   string
	minAvailable int32
}
