{{- /*
imagePullSecret
backingservicesRegistrySecret
deploymentName
tlssecretsnippet
backingservices.gke.backendConfig
podAffinity
tolerations
are copied from backingservices/templates/_supplemental.tpl because helm lint requires
charts to render standalone. See: https://github.com/helm/helm/issues/11260 for more details.
*/}}

{{- define "imagePullSecret" }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.docker.registry.url (printf "%s:%s" .Values.docker.registry.username .Values.docker.registry.password | b64enc) | b64enc }}
{{- end }}

{{- define "backingservicesRegistrySecret" }}
{{- $depName := printf "%s" (include "deploymentName" (dict "root" .root "defaultname" .defaultname )) -}}
{{- $depName -}}-registry-secret
{{- end }}

{{- define "deploymentName" }}{{ $deploymentNamePrefix := .defaultname }}{{ if (.root.deployment) }}{{ if (.root.deployment.name) }}{{ $deploymentNamePrefix = .root.deployment.name }}{{ end }}{{ end }}{{ if (.root.name) }}{{ $deploymentNamePrefix = .root.name }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "tlssecretsnippet" -}}
tls:
- hosts:
  - {{ include "domainName" (dict "node" .node) }}
  secretName: {{ .node.ingress.tls.secretName }}
{{- end -}}

{{- define "domainName" }}
  {{- if .node.ingress -}}
  {{- if .node.ingress.domain -}}
    {{ .node.ingress.domain }}
  {{- end -}}
  {{- else if .node.service.domain -}}
    {{ .node.service.domain }}
  {{- end -}}
{{- end }}

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

{{/*
Autopilot secret name - uses pre-existing secret or auto-generated one
*/}}
{{- define "autopilot.credentialsSecretName" -}}
{{- if .Values.providerCredentialsSecret -}}
{{- .Values.providerCredentialsSecret -}}
{{- else -}}
{{- $depName := include "deploymentName" (dict "root" .Values "defaultname" "autopilot") -}}
{{- printf "%s-provider-credentials" $depName -}}
{{- end -}}
{{- end -}}

{{/*
Check if any inline credentials are provided
*/}}
{{- define "autopilot.hasInlineCredentials" -}}
{{- if or .Values.azure.endpoint .Values.azure.apiKey .Values.aws.accessKeyId .Values.vertex.credentials .Values.vertex.applicationCredentialsFile -}}
true
{{- end -}}
{{- end -}}

{{/*
Check if any models config is provided (inline, existing configmap, or use default bundled file)
*/}}
{{- define "autopilot.hasCustomModels" -}}
{{- if and .Values.customModels (or .Values.customModels.existingConfigMap .Values.customModels.inline) -}}
true
{{- else if .Values.deployModelsConfigMap -}}
true
{{- end -}}
{{- end -}}

{{/*
Check if we need to create a ConfigMap (inline or default file, but NOT existing configmap)
*/}}
{{- define "autopilot.hasModelsConfig" -}}
{{- if and .Values.customModels .Values.customModels.inline -}}
true
{{- else if .Values.deployModelsConfigMap -}}
true
{{- end -}}
{{- end -}}

{{/*
Resolve the models ConfigMap name
*/}}
{{- define "autopilot.modelsConfigMapName" -}}
{{- if and .Values.customModels .Values.customModels.existingConfigMap -}}
{{- .Values.customModels.existingConfigMap -}}
{{- else -}}
{{- $depName := include "deploymentName" (dict "root" .Values "defaultname" "autopilot") -}}
{{- printf "%s-models" $depName -}}
{{- end -}}
{{- end -}}
