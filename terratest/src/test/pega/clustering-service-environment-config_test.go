package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestClusteringServiceEnvironmentConfig(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
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
					"hazelcast.clusteringServiceEnabled": "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/hazelcast/templates/clustering-service-environment-config.yaml"})
			VerifyClusteringServiceEnvironmentConfig(t,yamlContent, options)

		}
	}


}

func VerifyClusteringServiceEnvironmentConfig(t *testing.T, yamlContent string, options *helm.Options) {

	var clusteringServiceEnvConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &clusteringServiceEnvConfigMap)
			clusteringServiceEnvConfigData := clusteringServiceEnvConfigMap.Data
			require.Equal(t, clusteringServiceEnvConfigData["NAMESPACE"], "default")
			require.Equal(t, clusteringServiceEnvConfigData["JAVA_OPTS"], "-XX:MaxRAMPercentage=80.0 -XX:InitialRAMPercentage=80.0 -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/opt/hazelcast/logs/heapdump.hprof -XX:+UseG1GC -XX:NewRatio=3 -XshowSettings:vm -XX:InitiatingHeapOccupancyPercent=45 -Xlog:gc*,gc+phases=debug:file=/opt/hazelcast/logs/gc.log:time,pid,tags:filecount=5,filesize=3m")
			require.Equal(t, clusteringServiceEnvConfigData["SERVICE_NAME"], "clusteringservice-service")
			require.Equal(t, clusteringServiceEnvConfigData["MIN_CLUSTER_SIZE"], "3")
            require.Equal(t, clusteringServiceEnvConfigData["JMX_ENABLED"], "true")
            require.Equal(t, clusteringServiceEnvConfigData["HEALTH_MONITORING_LEVEL"], "OFF")
            require.Equal(t, clusteringServiceEnvConfigData["GROUP_NAME"], "prpchz")
            require.Equal(t, clusteringServiceEnvConfigData["GRACEFUL_SHUTDOWN_MAX_WAIT_SECONDS"], "600")
            require.Equal(t, clusteringServiceEnvConfigData["LOGGING_LEVEL"], "info")
            require.Equal(t, clusteringServiceEnvConfigData["DIAGNOSTICS_ENABLED"], "true")
            require.Equal(t, clusteringServiceEnvConfigData["DIAGNOSTICS_METRIC_LEVEL"], "info")
            require.Equal(t, clusteringServiceEnvConfigData["DIAGNOSTICS_FILE_COUNT"], "3")
            require.Equal(t, clusteringServiceEnvConfigData["DIAGNOSTIC_LOG_FILE_SIZE_MB"], "50")

		}
	}
}
