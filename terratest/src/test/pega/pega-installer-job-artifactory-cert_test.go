package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sbatch "k8s.io/api/batch/v1"
)

func TestPegaInstallerJobWithArtifactoryCert(t *testing.T) {

	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.customArtifactory.enableSSLVerification": "true",
						"global.customArtifactory.certificate": "self-signed-certificate.cer: |\n-----BEGIN CERTIFICATE-----\nMIIDdTCCAl2gAwIBAgIENdb1mTANBgkqhkiG9w0BAQsFADBrMQswCQYDVQQGEwJJ\nTjESMBAGA1UECBMJVGVsYW5nYW5hMRIwEAYDVQQHEwlIeWRlcmFiYWQxDTALBgNV\nBAoTBFBlZ2ExDTALBgNVBAsTBFBERFMxFjAUBgNVBAMTDTEwLjIyNS43MS4xNDMw\nHhcNMjIwMzIyMTYwNzQxWhcNMjMwMzE3MTYwNzQxWjBrMQswCQYDVQQGEwJJTjES\nMBAGA1UECBMJVGVsYW5nYW5hMRIwEAYDVQQHEwlIeWRlcmFiYWQxDTALBgNVBAoT\nBFBlZ2ExDTALBgNVBAsTBFBERFMxFjAUBgNVBAMTDTEwLjIyNS43MS4xNDMwggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCNYvEPbKRZ1y3u5fmPcvdBaQVK\nXsy3ioY7ToMP6vpNXGmRhI06t6jVzDVU85jtdSb4On2B4uyZwUSxO9cWfOFtI6wW\nnrdhRxmygFvXwinon6LXcoRfK9TjI72C/694UWu/UysaUp8yWyrHf2XfQK2qFqMC\nej57bpRSCME2maAXKC88IGNpeX2XjhICUKPPBrWDpK4Jq7NhcgV6z0hmtCWg8+lj\nmMZRXDoJKMNNhWOKMYhL/djDJZH+PMON7sOGVtE52U0UBLhff6Ee4ERBRummNgCv\nxt4MTmzgewsN1uKQ5MiBjtJduVEmKhhiIV38QetrCPpejAHOJLFe2l5VfKHDAgMB\nAAGjITAfMB0GA1UdDgQWBBTK0eVfaa41Vr4qXTww3RBTgFO78DANBgkqhkiG9w0B\nAQsFAAOCAQEAOAjezNJmMx9j0hnutOspnHC8iOqaFQjW8t6D9cWEQALd2PNPB5S9\nQxlEuaN3x/zbtNI55fxZW6ryP/AJ0DclTs8vwzEk7DJ1Yt7vMfFG6DxbIUlPY677\nDGB23K68BXl8MtSYvOLbDwXYjyMDUzcmojaIjS6RwW8C5yvXW34h2jjwVWQm1yti\n46xANKLHEVTp44LiG+gf/9TxfQjSQXpSdgdMbJB744tMmozyfbtulWE0T5dBvd8w\ncdbPKbgldsv4bc8EojOYRRasYu6nZqP+8Tw/4jHr4IB2kiuJ63gs6IlqnzDyzzQ7\nYc0a+hYe1cTSXQn23aL/c9v/901LUpdAYw==\n-----END CERTIFICATE-----\n",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")

				// If there are three slices, it means that it is a pega-upgrade-deploy job
				if len(yamlSplit) == 4 {
					for index, jobInfo := range yamlSplit {
						if index >= 1 && index <= 3 {
							assertJobArtifactoryCertVolumeAndMount(t, jobInfo, true)
						}
					}
				} else {
					if operation == "install" || operation == "install-deploy" {
						assertJobArtifactoryCertVolumeAndMount(t, yamlSplit[1], true)
					} else {
						assertJobArtifactoryCertVolumeAndMount(t, yamlSplit[1], true)
					}
				}
			}
	}
}

func TestPegaInstallerJobWithoutArtifactoryCert(t *testing.T) {

	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
						"global.customArtifactory.enableSSLVerification": "false",
						"global.customArtifactory.certificate": "self-signed-certificate.cer: |\n-----BEGIN CERTIFICATE-----\nMIIDdTCCAl2gAwIBAgIENdb1mTANBgkqhkiG9w0BAQsFADBrMQswCQYDVQQGEwJJ\nTjESMBAGA1UECBMJVGVsYW5nYW5hMRIwEAYDVQQHEwlIeWRlcmFiYWQxDTALBgNV\nBAoTBFBlZ2ExDTALBgNVBAsTBFBERFMxFjAUBgNVBAMTDTEwLjIyNS43MS4xNDMw\nHhcNMjIwMzIyMTYwNzQxWhcNMjMwMzE3MTYwNzQxWjBrMQswCQYDVQQGEwJJTjES\nMBAGA1UECBMJVGVsYW5nYW5hMRIwEAYDVQQHEwlIeWRlcmFiYWQxDTALBgNVBAoT\nBFBlZ2ExDTALBgNVBAsTBFBERFMxFjAUBgNVBAMTDTEwLjIyNS43MS4xNDMwggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCNYvEPbKRZ1y3u5fmPcvdBaQVK\nXsy3ioY7ToMP6vpNXGmRhI06t6jVzDVU85jtdSb4On2B4uyZwUSxO9cWfOFtI6wW\nnrdhRxmygFvXwinon6LXcoRfK9TjI72C/694UWu/UysaUp8yWyrHf2XfQK2qFqMC\nej57bpRSCME2maAXKC88IGNpeX2XjhICUKPPBrWDpK4Jq7NhcgV6z0hmtCWg8+lj\nmMZRXDoJKMNNhWOKMYhL/djDJZH+PMON7sOGVtE52U0UBLhff6Ee4ERBRummNgCv\nxt4MTmzgewsN1uKQ5MiBjtJduVEmKhhiIV38QetrCPpejAHOJLFe2l5VfKHDAgMB\nAAGjITAfMB0GA1UdDgQWBBTK0eVfaa41Vr4qXTww3RBTgFO78DANBgkqhkiG9w0B\nAQsFAAOCAQEAOAjezNJmMx9j0hnutOspnHC8iOqaFQjW8t6D9cWEQALd2PNPB5S9\nQxlEuaN3x/zbtNI55fxZW6ryP/AJ0DclTs8vwzEk7DJ1Yt7vMfFG6DxbIUlPY677\nDGB23K68BXl8MtSYvOLbDwXYjyMDUzcmojaIjS6RwW8C5yvXW34h2jjwVWQm1yti\n46xANKLHEVTp44LiG+gf/9TxfQjSQXpSdgdMbJB744tMmozyfbtulWE0T5dBvd8w\ncdbPKbgldsv4bc8EojOYRRasYu6nZqP+8Tw/4jHr4IB2kiuJ63gs6IlqnzDyzzQ7\nYc0a+hYe1cTSXQn23aL/c9v/901LUpdAYw==\n-----END CERTIFICATE-----\n",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-job.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")

				// If there are three slices, it means that it is a pega-upgrade-deploy job
				if len(yamlSplit) == 4 {
					for index, jobInfo := range yamlSplit {
						if index >= 1 && index <= 3 {
							assertJobArtifactoryCertVolumeAndMount(t, jobInfo, false)
						}
					}
				} else {
					if operation == "install" || operation == "install-deploy" {
						assertJobArtifactoryCertVolumeAndMount(t, yamlSplit[1], false)
					} else {
						assertJobArtifactoryCertVolumeAndMount(t, yamlSplit[1], false)
					}
				}
		}
	}
}

func assertJobArtifactoryCertVolumeAndMount(t *testing.T, jobYaml string, shouldHaveVol bool) {
	var jobObj k8sbatch.Job
	UnmarshalK8SYaml(t, jobYaml, &jobObj)

	jobSpec := jobObj.Spec.Template.Spec

	var foundVol = false
	for _, vol := range jobSpec.Volumes {
		if vol.Name == "pega-volume-custom-artifactory-certificate" {
			foundVol = true
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVol)

	var foundVolMount = false
	for _, container := range jobSpec.Containers {
			for _, volMount := range container.VolumeMounts {
				if volMount.Name == "pega-volume-custom-artifactory-certificate" {
					require.Equal(t, "/opt/pega/artifactory/cert", volMount.MountPath)
					foundVolMount = true
					break
				}
			}
			break
	}
	require.Equal(t, shouldHaveVol, foundVolMount)

}
