{{- define  "cdn.service" -}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .name }}
  labels:
    app: {{ .name }}
spec:
  type: {{ .root.Values.serviceType }}
  selector:
    app: {{ .name }}
  ports:
    - protocol: TCP
      port: {{ .root.Values.pegaStaticPort }}
      targetPort: {{ .root.Values.pegaStaticTargetPort }}
{{ end }}