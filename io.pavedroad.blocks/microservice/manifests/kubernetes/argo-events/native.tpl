{{ define "native.tpl" }}
apiVersion: argoproj.io/v1alpha1
kind: EventBus
metadata:
  name: default
spec:
  nats:
    native:
      # Optional, defaults to 3. If it is < 3, set it to 3, that is the minimal requirement.
{{end}}
      replicas: 3
      # Optional, authen strategy, "none" or "token", defaults to "none"
      auth: token
#      containerTemplate:
#        resources:
#          requests:
#            cpu: "10m"
#      metricsContainerTemplate:
#        resources:
#          requests:
#            cpu: "10m"
#      antiAffinity: false
#      persistence:
#        storageClassName: standard
#        accessMode: ReadWriteOnce
#        volumeSize: 10Gi
