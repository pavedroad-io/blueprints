{{define "login-stagging-repo.tpl"}}
#!/bin/bash
aws ecr get-login-password | docker login --username AWS --password-stdin 400276217548.dkr.ecr.us-west-1.amazonaws.com
{{end}}
