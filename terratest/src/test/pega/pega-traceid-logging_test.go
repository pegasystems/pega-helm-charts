package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaTraceIdLogging(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	t.Run("default values should not emit PEGA_LOG_CORRELATION_ID_ENABLED", func(t *testing.T) {
		var options = &helm.Options{
			SetValues: map[string]string{
				"global.provider":        "k8s",
				"global.actions.execute": "deploy",
			},
		}

		yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
		VerifyEnvNotPresent(t, yamlContent, "PEGA_LOG_CORRELATION_ID_ENABLED")
	})

	t.Run("correlationIdEnabled true should emit PEGA_LOG_CORRELATION_ID_ENABLED", func(t *testing.T) {
		var options = &helm.Options{
			SetValues: map[string]string{
				"global.provider":                "k8s",
				"global.actions.execute":         "deploy",
				"global.logging.correlationIdEnabled":  "true",
			},
		}

		yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
		VerifyEnvValue(t, yamlContent, "PEGA_LOG_CORRELATION_ID_ENABLED", "true")
	})

	t.Run("correlationIdEnabled false should not emit PEGA_LOG_CORRELATION_ID_ENABLED", func(t *testing.T) {
		var options = &helm.Options{
			SetValues: map[string]string{
				"global.provider":                "k8s",
				"global.actions.execute":         "deploy",
				"global.logging.correlationIdEnabled":  "false",
			},
		}

		yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
		VerifyEnvNotPresent(t, yamlContent, "PEGA_LOG_CORRELATION_ID_ENABLED")
	})

	t.Run("tier config contains prlog4j2.xml.tmpl key", func(t *testing.T) {
		var options = &helm.Options{
			SetValues: map[string]string{
				"global.provider":        "k8s",
				"global.actions.execute": "deploy",
			},
		}

		yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
		var pegaConfigMap k8score.ConfigMap
		configSlice := strings.Split(yamlContent, "---")
		for index, configData := range configSlice {
			if index >= 1 && index <= 3 {
				UnmarshalK8SYaml(t, configData, &pegaConfigMap)
				_, exists := pegaConfigMap.Data["prlog4j2.xml.tmpl"]
				require.True(t, exists, "prlog4j2.xml.tmpl key should exist in tier config")
				_, existsStatic := pegaConfigMap.Data["prlog4j2.xml"]
				require.True(t, existsStatic, "prlog4j2.xml key should exist in tier config for backward compatibility")
			}
		}
	})

	t.Run("template content includes Go conditional syntax", func(t *testing.T) {
		var options = &helm.Options{
			SetValues: map[string]string{
				"global.provider":        "k8s",
				"global.actions.execute": "deploy",
			},
		}

		yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
		var pegaConfigMap k8score.ConfigMap
		configSlice := strings.Split(yamlContent, "---")
		for index, configData := range configSlice {
			if index >= 1 && index <= 3 {
				UnmarshalK8SYaml(t, configData, &pegaConfigMap)
				tmplContent := pegaConfigMap.Data["prlog4j2.xml.tmpl"]
				require.Contains(t, tmplContent, `{{ if .Env.PEGA_LOG_CORRELATION_ID_ENABLED }}`)
				require.Contains(t, tmplContent, `[%X{ext-correlation-id}]`)
				require.Contains(t, tmplContent, `[%X{int-correlation-id}]`)
				require.Contains(t, tmplContent, `{{ end }}`)
			}
		}
	})
}

