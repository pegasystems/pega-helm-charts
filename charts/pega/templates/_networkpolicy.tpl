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

{{- /* Renders TCP NetworkPolicyPort entries for the built-in rules. Argument: list of ports. */ -}}
{{- define "networkPolicyPorts" -}}
{{- range . }}
- protocol: TCP
  port: {{ . }}
{{- end }}
{{- end -}}

{{- /*
Renders user-supplied NetworkPolicyEgressRule entries (networkPolicy.database, networkPolicy.kafka).
Arguments: name, rules, required.
*/ -}}
{{- define "networkPolicyEgressRules" -}}
{{- if not (kindIs "slice" (.rules | default list)) -}}
{{- fail (printf "%s must be a list of NetworkPolicy egress rules" .name) -}}
{{- end -}}
{{- if .rules -}}
{{- toYaml .rules -}}
{{- else if .required -}}
{{- fail (printf "%s must contain at least one NetworkPolicy egress rule" .name) -}}
{{- end -}}
{{- end -}}

{{- /* Egress rule to pods of this release. Arguments: selector (YAML string), ports. */ -}}
{{- define "networkPolicyPodEgressRule" -}}
- to:
  - podSelector:
      {{- .selector | nindent 6 }}
  ports:
  {{- include "networkPolicyPorts" .ports | trim | nindent 2 }}
{{- end -}}

{{- /* Ingress rule from pods of this release. Arguments: selectors (list of YAML strings), ports. */ -}}
{{- define "networkPolicyPodIngressRule" -}}
- from:
  {{- range .selectors }}
  - podSelector:
      {{- . | nindent 6 }}
  {{- end }}
  ports:
  {{- include "networkPolicyPorts" .ports | trim | nindent 2 }}
{{- end -}}
