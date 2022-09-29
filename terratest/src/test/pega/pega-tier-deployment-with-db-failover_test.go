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

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

            fmt.Println(vendor + "-" + operation)

            var options = &helm.Options{
                ValuesFiles: []string{"data/values_with_overidden_liveness_probe_config.yaml"},
                SetValues: map[string]string{
                    "global.provider":               vendor,
                    "global.actions.execute":        operation,
                    "installer.upgrade.upgradeType": "zero-downtime",
                },
            }

            yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
            yamlSplit := strings.Split(yamlContent, "---")

            //web tier uses defaults not in values.yaml (failure threshold: 3, period seconds: 30)
            assertLivenessProbeFailureThreshold(t, yamlSplit[1], 4, 30)

            //batch tier uses (failure threshold: 5, period seconds: 100)
            assertLivenessProbeFailureThreshold(t, yamlSplit[2], 6, 100)

            //stream tier uses (failure threshold: 7, period seconds: 300) -- retryTimeout has max of 180s
            //so we get into more interesting math...
            assertLivenessProbeFailureThreshold(t, yamlSplit[3], 12, 180)
		}
	}
}

func assertLivenessProbeFailureThreshold(t *testing.T, yaml string, expectedMaxRetries int, expectedRetryTimeout int) {
    if (strings.Contains(yaml, "kind: StatefulSet")) {
        var statefulsetObj appsv1beta2.StatefulSet
	    UnmarshalK8SYaml(t, yaml, &statefulsetObj)
	    VerifyPodSpecClassloaderRetrySettings(t, &statefulsetObj.Spec.Template.Spec, expectedMaxRetries, expectedRetryTimeout)
    } else {
	    var deploymentObj appsv1.Deployment
	    UnmarshalK8SYaml(t, yaml, &deploymentObj)
	    VerifyPodSpecClassloaderRetrySettings(t, &deploymentObj.Spec.Template.Spec,  expectedMaxRetries, expectedRetryTimeout)
    }
}

func VerifyPodSpecClassloaderRetrySettings(t *testing.T, pod *k8score.PodSpec, expectedMaxRetries int, expectedRetryTimeout int) {
	maxRetries := ""
	retryTimeout := ""

	for _, envItem := range pod.Containers[0].Env {
	    if (envItem.Name=="RETRY_TIMEOUT") {
            retryTimeout = envItem.Value
	    }
        if (envItem.Name=="MAX_RETRIES") {
            maxRetries = envItem.Value
        }
	}

	require.NotEqual(t, "", retryTimeout)
	require.NotEqual(t, "", maxRetries)

	require.Equal(t, strconv.Itoa(expectedRetryTimeout), retryTimeout)
	require.Equal(t, strconv.Itoa(expectedMaxRetries), maxRetries)
}

