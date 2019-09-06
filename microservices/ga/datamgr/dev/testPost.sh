{{define "testPost.sh"}}
#!/bin/bash

curl -H "Content-Type: application/json" \
     -X POST \
     -d @{{.Name}}.json \
     -v http://localhost:8083/api/v1/namespace/pavedroad.io/{{.Name}}
{{end}}
