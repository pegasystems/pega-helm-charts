{{- define "pegaTraefikConfigTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if and (.node.tlscertificates) (.node.tlscertificates.enabled) }}
kind: ServersTransport
apiVersion: traefik.containo.us/v1alpha1
metadata:
  name: {{ .name }}-servers-transport
  namespace: {{ .root.Release.Namespace }}
spec:
# set the below to true if the connection from traefik to the backend is to be encrypted but not validated using self-signed certificates
  insecureSkipVerify: false
  rootCAsSecrets:
    - {{ .depname }}-tomcat-certificates-secret
  serverName: pega.com
{{- end }}
{{- end}}