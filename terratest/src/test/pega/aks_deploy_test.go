package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
)

// TestPegaStandardTierDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestPegaAKSStandardTierDeployment(t *testing.T) {
	// set action execute to install
	var options = &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
			"global.provider":        "aks",
		},
	}

	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, options, []string{"wait-for-pegasearch", "wait-for-cassandra"})
}

// set action execute to install
var installDeployoptions = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "install-deploy",
		"global.provider":        "aks",
	},
}

// TestPegaAKSInstallDeployDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestPegaAKSInstallDeployDeployment(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, installDeployoptions, []string{"wait-for-pegainstall", "wait-for-pegasearch", "wait-for-cassandra"})
}

// set action execute to install
var aksUpgradeDeployOptions = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "upgrade-deploy",
		"global.provider":        "aks",
	},
}

// TestPegaAKSUpgradeDeployDeployment - Test case to verify the upgrade-deploy on AKS provider.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
// Special case in AKS during rolling restart to verify environments variables that are specific to AKS cluster - aksSpecificUpgraderDeployEnvs() method
func TestPegaAKSUpgradeDeployDeployment(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, aksUpgradeDeployOptions, []string{"wait-for-pegaupgrade"})
}
