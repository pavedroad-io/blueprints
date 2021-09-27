{{ define "manifests/kubernetes/stag/db/disk1.sh"}}
aws ec2 create-volume \
  --region us-west-1 \
  --availability-zone us-west-1b \
  --size 10 \
  --volume-type gp3 \
  --tag-specifications 'ResourceType=volume,Tags=[{Key=env,Value=stagging},{Key=Name,Value=staging-db}]' 
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
