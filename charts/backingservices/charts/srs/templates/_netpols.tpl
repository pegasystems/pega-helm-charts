{{/*
Network policy: kube-dns
*/}}
{{- define "srs.netpol.kube-dns" -}}
- namespaceSelector:
    matchLabels:
      name: kube-system
- podSelector:
    matchExpressions:
      - key: k8s-app
        operator: In
        values: ["kube-dns", "coredns"]
ports:
- protocol: TCP
  port: 53
- protocol: TCP
  port: 1053
- protocol: TCP
  port: 80
- protocol: TCP
  port: 8080
{{- end -}}