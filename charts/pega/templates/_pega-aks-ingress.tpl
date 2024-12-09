{{- define "pega.aks.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
{{ include "ingressApiVersion" . }}
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations:
{{- if not (.node.ingress.ingressClassName) }}
    # Ingress class used is 'azure/application-gateway'
    kubernetes.io/ingress.class: azure/application-gateway
{{- end }}
    # Ingress annotations for aks
    appgw.ingress.kubernetes.io/cookie-based-affinity: "true"
{{- if .node.ingress.annotations }}
{{ toYaml .node.ingress.annotations | indent 4 }}
{{- end }}
{{ if ( include "ingressTlsEnabled" . ) }}
    # HTTP to HTTPS Redirect
    appgw.ingress.kubernetes.io/ssl-redirect: "true"
{{ end }}
{{- if ((.node.service).tls).enabled }}
    # TLS certificate used for the ingress
    appgw.ingress.kubernetes.io/backend-protocol: https
{{- end }}
spec:
{{- if .node.ingress.ingressClassName }}
  ingressClassName: {{ .node.ingress.ingressClassName }}
{{- end }}
{{ if ( include "ingressTlsEnabled" . ) }}
{{- if .node.ingress.tls.secretName }}
{{ include "tlssecretsnippet" . }}
{{ end }}
{{ end }}
  rules:
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, traefik http port also has to be provided in the URL.
  - host: {{ template "domainName" dict "node" .node }}
    http:
      paths:
# protocol will be set to https only when either ingress is enabled or domain is set
{{ include "defaultIngressRule" . | indent 6 }}
---     
{{- end }}
