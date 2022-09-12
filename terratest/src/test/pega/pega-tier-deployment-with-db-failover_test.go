package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"strconv"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
)

func TestPegaTierDeploymentWithDBFailover(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

    var dbFailoverConfigs = []dbFailoverConfig {
                dbFailoverConfig {3, 4, 3},
                dbFailoverConfig{3, 40, 5},
                dbFailoverConfig {4, 3, 3},
                dbFailoverConfig{40, 3, 5},
        }

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, dbFailoverConfig := range dbFailoverConfigs {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"global.classloading.maxRetries": strconv.Itoa(dbFailoverConfig.maxRetries),
						"global.classloading.retryTimeout": strconv.Itoa(dbFailoverConfig.retryTimeout),
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				assertLivenessProbeFailureThreshold(t, yamlSplit[1], dbFailoverConfig.expectedFailureThreshold)
				assertLivenessProbeFailureThreshold(t, yamlSplit[2], dbFailoverConfig.expectedFailureThreshold)
				assertLivenessProbeFailureThreshold(t, yamlSplit[3], dbFailoverConfig.expectedFailureThreshold)
			}
		}
	}
}

func assertLivenessProbeFailureThreshold(t *testing.T, yaml string, expectedFailureThreshold int) {
    if (strings.Contains(yaml, "kind: StatefulSet")) {
        var statefulsetObj appsv1beta2.StatefulSet
	    UnmarshalK8SYaml(t, yaml, &statefulsetObj)
	    VerifyPodSpecLivenessFailureThreshold(t, &statefulsetObj.Spec.Template.Spec, expectedFailureThreshold)
    } else {
	    var deploymentObj appsv1.Deployment
	    UnmarshalK8SYaml(t, yaml, &deploymentObj)
	    VerifyPodSpecLivenessFailureThreshold(t, &deploymentObj.Spec.Template.Spec,  expectedFailureThreshold)
    }
}

func VerifyPodSpecLivenessFailureThreshold(t *testing.T, pod *k8score.PodSpec, expectedFailureThreshold int) {
	require.Equal(t, int32(expectedFailureThreshold), pod.Containers[0].LivenessProbe.FailureThreshold)
}

type dbFailoverConfig struct {
	maxRetries                  int
	retryTimeout                int
	expectedFailureThreshold    int
}
