
#!/bin/bash

K8SCONFIG="microk8s.config"
K8SCONTEXT="microk8s"

clear
echo "Saving ${K8SCONFIG} to ${HOME}/.kube/config"
echo "gathering microk8s configuration "

clusterName=`${K8SCONFIG} | yq r - clusters[0].name`
clusterServer=`${K8SCONFIG} | yq r - clusters[0].cluster.server`
clusterCertificate=`${K8SCONFIG} | yq r - clusters[0].cluster.certificate-authority-data`
userName=`microk8s.config | yq r - users[0].name`
userToken=`${K8SCONFIG} | yq r - users[0].user.token`

# Add/update the server and the user
echo "Saving / updating ${K8SCONTEXT}"
kubectl config set-cluster ${clusterName} --server=${clusterServer}
kubectl config set clusters.${clusterName}.certificate-authority-data ${clusterCertificate}

kubectl config set-credentials ${userName}
kubectl config set users.${userName}.token ${userToken}

# Add context and use it
kubectl config set-context microk8s --cluster=${clusterName} --user=${userName}
kubectl config set current-context microk8s

