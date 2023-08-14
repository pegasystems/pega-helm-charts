{{- /*
deploymentName
pegaCredentialsSecret
pegaRegistrySecret
imagePullSecrets
pegaCredentialVolumeTemplate
pegaVolumeCredentials
customArtifactorySSLVerificationEnabled
performDeployment
performInstallAndDeployment
performUpgradeAndDeployment are copied from pega/templates/_helpers.tpl because helm lint requires
charts to render standalone. See: https://github.com/helm/helm/issues/11260 for more details.
*/}}


{{- define "deploymentName" }}{{ $deploymentNamePrefix := "pega" }}{{ if (.Values.global.deployment) }}{{ if (.Values.global.deployment.name) }}{{ $deploymentNamePrefix = .Values.global.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "pegaVolumeCredentials" }}pega-volume-credentials{{- end }}

{{- define "pegaCredentialsSecret" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-credentials-secret
{{- end }}

{{- define "pegaCredentialVolumeTemplate" }}
- name: {{ template "pegaVolumeCredentials" }}
  projected:
    defaultMode: 420
    sources:
    - secret:
        name: {{ template "pegaCredentialsSecret" $ }}
  {{ if ((.Values.global.jdbc).external_secret_name) }}
    - secret:
        name: {{ .Values.global.jdbc.external_secret_name }}
  {{- end }}
  {{ if ((.Values.hazelcast).external_secret_name)}}
    - secret:
        name: {{ .Values.hazelcast.external_secret_name }}
  {{- end }}
  {{ if ((.Values.global.customArtifactory.authentication).external_secret_name) }}
    - secret:
        name: {{ .Values.global.customArtifactory.authentication.external_secret_name }}
  {{- end }}
  {{ if ((.Values.dds).external_secret_name)}}
    - secret:
        name: {{ .Values.dds.external_secret_name }}
  {{- end }}
  {{ if ((.Values.stream).external_secret_name)}}
    - secret:
        name: {{ .Values.stream.external_secret_name }}
  {{- end }}
{{- end}}

{{- define "customArtifactorySSLVerificationEnabled" }}
{{- if (.Values.global.customArtifactory) }}
{{- if (.Values.global.customArtifactory.enableSSLVerification) }}
{{- if (eq .Values.global.customArtifactory.enableSSLVerification true) -}}
true
{{- else -}}
false
{{- end }}
{{- end }}
{{- end }}
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


{{- define "performDeployment" }}
  {{- if or (eq .Values.global.actions.execute "deploy") (eq .Values.global.actions.execute "install-deploy") (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performInstallAndDeployment" }}
  {{- if (eq .Values.global.actions.execute "install-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performUpgradeAndDeployment" }}
  {{- if (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}