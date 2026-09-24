{{- /*
deploymentName
networkPolicyName
pegaRegistrySecret
imagePullSecrets
pegaVolumeCredentials
customArtifactorySSLVerificationEnabled
performDeployment
performInstallAndDeployment
performUpgradeAndDeployment
pega-db-secret-name
pega-hz-secret-name
deployDBSecret
deployNonExtDBSecret
podAffinity
tolerations
secretResolver are copied from pega/templates/_helpers.tpl because helm lint requires
charts to render standalone. See: https://github.com/helm/helm/issues/11260 for more details.
*/}}


{{- define "deploymentName" }}{{ $deploymentNamePrefix := "pega" }}{{ if (.Values.global.deployment) }}{{ if (.Values.global.deployment.name) }}{{ $deploymentNamePrefix = .Values.global.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "networkPolicyName" -}}
{{- $fullName := printf "%s-networkpolicy-%s" .deploymentName .suffix -}}
{{- if le (len $fullName) 63 -}}
{{- $fullName -}}
{{- else -}}
{{- printf "%s-%s" (printf "%s-networkpolicy" .deploymentName | trunc 54 | trimSuffix "-") (sha256sum $fullName | trunc 8) -}}
{{- end -}}
{{- end -}}

{{- define "pegaServiceAccountName" -}}
{{- if .Values.global.serviceAccount.enabled -}}
{{- if .Values.global.serviceAccount.name -}}
{{- .Values.global.serviceAccount.name -}}
{{- else -}}
{{- printf "%s-serviceaccount" (include "deploymentName" .) | trunc 253 | trimSuffix "-" -}}
{{- end -}}
{{- else -}}
{{- printf "%s-serviceaccount" (include "deploymentName" .) | trunc 253 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "pegaInstallerServiceAccountName" -}}
{{- if .useGlobal -}}
{{- include "pegaServiceAccountName" .root -}}
{{- else if .config.name -}}
{{- .config.name -}}
{{- else -}}
{{- printf "%s-installer-serviceaccount" (include "deploymentName" .root) | trunc 253 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "validateServiceAccountName" -}}
{{- if not (regexMatch "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$" .name) -}}
{{- fail (printf "%s must be a valid DNS subdomain" .message) -}}
{{- end -}}
{{- if gt (len .name) 253 -}}
{{- fail (printf "%s must be no longer than 253 characters" .message) -}}
{{- end -}}
{{- end -}}

{{- define "validateServiceAccountBoolean" -}}
{{- if not (kindIs "bool" .value) -}}
{{- fail (printf "%s must be a boolean" .message) -}}
{{- end -}}
{{- end -}}

{{- define "validateNetworkPolicyPorts" -}}
{{- range .ports }}
{{- $port := toString . }}
{{- if not (regexMatch "^[0-9]+$" $port) }}
{{- fail (printf "%s ports must be integers between 1 and 65535" $.name) }}
{{- end }}
{{- if or (lt (atoi $port) 1) (gt (atoi $port) 65535) }}
{{- fail (printf "%s ports must be integers between 1 and 65535" $.name) }}
{{- end }}
{{- end }}
{{- end }}

{{- define "pegaVolumeCredentials" }}pega-volume-credentials{{- end }}

{{- define "initContainerResources" }}
  resources:
    # Resources requests/limits for initContainers
    requests:
      cpu: 50m
      memory: 64Mi
    limits:
      cpu: 50m
      memory: 64Mi
{{- end }}

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
{{- if and .Values.global.docker.registry (not .Values.global.docker.imagePullSecretNames) }}
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

{{- define "pega-upgrade-rest-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-upgrade-rest-secret
{{- end -}}

{{- define "pega-hz-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-hz-secret
{{- end -}}

{{- define "deployDBSecret" -}} 
true
{{- end }}

{{- define "deployNonExtDBSecret" }}
{{- if and (eq (include "deployDBSecret" .) "true") (not (.Values.global.jdbc).external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{- define "deployPegaRESTSecret" -}}
true
{{- end }}

{{- define "deployNonExtPegaRESTSecret" }}
{{- if and (eq (include "deployPegaRESTSecret" .) "true") (not (.Values.upgrade).pega_rest_external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{- define "secretResolver" }}
{{- if (eq (include .deploySecret .context) "true") }}
- secret:
{{- if (eq (include .deployNonExtsecret .context) "true") }}
    name: {{ include .nonExtSecretName .context}}
{{- else }}
    name: {{ .extSecretName }}
{{- end -}}
{{- end -}}
{{- end  -}}

{{- define "podAffinity" }}
{{- if .affinity }}
affinity:
{{- toYaml .affinity | nindent 2 }}
{{- end }}
{{ end }}

{{- define "tolerations" }}
{{- if .tolerations }}
tolerations:
{{- toYaml .tolerations | nindent 2 }}
{{- end }}
{{ end }}