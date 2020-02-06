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
    {{- if (eq .node.ingress.tls.enabled true) }}
    # specifies the ports that ALB used to listen on
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS": 443}]'
    # set the redirect action to redirect http traffic into https
    alb.ingress.kubernetes.io/actions.ssl-redirect: '{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}'
    {{- if (.node.ingress.tls.ssl_annotation) -}}
    {{ toYaml .node.ingress.tls.ssl_annotation }}
    {{- end }}
    {{- else }}
    # specifies the ports that ALB used to listen on
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}]'
    {{ end }}
    # override the default scheme internal as ALB should be internet-facing 
    alb.ingress.kubernetes.io/scheme: internet-facing
    # enable sticky sessions on target group
    alb.ingress.kubernetes.io/target-group-attributes: stickiness.enabled=true,stickiness.lb_cookie.duration_seconds={{ include "lbSessionCookieStickiness" . }}
    # set to ip mode to route traffic directly to the pods ip
    alb.ingress.kubernetes.io/target-type: ip
spec:
  rules:
  {{- if (eq .node.ingress.tls.enabled true) }}
  - http:
      paths:
      - backend:
          serviceName: ssl-redirect
          servicePort: use-annotation
  {{ end }}
  # The calls will be redirected from {{ .node.domain }} to below mentioned backend serviceName and servicePort.
  # To access the below service, along with {{ .node.domain }}, alb http port also has to be provided in the URL.
  - host: {{ .node.ingress.domain }}
    http:
      paths:
      - backend:
          serviceName: {{ .name }}
          servicePort: {{ .node.service.port }}
---
{{- end }}
