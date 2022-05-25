{{- define "pegaTomcatKeystoreSecretTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if .node.tlscertificates }}
{{- if .node.tlscertificates.enabled }}
kind: Secret
apiVersion: v1
type: Opaque
metadata:
  name: {{ .name }}-tomcat-keystore-secret
  namespace: {{ .root.Release.Namespace }}
stringData:
  # cert Files
  # Base64 encoded password for enabling TLS in tomcat
{{- if .node.tlscertificates.keystorepassword }}
  CERT_PASSWORD: {{ .node.tlscertificates.keystorepassword | quote}}
{{- else }}
  CERT_PASSWORD: "123456"
{{- end }}
{{- if .node.tlscertificates.keystore }}
  CERT_CONTENT: {{ .node.tlscertificates.keystore | quote -}}
{{- else }}
  CERT_CONTENT: {{ .root.Files.Get "config/certs/pegakeystore.jks" | b64enc | indent 2 }}
{{- end }}
# this field is used for traefik, it expects the root CA certificate in a secret under the field ca.crt
{{- if .node.tlscertificates.cacertificate }}
  ca.crt: {{ .node.tlscertificates.cacertificate | b64dec | quote -}}
{{- else }}
  ca.crt: {{ .root.Files.Get "config/certs/pegaca.crt" | quote | indent 2 }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
