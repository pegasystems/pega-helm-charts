{{- define "pegaRegistryCredentialsSecretTemplate" }}
kind: Secret
apiVersion: v1
metadata:
  name: {{ template "pegaRegistrySecret" $ }}
  namespace: {{ .Release.Namespace }}
  annotations:
    "helm.sh/hook": pre-install, pre-upgrade
    "helm.sh/hook-weight": "0"
    "helm.sh/hook-delete-policy": before-hook-creation
data:
  .dockerconfigjson: {{ template "imagePullSecret" . }}
type: kubernetes.io/dockerconfigjson
{{- end }}