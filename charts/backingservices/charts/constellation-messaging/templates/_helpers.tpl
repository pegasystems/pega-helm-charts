{{- define "podAffinity" }}
{{- if .affinity }}
# Pod Affinity
affinity:
{{- toYaml .affinity | nindent 2 }}
{{- end }}
{{ end }}