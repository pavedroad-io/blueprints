{{define "tilt.tpl"}}
k8s_yaml(kustomize('manifests/kubernetes/dev'))

docker_build('localhost:32000/{{.Info.GitHubOrg}}/{{.Info.Name | ToLower}}', '.', dockerfile='manifests/Dockerfile')
docker_build('localhost:32000/{{.Info.GitHubOrg}}/{{.Info.Name | ToLower}}initdb', '.', dockerfile='manifests/InitDbDockerFile')
{{end}}
