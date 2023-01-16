package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	k8sbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8srbac "k8s.io/api/rbac/v1"
	"testing"
	"strings"
	"fmt"
)

func TestClusteringServiceMigration(t *testing.T) {

	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
    var supportedOperations =  []string{"deploy","install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

		    fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"hazelcast.migration.initiateMigration":  "true",
					"hazelcast.clusteringServiceEnabled":  "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/clustering-service-migration.yaml"})
			yamlSplit := strings.Split(yamlContent, "---")
			assertServiceAccount(t, yamlSplit[1], options)
            assertRole(t, yamlSplit[2], options)
            assertRoleBinding(t, yamlSplit[3], options)
			assertMigrationJob(t, yamlSplit[4], options)

		}
	}
}
func assertServiceAccount(t *testing.T, serviceAccountYaml string, options *helm.Options) {
	var serviceAccount *corev1.ServiceAccount
	UnmarshalK8SYaml(t, serviceAccountYaml, &serviceAccount)

	require.Equal(t, serviceAccount.ObjectMeta.Name, "clusteringservice-migration-sa")
	require.Equal(t, serviceAccount.ObjectMeta.Namespace, "default")
}

func assertRole(t *testing.T, roleYaml string, options *helm.Options) {
	var roleObj k8srbac.Role
	UnmarshalK8SYaml(t, roleYaml, &roleObj)

    require.Equal(t, roleObj.ObjectMeta.Name, "clusteringservice-migration-role")
    require.Equal(t, roleObj.ObjectMeta.Namespace, "default")
	require.Equal(t, roleObj.Rules[0].APIGroups, []string{""})
	require.Equal(t, roleObj.Rules[0].Resources, []string{"pods"})
	require.Equal(t, roleObj.Rules[0].Verbs, []string{"get", "list"})
    require.Equal(t, roleObj.Rules[1].APIGroups, []string{""})
    require.Equal(t, roleObj.Rules[1].Resources, []string{"pods/exec"})
    require.Equal(t, roleObj.Rules[1].Verbs, []string{"create"})
}

func assertRoleBinding(t *testing.T, roleBinding string, options *helm.Options) {
	var roleBindingObj k8srbac.RoleBinding
	UnmarshalK8SYaml(t, roleBinding, &roleBindingObj)

    require.Equal(t, roleBindingObj.ObjectMeta.Name, "clusteringservice-migration-role-binding")
    require.Equal(t, roleBindingObj.ObjectMeta.Namespace, "default")
	require.Equal(t, roleBindingObj.Subjects[0].Kind, "ServiceAccount")
	require.Equal(t, roleBindingObj.Subjects[0].Name, "clusteringservice-migration-sa")
    require.Equal(t, roleBindingObj.RoleRef.APIGroup, "rbac.authorization.k8s.io")
    require.Equal(t, roleBindingObj.RoleRef.Kind, "Role")
    require.Equal(t, roleBindingObj.RoleRef.Name, "clusteringservice-migration-role")
}

func assertMigrationJob(t *testing.T, jobYaml string, options *helm.Options) {
	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, jobYaml, &jobObj)

	jobSpec := jobObj.Spec.Template.Spec

    require.Equal(t, jobObj.ObjectMeta.Name, "clusteringservice-migration-job")
    require.Equal(t, jobObj.ObjectMeta.Namespace, "default")
    require.Equal(t, jobObj.Spec.Template.ObjectMeta.Name, "clusteringservice-migration-job")
	require.Equal(t, jobSpec.ServiceAccountName, "clusteringservice-migration-sa")
}


