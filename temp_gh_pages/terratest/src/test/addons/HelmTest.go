package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

type HelmTest struct {
	T           *testing.T
	ChartPath   string
	HelmOptions *helm.Options
}

func NewHelmTest(t *testing.T, chartRelativePath string, options map[string]string) *HelmTest {
	t.Parallel()

	path, err := filepath.Abs(chartRelativePath)
	require.NoError(t, err)

	return &HelmTest{
		T:           t,
		ChartPath:   path,
		HelmOptions: &helm.Options{SetValues: options},
	}
}
