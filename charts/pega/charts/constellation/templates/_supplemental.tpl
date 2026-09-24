{{- /*
deploymentName
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

The ServiceAccount helpers (pegaServiceAccountName ... validateServiceAccountBoolean) are shared
by the pega chart and all subcharts. This file is the master copy; run sync_supplementals.sh after
editing it.
*/}}


{{- define "deploymentName" }}{{ $deploymentNamePrefix := "pega" }}{{ if (.Values.global.deployment) }}{{ if (.Values.global.deployment.name) }}{{ $deploymentNamePrefix = .Values.global.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "pegaServiceAccountName" -}}
{{- $config := (.Values.global).serviceAccount | default dict -}}
{{- $config.name | default (printf "%s-serviceaccount" (include "deploymentName" .)) -}}
{{- end -}}

{{- define "pegaInstallerServiceAccountDefaultName" -}}
{{- printf "%s-installer-serviceaccount" (include "deploymentName" .) -}}
{{- end -}}

{{- /*
pegaResolveServiceAccountName returns the ServiceAccount a pod should use, or an empty
string when the pod keeps the namespace default. An explicit override always wins.
Arguments: override, config, configPath, generatedName.
*/ -}}
{{- define "pegaResolveServiceAccountName" -}}
{{- if .override -}}
{{- .override -}}
{{- else if .config.enabled -}}
{{- if .config.name -}}
{{- .config.name -}}
{{- else if .config.create -}}
{{- .generatedName -}}
{{- else -}}
{{- fail (printf "%s.enabled requires %s.name or %s.create=true" .configPath .configPath .configPath) -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- /*
pegaServiceAccountSpec renders serviceAccountName and automountServiceAccountToken for a pod spec.
automountServiceAccountToken is only set for chart-managed accounts, never for explicit overrides.
Arguments: override, overridePath, config, configPath, generatedName.
*/ -}}
{{- define "pegaServiceAccountSpec" -}}
{{- $name := include "pegaResolveServiceAccountName" . -}}
{{- if $name -}}
{{- if .override -}}
{{- include "validateServiceAccountName" (dict "name" $name "message" .overridePath) -}}
serviceAccountName: {{ $name }}
{{- else -}}
{{- include "validateServiceAccountName" (dict "name" $name "message" (printf "%s.name" .configPath)) -}}
{{- include "validateServiceAccountBoolean" (dict "value" .config.automountServiceAccountToken "message" (printf "%s.automountServiceAccountToken" .configPath)) -}}
serviceAccountName: {{ $name }}
automountServiceAccountToken: {{ .config.automountServiceAccountToken }}
{{- end -}}
{{- end -}}
{{- end -}}

{{- /* Pod spec ServiceAccount fields for workloads that use the shared global.serviceAccount. */ -}}
{{- define "pegaWorkloadServiceAccountSpec" -}}
{{- include "pegaServiceAccountSpec" (dict
      "override" (.override | default "")
      "overridePath" .overridePath
      "config" ((.root.Values.global).serviceAccount | default dict)
      "configPath" "global.serviceAccount"
      "generatedName" (include "pegaServiceAccountName" .root)) -}}
{{- end -}}

{{- /*
pegaInstallerServiceAccountArgs returns (as JSON) the pegaServiceAccountSpec arguments for installer jobs.
Must be called with the installer chart context. installer.serviceAccount takes precedence over
global.serviceAccount; installer.serviceAccountName overrides both.
*/ -}}
{{- define "pegaInstallerServiceAccountArgs" -}}
{{- $args := dict
      "override" (.Values.serviceAccountName | default "")
      "overridePath" "installer.serviceAccountName"
      "config" (.Values.serviceAccount | default dict)
      "configPath" "installer.serviceAccount"
      "generatedName" (include "pegaInstallerServiceAccountDefaultName" .) -}}
{{- $global := (.Values.global).serviceAccount | default dict -}}
{{- if and (not $args.config.enabled) $global.enabled -}}
{{- $_ := set $args "config" $global -}}
{{- $_ := set $args "configPath" "global.serviceAccount" -}}
{{- $_ := set $args "generatedName" (include "pegaServiceAccountName" .) -}}
{{- end -}}
{{- toJson $args -}}
{{- end -}}

{{- /*
pegaServiceAccountResource renders a ServiceAccount when config.enabled and config.create are true.
Arguments: config, configPath, generatedName, namespace.
*/ -}}
{{- define "pegaServiceAccountResource" -}}
{{- $config := .config | default dict -}}
{{- if and $config.create (not $config.enabled) -}}
{{- fail (printf "%s.create requires %s.enabled=true" .configPath .configPath) -}}
{{- end -}}
{{- if $config.enabled -}}
{{- include "validateServiceAccountBoolean" (dict "value" $config.automountServiceAccountToken "message" (printf "%s.automountServiceAccountToken" .configPath)) -}}
{{- end -}}
{{- if and $config.enabled $config.create -}}
{{- $name := $config.name | default .generatedName -}}
{{- include "validateServiceAccountName" (dict "name" $name "message" (printf "%s.name" .configPath)) -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ $name }}
  namespace: {{ .namespace }}
  {{- with $config.labels }}
  labels:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with $config.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
automountServiceAccountToken: {{ $config.automountServiceAccountToken }}
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