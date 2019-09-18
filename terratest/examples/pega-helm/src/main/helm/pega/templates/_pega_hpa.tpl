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
    apiVersion: extensions/v1beta1
    kind: Deployment
    name: {{ .deploymentName | quote }}
  minReplicas: {{ .hpa.minReplicas }}
  maxReplicas: {{ .hpa.maxReplicas }}
  metrics:
  - type: Resource
    resource:
      name: cpu
      targetAverageUtilization: {{ .hpa.targetAverageCPUUtilization }}
  - type: Resource
    resource:
      name: memory
      targetAverageUtilization: {{ .hpa.targetAverageMemoryUtilization }}
{{- end -}}
{{- end -}}