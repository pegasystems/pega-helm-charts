{{- define "searchName" -}}
pega-search
{{- end -}}
{{- define "isExternalSearch" }}
  {{- if or (and (.Values.externalURL) (ne .Values.externalURL "http://pega-search")) .Values.externalSearchService -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }} 