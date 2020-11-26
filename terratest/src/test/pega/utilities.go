package pega

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

var volumeDefaultMode int32 = 420
var volumeDefaultModePtr = &volumeDefaultMode


// compareConfigMapData - Compares the config map deployed for each kind of tier with the excepted xml's
func compareConfigMapData(t *testing.T, actualFileData string, expectedFileName string) {
	expectedFile, err := ioutil.ReadFile(expectedFileName)
	require.Empty(t, err)
	expectedFileData := string(expectedFile)
	expectedFileData = strings.Replace(expectedFileData, "\r", "", -1)

	equal := false
	if expectedFileData == actualFileData {
		equal = true
	}
	require.Equal(t, true, equal)
}

//aksSpecificUpgraderDeployEnvs - Test aks specific upgrade-deploy environmnet variables in case of upgrade-deploy
func aksSpecificUpgraderDeployEnvs(t *testing.T, options *helm.Options, container k8score.Container) {
	if options.SetValues["global.provider"] == "aks" && options.SetValues["global.actions.execute"] == "upgrade-deploy" {
		require.Equal(t, container.Env[0].Name, "KUBERNETES_SERVICE_HOST")
		require.Equal(t, container.Env[0].Value, "API_SERVICE_ADDRESS")
		require.Equal(t, container.Env[1].Name, "KUBERNETES_SERVICE_PORT_HTTPS")
		require.Equal(t, container.Env[1].Value, "SERVICE_PORT_HTTPS")
		require.Equal(t, container.Env[2].Name, "KUBERNETES_SERVICE_PORT")
		require.Equal(t, container.Env[2].Value, "SERVICE_PORT_HTTPS")
	}
}

// VerifyInitContinerData - Verifies any possible initContainer that can occur in pega helm chart deployments
func VerifyInitContinerData(t *testing.T, containers []k8score.Container, options *helm.Options) {

	if len(containers) == 0 {
		println("no init containers")
	}

	for i := 0; i < len(containers); i++ {
		container := containers[i]
		name := container.Name
		if name == "wait-for-pegainstall" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-db-install"}, container.Args)
		} else if name == "wait-for-pegasearch" {
			require.Equal(t, "busybox:1.31.0", container.Image)
			require.Equal(t, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-cassandra" {
			require.Equal(t, "cassandra:3.11.3", container.Image)
			require.Equal(t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" pega-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-cassandra" {
			require.Equal(t, "cassandra:3.11.3", container.Image)
			require.Equal(t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" pega-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-pegaupgrade" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-db-upgrade"}, container.Args)
			aksSpecificUpgraderDeployEnvs(t, options, container)
		} else if name == "wait-for-pre-dbupgrade" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-pre-upgrade"}, container.Args)
		} else if name == "wait-for-rolling-updates" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"sh", "-c", " kubectl rollout status deployment/pega-web --namespace default && kubectl rollout status deployment/pega-batch --namespace default && kubectl rollout status statefulset/pega-stream --namespace default"}, container.Command)
		} else if name == "wait-for-hazelcast" {
			require.Equal(t, "busybox:1.31.0", container.Image)
			require.Equal(t, []string{"sh", "-c", "set -e counter=0; while [ -z $(wget -S \"http://pega-hazelcast-service.default:5701/hazelcast/health/cluster-state\" 2>&1 | grep \"HTTP/\" | awk '{print $2}') ] || [ $(wget -q -O - \"http://pega-hazelcast-service.default:5701/hazelcast/health/cluster-size\" /dev/null) -ne 3 ] || [ $(wget -S \"http://pega-hazelcast-service.default:5701/hazelcast/health/cluster-state\" 2>&1 | grep \"HTTP/\" | awk '{print $2}') -ne 200 ]; do echo \"waiting for hazelcast pods to start and join the cluster...\" ; counter=$(($counter+5)); sleep 5; if [ $counter -gt 150 ]; then echo \"Timeout Reached. Hazelcast pods failed to join the cluster\"; exit 1; fi done; echo \"Hazelcast cluster is up now\"; exit 0;\n"}, container.Args)
		} else {
			fmt.Println("invalid init containers found.. please check the list", name)
			t.Fail()
		}
	}
}
