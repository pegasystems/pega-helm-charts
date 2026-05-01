package pega

import (
	"path/filepath"
	"testing"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaInstallerLibDownloadScriptConfig(t *testing.T) {
    helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
        SetValues: map[string]string{
            "global.deployment.name":        "pega",
            "global.provider":               "k8s",
            "global.actions.execute":        "install",
            "global.downloadContainer.image": "ICDOWNLOAD_IMAGE:1.0",
        },
    }
    yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-lib-download-script-config.yaml"})

    var envConfigMap k8score.ConfigMap
    UnmarshalK8SYaml(t, yamlContent, &envConfigMap)

    require.Equal(t, "pega-installer-lib-download-script-config", envConfigMap.ObjectMeta.Name)
    VerifyEnvPresent(t, &envConfigMap, "download-jdbc-lib.sh")
}


func TestPegaInstallerLibDownloadScriptConfigRenderedForActions(t *testing.T) {
    var supportedOperations = []string{"install", "upgrade", "install-deploy", "upgrade-deploy"}

    helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

    for _, operation := range supportedOperations {

        var options = &helm.Options{
            SetValues: map[string]string{
                "global.deployment.name":        "pega",
                "global.provider":               "k8s",
                "global.actions.execute":        operation,
                "installer.upgrade.upgradeType": "out-of-place",
                "global.downloadContainer.image": "ICDOWNLOAD_IMAGE:1.0",
            },
        }
        yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-lib-download-script-config.yaml"})

        var envConfigMap k8score.ConfigMap
        UnmarshalK8SYaml(t, yamlContent, &envConfigMap)

        require.Equal(t, "pega-installer-lib-download-script-config", envConfigMap.ObjectMeta.Name)
        VerifyEnvPresent(t, &envConfigMap, "download-jdbc-lib.sh")
    }
}