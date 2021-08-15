{{ define "Namespace.tpl"}}
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Project.Kubernetes.Namespace}}
{{end}}
