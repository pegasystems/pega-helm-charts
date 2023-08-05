{{- /*
deploymentName
pegaRegistrySecret
imagePullSecrets
pegaVolumeCredentials
customArtifactorySSLVerificationEnabled
performDeployment
performInstallAndDeployment
performUpgradeAndDeployment
genericSecretResolver
pegaDBSecretName
pegaHZSecretName are copied from pega/templates/_helpers.tpl because helm lint requires
charts to render standalone. See: https://github.com/helm/helm/issues/11260 for more details.
*/}}


{{- define "deploymentName" }}{{ $deploymentNamePrefix := "pega" }}{{ if (.Values.global.deployment) }}{{ if (.Values.global.deployment.name) }}{{ $deploymentNamePrefix = .Values.global.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "pegaVolumeCredentials" }}pega-volume-credentials{{- end }}

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

{{- define "pega-db-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-db-secret
{{- end -}}

{{- define "pega-hz-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-hz-secret
{{- end -}}

{{- define "genericSecretResolver" }}
    {{- if (.externalSecretName) -}}{{ .externalSecretName }}{{- else -}}{{ (include .valuesSecretName .context) }}{{- end -}}
{{- end  -}}