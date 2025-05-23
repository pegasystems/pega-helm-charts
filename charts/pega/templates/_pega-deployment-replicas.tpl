{{- define  "pega.deployment.replicas" -}}
{{- if or (not .node.hpa) (eq .node.hpa.enabled false) }}
# Replicas specify the number of copies for {{ .node.name }} when HPA is not in use
replicas: {{ .node.replicas }}
{{- end }}
{{- end }}