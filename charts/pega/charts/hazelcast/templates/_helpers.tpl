{{- define "hazelcastName" -}} pega-hazelcast {{- end -}}
{{- define "hazelcastEnvironmentConfig" -}} pega-hz-env-config {{- end -}}

{{- define "clusteringServiceName" -}} clusteringservice {{- end -}}
{{- define "clusteringServiceEnvironmentConfig" -}} clusteringservice-env-config {{- end -}}


{{- define "isHazelcastEnabled" }}
 {{- if .Values.enabled -}}
  true
 {{- else -}}
  false
 {{- end -}}
{{- end }}

{{- define "isClusteringServiceEnabled" }}
 {{- if .Values.clusteringServiceEnabled -}}
  true
 {{- else -}}
  false
 {{- end -}}
{{- end }}

{{- define "isMigrationInit" }}
 {{- if .Values.migration.initiateMigration -}}
  true
 {{- else -}}
  false
 {{- end -}}
{{- end }}

{{- define "isEncryptionEnabled" }}
  {{- if .Values.encryption.enabled  -}}
   true
  {{- else -}}
   false
  {{- end -}}
{{- end -}}

{{- define "isHighlySecureCryptoModeEnabled" }}
  {{- if and ( .Values.encryption.enabled ) ( .Values.global.highlySecureCryptoModeEnabled) -}}
   true
  {{- else -}}
   false
  {{- end -}}
{{- end -}}

{{- define "hazelcastVolumeCredentials" }}hazelcast-volume-credentials{{- end }}

{{- define "hazelcastVolumeTemplate" }}
- name: {{ template "hazelcastVolumeCredentials" }}
  projected:
    defaultMode: 420
    sources:
    {{- $d := dict "deploySecret" "deployHzServerSecret" "deployNonExtsecret" "deployNonExtHzServerSecret" "extSecretName" .Values.external_secret_name "nonExtSecretName" "pega-hz-secret-name" "context" $ -}}
    {{ include "secretResolver" $d | indent 4}}
{{- end}}

{{- define "deployHzServerSecret" -}}
true
{{- end }}

{{- define "deployNonExtHzServerSecret" }}
{{- if and (eq (include "deployHzServerSecret" .) "true") (not (.Values).external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

# Override this template to generate additional pod annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedHazelcastServicePodAnnotations" }}
{{- end }}

# Override this template to generate additional service annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedHazelcastServiceAnnotations" }}
{{- end }}

# Override this template to generate additional pod annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedClusteringServicePodAnnotations" }}
{{- end }}

# Override this template to generate additional pod labels that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedHazelcastServicePodLabels" }}
{{- end }}

# Override this template to generate additional service annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedClusteringServiceAnnotations" }}
{{- end }}

# Override this template to generate additional pod labels that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedClusteringServicePodLabels" }}
{{- end }}

# Override this template to generate additional Clustering Service StatefulSet annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedClusteringServiceStatefulSetAnnotations" }}
{{- end }}

{{- /*
Renders a map of annotations as quoted strings. Values are evaluated with `tpl`, so they can reference other values.
- null value: key is omitted (use in a stage override file to remove a key)
- template value that renders to an empty string: key is omitted (use for conditional keys)
- literal "" value: key is kept with an empty value
*/}}
{{- define "renderAnnotations" }}
{{- $root := .root }}
{{- range $key, $value := .annotations }}
{{- if not (kindIs "invalid" $value) }}
{{- $rendered := tpl (toString $value) $root }}
{{- if or $rendered (eq (toString $value) "") }}
{{ $key }}: {{ $rendered | quote }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{- define "clusteringServiceStatefulSetAnnotations" }}
{{- include "renderAnnotations" (dict "root" . "annotations" .Values.statefulsetAnnotations) }}
{{- include "generatedClusteringServiceStatefulSetAnnotations" . }}
{{- end }}

# Override this template to generate additional Hazelcast StatefulSet annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedHazelcastStatefulSetAnnotations" }}
{{- end }}

{{- define "hazelcastStatefulSetAnnotations" }}
{{- include "renderAnnotations" (dict "root" . "annotations" .Values.statefulsetAnnotations) }}
{{- include "generatedHazelcastStatefulSetAnnotations" . }}
{{- end }}

{{- /*
Returns a CPU quantity (e.g. 1, "1.5", "500m") scaled by a percentage, in millicores.
Usage: {{ include "cpuRequestPercent" (dict "cpu" .Values.resources.requests.cpu "percent" 150) }} -> 1500m
*/}}
{{- define "cpuRequestPercent" }}
{{- $cpu := toString .cpu }}
{{- $millis := 0.0 }}
{{- if hasSuffix "m" $cpu }}
{{- $millis = float64 (trimSuffix "m" $cpu) }}
{{- else }}
{{- $millis = mulf (float64 $cpu) 1000 }}
{{- end }}
{{- printf "%dm" (int (floor (divf (mulf $millis (float64 .percent)) 100))) }}
{{- end }}


{{- define "performDeployment" }}
  {{- if or (eq .Values.global.actions.execute "deploy") (eq .Values.global.actions.execute "install-deploy") (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}