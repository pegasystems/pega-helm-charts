package backingservices

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestConstellationIngressDisabled(t *testing.T) {

	var supportedVendors = []string{"k8s", "gke", "aws"}

	for _, vendor := range supportedVendors {

		fmt.Println(vendor)
		/* Run subtest in parallel */
		t.Run(vendor, func(subtest *testing.T) {
			var mapconstellation = map[string]string{
				"constellation.enabled":        		"true",
				"constellation.deployment.name":        "constellation-test",
				"constellation.provider":               vendor,
				"constellation.ingress.enabled":        "false",
			};
	
			var helmtest = NewHelmTest(subtest, helmChartRelativePath, mapconstellation);
			helmChartParser := NewHelmConfigParser(
				helmtest,
			)
	
			for _, i := range constellationGKEResourcesForIngress {
				require.False(subtest, helmChartParser.Contains(SearchResourceOption{
					Name: i.Name,
					Kind: i.Kind,
				}))
			}
		})
	}

}

func TestConstellationMessagingIngressDisabled(t *testing.T) {

	var supportedVendors = []string{"k8s", "gke", "aws"}

	for _, vendor := range supportedVendors {

		fmt.Println(vendor)
		/* Run subtest in parallel */
		t.Run(vendor, func(subtest *testing.T) {
			var mapconstellation = map[string]string{
				"constellation.enabled":        		"true",
				"constellation.deployment.name":        "constellation-messaging-test",
				"constellation.provider":               vendor,
				"constellation.ingress.enabled":        "false",
			};
	
			var helmtest = NewHelmTest(subtest, helmChartRelativePath, mapconstellation);
			helmChartParser := NewHelmConfigParser(
				helmtest,
			)
	
			for _, i := range constellationMessagingGKEResourcesForIngress {
				require.False(subtest, helmChartParser.Contains(SearchResourceOption{
					Name: i.Name,
					Kind: i.Kind,
				}))
			}
		})
	}

}

func TestConstellationGKEIngressEnabled(t *testing.T) {

	var mapconstellation = map[string]string{
		"constellation.enabled":        		"true",
		"constellation.deployment.name":        "constellation-test",
		"constellation.provider":               "gke",
		"constellation.ingress.enabled":        "true",
	};

	var helmtest = NewHelmTest(t, helmChartRelativePath, mapconstellation);
	helmChartParser := NewHelmConfigParser(
		helmtest,
	)

	for _, i := range constellationGKEResourcesForIngress {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}

}

func TestConstellationMessagingGKEIngressEnabled(t *testing.T) {

	var mapconstellation = map[string]string{
		"constellation.enabled":        		"true",
		"constellation.deployment.name":        "constellation-messaging-test",
		"constellation.provider":               "gke",
		"constellation.ingress.enabled":        "true",
	};

	var helmtest = NewHelmTest(t, helmChartRelativePath, mapconstellation);
	helmChartParser := NewHelmConfigParser(
		helmtest,
	)

	for _, i := range constellationMessagingGKEResourcesForIngress {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}

}

func TestConstellationIngressEnabled(t *testing.T) {

	var supportedVendors = []string{"k8s", "aws"}

	for _, vendor := range supportedVendors {

		fmt.Println(vendor)
		/* Run subtest in parallel */
		t.Run(vendor, func(subtest *testing.T) {
			var mapconstellation = map[string]string{
				"constellation.enabled":        		"true",
				"constellation.deployment.name":        "constellation-test",
				"constellation.provider":               vendor,
				"constellation.ingress.enabled":        "true",
			};
	
			var helmtest = NewHelmTest(subtest, helmChartRelativePath, mapconstellation);
			helmChartParser := NewHelmConfigParser(
				helmtest,
			)
	
			for _, i := range constellationResourcesForIngress {
				require.True(subtest, helmChartParser.Contains(SearchResourceOption{
					Name: i.Name,
					Kind: i.Kind,
				}))
			}
		})
	}

}

func TestConstellationMessagingIngressEnabled(t *testing.T) {

	var supportedVendors = []string{"k8s", "aws"}

	for _, vendor := range supportedVendors {

		fmt.Println(vendor)
		/* Run subtest in parallel */
		t.Run(vendor, func(subtest *testing.T) {
			var mapconstellation = map[string]string{
				"constellation.enabled":        		"true",
				"constellation.deployment.name":        "constellation-messaging-test",
				"constellation.provider":               vendor,
				"constellation.ingress.enabled":        "true",
			};
	
			var helmtest = NewHelmTest(subtest, helmChartRelativePath, mapconstellation);
			helmChartParser := NewHelmConfigParser(
				helmtest,
			)
	
			for _, i := range constellationMessagingResourcesForIngress {
				require.True(subtest, helmChartParser.Contains(SearchResourceOption{
					Name: i.Name,
					Kind: i.Kind,
				}))
			}
		})
	}

}

var constellationGKEResourcesForIngress = []SearchResourceOption{
	{
		Name: "constellation-test",
		Kind: "Ingress",
	},
	{
		Name: "constellation-test",
		Kind: "BackendConfig",
	},
}

var constellationMessagingGKEResourcesForIngress = []SearchResourceOption{
	{
		Name: "constellation-messaging-test",
		Kind: "Ingress",
	},
	{
		Name: "constellation-messaging-test",
		Kind: "BackendConfig",
	},
}

var constellationResourcesForIngress = []SearchResourceOption{
	{
		Name: "constellation-test",
		Kind: "Ingress",
	},
}

var constellationMessagingResourcesForIngress = []SearchResourceOption{
	{
		Name: "constellation-messaging-test",
		Kind: "Ingress",
	},
}