package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	apisv1beta2 "github.com/GoogleCloudPlatform/gke-managed-certs/pkg/apis/networking.gke.io/v1beta1"
	"path/filepath"
	"testing"
	"strings"
)


func TestPegaGKEManagedCertificate(t *testing.T) {
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testsPath,err := filepath.Abs(PegaHelmChartTestsPath)
	require.NoError(t, err)

	for _,operation := range supportedOperations{

			var options = &helm.Options{			
				SetValues: map[string]string{
					"global.provider":        "gke",
					"global.actions.execute": operation,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-gke-managed-certificate.yaml"}, "--values" , testsPath + "/data/values_gke_managedcertificate.yaml")
			yamlSplit := strings.Split(yamlContent, "---")
			verifyManagedCertificateForWeb(t, yamlSplit[1])			
			verifyManagedCertificateForStream(t, yamlSplit[2])
		}

}

func verifyManagedCertificateForWeb(t *testing.T, yamlContent string) {
	var managedCertificate apisv1beta2.ManagedCertificate
	UnmarshalK8SYaml(t, yamlContent, &managedCertificate)
	require.Equal(t, "managed-certificate-web", managedCertificate.Name)
}

func verifyManagedCertificateForStream(t *testing.T, yamlContent string) {
	var managedCertificate apisv1beta2.ManagedCertificate
	UnmarshalK8SYaml(t, yamlContent, &managedCertificate)
	require.Equal(t, "managed-certificate-stream", managedCertificate.Name)
}
