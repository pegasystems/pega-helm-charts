{{- define  "pega.service" -}}
# Service instance for {{ .name }}
kind: Service
apiVersion: v1
metadata:
  # Name of the service for
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
{{- if .node.service.annotations }}
  annotations: 
    # Custom annotations
{{ toYaml .node.service.annotations | indent 4 }}
{{- else }}
  {{- if (eq .root.Values.global.provider "k8s") }}
  annotations:
    # Override the default wrr load balancer algorithm.
    traefik.ingress.kubernetes.io/load-balancer-method: drr
    # Sets the maximum number of simultaneous connections to the backend
    # Must be used in conjunction with the label below to take effect
    traefik.ingress.kubernetes.io/max-conn-amount: '10'
    # Enable session affinity
    traefik.ingress.kubernetes.io/service.sticky.cookie: "true"
    traefik.ingress.kubernetes.io/service.sticky.cookie.name: UNIQUE-PEGA-COOKIE-NAME
{{- if ((.node.service).tls).enabled }}
{{- if (.node.service.tls.traefik).enabled }}
    # Sets serversTreansport that has config in order to verify rootCA
    traefik.ingress.kubernetes.io/service.serverstransport: {{ .root.Release.Namespace }}-{{ .name }}-servers-transport@kubernetescrd
{{- end }}
{{- end }}
  {{- else if (eq .root.Values.global.provider "gke") }}
  annotations:
    cloud.google.com/neg: '{"ingress": true}'
    cloud.google.com/app-protocols: '{"https":"HTTPS","http":"HTTP"}'
    {{ if (semverCompare "< 1.16.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}beta.{{ end -}}cloud.google.com/backend-config: '{"ports": {"{{ .node.service.port }}": "{{ .name }}"}}'
  {{ end }}
{{- end }}
spec:
  type:
  {{- if (.node.service.serviceType) -}}
  {{ indent 1 (.node.service.serviceType) }}
  {{- else if or (eq .root.Values.global.provider "gke") (eq .root.Values.global.provider "eks") -}}
  {{ indent 1 "NodePort" }}
  {{- else -}}
  {{ indent 1 "LoadBalancer" }}
  {{- end }}
  {{- if and ( and (.node.service.serviceType) (eq (toString .node.service.serviceType) "LoadBalancer")) (.node.service.loadBalancerSourceRanges) }}
  loadBalancerSourceRanges:
  {{- range .node.service.loadBalancerSourceRanges }}
    - "{{ . }}"
  {{- end }}
  {{- end }}
  # Specification of on which port the service is enabled
  ports:
{{- if or (not (hasKey .node.service "httpEnabled")) (.node.service.httpEnabled) }}
  - name: http
    port: {{ .node.service.port }}
    targetPort: {{ .node.service.targetPort }}
{{- end }}
{{- if (.node.service.tls).enabled }}
  - name: https
    port: {{ .node.service.tls.port }}
    targetPort: {{ .node.service.tls.targetPort }}
{{- end }}
{{- if .node.service.customServicePorts }}
{{ toYaml .node.service.customServicePorts | indent 2 }}
{{- end }}
  selector:
    app: {{ .name }}
---
{{- end -}}
