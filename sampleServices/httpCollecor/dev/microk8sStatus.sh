
#!/bin/bash

STATUS=`microk8s.status | sed  's/.*\(\bnot running\b\).*/\1/'`

if [ "$STATUS" = "not running" ]
then
	echo "down"
else
	echo "up"
fi
