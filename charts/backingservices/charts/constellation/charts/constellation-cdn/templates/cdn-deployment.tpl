{{- define  "cdn.deployment" -}}
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: {{ .name }}
  labels:
    app: {{ .name }}
spec:
  replicas: {{ .node.replicas }}
  selector:
    matchLabels:
      app: {{ .name }}
  template:
    metadata:
      labels:
        app: {{ .name }}
    spec:
      imagePullSecrets:
      {{- range .root.Values.imagePullSecretNames }}
        - name: {{ . }}
      {{- end }}
      containers:
        - name: {{ .name }}
          imagePullPolicy: {{ .root.Values.imagePullPolicy }}
          image: {{ .node.image }}
          ports:
            - containerPort: {{ .root.Values.pegaStaticTargetPort }}
{{ end }}