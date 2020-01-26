
#!/bin/bash

STATUS=`microk8s.status | head -1 | sed  's/.*\(\bis running\b\).*/\1/'`

if [ "$STATUS" = "is running" ]
then
	echo "stopping microk8s"
	microk8s.stop
else
	echo "microk8s is not running"
fi
