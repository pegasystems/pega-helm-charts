{{- define  "pega.service" -}}
# Service instance for {{ .name }}
kind: Service
apiVersion: v1
metadata:
  # Name of the service for
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations: 
{{- if .node.service.annotations }}
    # Custom annotations
{{ toYaml .node.service.annotations | indent 4 }}
{{- else }}
  {{- if and (ne .root.Values.global.provider "eks") (ne .root.Values.global.provider "openshift") }}
    # Enable backend sticky sessions
    traefik.ingress.kubernetes.io/affinity: 'true'
    # Override the default wrr load balancer algorithm.
    traefik.ingress.kubernetes.io/load-balancer-method: drr
    # Sets the maximum number of simultaneous connections to the backend
    # Must be used in conjunction with the label below to take effect
    traefik.ingress.kubernetes.io/max-conn-amount: '10'
    # Manually set the cookie name for sticky sessions
    traefik.ingress.kubernetes.io/session-cookie-name: UNIQUE-PEGA-COOKIE-NAME
  {{ end }}
{{- end }}
spec:
  type: {{ .node.service.serviceType | default "LoadBalancer" }}
  # Specification of on which port the service is enabled
  ports:
  - name: http
    port: {{ .node.service.port }}
    targetPort: {{ .node.service.targetPort }}
  selector:
    app: {{ .name }}
---
{{- end -}}
