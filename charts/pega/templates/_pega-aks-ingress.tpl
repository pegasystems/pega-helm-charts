{{- define "pega.aks.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations:
    # Ingress class used is 'azure/application-gateway'
    kubernetes.io/ingress.class: azure/application-gateway
    # Ingress annotations for aks
    appgw.ingress.kubernetes.io/cookie-based-affinity: "true"
{{- if .node.ingress.annotations }}
{{ toYaml .node.ingress.annotations | indent 4 }}
{{- end }}
{{ if ( include "ingressTlsEnabled" . ) }}
    # HTTP to HTTPS Redirect
    appgw.ingress.kubernetes.io/ssl-redirect: "true"
{{ end }}
spec:
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
      - backend:
          serviceName: {{ .name }}
          servicePort: {{ .node.service.port }}
---     
{{- end }}
