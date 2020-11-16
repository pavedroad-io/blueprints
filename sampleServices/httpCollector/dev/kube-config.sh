
#!/bin/bash

CMD0="microk8s.config"
CMD1=`microk8s.config > $HOME/.kube/config`

clear
echo "Writing microk8s.config to $HOME/.kube/config"
echo "microk8s.config > $HOME/.kube/config"
$CMD1
