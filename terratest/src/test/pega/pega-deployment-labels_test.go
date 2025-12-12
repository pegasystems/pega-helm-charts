package pega

import (
	"testing"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
)


func TestPegaTierDeploymentWithDeploymentLabels(t *testing.T) {
 var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
 var supportedOperations = []string{"deploy"}
 var deploymentNames = []string{"pega", "myapp-dev"}
 helmChartPath, err := filepath.Abs(PegaHelmChartPath)
 require.NoError(t, err)
 var depObj appsv1.Deployment
 for _, vendor := range supportedVendors {
  for _, operation := range supportedOperations {
   for _, deploymentName := range deploymentNames {
    options := &helm.Options{
     SetValues: map[string]string{
      "global.provider":               vendor,
      "global.actions.execute":        operation,
      "global.deployment.name":        deploymentName,
      "global.tier[0].name":           "web",
      "global.tier[0].deploymentLabels.foo": "bar",
      "global.tier[0].deploymentLabels.env": "prod",
     },
    }
    yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
    yamlSplit := strings.Split(yamlContent, "---")
    UnmarshalK8SYaml(t, yamlSplit[1], &depObj)
    labels := depObj.ObjectMeta.Labels
    require.Equal(t, "bar", labels["foo"])
    require.Equal(t, "prod", labels["env"])
   }
  }
 }
}