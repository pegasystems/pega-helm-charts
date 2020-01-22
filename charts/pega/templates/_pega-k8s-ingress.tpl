{{- define "pega.k8s.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations:
{{- if .node.ingress }}
{{- if .node.ingress.annotations }}
    # Custom annotations
{{ toYaml .node.ingress.annotations | indent 4 }}
{{- end }}
{{- else }}
    # Ingress class used is 'traefik'
    kubernetes.io/ingress.class: {{ include "ingressClass" . }}
    {{- if (eq .root.Values.global.provider "aks") }}
    # Ingress annotations for aks
    appgw.ingress.kubernetes.io/cookie-based-affinity: "true"
    {{ end }}
{{- end }}
spec:
  rules:
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, traefik http port also has to be provided in the URL.
  - host: {{ .node.service.domain }}
    http:
      paths:
      - backend:
          serviceName: {{ .name }}
          servicePort: {{ .node.service.port }}
---     
{{- end }}
