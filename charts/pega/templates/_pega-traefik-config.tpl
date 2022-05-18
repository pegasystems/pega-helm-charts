{{- define "pegaTraefikConfigTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if .node.tlscertificates }}
{{- if .node.tlscertificates.enabled }}
{{- if .node.tlscertificates.traefik.enabled }}
kind: ServersTransport
apiVersion: traefik.containo.us/v1alpha1
metadata:
  name: {{ .name }}-servers-transport
  namespace: {{ .root.Release.Namespace }}
spec:
# set the below to true if the connection from traefik to the backend is to be encrypted but not validated using self-signed certificates
{{- if .node.tlscertificates.traefik.insecureSkipVerify }}
  insecureSkipVerify: true
{{- else }}
  insecureSkipVerify: false
{{- end }}
  rootCAsSecrets:
    - {{ .depname }}-tomcat-certificates-secret
  serverName: {{ .node.tlscertificates.traefik.serverName -}}
{{- if .node.tlscertificates.traefik.serverName }}
  serverName: {{ .node.tlscertificates.traefik.serverName -}}
{{- else }}
  serverName: www.pega.com
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end}}