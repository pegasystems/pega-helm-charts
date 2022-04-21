package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaInstallerConnectionPropsConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade", "upgrade-deploy"}
	var supportedDbs = []string{"postgres", "mssql", "oracledate", "udb", "db2zos"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, dbPlatform := range supportedDbs {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                  vendor,
						"global.actions.execute":           operation,
						"global.jdbc.dbType":               dbPlatform,
						"global.jdbc.connectionProperties": "prop1=value1;prop2=value2",
						"installer.upgrade.upgradeType":    "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
				assertInstallerConnectionPropsConfig(t, yamlContent, dbPlatform)
			}
		}
	}
}

func assertInstallerConnectionPropsConfig(t *testing.T, configYaml string, dbPlatform string) {
	var installConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, configYaml, &installConfigMap)
	installConfigData := installConfigMap.Data

	compareConfigMapData(t, installConfigData["prconfig.xml.tmpl"], "data/expectedPrconfig.xml")
	compareConfigMapData(t, installConfigData["setupDatabase.properties.tmpl"], "data/expectedSetupdatabase.properties")
	compareConfigMapData(t, installConfigData["prbootstrap.properties.tmpl"], "data/expectedPRbootstrap.properties")
	compareConfigMapData(t, installConfigData["migrateSystem.properties.tmpl"], "data/expectedMigrateSystem.properties.tmpl")
	compareConfigMapData(t, installConfigData["prlog4j2.xml"], "data/expectedPRlog4j2.xml")
	compareConfigMapData(t, installConfigData["prpcUtils.properties.tmpl"], "data/expectedPRPCUtils.properties.tmpl")
	compareConfigMapData(t, installConfigData[fmt.Sprintf("%s.conf", dbPlatform)], fmt.Sprintf("data/expected%s.conf", dbPlatform))
	if dbPlatform == "db2zos" {
		compareConfigMapData(t, installConfigData["DB2SiteDependent.properties"], "data/expectedDB2SiteDependent.properties")
	} else {
		assert.Equal(t, installConfigData["DB2SiteDependent.properties"], "")
	}
}
