package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"io"
	k8score "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierConfigWithWeb(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	webPath := filepath.Join(helmChartPath, "config", "deploy", "web.xml")

	err = CopyFile("data/expectedInstallDeployWeb.xml", webPath)
	require.NoError(t, err)
	defer os.Remove(webPath)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":               vendor,
					"global.actions.execute":        operation,
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
			VerifyTierConfgWithWeb(t, yamlContent, options)

		}
	}

}

func CopyFile(src, dest string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

// VerifyTierConfg - Performs the tier specific configuration assetions with the values as provided in default values.yaml
func VerifyTierConfgWithWeb(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(yamlContent, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			pegaConfigMapData := pegaConfigMap.Data
			compareConfigMapData(t, pegaConfigMapData["prconfig.xml"], "data/expectedInstallDeployPrconfig.xml")
			compareConfigMapData(t, pegaConfigMapData["context.xml.tmpl"], "data/expectedInstallDeployContext.xml.tmpl")
			compareConfigMapData(t, pegaConfigMapData["prlog4j2.xml"], "data/expectedInstallDeployPRlog4j2.xml")
			compareConfigMapData(t, pegaConfigMapData["server.xml.tmpl"], "data/expectedInstallDeployServer.xml.tmpl")
			compareConfigMapData(t, pegaConfigMapData["web.xml"], "data/expectedInstallDeployWeb.xml")
		}
	}
}
