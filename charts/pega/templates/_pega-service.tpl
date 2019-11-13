{{- define  "pega.service" -}}
# Service instance for {{ .name }}
kind: Service
apiVersion: v1
metadata:
  # Name of the service for
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  {{- if and (ne .root.Values.global.provider "eks") (ne .root.Values.global.provider "openshift") }}
  annotations: 
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
spec:
  type: LoadBalancer
  # Specification of on which port the service is enabled
  ports:
  - name: http
    port: {{ .port }}
    targetPort: {{ .targetPort }}
  selector:
    app: {{ .name }}
---
{{- end -}}
