{{ define "install.tpl"}}
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: eventbus.argoproj.io
spec:
  group: argoproj.io
  names:
    kind: EventBus
    listKind: EventBusList
    plural: eventbus
    shortNames:
    - eb
    singular: eventbus
  scope: Namespaced
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            type: object
            x-kubernetes-preserve-unknown-fields: true
          status:
            type: object
            x-kubernetes-preserve-unknown-fields: true
        required:
        - metadata
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: eventsources.argoproj.io
spec:
  group: argoproj.io
  names:
    kind: EventSource
    listKind: EventSourceList
    plural: eventsources
    shortNames:
    - es
    singular: eventsource
  scope: Namespaced
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            type: object
            x-kubernetes-preserve-unknown-fields: true
          status:
            type: object
            x-kubernetes-preserve-unknown-fields: true
        required:
        - metadata
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: sensors.argoproj.io
spec:
  group: argoproj.io
  names:
    kind: Sensor
    listKind: SensorList
    plural: sensors
    shortNames:
    - sn
    singular: sensor
  scope: Namespaced
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            type: object
            x-kubernetes-preserve-unknown-fields: true
          status:
            type: object
            x-kubernetes-preserve-unknown-fields: true
        required:
        - metadata
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: argo-events-sa
  namespace: argo-events
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
  name: argo-events-aggregate-to-admin
rules:
- apiGroups:
  - argoproj.io
  resources:
  - sensors
  - sensors/finalizers
  - sensors/status
  - eventsources
  - eventsources/finalizers
  - eventsources/status
  - eventbus
  - eventbus/finalizers
  - eventbus/status
  verbs:
  - create
  - delete
  - deletecollection
  - get
  - list
  - patch
  - update
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
  name: argo-events-aggregate-to-edit
rules:
- apiGroups:
  - argoproj.io
  resources:
  - sensors
  - sensors/finalizers
  - sensors/status
  - eventsources
  - eventsources/finalizers
  - eventsources/status
  - eventbus
  - eventbus/finalizers
  - eventbus/status
  verbs:
  - create
  - delete
  - deletecollection
  - get
  - list
  - patch
  - update
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
  name: argo-events-aggregate-to-view
rules:
- apiGroups:
  - argoproj.io
  resources:
  - sensors
  - sensors/finalizers
  - sensors/status
  - eventsources
  - eventsources/finalizers
  - eventsources/status
  - eventbus
  - eventbus/finalizers
  - eventbus/status
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: argo-events-role
rules:
- apiGroups:
  - argoproj.io
  resources:
  - sensors
  - sensors/finalizers
  - sensors/status
  - eventsources
  - eventsources/finalizers
  - eventsources/status
  - eventbus
  - eventbus/finalizers
  - eventbus/status
  verbs:
  - create
  - delete
  - deletecollection
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - pods
  - pods/exec
  - configmaps
  - secrets
  - services
  - persistentvolumeclaims
  verbs:
  - create
  - get
  - list
  - watch
  - update
  - patch
  - delete
- apiGroups:
  - apps
  resources:
  - deployments
  - statefulsets
  verbs:
  - create
  - get
  - list
  - watch
  - update
  - patch
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: argo-events-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: argo-events-role
subjects:
- kind: ServiceAccount
  name: argo-events-sa
  namespace: argo-events
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: eventbus-controller
  namespace: argo-events
spec:
  replicas: 1
  selector:
    matchLabels:
      app: eventbus-controller
  template:
    metadata:
      labels:
        app: eventbus-controller
    spec:
      containers:
      - args:
        - eventbus-controller
        env:
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: NATS_STREAMING_IMAGE
          value: nats-streaming:0.17.0
        - name: NATS_METRICS_EXPORTER_IMAGE
          value: synadia/prometheus-nats-exporter:0.6.2
        image: quay.io/argoproj/argo-events:v1.4.0
        imagePullPolicy: Always
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 3
          periodSeconds: 3
        name: eventbus-controller
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 3
          periodSeconds: 3
      securityContext:
        runAsNonRoot: true
        runAsUser: 9731
      serviceAccountName: argo-events-sa
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: eventsource-controller
  namespace: argo-events
spec:
  replicas: 1
  selector:
    matchLabels:
      app: eventsource-controller
  template:
    metadata:
      labels:
        app: eventsource-controller
    spec:
      containers:
      - args:
        - eventsource-controller
        env:
        - name: EVENTSOURCE_IMAGE
          value: quay.io/argoproj/argo-events:v1.4.0
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/argoproj/argo-events:v1.4.0
        imagePullPolicy: Always
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 3
          periodSeconds: 3
        name: eventsource-controller
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 3
          periodSeconds: 3
      securityContext:
        runAsNonRoot: true
        runAsUser: 9731
      serviceAccountName: argo-events-sa
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sensor-controller
  namespace: argo-events
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sensor-controller
  template:
    metadata:
      labels:
        app: sensor-controller
    spec:
      containers:
      - args:
        - sensor-controller
        env:
        - name: SENSOR_IMAGE
          value: quay.io/argoproj/argo-events:v1.4.0
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/argoproj/argo-events:v1.4.0
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 3
          periodSeconds: 3
        name: sensor-controller
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 3
          periodSeconds: 3
      securityContext:
        runAsNonRoot: true
        runAsUser: 9731
      serviceAccountName: argo-events-sa
{{end}}
