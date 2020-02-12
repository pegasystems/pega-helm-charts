{{- define "pega.k8s.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations:
    # Ingress class used is 'traefik'
    kubernetes.io/ingress.class: traefik
spec:
  rules:
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, traefik http port also has to be provided in the URL.
  - host: {{ .node.service.domain }}
    http:
      paths: 
      {{ if and .root.Values.constellation (eq .root.Values.constellation.enabled true)}}
      - path: /prweb/constellation     
        backend:
          serviceName: constellation
          servicePort: 3000
      {{ end }}
      - backend: 
          serviceName: {{ .name }} 
          servicePort: {{ .node.service.port }}
---     
{{- end }}
