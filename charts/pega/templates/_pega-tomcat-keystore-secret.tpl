{{- define "pegaTomcatKeystoreSecretTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if ((.node.service).tls).enabled }}
kind: Secret
apiVersion: v1
type: Opaque
metadata:
  name: {{ .name }}-tomcat-keystore-secret
  namespace: {{ .root.Release.Namespace }}
data:
# supports either keystore and password combo or certificate, chain and private key files in PEM format
{{- if and .node.service.tls.certificateFile .node.service.tls.certificateKeyFile }}
  ca.crt: {{ .node.service.tls.cacertificate  | indent 4 }}
  TOMCAT_CERTIFICATE_FILE: {{ .node.service.tls.certificateFile  | indent 4 }}
  TOMCAT_CERTIFICATE_KEY_FILE: {{ .node.service.tls.certificateKeyFile  | indent 4 }}
  TOMCAT_CERTIFICATE_CHAIN_FILE: {{ .node.service.tls.cacertificate  | quote | indent 4 -}}
{{- else }}
  # cert Files
  # Base64 encoded password for enabling TLS in tomcat
{{- if .node.service.tls.keystorepassword }}
  TOMCAT_KEYSTORE_PASSWORD: {{ .node.service.tls.keystorepassword | b64enc | quote}}
{{- else }}
  TOMCAT_KEYSTORE_PASSWORD: {{ "123456" | b64enc }}
{{- end }}
{{- if .node.service.tls.keystore }}
  TOMCAT_KEYSTORE_CONTENT: {{ .node.service.tls.keystore | quote -}}
{{- else }}
  TOMCAT_KEYSTORE_CONTENT: {{ .root.Files.Get "config/certs/pegakeystore.jks"  | b64enc | indent 2 }}
{{- end }}
# this field is used for traefik, it expects the root CA certificate in a secret under the field ca.crt
{{- if .node.service.tls.cacertificate }}
  ca.crt: {{ .node.service.tls.cacertificate  | quote -}}
{{- else }}
  ca.crt: {{ .root.Files.Get "config/certs/pegaca.crt" | b64enc | quote | indent 2 }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
