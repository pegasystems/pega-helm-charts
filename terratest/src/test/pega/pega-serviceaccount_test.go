package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPegaServiceAccountsDisabledByDefault(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	options := &helm.Options{SetValues: map[string]string{
		"global.provider":        "k8s",
		"global.actions.execute": "install-deploy",
	}}
	for _, template := range []string{"templates/pega-serviceaccount.yaml", "charts/installer/templates/pega-installer-serviceaccount.yaml"} {
		_, err := RenderTemplateE(t, options, helmChartPath, []string{template})
		require.Error(t, err, template)
		require.Contains(t, err.Error(), "could not find template", template)
	}

	specs := renderPodSpecs(t, options, helmChartPath)
	require.NotEmpty(t, specs)
	for name, spec := range specs {
		require.Empty(t, spec.ServiceAccountName, name)
		require.Nil(t, spec.AutomountServiceAccountToken, name)
	}
}

func TestPegaServiceAccountsCreated(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	options := &helm.Options{SetValues: map[string]string{
		"global.provider":               "k8s",
		"global.actions.execute":        "install-deploy",
		"global.serviceAccount.enabled": "true",
		"global.serviceAccount.create":  "true",
		"global.serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn": "arn:aws:iam::1:role/pega",
		"installer.serviceAccount.automountServiceAccountToken":            "false",
		"constellation.enabled":                    "true",
		"hazelcast.clusteringServiceEnabled":       "true",
		"installer.serviceAccount.enabled":         "true",
		"installer.serviceAccount.create":          "true",
		"global.tier[0].name":                      "web",
		"global.tier[0].custom.serviceAccountName": "web-account",
		"global.tier[1].name":                      "batch",
	}}

	var runtime corev1.ServiceAccount
	UnmarshalK8SYaml(t, RenderTemplate(t, options, helmChartPath, []string{"templates/pega-serviceaccount.yaml"}), &runtime)
	require.Equal(t, "pega-serviceaccount", runtime.Name)
	require.Equal(t, "arn:aws:iam::1:role/pega", runtime.Annotations["eks.amazonaws.com/role-arn"])
	require.True(t, *runtime.AutomountServiceAccountToken)

	var installer corev1.ServiceAccount
	UnmarshalK8SYaml(t, RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-serviceaccount.yaml"}), &installer)
	require.Equal(t, "pega-installer-serviceaccount", installer.Name)
	require.False(t, *installer.AutomountServiceAccountToken)

	specs := renderPodSpecs(t, options, helmChartPath)
	require.Equal(t, "web-account", specs["pega-web"].ServiceAccountName)
	require.Nil(t, specs["pega-web"].AutomountServiceAccountToken, "explicit overrides keep their own automount setting")
	require.Equal(t, "pega-serviceaccount", specs["pega-batch"].ServiceAccountName)
	require.True(t, *specs["pega-batch"].AutomountServiceAccountToken)
	for _, name := range []string{"pega-search", "pega-hazelcast", "clusteringservice", "constellation"} {
		require.Equal(t, "pega-serviceaccount", specs[name].ServiceAccountName, name)
		require.True(t, *specs[name].AutomountServiceAccountToken, name)
	}
	require.Equal(t, "pega-installer-serviceaccount", specs["pega-db-install"].ServiceAccountName)
	require.False(t, *specs["pega-db-install"].AutomountServiceAccountToken)
}

func TestPegaServiceAccountsExisting(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	options := &helm.Options{SetValues: map[string]string{
		"global.provider":                  "k8s",
		"global.actions.execute":           "install-deploy",
		"global.serviceAccount.enabled":    "true",
		"global.serviceAccount.name":       "existing-runtime",
		"installer.serviceAccountName":     "existing-installer",
		"hazelcast.serviceAccountName":     "existing-hazelcast",
		"pegasearch.serviceAccountName":    "existing-search",
		"constellation.enabled":            "true",
		"constellation.serviceAccountName": "existing-constellation",
	}}

	for _, template := range []string{"templates/pega-serviceaccount.yaml", "charts/installer/templates/pega-installer-serviceaccount.yaml"} {
		_, err = RenderTemplateE(t, options, helmChartPath, []string{template})
		require.Error(t, err, template)
		require.Contains(t, err.Error(), "could not find template", template)
	}

	specs := renderPodSpecs(t, options, helmChartPath)
	require.Equal(t, "existing-runtime", specs["pega-web"].ServiceAccountName)
	require.Nil(t, specs["pega-web"].AutomountServiceAccountToken, "existing accounts keep their own automount setting")
	for name, account := range map[string]string{"pega-hazelcast": "existing-hazelcast", "pega-search": "existing-search", "constellation": "existing-constellation"} {
		require.Equal(t, account, specs[name].ServiceAccountName, name)
		require.Nil(t, specs[name].AutomountServiceAccountToken, name)
	}
	require.Equal(t, "existing-installer", specs["pega-db-install"].ServiceAccountName)
}

func TestPegaServiceAccountsValidation(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		values   map[string]string
		expected string
	}{
		{"create without enabled", map[string]string{"global.serviceAccount.create": "true"}, "global.serviceAccount.create requires global.serviceAccount.enabled=true"},
		{"enabled without name or create", map[string]string{"global.serviceAccount.enabled": "true"}, "global.serviceAccount.enabled requires global.serviceAccount.name or global.serviceAccount.create=true"},
		{"installer create without enabled", map[string]string{"installer.serviceAccount.create": "true"}, "installer.serviceAccount.create requires installer.serviceAccount.enabled=true"},
		{"installer create without enabled on deploy", map[string]string{"global.actions.execute": "deploy", "installer.serviceAccount.create": "true"}, "installer.serviceAccount.create requires installer.serviceAccount.enabled=true"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			values := map[string]string{"global.provider": "k8s", "global.actions.execute": "install-deploy"}
			for key, value := range testCase.values {
				values[key] = value
			}
			_, err := helm.RenderTemplateE(t, &helm.Options{SetValues: values}, helmChartPath, PegaHelmRelease, []string{})
			require.Error(t, err)
			require.Contains(t, err.Error(), testCase.expected)
		})
	}
}

// renderPodSpecs returns the pod specs of all Deployments, StatefulSets and Jobs keyed by workload name.
func renderPodSpecs(t *testing.T, options *helm.Options, chartPath string) map[string]corev1.PodSpec {
	output := RenderTemplate(t, options, chartPath, []string{})
	specs := make(map[string]corev1.PodSpec)
	for _, document := range strings.Split(output, "\n---") {
		var typeMeta metav1.TypeMeta
		if err := helm.UnmarshalK8SYamlE(t, document, &typeMeta); err != nil {
			continue
		}
		switch typeMeta.Kind {
		case "Deployment":
			var obj appsv1.Deployment
			UnmarshalK8SYaml(t, document, &obj)
			specs[obj.Name] = obj.Spec.Template.Spec
		case "StatefulSet":
			var obj appsv1.StatefulSet
			UnmarshalK8SYaml(t, document, &obj)
			specs[obj.Name] = obj.Spec.Template.Spec
		case "Job":
			var obj batchv1.Job
			UnmarshalK8SYaml(t, document, &obj)
			specs[obj.Name] = obj.Spec.Template.Spec
		}
	}
	return specs
}
