{{- define "pegaTraefikConfigTemplate" }}
# Secret used for tls certificates to configure https to tomcat
{{- if ((.node.service).tls).enabled }}
{{- if (.node.service.tls.traefik).enabled }}
kind: ServersTransport
apiVersion: traefik.containo.us/v1alpha1
metadata:
  name: {{ .name }}-servers-transport
  namespace: {{ .root.Release.Namespace }}
spec:
# set the below to true if the connection from traefik to the backend is to be encrypted but not validated using self-signed certificates
{{- if .node.service.tls.traefik.insecureSkipVerify }}
  insecureSkipVerify: true
{{- else }}
  insecureSkipVerify: false
{{- end }}
#For traefik, it expects the root CA certificate in a secret under the field ca.crt
  rootCAsSecrets:
  {{- if .node.service.tls.external_secret_name }}
    - {{ .node.service.tls.external_secret_name }}
  {{- else }}
    - {{ .depname }}-tomcat-keystore-secret
  {{- end }}
{{- if .node.service.tls.traefik.serverName }}
  serverName: {{ .node.service.tls.traefik.serverName -}}
{{- else }}
  serverName: www.pega.com
{{- end }}
{{- end }}
{{- end }}
{{- end}}