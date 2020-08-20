{{- define "pega.gke.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
apiVersion: extensions/v1beta1
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
    {{ toYaml .node.ingress.tls.ssl_annotation }}
  {{ end }}
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
  backend:
    serviceName: {{ .name }}
    servicePort: {{ .node.service.port }}
  rules:
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, http/https port also has to be provided in the URL.
  - host: {{ template "domainName" dict "node" .node }}
    http:
      paths: 
      {{ if and .root.Values.constellation (eq .root.Values.constellation.enabled true) }}
      - path: /c11n     
        backend:
          serviceName: constellation
          servicePort: 3000
      {{ end }}
      - backend: 
          serviceName: {{ .name }} 
          servicePort: {{ .node.service.port }}
---
{{- end }}