{{define "dev/apifix.sh"}}#!/bin/bash
  
fixrequired=`cat docs/api.json | jq '.definitions' `

if [[ "$fixrequired" == "null" ]]
then
        echo "Adding definitions to docs/api.json"
        cat docs/api.json | jq '. + {definitions: {}}' > docs/fix.json
        mv docs/fix.json docs/api.json
fi
{{end}}
