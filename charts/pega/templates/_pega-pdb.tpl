{{- define "pega.pdb" -}}
{{- if .pdb.enabled -}}
# The Pod Disruption Budget for {{ .name }} deployment
{{- if (semverCompare ">= 1.21.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
apiVersion: policy/v1
{{- else }}
apiVersion: policy/v1beta1
{{- end }}
kind: PodDisruptionBudget
metadata:
  name: {{ .name }}-pdb
  namespace: {{ .root.Release.Namespace }}
spec:
  {{- if .pdb.minAvailable }}
  minAvailable: {{ .pdb.minAvailable }}
  {{- else if .pdb.maxUnavailable }}
  maxUnavailable: {{ .pdb.maxUnavailable }}
  {{- end }}
  selector:
    matchLabels:
      app: {{ .name }}
      
---
{{- end -}}
{{- end -}}
