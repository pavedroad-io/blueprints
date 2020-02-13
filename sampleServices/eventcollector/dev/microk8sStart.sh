
#!/bin/bash

STATUS=`microk8s.status | sed  's/.*\(\bnot running\b\).*/\1/'`

if [ "$STATUS" = "not running" ]
then
	echo "microk8s is not running"
	echo "starting microk8s"
	microk8s.start
	echo "waiting for ready "
	microk8s.status --wait-ready
	echo "updating $HOME/.kube/config"
	microk8s.config > $HOME/.kube/config
	echo "enabling required services"
	microk8s.enable dns registry
else
	echo "microk8s is already running"

fi
