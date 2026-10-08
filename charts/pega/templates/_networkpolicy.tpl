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

{{- /* Matches the labels of the cassandra subchart (cassandra.name). */ -}}
{{- define "networkPolicyCassandraSelector" -}}
matchLabels:
  app: {{ default "cassandra" (.Values.cassandra).nameOverride | trunc 63 | trimSuffix "-" }}
  release: {{ .Release.Name }}
{{- end -}}

{{- /* CQL container port of the cassandra subchart. */ -}}
{{- define "networkPolicyCassandraPort" -}}
{{- default 9042 (((.Values.cassandra).config).ports).cql -}}
{{- end -}}

{{- /* TCP ports of the tier pods: HTTP(S), embedded Hazelcast and TCP tier[].custom.ports. Returns a comma-separated list. */ -}}
{{- define "networkPolicyTierPorts" -}}
{{- $ports := list 8080 8443 5701 -}}
{{- range .Values.global.tier -}}
{{- range ((.custom).ports | default list) -}}
{{- if and .containerPort (eq (.protocol | default "TCP") "TCP") -}}
{{- $ports = append $ports (.containerPort | int) -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- $ports | uniq | join "," -}}
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
{{- range .rules -}}
{{- if not (and (kindIs "map" .) (or .to .ports)) -}}
{{- fail (printf "%s entries must be NetworkPolicy egress rules with to and/or ports; an empty rule would allow all egress" $.name) -}}
{{- end -}}
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
