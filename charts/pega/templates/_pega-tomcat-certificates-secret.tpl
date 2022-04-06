{{- define "pegaTomcatCertificatesSecretTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if and (.node.tlscertificates) (.node.tlscertificates.enabled) }}
kind: Secret
apiVersion: v1
type: Opaque
metadata:
  name: {{ .name }}-tomcat-certificates-secret
  namespace: {{ .root.Release.Namespace }}
stringData:
  # cert Files
{{- if and (.node.tlscertificates) (.node.tlscertificates.enabled) }}
  # Base64 encoded password for enabling TLS in tomcat
  CERT_PASSWORD: {{ .node.tlscertificates.tlspassword | quote}}
  CERT_CONTENT: {{ .node.tlscertificates.certificate -}}
{{- end }}
{{- end }}
{{- end}}