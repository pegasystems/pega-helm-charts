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

{{- define "hazelcastVolumeCredentials" }}hazelcast-volume-credentials{{- end }}

{{- define "hazelcastVolumeTemplate" }}
- name: {{ template "hazelcastVolumeCredentials" }}
  projected:
    defaultMode: 420
    sources:
    - secret:
        name: {{ template "genericSecretResolver" dict "externalSecretName" .Values.external_secret_name "valuesSecretName" "pega-hz-secret-name" }}
{{- end}}


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


{{- define "performDeployment" }}
  {{- if or (eq .Values.global.actions.execute "deploy") (eq .Values.global.actions.execute "install-deploy") (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}