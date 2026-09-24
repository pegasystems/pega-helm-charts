{{- /*
Helpers for pega-networkpolicy.yaml and pega-networkpolicy-custom.yaml.
Helpers that render YAML lists return unindented content; callers trim and indent it.
*/ -}}

{{- define "networkPolicyName" -}}
{{- $deploymentName := include "deploymentName" .root -}}
{{- $fullName := printf "%s-networkpolicy-%s" $deploymentName .suffix -}}
{{- if le (len $fullName) 63 -}}
{{- $fullName -}}
{{- else -}}
{{- printf "%s-%s" (printf "%s-networkpolicy" $deploymentName | trunc 54 | trimSuffix "-") (sha256sum $fullName | trunc 8) -}}
{{- end -}}
{{- end -}}

{{- /* Suffixes used by built-in policies. customPolicies must not reuse them. */ -}}
{{- define "networkPolicyBuiltInSuffixes" -}}
tiers,installer,constellation,hazelcast,clusteringservice,search,cassandra
{{- end -}}

{{- define "networkPolicyMetadata" -}}
metadata:
  name: {{ include "networkPolicyName" . }}
  namespace: {{ .root.Release.Namespace }}
{{- end -}}

{{- /* Selects the pods of all Pega tiers of this release by their app label. */ -}}
{{- define "networkPolicyTierSelector" -}}
matchExpressions:
- key: app
  operator: In
  values:
  {{- range .Values.global.tier }}
  - {{ printf "%s-%s" (include "deploymentName" $) .name | quote }}
  {{- end }}
{{- end -}}

{{- define "networkPolicyInstallerSelector" -}}
matchLabels:
  app: installer
{{- end -}}

{{- define "networkPolicyHazelcastSelector" -}}
matchLabels:
  app: {{ include "hazelcastName" . }}
  component: Hazelcast
{{- end -}}

{{- define "networkPolicyClusteringServiceSelector" -}}
matchLabels:
  app: {{ include "clusteringServiceName" . }}
  component: Hazelcast
{{- end -}}

{{- define "networkPolicySearchSelector" -}}
matchLabels:
  app: {{ include "searchName" . }}
  component: Search
{{- end -}}

{{- define "networkPolicyConstellationSelector" -}}
matchLabels:
  app: constellation
{{- end -}}

{{- define "networkPolicyCassandraSelector" -}}
matchLabels:
  app: cassandra
  release: {{ .Release.Name }}
{{- end -}}

{{- define "networkPolicyDeployTiers" -}}
{{- include "performDeployment" . | trim -}}
{{- end -}}

{{- define "networkPolicyRunInstaller" -}}
{{- if or (eq (include "performInstall" .) "true") (eq (include "performUpgrade" .) "true") -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{- /* Mirrors the isExternalSearch logic of the pegasearch subchart, evaluated from the parent chart. */ -}}
{{- define "networkPolicyInternalSearch" -}}
{{- $externalURL := .Values.pegasearch.externalURL | default "" -}}
{{- if .Values.pegasearch.externalSearchService -}}
false
{{- else if or (empty $externalURL) (eq $externalURL "http://pega-search") (eq $externalURL (include "defaultSearchURL" .)) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{- /*
Renders NetworkPolicyPort entries. Each port is either an integer (TCP) or a map with
protocol (TCP, UDP or SCTP) and port. Arguments: name, ports.
*/ -}}
{{- define "networkPolicyPorts" -}}
{{- if empty .ports -}}
{{- fail (printf "%s.ports must contain at least one port" .name) -}}
{{- end -}}
{{- range .ports }}
{{- $protocol := "TCP" }}
{{- $port := . }}
{{- if kindIs "map" . }}
{{- $protocol = .protocol | default "TCP" | toString | upper }}
{{- $port = .port }}
{{- end }}
{{- if not (has $protocol (list "TCP" "UDP" "SCTP")) }}
{{- fail (printf "%s.ports protocol must be TCP, UDP or SCTP" $.name) }}
{{- end }}
{{- $portString := toString $port }}
{{- if or (not (regexMatch "^[0-9]+$" $portString)) (lt (atoi $portString) 1) (gt (atoi $portString) 65535) }}
{{- fail (printf "%s.ports must be integers between 1 and 65535" $.name) }}
{{- end }}
- protocol: {{ $protocol }}
  port: {{ atoi $portString }}
{{- end }}
{{- end -}}

{{- /* Renders NetworkPolicyPeer entries from cidrs, namespaceSelector and podSelector. Empty selectors are ignored. */ -}}
{{- define "networkPolicyPeers" -}}
{{- range .cidrs }}
- ipBlock:
    cidr: {{ . | quote }}
{{- end }}
{{- if .namespaceSelector }}
- namespaceSelector:
    {{- toYaml .namespaceSelector | nindent 4 }}
  {{- with .podSelector }}
  podSelector:
    {{- toYaml . | nindent 4 }}
  {{- end }}
{{- else if .podSelector }}
- podSelector:
    {{- toYaml .podSelector | nindent 4 }}
{{- end }}
{{- end -}}

{{- /*
Renders one egress rule for a configured destination (cidrs/namespaceSelector/podSelector + ports).
Renders nothing for an unconfigured optional destination. Arguments: name, config, required.
*/ -}}
{{- define "networkPolicyDestinationRule" -}}
{{- $config := .config | default dict -}}
{{- range $config.cidrs -}}
{{- if not (regexMatch "^[0-9a-fA-F:.]+/[0-9]{1,3}$" (toString .)) -}}
{{- fail (printf "%s.cidrs entries must be CIDR blocks, for example 10.0.0.0/16" $.name) -}}
{{- end -}}
{{- end -}}
{{- $peers := include "networkPolicyPeers" $config | trim -}}
{{- if $peers -}}
- to:
  {{- $peers | nindent 2 }}
  ports:
  {{- include "networkPolicyPorts" (dict "name" .name "ports" $config.ports) | trim | nindent 2 }}
{{- else if .required -}}
{{- fail (printf "%s requires at least one of cidrs, namespaceSelector or podSelector" .name) -}}
{{- else if $config.ports -}}
{{- fail (printf "%s.ports is set but no cidrs, namespaceSelector or podSelector is configured" .name) -}}
{{- end -}}
{{- end -}}

{{- /*
DNS egress to any destination, added to every built-in policy. Port 5353 covers DNS pods that
listen on it behind the port-53 Service (for example openshift-dns).
*/ -}}
{{- define "networkPolicyDnsRule" -}}
- ports:
  - protocol: UDP
    port: 53
  - protocol: TCP
    port: 53
  - protocol: UDP
    port: 5353
  - protocol: TCP
    port: 5353
{{- end -}}

{{- /* Egress rule to pods of this release. Arguments: selector (YAML string), ports. */ -}}
{{- define "networkPolicyPodEgressRule" -}}
- to:
  - podSelector:
      {{- .selector | nindent 6 }}
  ports:
  {{- include "networkPolicyPorts" (dict "name" "networkPolicy" "ports" .ports) | trim | nindent 2 }}
{{- end -}}

{{- /* Ingress rule from pods of this release. Arguments: selectors (list of YAML strings), ports. */ -}}
{{- define "networkPolicyPodIngressRule" -}}
- from:
  {{- range .selectors }}
  - podSelector:
      {{- . | nindent 6 }}
  {{- end }}
  ports:
  {{- include "networkPolicyPorts" (dict "name" "networkPolicy" "ports" .ports) | trim | nindent 2 }}
{{- end -}}
