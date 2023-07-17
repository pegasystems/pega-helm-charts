{{- /*
deploymentName and pegaCredentialsSecret are copied from pega/templates/_helpers.tpl because helm lint requires
charts to render standalone. See: https://github.com/helm/helm/issues/11260 for more details.
*/}}

{{- define "deploymentName" }}{{ $deploymentNamePrefix := "pega" }}{{ if (.Values.global.deployment) }}{{ if (.Values.global.deployment.name) }}{{ $deploymentNamePrefix = .Values.global.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "pegaCredentialsSecret" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-credentials-secret
{{- end }}

{{- define "pegaRegistrySecret" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-registry-secret
{{- end }}

{{- define "imagePullSecrets" }}
{{- if .Values.global.docker.registry }}
- name: {{ template "pegaRegistrySecret" $ }}
{{- end }}
{{- if (.Values.global.docker.imagePullSecretNames) }}
{{- range .Values.global.docker.imagePullSecretNames }}
- name: {{ . }}
{{- end -}}
{{- end -}}
{{- end -}}


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
        name: {{ template "pegaCredentialsSecret" $ }}
  {{ if ((.Values.global.jdbc).external_secret_name) }}
    - secret:
        name: {{ .Values.global.jdbc.external_secret_name }}
  {{- end }}
  {{ if (.Values.external_secret_name)}}
    - secret:
        name: {{ .Values.external_secret_name }}
  {{- end }}
  {{ if ((.Values.global.customArtifactory.authentication).external_secret_name) }}
    - secret:
        name: {{ .Values.global.customArtifactory.authentication.external_secret_name }}
  {{- end }}
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
