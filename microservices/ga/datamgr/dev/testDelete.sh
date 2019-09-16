{{define "testDelete.sh"}}
#!/bin/bash

curl -X DELETE \
     -H "Content-Type: application/json" \
     -v http://localhost:8083/api/v1/namespace/pavedroad.io/{{.Name}}/1504d9f7-d791-4342-8fc2-10618a44a749

{{end}}

