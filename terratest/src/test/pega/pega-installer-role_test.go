package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8srbac "k8s.io/api/rbac/v1"
)

func TestPegaInstallerRole(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "upgrade-deploy"}

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
			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-role.yaml"})
			assertInstallerRole(t, yamlContent)
		}
	}
}

func assertInstallerRole(t *testing.T, roleYaml string) {
	var deployRoleObj k8srbac.Role
	UnmarshalK8SYaml(t, roleYaml, &deployRoleObj)
	require.Equal(t, deployRoleObj.Rules[0].APIGroups, []string{"apps"})
	require.Equal(t, deployRoleObj.Rules[0].Resources, []string{"deployments", "statefulsets"})
	require.Equal(t, deployRoleObj.Rules[0].Verbs, []string{"get", "watch", "list"})
	require.Equal(t, deployRoleObj.Rules[1].APIGroups, []string{"batch"})
	require.Equal(t, deployRoleObj.Rules[1].Resources, []string{"jobs"})
	require.Equal(t, deployRoleObj.Rules[1].Verbs, []string{"get", "watch", "list"})
	require.Equal(t, deployRoleObj.Rules[2].APIGroups, []string{"extensions"})
	require.Equal(t, deployRoleObj.Rules[2].Resources, []string{"deployments"})
	require.Equal(t, deployRoleObj.Rules[2].Verbs, []string{"get", "watch", "list"})
}
