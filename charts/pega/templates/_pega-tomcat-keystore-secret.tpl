{{- define "pegaTomcatKeystoreSecretTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if ((.node.service).tls).enabled }}
kind: Secret
apiVersion: v1
type: Opaque
metadata:
  name: {{ .name }}-tomcat-keystore-secret
  namespace: {{ .root.Release.Namespace }}
stringData:
  # cert Files
  # Base64 encoded password for enabling TLS in tomcat
{{- if .node.service.tls.keystorepassword }}
  TOMCAT_KEYSTORE_PASSWORD: {{ .node.service.tls.keystorepassword | quote}}
{{- else }}
  TOMCAT_KEYSTORE_PASSWORD: "123456"
{{- end }}
{{- if .node.service.tls.keystore }}
  TOMCAT_KEYSTORE_CONTENT: {{ .node.service.tls.keystore | quote -}}
{{- else }}
  TOMCAT_KEYSTORE_CONTENT: {{ .root.Files.Get "config/certs/pegakeystore.jks" | b64enc | indent 2 }}
{{- end }}
# this field is used for traefik, it expects the root CA certificate in a secret under the field ca.crt
{{- if .node.service.tls.cacertificate }}
  ca.crt: {{ .node.service.tls.cacertificate | b64dec | quote -}}
{{- else }}
  ca.crt: {{ .root.Files.Get "config/certs/pegaca.crt" | quote | indent 2 }}
{{- end }}
{{- end }}
{{- end }}
