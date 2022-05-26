{{- define "pegaTomcatKeystoreSecretTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if .node.tls }}
{{- if .node.tls.enabled }}
kind: Secret
apiVersion: v1
type: Opaque
metadata:
  name: {{ .name }}-tomcat-keystore-secret
  namespace: {{ .root.Release.Namespace }}
stringData:
  # cert Files
  # Base64 encoded password for enabling TLS in tomcat
{{- if .node.tls.keystorepassword }}
  TOMCAT_KEYSTORE_PASSWORD: {{ .node.tls.keystorepassword | quote}}
{{- else }}
  TOMCAT_KEYSTORE_PASSWORD: "123456"
{{- end }}
{{- if .node.tls.keystore }}
  TOMCAT_KEYSTORE_CONTENT: {{ .node.tls.keystore | quote -}}
{{- else }}
  TOMCAT_KEYSTORE_CONTENT: {{ .root.Files.Get "config/certs/pegakeystore.jks" | b64enc | indent 2 }}
{{- end }}
# this field is used for traefik, it expects the root CA certificate in a secret under the field ca.crt
{{- if .node.tls.cacertificate }}
  ca.crt: {{ .node.tls.cacertificate | b64dec | quote -}}
{{- else }}
  ca.crt: {{ .root.Files.Get "config/certs/pegaca.crt" | quote | indent 2 }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
