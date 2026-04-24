package pega

import (
	"path/filepath"
	"testing"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaLibDownloadScriptConfig(t *testing.T) {
    helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
        SetValues: map[string]string{
            "global.deployment.name":        "pega",
            "global.provider":               "k8s",
            "global.actions.execute":        "deploy",
            "global.downloadContainer.image": "ICDOWNLOAD_IMAGE:1.0",
        },
    }
    yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-lib-download-script-config.yaml"})

    var envConfigMap k8score.ConfigMap
    UnmarshalK8SYaml(t, yamlContent, &envConfigMap)

    require.Equal(t, "pega-runtime-lib-download-script-config", envConfigMap.ObjectMeta.Name)
    VerifyEnvPresent(t, &envConfigMap, "download-jdbc-lib.sh")
}

func VerifyEnvPresent(t *testing.T, envConfigMap *k8score.ConfigMap, entry string) {
	envConfigData := envConfigMap.Data
	_, previouslySet := envConfigData[entry]
	require.Equal(t, true, previouslySet)
}