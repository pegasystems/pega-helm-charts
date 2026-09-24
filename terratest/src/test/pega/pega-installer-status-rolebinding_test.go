package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8srbac "k8s.io/api/rbac/v1"
)

func TestPegaInstallerStatusRoleBinding(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":               vendor,
					"global.actions.execute":        operation,
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}
			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-status-rolebinding.yaml"})
			assertInstallerRoleBinding(t, yamlContent)
		}
	}
}

func TestPegaInstallerStatusRoleBindingUsesRuntimeAndInstallerAccounts(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		values        map[string]string
		expectedNames []string
	}{
		{
			name: "shared runtime account",
			values: map[string]string{
				"global.provider":                  "k8s",
				"global.actions.execute":           "install-deploy",
				"global.serviceAccount.enabled":    "true",
				"global.serviceAccount.create":     "true",
				"global.serviceAccount.name":       "runtime-account",
				"installer.serviceAccount.enabled": "false",
				"installer.serviceAccount.create":  "false",
			},
			expectedNames: []string{"runtime-account"},
		},
		{
			name: "separate installer account",
			values: map[string]string{
				"global.provider":                  "k8s",
				"global.actions.execute":           "install-deploy",
				"global.serviceAccount.enabled":    "true",
				"global.serviceAccount.create":     "true",
				"global.serviceAccount.name":       "runtime-account",
				"installer.serviceAccount.enabled": "true",
				"installer.serviceAccount.create":  "true",
				"installer.serviceAccount.name":    "installer-account",
			},
			expectedNames: []string{"runtime-account", "installer-account"},
		},
		{
			name: "tier override without shared account",
			values: map[string]string{
				"global.provider":                          "k8s",
				"global.actions.execute":                   "install-deploy",
				"global.tier[0].name":                      "web",
				"global.tier[0].custom.serviceAccountName": "web-account",
				"global.tier[1].name":                      "batch",
			},
			expectedNames: []string{"web-account", "default"},
		},
		{
			name: "existing installer account",
			values: map[string]string{
				"global.provider":               "k8s",
				"global.actions.execute":        "upgrade-deploy",
				"installer.upgrade.upgradeType": "zero-downtime",
				"installer.serviceAccountName":  "existing-installer",
			},
			expectedNames: []string{"default", "existing-installer"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			options := &helm.Options{SetValues: testCase.values}
			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-status-rolebinding.yaml"})

			var roleBinding k8srbac.RoleBinding
			helm.UnmarshalK8SYaml(t, yamlContent, &roleBinding)
			require.Equal(t, testCase.expectedNames, subjectNames(roleBinding))
		})
	}
}

func subjectNames(roleBinding k8srbac.RoleBinding) []string {
	names := make([]string, 0, len(roleBinding.Subjects))
	for _, subject := range roleBinding.Subjects {
		names = append(names, subject.Name)
	}
	return names
}

func assertInstallerRoleBinding(t *testing.T, roleBindingYaml string) {
	var installerRoleBindingObj k8srbac.RoleBinding
	helm.UnmarshalK8SYaml(t, roleBindingYaml, &installerRoleBindingObj)
	require.Equal(t, installerRoleBindingObj.RoleRef.APIGroup, "rbac.authorization.k8s.io")
	require.Equal(t, installerRoleBindingObj.RoleRef.Kind, "Role")
	require.Equal(t, installerRoleBindingObj.RoleRef.Name, "jobs-reader")

	require.Equal(t, installerRoleBindingObj.Subjects[0].Kind, "ServiceAccount")
	require.Equal(t, installerRoleBindingObj.Subjects[0].Name, "default")
	require.Equal(t, installerRoleBindingObj.Subjects[0].Namespace, "default")
}
