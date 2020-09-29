{{- define "pega.hpa" -}}
{{- if .hpa.enabled -}}
# The Horizontal Pod Autoscaler for {{ .deploymentName }} deployment
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: {{ .name | quote}}
  namespace: {{ .root.Release.Namespace }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ .deploymentName | quote }}
  {{- if .hpa.minReplicas }}
  minReplicas: {{ .hpa.minReplicas }}
  {{- else }}
  minReplicas: 1
  {{- end }}
  {{- if .hpa.maxReplicas }}
  maxReplicas: {{ .hpa.maxReplicas }}
  {{- else }}
  maxReplicas: 5
  {{- end }} 
  metrics:
  - type: Resource
    resource:
      name: cpu
      {{- if .hpa.targetAverageCPUUtilization }}
      targetAverageUtilization: {{ .hpa.targetAverageCPUUtilization }}
      {{- else }}
      targetAverageUtilization: 70
      {{- end }}  
  - type: Resource
    resource:
      name: memory
      {{- if .hpa.targetAverageMemoryUtilization }}
      targetAverageUtilization: {{ .hpa.targetAverageMemoryUtilization }}
      {{- else }}
      targetAverageUtilization: 85
      {{- end }}
---
{{- end -}}
{{- end -}}
