{{define "tilt.tpl"}}
k8s_yaml(kustomize('manifests/kubernetes/dev'))

docker_build('localhost:32000/{{.Info.Organization}}/{{.Info.Name}}', '.', dockerfile='manifests/Dockerfile')
docker_build('localhost:32000/{{.Info.Organization}}/{{.Info.Name}}initdb', '.', dockerfile='manifests/InitDbDockerFile')
{{end}}
