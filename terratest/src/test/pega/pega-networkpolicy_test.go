package pega

import (
	"path/filepath"
	"reflect"
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
		"global.provider":                              "k8s",
		"global.actions.execute":                       "deploy",
		"cassandra.enabled":                            "false",
		"networkPolicy.enabled":                        "true",
		"networkPolicy.database[0].to[0].ipBlock.cidr": "10.0.0.0/24",
		"networkPolicy.database[0].ports[0].protocol":  "TCP",
		"networkPolicy.database[0].ports[0].port":      "5432",
		"networkPolicy.kafka[0].to[0].ipBlock.cidr":    "10.2.0.0/24",
		"networkPolicy.kafka[0].ports[0].protocol":     "TCP",
		"networkPolicy.kafka[0].ports[0].port":         "9092",
	}
	// Deletions run first so that a "null" override never removes a key set by another override.
	for key, value := range overrides {
		if value == "null" {
			for existing := range values {
				if strings.HasPrefix(existing, key) {
					delete(values, existing)
				}
			}
		}
	}
	for key, value := range overrides {
		if value != "null" {
			values[key] = value
		}
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
	require.True(t, hasEgressToPod(tiers, "app", "pega-search", 9200))
	require.True(t, hasEgressToPod(tiers, "app", "pega-hazelcast", 5701))
	require.True(t, hasEgressToCIDR(tiers, "10.2.0.0/24", 9092))
	for _, rule := range tiers.Spec.Ingress {
		for _, peer := range rule.From {
			require.Nil(t, peer.NamespaceSelector, "external ingress sources must come from customPolicies")
			require.Nil(t, peer.IPBlock, "external ingress sources must come from customPolicies")
		}
	}

	hazelcast := policies["pega-networkpolicy-hazelcast"]
	require.Equal(t, map[string]string{"app": "pega-hazelcast", "component": "Hazelcast"}, hazelcast.Spec.PodSelector.MatchLabels)
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
		require.False(t, hasEgressPort(policy, corev1.ProtocolUDP, 53), "DNS must come from customPolicies: "+name)
	}

	for _, name := range []string{"pega-networkpolicy-tiers", "pega-networkpolicy-installer"} {
		require.True(t, hasEgressToCIDR(policies[name], "10.0.0.0/24", 5432), name)
		require.True(t, hasEgressToPod(policies[name], "app", "cassandra", 9042), name)
	}
	require.True(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "constellation", 3000))
	require.True(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "clusteringservice", 5701))
	require.Equal(t, map[string]string{"app": "installer"}, policies["pega-networkpolicy-installer"].Spec.PodSelector.MatchLabels)

	tiers := policies["pega-networkpolicy-tiers"]
	require.True(t, hasIngressFromPod(tiers, "app", "installer", 8080))
	require.True(t, hasIngressFromPod(tiers, "app", "constellation", 8080))
	require.False(t, hasIngressFromPod(tiers, "app", "pega-search", 8080))
	require.True(t, hasIngressFromPod(policies["pega-networkpolicy-hazelcast"], "app", "pega-hazelcast", 5701))
	require.True(t, hasIngressFromPod(policies["pega-networkpolicy-clusteringservice"], "app", "clusteringservice", 5701))
	require.True(t, hasIngressFromPod(policies["pega-networkpolicy-search"], "app", "pega-search", 9300))
	cassandra := policies["pega-networkpolicy-cassandra"]
	require.True(t, hasIngressFromPod(cassandra, "app", "installer", 9042))
	require.True(t, hasIngressFromPod(cassandra, "app", "cassandra", 7000))
	require.Empty(t, policies["pega-networkpolicy-installer"].Spec.Ingress)
	require.Equal(t, []networkingv1.PolicyType{networkingv1.PolicyTypeEgress}, policies["pega-networkpolicy-installer"].Spec.PolicyTypes)
	require.True(t, hasIngressFromTierSelectorPort(tiers, 5701))
	for name, port := range map[string]int{"pega-networkpolicy-hazelcast": 5701, "pega-networkpolicy-search": 9200, "pega-networkpolicy-constellation": 3000} {
		require.True(t, hasIngressFromTierSelectorPort(policies[name], port), name)
	}
}

func TestPegaNetworkPoliciesActions(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testCases := map[string][]string{
		"install":        {"pega-networkpolicy-installer"},
		"upgrade":        {"pega-networkpolicy-installer"},
		"upgrade-deploy": {"pega-networkpolicy-tiers", "pega-networkpolicy-installer", "pega-networkpolicy-hazelcast", "pega-networkpolicy-search"},
	}
	for action, expected := range testCases {
		t.Run(action, func(t *testing.T) {
			policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
				"global.actions.execute":        action,
				"installer.upgrade.upgradeType": getUpgradeTypeForUpgradeAction(action),
			})}, helmChartPath)

			require.ElementsMatch(t, expected, policyNames(policies))
			installer := policies["pega-networkpolicy-installer"]
			require.True(t, hasEgressToCIDR(installer, "10.0.0.0/24", 5432))
			require.Equal(t, action == "upgrade-deploy", hasEgressToTierSelector(installer), "installer reaches tiers only when they are deployed")
		})
	}
}

func TestPegaNetworkPoliciesOptionalComponents(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"stream.enabled":      "false",
		"networkPolicy.kafka": "null",
		"hazelcast.enabled":   "false",
	})}, helmChartPath)

	require.ElementsMatch(t, []string{"pega-networkpolicy-tiers", "pega-networkpolicy-search"}, policyNames(policies))
	tiers := policies["pega-networkpolicy-tiers"]
	require.False(t, hasEgressToPod(tiers, "app", "pega-hazelcast", 5701))
	require.False(t, hasEgressToCIDR(tiers, "10.2.0.0/24", 9092))
}

func TestPegaNetworkPoliciesCassandraOverrides(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"cassandra.enabled":          "true",
		"cassandra.nameOverride":     "dds",
		"cassandra.config.ports.cql": "9142",
	})}, helmChartPath)

	cassandra := policies["pega-networkpolicy-cassandra"]
	require.Equal(t, "dds", cassandra.Spec.PodSelector.MatchLabels["app"])
	require.True(t, hasIngressFromTierSelectorPort(cassandra, 9142))
	require.True(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "dds", 9142))
	require.False(t, hasEgressToPod(policies["pega-networkpolicy-tiers"], "app", "dds", 9042))
}

func TestPegaNetworkPoliciesTierCustomPorts(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"global.tier[0].name":                          "web",
		"global.tier[0].custom.ports[0].containerPort": "9100",
		"global.tier[0].custom.ports[1].containerPort": "9101",
		"global.tier[0].custom.ports[1].protocol":      "UDP",
		"global.tier[1].name":                          "batch",
		"global.tier[1].custom.ports[0].containerPort": "8080",
	})}, helmChartPath)

	tiers := policies["pega-networkpolicy-tiers"]
	require.True(t, hasIngressFromTierSelectorPort(tiers, 9100))
	require.True(t, hasEgressToPod(tiers, "app", "pega-hazelcast", 5701))
	require.False(t, hasIngressFromTierSelectorPort(tiers, 9101), "non-TCP custom ports must come from customPolicies")
	require.Len(t, tiers.Spec.Ingress[0].Ports, 4, "tier ports must be de-duplicated")
	require.Len(t, tiers.Spec.Egress[0].Ports, 4)
}

func TestPegaNetworkPoliciesDeploymentName(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"global.deployment.name": "prod",
	})}, helmChartPath)
	require.ElementsMatch(t, []string{"prod-networkpolicy-tiers", "prod-networkpolicy-hazelcast", "prod-networkpolicy-search"}, policyNames(policies))
	require.Equal(t, []string{"prod-web", "prod-batch"}, policies["prod-networkpolicy-tiers"].Spec.PodSelector.MatchExpressions[0].Values)

	longName := strings.Repeat("a", 60)
	policies = renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"global.deployment.name": longName,
	})}, helmChartPath)
	require.Len(t, policies, 3)
	for name := range policies {
		require.LessOrEqual(t, len(name), 63, name)
		require.True(t, strings.HasPrefix(name, longName[:40]), name)
	}
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

func TestPegaNetworkPoliciesExternalRulesPassThrough(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	policies := renderNetworkPolicies(t, &helm.Options{SetValues: networkPolicyValues(map[string]string{
		"global.actions.execute":                                          "install-deploy",
		"networkPolicy.database[0].to[0].ipBlock.except[0]":               "10.0.0.128/25",
		"networkPolicy.database[0].ports[1].protocol":                     "UDP",
		"networkPolicy.database[0].ports[1].port":                         "5433",
		"networkPolicy.kafka[0].to[0].ipBlock.cidr":                       "null",
		"networkPolicy.kafka[0].to[0].namespaceSelector.matchLabels.team": "streaming",
	})}, helmChartPath)

	for _, name := range []string{"pega-networkpolicy-tiers", "pega-networkpolicy-installer"} {
		rule := egressRuleToCIDR(policies[name], "10.0.0.0/24")
		require.NotNil(t, rule, name)
		require.Equal(t, []string{"10.0.0.128/25"}, rule.To[0].IPBlock.Except, name)
		require.True(t, hasPort(rule.Ports, corev1.ProtocolTCP, 5432), name)
		require.True(t, hasPort(rule.Ports, corev1.ProtocolUDP, 5433), name)
	}
	require.True(t, hasEgressToNamespace(policies["pega-networkpolicy-tiers"], "team", "streaming", 9092))
	require.False(t, hasEgressToNamespace(policies["pega-networkpolicy-installer"], "team", "streaming", 9092), "kafka rules apply to tiers only")
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
		{"missing database", map[string]string{"networkPolicy.database": "null"}, "networkPolicy.database must contain at least one NetworkPolicy egress rule"},
		{"database not a list", map[string]string{"networkPolicy.database": "null", "networkPolicy.database.cidrs[0]": "10.0.0.0/24"}, "networkPolicy.database must be a list of NetworkPolicy egress rules"},
		{"database rule without to or ports", map[string]string{"networkPolicy.database": "null", "networkPolicy.database[0].description": "any"}, "networkPolicy.database entries must be NetworkPolicy egress rules"},
		{"scalar kafka rule", map[string]string{"networkPolicy.kafka": "null", "networkPolicy.kafka[0]": "10.2.0.0/24"}, "networkPolicy.kafka entries must be NetworkPolicy egress rules"},
		{"kafka required with stream", map[string]string{"networkPolicy.kafka": "null"}, "networkPolicy.kafka must contain at least one NetworkPolicy egress rule"},
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

func egressRuleToCIDR(policy networkingv1.NetworkPolicy, cidr string) *networkingv1.NetworkPolicyEgressRule {
	for i, rule := range policy.Spec.Egress {
		for _, peer := range rule.To {
			if peer.IPBlock != nil && peer.IPBlock.CIDR == cidr {
				return &policy.Spec.Egress[i]
			}
		}
	}
	return nil
}

func hasEgressToNamespace(policy networkingv1.NetworkPolicy, key string, value string, port int) bool {
	for _, rule := range policy.Spec.Egress {
		for _, peer := range rule.To {
			if peer.NamespaceSelector != nil && peer.NamespaceSelector.MatchLabels[key] == value && hasPort(rule.Ports, corev1.ProtocolTCP, port) {
				return true
			}
		}
	}
	return false
}

func hasIngressFromPod(policy networkingv1.NetworkPolicy, key string, value string, port int) bool {
	for _, rule := range policy.Spec.Ingress {
		for _, peer := range rule.From {
			if peer.PodSelector != nil && peer.NamespaceSelector == nil && peer.PodSelector.MatchLabels[key] == value && hasPort(rule.Ports, corev1.ProtocolTCP, port) {
				return true
			}
		}
	}
	return false
}

// isTierSelector matches the selector of the default tiers (pega-web, pega-batch).
func isTierSelector(selector *metav1.LabelSelector) bool {
	return selector != nil && len(selector.MatchLabels) == 0 && reflect.DeepEqual(selector.MatchExpressions, []metav1.LabelSelectorRequirement{{
		Key:      "app",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"pega-web", "pega-batch"},
	}})
}

func hasIngressFromTierSelectorPort(policy networkingv1.NetworkPolicy, port int) bool {
	for _, rule := range policy.Spec.Ingress {
		for _, peer := range rule.From {
			if isTierSelector(peer.PodSelector) && hasPort(rule.Ports, corev1.ProtocolTCP, port) {
				return true
			}
		}
	}
	return false
}

func hasEgressToTierSelector(policy networkingv1.NetworkPolicy) bool {
	for _, rule := range policy.Spec.Egress {
		for _, peer := range rule.To {
			if isTierSelector(peer.PodSelector) {
				return true
			}
		}
	}
	return false
}
