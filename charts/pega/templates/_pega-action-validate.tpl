{{- define  "pega.actionvalidate" -}}
{{- $validActions := list "install" "deploy" "install-deploy" "upgrade" "upgrade-deploy" }}
{{- if not (has .root.Values.global.actions.execute $validActions) }}
{{- fail "Action value is not correct. The valid values are 'install' 'deploy' 'install-deploy' 'upgrade' 'upgrade-deploy'" }}
{{- end }}
{{- end }}