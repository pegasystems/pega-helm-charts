package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

var networkPolicyTemplates = []string{
	"templates/pega-networkpolicy.yaml",
	"templates/pega-networkpolicy-custom.yaml",
}

func networkPolicyValues(overrides map[string]string) map[string]string {
	values := map[string]string{
		"global.provider":                 "k8s",
		"global.actions.execute":          "deploy",
		"cassandra.enabled":               "false",
		"networkPolicy.enabled":           "true",
		"networkPolicy.database.cidrs[0]": "10.0.0.0/24",
		"networkPolicy.database.ports[0]": "5432",
		"networkPolicy.kafka.cidrs[0]":    "10.2.0.0/24",
		"networkPolicy.kafka.ports[0]":    "9092",
	}
	for key, value := range overrides {
		if value == "null" {
			delete(values, key+"[0]")
		}
		values[key] = value
	}
	return values
}

func TestPegaNetworkPoliciesDisabledRenderNothing(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	options := &helm.Options{SetValues: map[string]string{"global.provider": "k8s"}}
	for _, template := range networkPolicyTemplates {
		_, err := RenderTemplateE(t, options, helmChartPath, []string{template})
		require.Error(t, err, template)
		require.Contains(t, err.Error(), "could not find template", template)
	}
}

func TestPegaNetworkPoliciesDeploy(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(nil)}, helmChartPath)
	require.ElementsMatch(t, []string{
		"pega-networkpolicy-tiers",
		"pega-networkpolicy-hazelcast",
		"pega-networkpolicy-search",
	}, policyNames(policies))

	tiers := policies["pega-networkpolicy-tiers"]
	require.Empty(t, tiers.Spec.PodSelector.MatchLabels)
	require.Equal(t, []metav1.LabelSelectorRequirement{{
		Key:      "app",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"pega-web", "pega-batch"},
	}}, tiers.Spec.PodSelector.MatchExpressions)
	require.ElementsMatch(t, []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress}, tiers.Spec.PolicyTypes)
	require.True(t, hasEgressToCIDR(tiers, "10.0.0.0/24", 5432))
	require.True(t, hasEgressPort(tiers, corev1.ProtocolUDP, 53), "every component policy must allow DNS")
	require.True(t, hasEgressToPod(tiers, "app", "pega-search", 9200))
	require.True(t, hasEgressToPod(tiers, "app", "pega-hazelcast", 5701))
	require.True(t, hasEgressPort(tiers, corev1.ProtocolUDP, 5353))
	require.True(t, hasEgressToCIDR(tiers, "10.2.0.0/24", 9092))
	for _, rule := range tiers.Spec.Ingress {
		for _, peer := range rule.From {
			require.Nil(t, peer.NamespaceSelector, "external ingress sources must come from customPolicies")
			require.Nil(t, peer.IPBlock, "external ingress sources must come from customPolicies")
		}
	}

	hazelcast := policies["pega-networkpolicy-hazelcast"]
	require.Equal(t, map[string]string{"app": "pega-hazelcast", "component": "Hazelcast"}, hazelcast.Spec.PodSelector.MatchLabels)
	require.True(t, hasEgressPort(hazelcast, corev1.ProtocolUDP, 53))
}

func TestPegaNetworkPoliciesAllComponents(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"global.actions.execute":             "install-deploy",
		"cassandra.enabled":                  "true",
		"constellation.enabled":              "true",
		"hazelcast.clusteringServiceEnabled": "true",
	})}, helmChartPath)

	require.ElementsMatch(t, []string{
		"pega-networkpolicy-tiers",
		"pega-networkpolicy-installer",
		"pega-networkpolicy-constellation",
		"pega-networkpolicy-hazelcast",
		"pega-networkpolicy-clusteringservice",
		"pega-networkpolicy-search",
		"pega-networkpolicy-cassandra",
	}, policyNames(policies))
	for name, policy := range policies {
		require.NotEqual(t, metav1.LabelSelector{}, policy.Spec.PodSelector, name)
		require.True(t, hasEgressPort(policy, corev1.ProtocolUDP, 53), name)
	}

	for _, name := range []string{"pega-networkpolicy-tiers", "pega-networkpolicy-installer"} {
		require.True(t, hasEgressToCIDR(policies[name], "10.0.0.0/24", 5432), name)
		require.True(t, hasEgressToPod(policies[name], "app", "cassandra", 9042), name)
	}
	require.True(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "constellation", 3000))
	require.True(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "clusteringservice", 5701))
	require.Equal(t, map[string]string{"app": "installer"}, policies["pega-networkpolicy-installer"].Spec.PodSelector.MatchLabels)
}

func TestPegaNetworkPoliciesMigrationJob(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"hazelcast.clusteringServiceEnabled":    "true",
		"hazelcast.migration.initiateMigration": "true",
	})}, helmChartPath)

	require.ElementsMatch(t, []string{
		"pega-networkpolicy-tiers",
		"pega-networkpolicy-hazelcast",
		"pega-networkpolicy-clusteringservice",
		"pega-networkpolicy-search",
	}, policyNames(policies), "the migration job must not be isolated")
}

func TestPegaNetworkPoliciesExternalSearch(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testCases := map[string]map[string]string{
		"external URL": {"pegasearch.externalURL": "https://search.example.com"},
		"SRS":          {"pegasearch.externalSearchService": "true", "pegasearch.externalURL": "http://srs:8080"},
	}
	for name, values := range testCases {
		t.Run(name, func(t *testing.T) {
			policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(values)}, helmChartPath)

			require.NotContains(t, policies, "pega-networkpolicy-search")
			require.False(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "pega-search", 9200))
		})
	}
}

func TestPegaNetworkPoliciesPortFormats(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"networkPolicy.kafka.ports[0].protocol": "tcp",
		"networkPolicy.kafka.ports[0].port":     "9142",
	})}, helmChartPath)

	require.True(t, hasEgressPort(policies["pega-networkpolicy-tiers"], corev1.ProtocolTCP, 9142))
}

func TestPegaNetworkPoliciesCustomPolicies(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"networkPolicy.customPolicies[0].name":                        "allow-agent",
		"networkPolicy.customPolicies[0].podSelector.matchLabels.app": "pega-web",
		"networkPolicy.customPolicies[0].policyTypes[0]":              "Egress",
	})}, helmChartPath, networkPolicyTemplates...)

	custom := policies["pega-networkpolicy-allow-agent"]
	require.Equal(t, map[string]string{"app": "pega-web"}, custom.Spec.PodSelector.MatchLabels)
	require.Equal(t, []networkingv1.PolicyType{networkingv1.PolicyTypeEgress}, custom.Spec.PolicyTypes)
}

func TestPegaNetworkPoliciesValidation(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		values   map[string]string
		expected string
	}{
		{"missing database", map[string]string{"networkPolicy.database.cidrs": "null"}, "networkPolicy.database requires at least one of"},
		{"database without ports", map[string]string{"networkPolicy.database.ports": "null"}, "networkPolicy.database.ports must contain at least one port"},
		{"invalid CIDR", map[string]string{"networkPolicy.database.cidrs[0]": "db.example.com"}, "networkPolicy.database.cidrs entries must be CIDR blocks"},
		{"invalid port", map[string]string{"networkPolicy.database.ports[0]": "70000"}, "networkPolicy.database.ports must be integers between 1 and 65535"},
		{"invalid protocol", map[string]string{"networkPolicy.kafka.ports[0].protocol": "ICMP", "networkPolicy.kafka.ports[0].port": "9092"}, "networkPolicy.kafka.ports protocol must be TCP, UDP or SCTP"},
		{"ports without peers", map[string]string{"stream.enabled": "false", "networkPolicy.kafka.cidrs": "null"}, "networkPolicy.kafka.ports is set but no cidrs"},
		{"kafka required with stream", map[string]string{"networkPolicy.kafka.cidrs": "null"}, "networkPolicy.kafka requires at least one of"},
		{"reserved custom name", map[string]string{"networkPolicy.customPolicies[0].name": "installer", "networkPolicy.customPolicies[0].podSelector.matchLabels.app": "x"}, "is reserved by a built-in policy"},
		{"duplicate custom name", map[string]string{
			"networkPolicy.customPolicies[0].name": "extra", "networkPolicy.customPolicies[0].podSelector.matchLabels.app": "x",
			"networkPolicy.customPolicies[1].name": "extra", "networkPolicy.customPolicies[1].podSelector.matchLabels.app": "y",
		}, "contains duplicate name"},
		{"custom without podSelector", map[string]string{"networkPolicy.customPolicies[0].name": "extra"}, "must define podSelector"},
		{"invalid custom name", map[string]string{"networkPolicy.customPolicies[0].name": "Extra_1", "networkPolicy.customPolicies[0].podSelector.matchLabels.app": "x"}, "must be a DNS-1123 label"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			options := &helm.Options{SetValues: networkPolicyValues(testCase.values)}
			_, err := RenderTemplateE(t, options, helmChartPath, []string{})
			require.Error(t, err)
			require.Contains(t, err.Error(), testCase.expected)
		})
	}
}

func renderNetworkPolicies(t *testing.T, options *helm.Options, chartPath string, templates ...string) map[string]networkingv1.NetworkPolicy {
	if len(templates) == 0 {
		templates = networkPolicyTemplates[:1]
	}
	yamlContent := RenderTemplate(t, options, chartPath, templates)
	policies := make(map[string]networkingv1.NetworkPolicy)

	for _, document := range strings.Split(yamlContent, "\n---") {
		if strings.TrimSpace(strings.ReplaceAll(document, "---", "")) == "" {
			continue
		}
		var policy networkingv1.NetworkPolicy
		helm.UnmarshalK8SYaml(t, document, &policy)
		if policy.Kind != "NetworkPolicy" {
			continue
		}
		require.NotContains(t, policies, policy.Name, "duplicate NetworkPolicy name")
		policies[policy.Name] = policy
	}

	return policies
}

func policyNames(policies map[string]networkingv1.NetworkPolicy) []string {
	names := make([]string, 0, len(policies))
	for name := range policies {
		names = append(names, name)
	}
	return names
}

func hasPort(ports []networkingv1.NetworkPolicyPort, protocol corev1.Protocol, port int) bool {
	expected := intstr.FromInt(port)
	for _, candidate := range ports {
		if candidate.Protocol != nil && *candidate.Protocol == protocol && candidate.Port != nil && *candidate.Port == expected {
			return true
		}
	}
	return false
}

func hasEgressPort(policy networkingv1.NetworkPolicy, protocol corev1.Protocol, port int) bool {
	for _, rule := range policy.Spec.Egress {
		if hasPort(rule.Ports, protocol, port) {
			return true
		}
	}
	return false
}

func hasEgressToCIDR(policy networkingv1.NetworkPolicy, cidr string, port int) bool {
	for _, rule := range policy.Spec.Egress {
		for _, peer := range rule.To {
			if peer.IPBlock != nil && peer.IPBlock.CIDR == cidr && hasPort(rule.Ports, corev1.ProtocolTCP, port) {
				return true
			}
		}
	}
	return false
}

func hasEgressToPod(policy networkingv1.NetworkPolicy, key string, value string, port int) bool {
	for _, rule := range policy.Spec.Egress {
		for _, peer := range rule.To {
			if peer.PodSelector != nil && peer.NamespaceSelector == nil && peer.PodSelector.MatchLabels[key] == value && hasPort(rule.Ports, corev1.ProtocolTCP, port) {
				return true
			}
		}
	}
	return false
}
