{{define "sonar-project.properties.tpl"}}
sonar.organization={{.Info.SonarCloudOrg}}
# must be unique in a given SonarQube instance
sonar.projectKey={{.Info.SonarCloudOrg}}_{{.Info.Name}}
# this is the name and version displayed in the SonarQube UI. Was mandatory prior to SonarQube 6.1.
sonar.projectName={{.Info.Name}}
sonar.projectVersion={{.Info.Version}}
sonar.go.coverage.reportPaths=artifacts/coverage.out
sonar.go.golint.reportPaths=artifacts/lint.out 
# Path is relative to the sonar-project.properties file. Replace "\" by "/" on Windows.
# This property is optional if sonar.modules is set. 
sonar.sources=.
sonar.host.url=https://sonarcloud.io
sonar.login={{SonarLogin}}

# Encoding of the source code. Default is default system encoding
#sonar.sourceEncoding=UTF-8
{{end}}
