package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	v1 "k8s.io/api/core/v1"
)

func TestPegaTierHPA(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}
    var deploymentNames = []string{"pega","myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
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
						"installer.upgrade.upgradeType": "zero-downtime",
                    },
                }

                yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"})
                verifyPegaHPAs(t, yamlContent, options, []hpa{
                    {
                        name:          getObjName(options, "-web-hpa"),
                        targetRefName: getObjName(options, "-web"),
                        kind:          "Deployment",
                        apiversion:    "apps/v1",
                        cpu:           true,
                        cpuValue:      parseResourceValue(t, "2.55"),
                    },
                    {
                        name:          getObjName(options, "-batch-hpa"),
                        targetRefName: getObjName(options, "-batch"),
                        kind:          "Deployment",
                        apiversion:    "apps/v1",
                        cpu:           true,
                        cpuValue:      parseResourceValue(t, "2.55"),
                    },
                })
            }
		}
	}
}


func TestPegaTierHPADisableTarget(t *testing.T) {
    var supportedVendors = []string{"k8s", "openshift", "eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}
    var deploymentNames = []string{"pega","myapp-dev"}

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
						"installer.upgrade.upgradeType": "zero-downtime",
                    },
                }


			    yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"}, "--values", testsPath+"/data/values_hpa_disabletarget.yaml")
                verifyPegaHPAs(t, yamlContent, options, []hpa{
                    {
                        name:          getObjName(options, "-web-hpa"),
                        targetRefName: getObjName(options, "-web"),
                        kind:          "Deployment",
                        apiversion:    "apps/v1",
                        mem:           true,
                        memPercent:    85,
                    },
                    {
                        name:          getObjName(options, "-batch-hpa"),
                        targetRefName: getObjName(options, "-batch"),
                        kind:          "Deployment",
                        apiversion:    "apps/v1",
                        cpu:           true,
                        cpuValue:      parseResourceValue(t, "2.55"),
                    },
                })
            }
		}
	}
}


func TestPegaTierOverrideValues(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
    var deploymentNames = []string{"pega","myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testsPath, err := filepath.Abs(PegaHelmChartTestsPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

            for _, depName := range deploymentNames {
                fmt.Println(vendor + "-" + operation + "-" + depName)

                var options = &helm.Options{
                    SetValues: map[string]string{
                        "global.provider":        vendor,
                        "global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
                    },
                }

                yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"}, "--values", testsPath+"/data/values_hpa_overridevalues.yaml")
                verifyPegaHPAs(t, yamlContent, options, []hpa{
                    {
                        name:          getObjName(options, "-web-hpa"),
                        targetRefName: getObjName(options, "-web"),
                        kind:          "Deployment",
                        apiversion:    "apps/v1",
                        cpu:           true,
                        cpuValue:      parseResourceValue(t, "4.13"),
                        mem:           true,
                        memPercent:    42,
                    },
                    {
                        name:          getObjName(options, "-batch-hpa"),
                        targetRefName: getObjName(options, "-batch"),
                        kind:          "Deployment",
                        apiversion:    "apps/v1",
                        cpu:           true,
                        cpuPercent:    24,
                    },
                })
			}
		}
	}
}

// verifyPegaHPAs - Splits the HPA object from the rendered template and asserts each HPA object
func verifyPegaHPAs(t *testing.T, yamlContent string, options *helm.Options, expectedHpas []hpa) {
	var pegaHpaObj autoscaling.HorizontalPodAutoscaler
	hpaSlice := strings.SplitAfter(yamlContent, "---")
	hpaSlice = hpaSlice[1:]

	require.Equal(t, len(expectedHpas), len(hpaSlice))

	for i, hpa := range hpaSlice {
		UnmarshalK8SYaml(t, hpa, &pegaHpaObj)
		verifyPegaHpa(t, &pegaHpaObj, expectedHpas[i])
	}
}

// verifyPegaHpa - Performs Pega HPA assertions with the values as provided in default values.yaml
func verifyPegaHpa(t *testing.T, hpaObj *autoscaling.HorizontalPodAutoscaler, expectedHpa hpa) {
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Name, expectedHpa.targetRefName)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Kind, expectedHpa.kind)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.APIVersion, expectedHpa.apiversion)
	currentMetricIndex := 0

	require.Equal(t, expectedHpa.expectedMetricCount(), len(hpaObj.Spec.Metrics))

	if expectedHpa.cpu {
		require.Equal(t, v1.ResourceName("cpu"), hpaObj.Spec.Metrics[currentMetricIndex].Resource.Name)
		if expectedHpa.cpuValue != (resource.Quantity{}) {
			require.Equal(t, autoscaling.MetricTargetType("Value"), hpaObj.Spec.Metrics[currentMetricIndex].Resource.Target.Type)
			require.Equal(t, expectedHpa.cpuValue, *hpaObj.Spec.Metrics[currentMetricIndex].Resource.Target.AverageValue)
		}
		if expectedHpa.cpuPercent != 0 {
			require.Equal(t, autoscaling.MetricTargetType("Utilization"), hpaObj.Spec.Metrics[currentMetricIndex].Resource.Target.Type)
			require.Equal(t, expectedHpa.cpuPercent, *hpaObj.Spec.Metrics[currentMetricIndex].Resource.Target.AverageUtilization)
		}
		currentMetricIndex++
	}
	if expectedHpa.mem {
		require.Equal(t, v1.ResourceName("memory"), hpaObj.Spec.Metrics[currentMetricIndex].Resource.Name)
		require.Equal(t, autoscaling.MetricTargetType("Utilization"), hpaObj.Spec.Metrics[currentMetricIndex].Resource.Target.Type)
		require.Equal(t, expectedHpa.memPercent, *hpaObj.Spec.Metrics[currentMetricIndex].Resource.Target.AverageUtilization)
		currentMetricIndex++
	}

	require.Equal(t, int32(5), hpaObj.Spec.MaxReplicas)
}

type hpa struct {
	name          string
	targetRefName string
	kind          string
	apiversion    string
	cpu           bool
	cpuValue      resource.Quantity
	cpuPercent    int32
	mem           bool
	memPercent    int32
}

func (h hpa) expectedMetricCount() int {
	result := 0
	if h.cpu {
		result++
	}

	if h.mem {
		result++
	}

	return result
}

func parseResourceValue(t *testing.T, resourceString string) resource.Quantity {
	value, err := resource.ParseQuantity(resourceString)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return value
}
