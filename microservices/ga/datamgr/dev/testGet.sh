{{define "testGet.sh"}}
#!/bin/bash

curl -H "Content-Type: application/json" \
     -v http://localhost:8083/api/v1/namespace/pavedroad.io/{{.Name}}/fool | jq '.'
{{end}}
