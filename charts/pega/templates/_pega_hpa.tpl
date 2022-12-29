{{- define "pega.hpa" -}}
{{- if .hpa.enabled -}}
# The Horizontal Pod Autoscaler for {{ .deploymentName }} deployment
{{- if(semverCompare ">= 1.23.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
apiVersion: autoscaling/v2
{{- else }}
apiVersion: autoscaling/v2beta2
{{- end }}
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
  {{- if (hasKey .hpa "enableCpuTarget" | ternary .hpa.enableCpuTarget true) }}
  - type: Resource
    resource:
      name: cpu
      target:
        {{- if .hpa.targetAverageCPUValue }} 
        type: Value
        averageValue: {{ .hpa.targetAverageCPUValue }}
        {{- else if .hpa.targetAverageCPUUtilization }}
        type: Utilization
        averageUtilization: {{ .hpa.targetAverageCPUUtilization }}
        {{- else }}
        type: Value
        averageValue: 2.55
        {{- end }}
  {{- end }}
  {{- if (hasKey .hpa "enableMemoryTarget" | ternary .hpa.enableMemoryTarget false) }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        {{- if .hpa.targetAverageMemoryUtilization }}
        averageUtilization: {{ .hpa.targetAverageMemoryUtilization }}
        {{- else }}
        averageUtilization: 85
        {{- end }}
  {{- end }}
  
---
{{- end -}}
{{- end -}}
