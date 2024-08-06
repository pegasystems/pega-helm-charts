{{- /*
imagePullSecret
backingservicesRegistrySecret
deploymentName
tlssecretsnippet
backingservices.gke.backendConfig
podAffinity
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


{{- define  "backingservices.gke.backendConfig" -}}
apiVersion: cloud.google.com/v1
kind: BackendConfig
metadata:
  name: {{ .name }}
spec:
  timeoutSec: 40
  connectionDraining:
    drainingTimeoutSec: 60
  healthCheck:
    checkIntervalSec: 5
    healthyThreshold: 1
    port: {{ .healthCheckPort }}
    requestPath: {{ .requestPath }}
    timeoutSec: 5
    type: HTTP
    unhealthyThreshold: 2
---
{{ end }}

{{- define "podAffinity" }}
{{- if .affinity }}
affinity:
{{- toYaml .affinity | nindent 2 }}
{{- end }}
{{ end }}