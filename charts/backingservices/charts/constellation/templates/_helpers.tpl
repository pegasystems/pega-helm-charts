{{- define "imagePullSecret" }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.docker.registry.url (printf "%s:%s" .Values.docker.registry.username .Values.docker.registry.password | b64enc) | b64enc }}
{{- end }}

{{- define "backingservicesRegistrySecret" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-registry-secret
{{- end }}

{{- define "deploymentName" }}{{ $deploymentNamePrefix := "constellation" }}{{ if (.Values.deployment) }}{{ if (.Values.deployment.name) }}{{ $deploymentNamePrefix = .Values.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "tlssecretsnippet" }}
  tls:
   - hosts:
     - {{ template "domainName" dict "node" .node }}
     secretName: {{ .node.ingress.tls.secretName }}
{{- end }}

{{- define "domainName" }}
  {{- if .node.ingress -}}
  {{- if .node.ingress.domain -}}
    {{ .node.ingress.domain }}
  {{- end -}}
  {{- else if .node.service.domain -}}
    {{ .node.service.domain }}
  {{- end -}}
{{- end }}