{{- define "pegaRegistryCredentialsSecretTemplate" }}
kind: Secret
apiVersion: v1
metadata:
  name: {{ template "pegaRegistrySecret" }}
  namespace: {{ .Release.Namespace }}
data:
  .dockerconfigjson: {{ template "imagePullSecret" . }}
type: kubernetes.io/dockerconfigjson
{{- end }}