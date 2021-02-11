{{- define  "pega.gke.backendConfig" -}}
apiVersion: cloud.google.com/v1
kind: BackendConfig
metadata:
  name: {{ .name }}
spec:
  timeoutSec: 40
  connectionDraining:
    drainingTimeoutSec: 60
  sessionAffinity:
    affinityType: "GENERATED_COOKIE"
    affinityCookieTtlSec: {{ template "lbSessionCookieStickiness" }}
  healthCheck:
    checkIntervalSec: 5
    healthyThreshold: 1
    port: 8080
    requestPath: /{{ template "pega.applicationContextPath" . }}/PRRestService/monitor/pingService/ping
    timeoutSec: 5
    type: HTTP
    unhealthyThreshold: 2
---
{{ end }}