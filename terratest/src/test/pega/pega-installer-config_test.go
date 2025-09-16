package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaInstallerConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}

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

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
			assertInstallerConfig(t, yamlContent)
		}
	}
}

func assertInstallerConfig(t *testing.T, configYaml string) {
	var installConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, configYaml, &installConfigMap)
	installConfigData := installConfigMap.Data

	compareConfigMapData(t, installConfigData["setupDatabase.properties.tmpl"], "data/expectedSetupdatabase.properties")
	compareConfigMapData(t, installConfigData["migrateSystem.properties.tmpl"], "data/expectedMigrateSystem.properties.tmpl")
	compareConfigMapData(t, installConfigData["prlog4j2.xml"], "data/expectedPRlog4j2.xml")
	compareConfigMapData(t, installConfigData["prpcUtils.properties.tmpl"], "data/expectedPRPCUtils.properties.tmpl")
}
