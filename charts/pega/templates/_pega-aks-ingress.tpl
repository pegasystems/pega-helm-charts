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
