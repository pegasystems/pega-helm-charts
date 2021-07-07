{{- define "pega.pdb" -}}
{{- if .pdb.enabled -}}
# The Pod Disruption Budget for {{ .name }} deployment
apiVersion: policy/v1beta1
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
