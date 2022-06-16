{{- define "pega.gke.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
{{ include "ingressApiVersion" . }}
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
{{ if (.node.ingress) }}
{{ if (.node.ingress.tls) }}
{{ if (eq .node.ingress.tls.enabled true) }}
  annotations:
    kubernetes.io/ingress.allow-http: "false"
{{ if (.node.ingress.tls.useManagedCertificate) }}
    networking.gke.io/managed-certificates: managed-certificate-{{ .node.name }}
{{ end }}
{{ if (.node.ingress.tls.ssl_annotation) }}
{{ toYaml .node.ingress.tls.ssl_annotation | indent 4 }}
{{ end }}
{{- if .node.ingress.annotations }}
{{ toYaml .node.ingress.annotations | indent 4 }}
{{- end }}
{{ end }}
{{ end }}
{{ end }}
spec:
{{ if (.node.ingress) }}
{{ if (.node.ingress.tls) }}
{{ if (eq .node.ingress.tls.enabled true) }}
{{ if .node.ingress.tls.secretName }}
{{ include "tlssecretsnippet" . }}
{{ end }}
{{ end }}
{{ end }}
{{ end }}
{{- if (semverCompare ">= 1.22.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
  defaultBackend:
{{ else }}
  backend:
{{ end }}
{{ include "ingressService" . | indent 4 }}
  rules:
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, http/https port also has to be provided in the URL.
  - host: {{ template "domainName" dict "node" .node }}
    http:
      paths: 
      {{ if and .root.Values.constellation (eq .root.Values.constellation.enabled true) }}
      - path: /c11n     
        pathType: ImplementationSpecific
        backend:
{{ include "ingressServiceC11n" . | indent 10 }}
      {{ end }}
      - pathType: ImplementationSpecific
        backend: 
# protocol will be set to https only when either ingress is enabled or domain is set
{{ include "ingressBackend" . }}
---
{{- end }}
