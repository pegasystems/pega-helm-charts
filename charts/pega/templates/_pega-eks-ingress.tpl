{{- define "pega.eks.ingress" -}}
# Ingress to be used for {{ .name }}
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations:
    # Ingress class used is 'alb'
    kubernetes.io/ingress.class: alb
{{ if (.node.service.domain) }}
    {{ template "eksHttpsAnnotationSnippet" }}
{{ else }}
{{- if (.node.ingress) }}
{{- if (.node.ingress.tls) }}
{{- if (eq .node.ingress.tls.enabled true) }}
    {{ template "eksHttpsAnnotationSnippet" }}
{{ if (.node.ingress.tls.ssl_annotation) }}
{{ toYaml .node.ingress.tls.ssl_annotation | indent 4 }}
{{ end }}
{{ end }}
{{ end }}
{{ else }}
    # specifies the ports that ALB used to listen on
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}]'
{{ end }}
{{ end }}
{{- if .node.ingress.annotations }}
{{ toYaml .node.ingress.annotations | indent 4 }}
{{- else }}
    # override the default scheme internal as ALB should be internet-facing
    alb.ingress.kubernetes.io/scheme: internet-facing
    # set to ip mode to route traffic directly to the pods ip
    alb.ingress.kubernetes.io/target-type: ip
{{- end }}
{{- if not (and (.node.ingress.annotations) ( .node.ingress.annotations | quote | regexFind "alb.ingress.kubernetes.io/target-group-attributes:" ) ) }}
    # enable sticky sessions on target group and set alb routing algorithm
    alb.ingress.kubernetes.io/target-group-attributes: load_balancing.algorithm.type=least_outstanding_requests,stickiness.enabled=true,stickiness.lb_cookie.duration_seconds={{ include "lbSessionCookieStickiness" . }}
{{- end }}
spec:
  rules:
  {{ if (.node.service.domain) }}
  - http:
      paths:
      - backend:
          serviceName: ssl-redirect
          servicePort: use-annotation
  {{ else }}
  {{ if ( include "ingressTlsEnabled" . ) }}
  - http:
      paths:
      - backend:
          serviceName: ssl-redirect
          servicePort: use-annotation
  {{ end }}
  {{ end }}
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, alb http port also has to be provided in the URL.
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
