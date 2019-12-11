{{define "docs/cheatsheet.md"}}# Kubernetes

## Basic constructs

                daemonsets <---> jobs <---> side cars 

       service <---> deployment <---> replicate controller <---> pod

                        persistentvolumeclaims

           ---------------- Execute on ------------------

          nodes <--->nodes <--->nodes <--->nodes <--->nodes

             ---------------- Storage ------------------

                         persistentvolumes
             
           ---------------- Managed by  ------------------

                   master <---> master <---> master

- [daemonsets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) run on every node providing support for applications like monitoring

- [jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/) are a pattern that provides a service to deployments

- [side cars](https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/) are a pattern that provides a service to deployments

-

- [Service](https://kubernetes.io/docs/concepts/services-networking/service/) exposes pods and managed load balancing

- [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) controls deploying pods and managing replication and health check.

- [Replicate controller](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/) manages scaling pods to the desired number specified

- [pods](https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/) Support running stateless conainers

-

- [Persistent Volume Claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) Allow a pod to request storage

-

- [nodes](https://kubernetes.io/docs/concepts/architecture/nodes/) are the servers workloads are scheduled on 

-

- [persistentvolumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) are statically or dynamically generated storage

-

- [master](https://kubernetes.io/docs/concepts/architecture/master-node-communication/) manage workload, nodes, and provide core services like API and discovery


# Kubectl

[Online guide](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands)


# Resources

Can be accessed using their short name or full name.  Multiple resources can be request depending on the kubernetes service/control being used.

```bash
kubectl get po,deploy,svc 
```

Naming

| Short Name | Full Name | Short Name | Full Name	|
| ---	| ---	| ---	| ---	|
|csr  | certificatesigningrequests | cs |componentstatuses|
| cm  | configmaps  | ds | daemonsets |
| deploy | deployments  | ep | endpoints|
| ev | evnets  | hpa | horizontalpodautoscalers|
| ing | ingresses  | limits | limitranges|
| no | nodes  | ns | namespaces|
| pvc | persistentvolumeclaims  | pv | persistentvolumes|
| pdb | poddisruptionbudgets  | psp | podsecuritypolicies|
| rs | replicasets  | rc | replicationcontrollers|
| quota | resourcequotas  | sa | serviceaccounts|
| svc | services  |||

## Get information about a resource type using explain

```bash
$ kubectl explain node
KIND:     Node
VERSION:  v1

DESCRIPTION:
     Node is a worker node in Kubernetes. Each node will have a unique
     identifier in the cache (i.e. in etcd).
.....
   metadata	<Object>
     Standard object's metadata. More info:
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata

   spec	<Object>
     Spec defines the behavior of a node.
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status

```

### Ask for more details about sub-object

```bash
$ kubectl explain node.spec.taints
KIND:     Node
VERSION:  v1

RESOURCE: taints <[]Object>

DESCRIPTION:
     If specified, the node's taints.

     The node this Taint is attached to has the "effect" on any pod that does
     not tolerate the Taint.

FIELDS:
   effect	<string> -required-
     Required. The effect of the taint on pods that do not tolerate the taint.
     Valid effects are NoSchedule, PreferNoSchedule and NoExecute.

   key	<string> -required-
     Required. The taint key to be applied to a node.

   timeAdded	<string>
     TimeAdded represents the time at which the taint was added. It is only
     written for NoExecute taints.

   value	<string>
     Required. The taint value corresponding to the taint key.
```

# Creating a pod on the fly via yaml
This examples uses the run command with the --genertor options.  It asks for a type of run-pod/v1 using the container tutum/dnsutils.  The --dry-run option executes without deploying to the k8s cluster and we redirect the output to dnsutils.yaml

```bash
$ kubectl run dnsutils --generator=run-pod/v1 --image=tutum/dnsutils --dry-run -o yaml > dnsutil.yaml

$ cat dnsutil.yaml 
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: dnsutils
  name: dnsutils
spec:
  containers:
  - image: tutum/dnsutils
    name: dnsutils
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}

$ kubectl create -f dnsutil.yaml 
pod/dnsutils created

$ kubectl get po -o wide
NAME                        READY   STATUS    RESTARTS   AGE    IP           NODE     NOMINATED NODE   READINESS GATES
dnsutils-74bdc55779-zbf5r   1/1     Running   0          5m1s   10.1.1.167   ubuntu   <none>           <none>

$ kubectl exec -it dnsutils-74bdc55779-zbf5r /bin/bash
root@dnsutils-74bdc55779-zbf5r:/# dig .....
```

## Debug problem with container

# Just once
```bash
$ kubectl logs dnsutils-74bdc55779-zbf5r
```

# Attach and stream logs
```bash
$ kubectl logs -f dnsutils-74bdc55779-zbf5r

CTRL-C to exit
```

# Get logs from prior instance of container
```bash
$ kubectl logs -p dnsutils-74bdc55779-zbf5r
```

# Get detailed pod history
```bash
$ kubectl get po -o yaml
apiVersion: v1
items:
- apiVersion: v1
  kind: Pod
  metadata:
    creationTimestamp: "2019-12-06T14:35:26Z"
    generateName: dnsutils-74bdc55779-
    labels:
      pod-template-hash: 74bdc55779
      run: dnsutils
......
  status:
    conditions:
    - lastProbeTime: null
      lastTransitionTime: "2019-12-06T14:35:26Z"
      status: "True"
      type: Initialized
    - lastProbeTime: null
      lastTransitionTime: "2019-12-06T14:35:31Z"
      status: "True"
      type: Ready
    - lastProbeTime: null
      lastTransitionTime: "2019-12-06T14:35:31Z"
      status: "True"
      type: ContainersReady
    - lastProbeTime: null
      lastTransitionTime: "2019-12-06T14:35:26Z"
      status: "True"
      type: PodScheduled
    containerStatuses:
    - containerID: containerd://35f035e758ba069a41067306b93fd7e94e6d49959951ced1b2db655f3423a593
      image: docker.io/tutum/dnsutils:latest
      imageID: sha256:c60036323e6dd80a2b1b33cd9d3021c8513ef7ccbeb69b2d6a6e9de6e6efd7c9
      lastState: {}
      name: dnsutils
      ready: true
      restartCount: 0
      started: true
      state:
        running:
          startedAt: "2019-12-06T14:35:30Z"
```


# Get just the container status
```bash
$ kubectl get po -o json | jq ".items[0].status.containerStatuses"
[
  {
    "containerID": "containerd://35f035e758ba069a41067306b93fd7e94e6d49959951ced1b2db655f3423a593",
    "image": "docker.io/tutum/dnsutils:latest",
    "imageID": "sha256:c60036323e6dd80a2b1b33cd9d3021c8513ef7ccbeb69b2d6a6e9de6e6efd7c9",
    "lastState": {},
    "name": "dnsutils",
    "ready": true,
    "restartCount": 0,
    "started": true,
    "state": {
      "running": {
        "startedAt": "2019-12-06T14:35:30Z"
      }
    }
  }
]
```

# Get just the lifecyle events
```bash
$ kubectl get po -o json | jq ".items[0].status.conditions"
[
  {
    "lastProbeTime": null,
    "lastTransitionTime": "2019-12-06T14:35:26Z",
    "status": "True",
    "type": "Initialized"
  },
  {
    "lastProbeTime": null,
    "lastTransitionTime": "2019-12-06T14:35:31Z",
    "status": "True",
    "type": "Ready"
  },
  {
    "lastProbeTime": null,
    "lastTransitionTime": "2019-12-06T14:35:31Z",
    "status": "True",
    "type": "ContainersReady"
  },
  {
    "lastProbeTime": null,
    "lastTransitionTime": "2019-12-06T14:35:26Z",
    "status": "True",
    "type": "PodScheduled"
  }
]
```

## Getting specific colums and sorting

### Note -A options includes all namespaces
```bash
$ kubectl get po -o wide --sort-by=.spec.nodeName -A
NAMESPACE            NAME                                   READY   STATUS    RESTARTS   AGE   IP           NODE     NOMINATED NODE   READINESS GATES
container-registry   registry-6c99589dc-vtwwp               1/1     Running   97         73d   10.1.1.155   ubuntu   <none>           <none>
default              dnsutils-74bdc55779-zbf5r              1/1     Running   0          28m   10.1.1.167   ubuntu   <none>           <none>
kube-system          coredns-f7867546d-4fh7b                1/1     Running   97         73d   10.1.1.156   ubuntu   <none>           <none>
kube-system          hostpath-provisioner-dd4c58fdb-d54g4   1/1     Running   75         70d   10.1.1.157   ubuntu   <none>           <none>

```

## Match labels

```bash
$ kubectl get deployments -l run=dnsutils
NAME       READY   UP-TO-DATE   AVAILABLE   AGE
dnsutils   1/1     1            1           4h46m
```

## For nodes

```bash
kubectl get nodes -o wide -L beta.kubernetes.io/arch -L beta.kubernetes.io/os -L beta.kubernetes.io/instance-type -L  kops.k8s.io/instancegroup
NAME     STATUS   ROLES    AGE   VERSION   INTERNAL-IP       EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME    ARCH    OS      INSTANCE-TYPE   INSTANCEGROUP
ubuntu   Ready    <none>   77d   v1.16.3   192.168.150.146   <none>        Ubuntu 18.04.3 LTS   4.15.0-70-generic   containerd://1.2.5   amd64   linux

```

# Skaffold deployments

## Build and deploy once
Use run command with --tail to watch the deployment progress

### Note: You need to run make before executing skaffold run
### Run will build the docker containers.  That requies a copy
### of the films microserice in the local directory.  Make
### will ensure that happens

```bash
skaffold run -f manifests/skaffold.yaml --tail
......

You will see containers get build, tagged, and pushed to the local repository

Then the deployment to microk8s

Starting deploy...
 - configmap/cockroach-configmap-6k4c6dfb5b created
 - configmap/films-configmap-5ckf596tfd created
 - service/films created
 - service/roach-ui created
 - deployment.apps/films created
 - deployment.apps/roach-ui created
 - persistentvolumeclaim/roach-ui-claim0 created

Then output from the init containers
[films-5644bf8dfd-5l6h2 wait-for-cockroach] Server:    10.152.183.10
[films-5644bf8dfd-5l6h2 wait-for-cockroach] Address 1: 10.152.183.10 kube-dns.kube-system.svc.cluster.local

repeats for each container, then the db initialization

[films-5644bf8dfd-z42s7 filmsdbinit] ========================================
[films-5644bf8dfd-z42s7 filmsdbinit]  Initializing tables
[films-5644bf8dfd-z42s7 filmsdbinit]  Starting at : Fri Dec  6 20:23:52 UTC 2019
[films-5644bf8dfd-z42s7 filmsdbinit]  Using: cockroach sql --insecure --host=roach-ui:26257
[films-5644bf8dfd-z42s7 filmsdbinit] 

...

CTRL-C to exit
```

## Delete the deployment

```bash
$ skaffold delete -f manifests/skaffold.yaml
Cleaning up...
 - configmap "cockroach-configmap-6k4c6dfb5b" deleted
 - configmap "films-configmap-5ckf596tfd" deleted
 - service "films" deleted
 - service "roach-ui" deleted
 - deployment.apps "films" deleted
 - deployment.apps "roach-ui" deleted
 - persistentvolumeclaim "roach-ui-claim0" deleted
$
```

## Get summary information for deployment
Here we ask for status of the pods, deployments, and services.  This does
not show our backend database, it uses different labels and selectors.

You can ask for the service and deployments using just the name, aka films.
Since the pods each have a unique hash, it won't match.  Using the label
selector always works.


```bash
$ kubectl get po,svc,deploy -l pavedroad.service=films
NAME                        READY   STATUS    RESTARTS   AGE
pod/films-7b47d6bdb-gzpf4   1/1     Running   0          35s
pod/films-7b47d6bdb-k8bgs   1/1     Running   0          35s
pod/films-7b47d6bdb-sdtns   1/1     Running   0          35s

NAME            TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
service/films   NodePort   10.152.183.38   <none>        8081:32396/TCP   36s

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/films   3/3     3            3           35s
$
```
Note: The services is TYPE NodePort.  That means it can be accessed external
to the cluster using the cluster IP + the NodePort assigned.  In this example,
that is 32396.

You can get both of those using dev/ scripts provided.

```bash

$ dev/getk8sip.sh
192.168.150.146

$ dev/getNodePort.sh films
32396
```
## See the database using the -l pavedroad.service=roach-ui

Note: There is no service deploy.  That means access external
      to the cluster is not allowed.

```bash
$ kubectl get po,svc,deploy -l pavedroad.service=roach-ui
NAME                            READY   STATUS    RESTARTS   AGE
pod/roach-ui-5fd4c76975-hplk5   1/1     Running   0          10m

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/roach-ui   1/1     1            1           10m
```

## Debugging

Start by using describe on the deployment

```bash
$ kubectl describe deploy films
Name:                   films
Namespace:              default
CreationTimestamp:      Fri, 06 Dec 2019 12:29:21 -0800
....
Events:
  Type    Reason             Age   From                   Message
  ----    ------             ----  ----                   -------
  Normal  ScalingReplicaSet  16m   deployment-controller  Scaled up replica set films-7b47d6bdb
```

What to look for?

- Init Containers: show configmap mappings and status
- Containers: show configmap mappings and status the application
- Conditions: show status of deployment rollout

Next, use describe on one of the containers

```bash
$ kubectl describe po films-7b47d6bdb-gzpf4
Name:         films-7b47d6bdb-gzpf4
Namespace:    default
Priority:     0
Node:         ubuntu/192.168.150.146
...
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  17m   default-scheduler  Successfully assigned default/films-7b47d6bdb-gzpf4 to ubuntu
  Normal  Pulled     17m   kubelet, ubuntu    Container image "busybox:1.28" already present on machine
  Normal  Created    17m   kubelet, ubuntu    Created container wait-for-cockroach
  Normal  Started    17m   kubelet, ubuntu    Started container wait-for-cockroach
  Normal  Pulled     17m   kubelet, ubuntu    Container image "localhost:32000/acme-demo/filmsinitdb:v0.0alpha-dirty@sha256:1564a6d0abe202a39a5877951e304e4ac1541d97595a79e3199b0322fce8b753" already present on machine
  Normal  Created    17m   kubelet, ubuntu    Created container filmsdbinit
  Normal  Started    17m   kubelet, ubuntu    Started container filmsdbinit
  Normal  Pulled     17m   kubelet, ubuntu    Container image "localhost:32000/acme-demo/films:v0.0alpha-dirty@sha256:23c917f7e41097cad133529ab15458e69bb89697f9500351b509c631750c32cf" already present on machine
  Normal  Created    17m   kubelet, ubuntu    Created container films
  Normal  Started    17m   kubelet, ubuntu    Started container films
```

What to look for?

- Events: Will show you the provisioning status for the init and app containers

## Looking at logs

Note: See kubectl syntax for pods

    namespace/podsname:path/.../file

```bash
# See k8s logs
$ kubectl logs films-7b47d6bdb-gzpf4

$ kubectl cp default/films-7b47d6bdb-gzpf4:/pavedroad/logs/films.log /tmp/films.log

# Get the application logs
$ cat /tmp/films.log
2019/12/06 20:29:33 filmsMain.go:106: Logfile opened logs/films.log
2019/12/06 20:29:33 filmsApp.go:59: Listing at: 0.0.0.0:8081

# log into the pod

$ kubectl exec -it films-7b47d6bdb-gzpf4 /bin/bash

root@films-7b47d6bdb-gzpf4:/pavedroad# cd logs

root@films-7b47d6bdb-gzpf4:/pavedroad/logs# cat films.log
2019/12/06 20:29:33 filmsMain.go:106: Logfile opened logs/films.log
2019/12/06 20:29:33 filmsApp.go:59: Listing at: 0.0.0.0:8081

root@films-7b47d6bdb-gzpf4:/pavedroad/logs# CTRL-D
```
## Manifest

All docker, docker-compose, and kubernetes manifest are located in 
the manifests directory. The kubernetes folder has directories for
each defined enviornment: dev, staging, test, and production.

```
manifests/
├── docker-compose.yaml
├── docker-db-only.yaml
├── Dockerfile
├── InitDbDockerFile
├── kubernetes
│   ├── dev
│   │   ├── db
│   │   │   ├── kustomization.yaml
│   │   │   ├── roach-ui-claim0-persistentvolumeclaim.yaml
│   │   │   ├── roach-ui-deployment.yaml
│   │   │   └── roach-ui-service.yaml
│   │   ├── films
│   │   │   ├── films-deployment.yaml
│   │   │   ├── films-service.yaml
│   │   │   └── kustomization.yaml
│   │   └── kustomization.yaml
│   ├── production
│   ├── staging
│   └── test
└── skaffold.yaml
```

## Skaffold manifest
In addtion to api and kind, there aer two major sections build and deploy.

Build:

Build provides a list of repositories, and a list of artfacts. Each artifact
points to a docker file and includes the repository server to publish too.

It can also include a tagPolicy, if omitted, the git commit hash is used.

Deploy:

Deployments are made using kustomize and kubectl.  Details on kustomzie follow
below.

```bash
apiVersion: skaffold/v1beta9
kind: Config
build:
  insecureRegistries:
    - localhost:32000
  artifacts:
  - image: localhost:32000/acme-demo/films
    context: .
    docker:
      dockerfile: manifests/Dockerfile
  - image: localhost:32000/acme-demo/filmsinitdb
    context: .
    docker:
      dockerfile: manifests/InitDbDockerFile
deploy:
  kustomize:
    path: "manifests/kubernetes/dev"
```

# Kustomize
The top level of each enviornment contains a kustomization.yaml.  
The top kustomization.yaml file also include a list of other locations
needed to build this deployment.  Each of those locations can define
there own kustomization.yaml file.  Those locations can be:

- file system
- http
- git

In the above skaffold.yaml, the deploy section points to 
manifests/kubernetes/dev. That file includes the directories db and films
in its bases declaration.  Then defines labels and annotations that will be
applied to all manifests globally.

manifests/kubernetes/dev/kustomization.yaml
```bash
bases:
  - db
  - films

commonLabels:
  pavedroad.env: dev

commonAnnotations:
  pavedroad.kustomize.base: films/manifests/kubernetes/dev
  pavedroad.kustomize.bases: "films,db"
```

The db kustomization.yaml defines a list of resources.  These are the manifests
it will process.  Manifests will not be included, even if the are in the 
directory, if they are not listed under resources.

It then includes a configMapGenerater.  This generator greates a confipmap
using the name specfied by name:.  And includes key/value pairs defined in
the literals section.

manifests/kubernetes/dev/db/kustomization.yaml
```bash
resources:
  - roach-ui-claim0-persistentvolumeclaim.yaml
  - roach-ui-deployment.yaml
  - roach-ui-service.yaml

configMapGenerator:
- name: cockroach-configmap
  literals:
  - host-ip=roach-ui
```

The films kustomization.yaml is similar only adding additional
labels and annotations to its resource manifests.

manifests/kubernetes/dev/films/kustomization.yaml
```bash
resources:
  - films-deployment.yaml
  - films-service.yaml

commonLabels:
  pavedroad.service: films

commonAnnotations:
  pavedroad.roadctl.version: alphav1
  pavedroad.roadctl.web: www.pavedroad.io
  pavedroad.roadctl.support: support@pavedroad.io

configMapGenerator:
- name: films-configmap
  literals:
  - database-ip=roach-ui
  - ip=0.0.0.0
  - port=8081
```
## NOTE: 

**Configmap names in your kustomization.yaml and your manifests use 
the same name.  However, when deployed, all references to that
configmap include a hash.**

```bash
containers:
  - image: localhost:32000/acme-demo/films:0.0
    env:
    - name: HTTP_IP_ADDR
      valueFrom:
        configMapKeyRef:
          name: films-configmap <<< HERE
          key: ip
```

**Hashed configmap names from our deployment**
```bash
kubectl get cm
NAME                             DATA   AGE
cockroach-configmap-6k4c6dfb5b   1      78m
films-configmap-5ckf596tfd       3      78m
```


## DNS

### Naming

svc.namespace.cluster.local
example: films.default.cluster.local

### Debugging

[k8s doc](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)

Start a debug pod

```bash
$kubectl apply -f https://k8s.io/examples/admin/dns/busybox.yaml
pod/busybox created

# Wait for pod to be ready
$ kubectl get po busybox
NAME      READY   STATUS    RESTARTS   AGE
busybox   1/1     Running   0          97s

# Make sure dns is working
$ kubectl exec -ti busybox -- nslookup kubernetes.default
Server:    10.152.183.10
Address 1: 10.152.183.10 kube-dns.kube-system.svc.cluster.local

Name:      kubernetes.default
Address 1: 10.152.183.1 kubernetes.default.svc.cluster.local


# Check resolv.conf
$ kubectl exec busybox cat /etc/resolv.conf
search default.svc.cluster.local svc.cluster.local cluster.local localdomain
nameserver 10.152.183.10
options ndots:5

# Make sure DNS pod is running
$ kubectl get po -n kube-system -l k8s-app=kube-dns -o wide
NAME                      READY   STATUS    RESTARTS   AGE   IP           NODE     NOMINATED NODE   READINESS GATES
coredns-f7867546d-4fh7b   1/1     Running   97         73d   10.1.1.156   ubuntu   <none>           <none>

# Check logs for all running pods
$ for p in $(kubectl get pods --namespace=kube-system -l k8s-app=kube-dns -o name); do kubectl logs --namespace=kube-system $p; done
.:53
2019-12-05T22:59:33.139Z [INFO] plugin/reload: Running configuration MD5 = 5d839962c224ea2e9fb32222b6a237d1
2019-12-05T22:59:33.140Z [INFO] CoreDNS-1.5.0
2019-12-05T22:59:33.140Z [INFO] linux/amd64, go1.12.2, e3f9a80
CoreDNS-1.5.0
linux/amd64, go1.12.2, e3f9a80

# Make sure the svc is running
$ kubectl get svc -n kube-system kube-dns
NAME       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                  AGE
kube-dns   ClusterIP   10.152.183.10   <none>        53/UDP,53/TCP,9153/TCP   73d

# Make sure end points are exposed
$ kubectl get ep -n kube-system kube-dns
NAME       ENDPOINTS                                     AGE
kube-dns   10.1.1.156:53,10.1.1.156:53,10.1.1.156:9153   73d

# Enable logging
$ kubectl edit configmap -n kube-system coredns


apiVersion: v1
data:
  Corefile: |
    .:53 {
        log  << ADD THIS LINE AND SAVE YOUR CHANGES
        errors
        health
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
          pods insecure
          fallthrough in-addr.arpa ip6.arpa
        }
        prometheus :9153
        forward . 8.8.8.8 8.8.4.4
        cache 30
        loop
        reload
        loadbalance
    }

# To just see the current DNS configuration use

$ kubectl describe configmap -n kube-system coredns
Name:         coredns
Namespace:    kube-system
Labels:       addonmanager.kubernetes.io/mode=EnsureExists
              k8s-app=kube-dns
Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"v1","data":{"Corefile":".:53 {\n    errors\n    health\n    ready\n    kubernetes cluster.local in-addr.arpa ip6.arpa {\n  ...

Data
====
Corefile:
----
.:53 {
    errors
    health
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
      pods insecure
      fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . 8.8.8.8 8.8.4.4
    cache 30
    loop
    reload
    loadbalance
}

Events:  <none>

# Interactively connect to busybox pods and play

$ kubectl exec -ti busybox /bin/sh
/ #
nslookup www.pavedroad.io
Server:    10.152.183.10
Address 1: 10.152.183.10 kube-dns.kube-system.svc.cluster.local

Name:      www.pavedroad.io
Address 1: 185.199.108.153

CTRL-D to exit

# Delete the debug container
$ kubectl delete po busybox
pod "busybox" deleted

# Using dnsutils, get support for dig

$ kubectl run  dnsutils -it --image=tutum/dnsutils --restart=Never
If you don't see a command prompt, try pressing enter.
root@dnsutils:/#
root@dnsutils:/#
root@dnsutils:/# dig www.pavedroad.io

; <<>> DiG 9.9.5-3ubuntu0.2-Ubuntu <<>> www.pavedroad.io
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 46858
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;www.pavedroad.io.		IN	A

;; ANSWER SECTION:
www.pavedroad.io.	30	IN	CNAME	pavedroad.io.
pavedroad.io.		30	IN	A	185.199.108.153

;; Query time: 64 msec
;; SERVER: 10.152.183.10#53(10.152.183.10)
;; WHEN: Fri Dec 06 14:18:26 UTC 2019
;; MSG SIZE  rcvd: 115

root@dnsutils:/# 

CTRL-D to exit

# Note STATUS == Completed, we told it as we told it --restart=Never
$ kubectl get po
NAME       READY   STATUS      RESTARTS   AGE
dnsutils   0/1     Completed   0          2m23s

# Delete it
$ kubectl delete po dnsutils
pod "dnsutils" deleted

# Without --restart option
$ kubectl run  dnsutils -it --image=tutum/dnsutils

If you don't see a command prompt, try pressing enter.
root@dnsutils-74bdc55779-52n9j:/#
root@dnsutils-74bdc55779-52n9j:/# exit

# See the pod name is different now
$ kubectl get po
NAME                        READY   STATUS    RESTARTS   AGE
dnsutils-74bdc55779-52n9j   1/1     Running   1          34s

# To reconnect
$ kubectl exec -it dnsutils-74bdc55779-52n9j /bin/sh
#
CTRL-D

# Clean it up
$ kubectl delete po dnsutils-74bdc55779-52n9j
pod "dnsutils-74bdc55779-52n9j" deleted

```




# Docker

- [Docker CLI](https://docs.docker.com/engine/reference/commandline) CLI reference for docker

## Docker commands

Build a docker file

```bash
$ docker build . -f manifests/InitDbDockerFile
```

Start a container

```bash
$ docker run -it pavedroadio/cockroachdb-client:0.3
```bash

List running containers

```bash
$ docker ps
```

Listing images

Local

```bash
$ docker image ls
```

Remote

```bash
$ docker image ls -all localhost:32000/filmsinitdb*
```

With digest

```bash
$ docker image ls --all --digests localhost:32000/films*
```

Delet images by digest

```bash
$ docker image rm -f digest
```

Starting and stopping a local DB instance for testing

## Docker compose commands

```bash
# start
$ docker-compose -f manifests/docker-db-only.yaml up -d

#stop
$ docker-compose -f manifests/docker-db-only.yaml down
```

# jq

## Just make the output pretty

```bash
kubectl get cm -n kube-system coredns -o json | jq '.'
{
  "apiVersion": "v1",
  "data": {
    "Corefile": ".:53 {\n    errors\n    health\n    ready\n    kubernetes cluster.local in-addr.arpa ip6.arpa {\n      pods insecure\n      fallthrough in-addr.arpa ip6.arpa\n    }\n    prometheus :9153\n    forward . 8.8.8.8 8.8.4.4\n    cache 30\n    loop\n    reload\n    loadbalance\n}\n"
  },
  "kind": "ConfigMap",
.....
}
```

```bash
$ kubectl get cm -n kube-system coredns -o json | jq ".metadata.name"
"coredns"

```

```bash
kubectl get cm -n kube-system coredns -o json | jq '[.metadata.name, .metadata.namespace]'
[
  "coredns",
  "kube-system"
]
```

## Working with arrays

### Get the first element
```bash
$ kubectl get po -o json -A | jq '.items[0].metadata.name'
"registry-6c99589dc-vtwwp"
```

### Get the second

```bash
$ kubectl get po -o json -A | jq '.items[1].metadata.name'
"dnsutils-74bdc55779-zbf5r"
```


## Get all

```bash
$ kubectl get po -o json -A | jq '.items[].metadata.name'
"registry-6c99589dc-vtwwp"
"dnsutils-74bdc55779-zbf5r"
"coredns-f7867546d-4fh7b"
"hostpath-provisioner-dd4c58fdb-d54g4"


## Get a range and pipe to a second filter

.items[0-2] asks for a subrange of the array

You can not just append.metadata when slicing an array like this.

You can pipeline jq filters and ask the second filter to print the
desired data.


```bash
$ kubectl get po -o json -A | jq '.items[0:2] | .[].metadata.name'
"registry-6c99589dc-vtwwp"
"dnsutils-74bdc55779-zbf5r"
```

## Dealing with speical 
Add quotes or brackets

.["foo$"] or ."foo$"

# Project directory structure

```bash
films
├── artifacts
│   ├── coverage.out
│   ├── gosec.out
│   ├── govet.out
│   └── lint.out
├── assets
│   └── images
│       └── films.png
├── builds
│   ├── films-darwin-386
│   ├── films-darwin-amd64
│   └── films-linux-amd64
├── dev
│   ├── busybox-diag.yaml
│   ├── db
│   │   ├── acme-demoAdmin.sql
│   │   ├── acme-demoGrantAdmin.sql
│   │   ├── acme-demo.sql
│   │   ├── filmsCreateTable.sql
│   │   └── filmsExecuteAll.sh
│   ├── films.json
│   ├── filmsPutData.json
│   ├── filmsRepositoryClean.sh
│   ├── getk8sip.sh
│   ├── getNodePort.sh
│   ├── kube-config.sh
│   ├── microk8sStart.sh
│   ├── microk8sStatus.sh
│   ├── microk8sStop.sh
│   ├── sql.sh
│   ├── testAll.sh
│   ├── testDelete.sh
│   ├── testGetList.sh
│   ├── testGet.sh
│   ├── testPost.sh
│   └── testPut.sh
├── docs
│   ├── api.html
│   ├── api.json
│   └── films.html
├── env
│   └── conf
│       ├── dev
│       │   ├── cockroach.yaml
│       │   ├── cockroach.yaml.bak
│       │   ├── films-db.yaml
│       │   ├── films-http.yaml
│       │   ├── kustomization.bak
│       │   └── kustomization.yaml
│       ├── production
│       ├── README.md
│       ├── staging
│       └── test
├── films
├── filmsApp.go
├── filmsDoc.go
├── filmsMain.go
├── filmsModel.go
├── films_test.go
├── films.yaml
├── Gopkg.lock
├── Gopkg.toml
├── logs
├── Makefile
├── manifests
│   ├── docker-compose.yaml
│   ├── docker-db-only.yaml
│   ├── Dockerfile
│   ├── InitDbDockerFile
│   ├── kubernetes
│   │   ├── dev
│   │   │   ├── db
│   │   │   │   ├── kustomization.yaml
│   │   │   │   ├── roach-ui-claim0-persistentvolumeclaim.yaml
│   │   │   │   ├── roach-ui-deployment.yaml
│   │   │   │   └── roach-ui-service.yaml
│   │   │   ├── films
│   │   │   │   ├── films-deployment.yaml
│   │   │   │   ├── films-service.yaml
│   │   │   │   └── kustomization.yaml
│   │   │   └── kustomization.yaml
│   │   ├── production
│   │   ├── staging
│   │   └── test
│   └── skaffold.yaml
├── README.md
├── sonarcloud.sh
├── sonar-project.properties
└── svcfilms.yaml
```

# References
- [kubernetes.io cheatsheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/) More Kubectl cheatsheets
- [kubectl commands cheatsheet](https://medium.com/faun/kubectl-commands-cheatsheet-43ce8f13adfb) Kenichi Shibata - General commands for an ops person
- [Kubectl output options](https://gist.github.com/so0k/42313dbb3b547a0f51a547bb968696ba) so0k - JSON path, jq, custom columns
- [JQ introduction](http://www.compciv.org/recipes/cli/jq-for-parsing-json/) compciv.org - Basic introduction to parsing with jq
- [jq github](https://stedolan.github.io/jq/) stedolan - GitHub home page
- [jq playground](https://jqplay.org/) stedolan - jq playground
- [jq docs](https://stedolan.github.io/jq/manual/) JQ Manual
- [skaffold.dev](https://skaffold.dev/docs/quickstart/) Documentation and Tutorials{{end}}
