package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"testing"
)

func TestHazelcast(t *testing.T) {

	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy"}
    var deploymentNames = []string{"pega","myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
		    for _, depName := range deploymentNames {
                var hazelcastOptions = &helm.Options{
                    SetValues: map[string]string{
                        "global.deployment.name": depName,
                        "global.provider":        vendor,
                        "global.actions.execute": operation,
						"installer.upgrade.upgradeType": "zero-downtime",
                        "hazelcast.enabled":  "true",
                    },
                }
                deploymentYaml := RenderTemplate(t, hazelcastOptions, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
                yamlSplit := strings.Split(deploymentYaml, "---")
                assertWeb(t,yamlSplit[1],hazelcastOptions)
                assertBatch(t,yamlSplit[2],hazelcastOptions)
                //assertStream(t,yamlSplit[3],hazelcastOptions)
            }
		}
	}
}
