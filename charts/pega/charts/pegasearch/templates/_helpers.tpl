{{- define "searchName" -}}
pega-search
{{- end -}}
{{- define "isExternalSearch" }}
  {{- if and (.Values.externalURL) (ne .Values.externalURL "http://pega-search") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }} 