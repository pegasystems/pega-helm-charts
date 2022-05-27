{{- define "pegaTraefikConfigTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if (.node.tls).enabled }}
{{- if .node.tls.traefik.enabled }}
kind: ServersTransport
apiVersion: traefik.containo.us/v1alpha1
metadata:
  name: {{ .name }}-servers-transport
  namespace: {{ .root.Release.Namespace }}
spec:
# set the below to true if the connection from traefik to the backend is to be encrypted but not validated using self-signed certificates
{{- if .node.tls.traefik.insecureSkipVerify }}
  insecureSkipVerify: true
{{- else }}
  insecureSkipVerify: false
{{- end }}
#For traefik, it expects the root CA certificate in a secret under the field ca.crt
  rootCAsSecrets:
    - {{ .depname }}-tomcat-keystore-secret
{{- if .node.tls.traefik.serverName }}
  serverName: {{ .node.tls.traefik.serverName -}}
{{- else }}
  serverName: www.pega.com
{{- end }}
{{- end }}
{{- end }}
{{- end}}