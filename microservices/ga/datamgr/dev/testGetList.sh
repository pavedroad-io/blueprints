{{define "testGetList.sh"}}
#!/bin/bash

curl -H "Content-Type: application/json" \
     -v http://localhost:8083/api/v1/namespace/pavedroad.io/{{.Name}}LIST | jq '.'
{{end}}


