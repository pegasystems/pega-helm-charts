{{- define "pega.gke.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
spec:
  backend:
    serviceName: {{ .name }}
    servicePort: {{ .node.service.port }}
  rules:
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, http/https port also has to be provided in the URL.
  - host: {{ .node.service.domain }}
    http:
      paths:
      - backend:
          serviceName: {{ .name }}
          servicePort: {{ .node.service.port }}
---
{{- end }}