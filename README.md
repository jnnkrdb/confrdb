# ConfRDB
ConfRDB is a Kubernetes Operator Package, that creates and updates Secrets and ConfigMaps in the cluster. It is used, to read a GlobalConfig or a GlobalSecret from the cluster and replicate its data into ConfigMaps or Secrets in the desired namespaces and keep the resources updated. 

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/jnnkrdb/configrdb)
[![CodeFactor](https://www.codefactor.io/repository/github/jnnkrdb/configrdb/badge)](https://www.codefactor.io/repository/github/jnnkrdb/configrdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/jnnkrdb/configrdb)](https://goreportcard.com/report/github.com/jnnkrdb/configrdb)
[![GitHub issues](https://badgen.net/github/issues/jnnkrdb/configrdb/)](https://github.com/jnnkrdb/configrdb/issues/)

## Table of Contents

- [Installation](#installation)
  - [Defaults](#defaults)
    - [Namespace](#namespace)
    - [CustomResourceDefinition GlobalConfig](#customresourcedefinition-globalconfig)
    - [CustomResourceDefinition GlobalSecret](#customresourcedefinition-globalsecret)
  - [Operator](#operator)
    - [ServiceAccount](#serviceaccount)
    - [ClusterRole](#clusterrole)
    - [ClusterRoleBinding](#clusterrolebinding)
    - [Deployment](#deployment)
    - [CustomResourceDefinition](#customresourcedefinition)
  - [Example Deployments](#example-deployments)
    - [GlobalConfig](#globalconfig)
    - [GlobalSecret](#globalsecret)
- [Configuration](#configuration)
  - [Operator Environment Variables](#operator-environment-variables)
  - [UI-Controller Angular Config](#ui-controller-angular-config)
- [RoadMap or Planned](#roadmap-or-planned)
    
## Installation
  
This part is about the installation of the ConfRDB service. It contains the collection of the kubernetes manifests and a short explanation about the overall service configuration. To get this service running, you need to deploy the yaml-files to your kubernetes cluster. The deployment of the ConfigMaps/Secrets will be handled with the CRDs of this project. Deploy a GlobalConfig to rollout ConfigMaps into the configured namespaces or deploy a GlobalSecret to do so with Secrets.

To deploy the service to your cluster, there are the following manifests, which are recommended to run the service. The manifests are minimalistic and do only contain the minimum neccessary information.

### Defaults
All ConfRDB-Controllers need some default deployment-manifests, for example the CustomResourceDefinitions. Those default deployments are listed in this section.
- [Namespace](#namespace)
- [CustomResourceDefinition GlobalConfig](#customresourcedefinition-globalconfig)
- [CustomResourceDefinition GlobalSecret](#customresourcedefinition-globalsecret)

#### Namespace
```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: confrdb
  labels:
    app: confrdb
```

#### CustomResourceDefinition GlobalConfig
```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalconfigs.globals.jnnkrdb.de
spec:
  group: globals.jnnkrdb.de
  names:
    kind: GlobalConfig
    listKind: GlobalConfigList
    plural: globalconfigs
    shortNames:
    - gc
    - gcs
    singular: globalconfig
  scope: Namespaced
  versions:
  - name: v1beta2
    schema:
      openAPIV3Schema:
        description: GlobalConfig is the Schema for the globalconfigs API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: GlobalConfigSpec defines the desired state of GlobalConfig
            properties:
              data:
                additionalProperties:
                  type: string
                type: object
              namespaces:
                description: struct which contains the information about the namespace
                  regex
                properties:
                  avoidregex:
                    default:
                    - default
                    items:
                      type: string
                    type: array
                  matchregex:
                    default:
                    - default
                    items:
                      type: string
                    type: array
                required:
                - avoidregex
                - matchregex
                type: object
            required:
            - data
            - namespaces
            type: object
          status:
            description: GlobalConfigStatus defines the observed state of GlobalConfig
            properties:
              conditions:
                items:
                  description: "Condition contains details for one aspect of the current
                    state of this API Resource. --- This struct is intended for direct
                    use as an array at the field path .status.conditions.  For example,
                    \n type FooStatus struct{ // Represents the observations of a
                    foo's current state. // Known .status.conditions.type are: \"Available\",
                    \"Progressing\", and \"Degraded\" // +patchMergeKey=type // +patchStrategy=merge
                    // +listType=map // +listMapKey=type Conditions []metav1.Condition
                    `json:\"conditions,omitempty\" patchStrategy:\"merge\" patchMergeKey:\"type\"
                    protobuf:\"bytes,1,rep,name=conditions\"` \n // other fields }"
                  properties:
                    lastTransitionTime:
                      description: lastTransitionTime is the last time the condition
                        transitioned from one status to another. This should be when
                        the underlying condition changed.  If that is not known, then
                        using the time when the API field changed is acceptable.
                      format: date-time
                      type: string
                    message:
                      description: message is a human readable message indicating
                        details about the transition. This may be an empty string.
                      maxLength: 32768
                      type: string
                    observedGeneration:
                      description: observedGeneration represents the .metadata.generation
                        that the condition was set based upon. For instance, if .metadata.generation
                        is currently 12, but the .status.conditions[x].observedGeneration
                        is 9, the condition is out of date with respect to the current
                        state of the instance.
                      format: int64
                      minimum: 0
                      type: integer
                    reason:
                      description: reason contains a programmatic identifier indicating
                        the reason for the condition's last transition. Producers
                        of specific condition types may define expected values and
                        meanings for this field, and whether the values are considered
                        a guaranteed API. The value should be a CamelCase string.
                        This field may not be empty.
                      maxLength: 1024
                      minLength: 1
                      pattern: ^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$
                      type: string
                    status:
                      description: status of the condition, one of True, False, Unknown.
                      enum:
                      - "True"
                      - "False"
                      - Unknown
                      type: string
                    type:
                      description: type of condition in CamelCase or in foo.example.com/CamelCase.
                        --- Many .condition.type values are consistent across resources
                        like Available, but because arbitrary conditions can be useful
                        (see .node.status.conditions), the ability to deconflict is
                        important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
                      maxLength: 316
                      pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$
                      type: string
                  required:
                  - lastTransitionTime
                  - message
                  - reason
                  - status
                  - type
                  type: object
                type: array
              deployedconfigmaps:
                items:
                  type: object
                type: array
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
``` 

#### CustomResourceDefinition GlobalSecret
```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalsecrets.globals.jnnkrdb.de
spec:
  group: globals.jnnkrdb.de
  names:
    kind: GlobalSecret
    listKind: GlobalSecretList
    plural: globalsecrets
    shortNames:
    - gs
    - gss
    singular: globalsecret
  scope: Namespaced
  versions:
  - name: v1beta2
    schema:
      openAPIV3Schema:
        description: GlobalSecret is the Schema for the globalsecrets API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: GlobalSecretSpec defines the desired state of GlobalSecret
            properties:
              data:
                additionalProperties:
                  type: string
                type: object
              namespaces:
                description: struct which contains the information about the namespace
                  regex
                properties:
                  avoidregex:
                    default:
                    - default
                    items:
                      type: string
                    type: array
                  matchregex:
                    default:
                    - default
                    items:
                      type: string
                    type: array
                required:
                - avoidregex
                - matchregex
                type: object
              type:
                enum:
                - Opaque
                - kubernetes.io/service-account-token
                - kubernetes.io/dockercfg
                - kubernetes.io/dockerconfigjson
                - kubernetes.io/basic-auth
                - kubernetes.io/ssh-auth
                - kubernetes.io/tls
                - bootstrap.kubernetes.io/token
                type: string
            required:
            - data
            - namespaces
            - type
            type: object
          status:
            description: GlobalSecretStatus defines the observed state of GlobalSecret
            properties:
              conditions:
                items:
                  description: "Condition contains details for one aspect of the current
                    state of this API Resource. --- This struct is intended for direct
                    use as an array at the field path .status.conditions.  For example,
                    \n type FooStatus struct{ // Represents the observations of a
                    foo's current state. // Known .status.conditions.type are: \"Available\",
                    \"Progressing\", and \"Degraded\" // +patchMergeKey=type // +patchStrategy=merge
                    // +listType=map // +listMapKey=type Conditions []metav1.Condition
                    `json:\"conditions,omitempty\" patchStrategy:\"merge\" patchMergeKey:\"type\"
                    protobuf:\"bytes,1,rep,name=conditions\"` \n // other fields }"
                  properties:
                    lastTransitionTime:
                      description: lastTransitionTime is the last time the condition
                        transitioned from one status to another. This should be when
                        the underlying condition changed.  If that is not known, then
                        using the time when the API field changed is acceptable.
                      format: date-time
                      type: string
                    message:
                      description: message is a human readable message indicating
                        details about the transition. This may be an empty string.
                      maxLength: 32768
                      type: string
                    observedGeneration:
                      description: observedGeneration represents the .metadata.generation
                        that the condition was set based upon. For instance, if .metadata.generation
                        is currently 12, but the .status.conditions[x].observedGeneration
                        is 9, the condition is out of date with respect to the current
                        state of the instance.
                      format: int64
                      minimum: 0
                      type: integer
                    reason:
                      description: reason contains a programmatic identifier indicating
                        the reason for the condition's last transition. Producers
                        of specific condition types may define expected values and
                        meanings for this field, and whether the values are considered
                        a guaranteed API. The value should be a CamelCase string.
                        This field may not be empty.
                      maxLength: 1024
                      minLength: 1
                      pattern: ^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$
                      type: string
                    status:
                      description: status of the condition, one of True, False, Unknown.
                      enum:
                      - "True"
                      - "False"
                      - Unknown
                      type: string
                    type:
                      description: type of condition in CamelCase or in foo.example.com/CamelCase.
                        --- Many .condition.type values are consistent across resources
                        like Available, but because arbitrary conditions can be useful
                        (see .node.status.conditions), the ability to deconflict is
                        important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
                      maxLength: 316
                      pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$
                      type: string
                  required:
                  - lastTransitionTime
                  - message
                  - reason
                  - status
                  - type
                  type: object
                type: array
              deployedsecrets:
                items:
                  type: object
                type: array
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
```

### Operator
The Operator contains the core functionality of this controller package. The operator requests GlobalConfigs and GlobalSecrets in the cluster and creates the ConfigMaps and Secrets with their specifications.
The Controller needs some specific kubernetes manifests to show its full potential:
  - [ServiceAccount](#serviceaccount)
  - [ClusterRole](#clusterrole)
  - [ClusterRoleBinding](#clusterrolebinding)
  - [Deployment](#deployment)

#### ServiceAccount
```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: confrdb
  name: sa-confrdb
  labels:
    app: confrdb
```
 
#### ClusterRole
```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cr-confrdb
rules:
  # Get/List/Watch Namespaces
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list", "watch"]
  # Get/Create/Patch/Update/Delete Secrets
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "create", "patch", "update", "delete"]
  # Get/List/Watch/Create/Patch/Update/Delete ConfigMaps 
  # Receives more rights for COnfigMaps than for Secrets
  # because of the leader election
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: ["get", "list", "watch", "create", "patch", "update", "delete"]
  # Get/List/Watch/Create/Patch/Update/Delete GlobalConfigs and GlobalSecrets
- apiGroups: ["globals.jnnkrdb.de"]
  resources: ["globalconfigs", "globalsecrets"]
  verbs: ["get", "list", "watch", "create", "patch", "update", "delete"]  
  # Update GlobalConfigs and GlobalSecrets Finalizers
- apiGroups: ["globals.jnnkrdb.de"]
  resources: ["globalconfigs/finalizers", "globalsecrets/finalizers"]
  verbs: ["update"]  
  # Get/Patch/Update GlobalConfigs and GlobalSecrets Status
- apiGroups: ["globals.jnnkrdb.de"]
  resources: ["globalconfigs/status", "globalsecrets/status"]
  verbs: ["get", "patch", "update"]
  # Get/List/Watch/Create/Patch/Update/Delete Leases for LeaderElection 
- apiGroups: ["coordination.k8s.io"]
  resources: ["leases"]
  verbs: ["get", "list", "watch", "create", "patch", "update", "delete"]
  # Create/Patch Events
- apiGroups: [""]
  resources: ["events"]
  verbs: ["create", "patch"]
```  

#### ClusterRoleBinding
```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: crb-confrdb
  labels:
    app: confrdb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cr-confrdb
subjects:
  - kind: ServiceAccount
    name: sa-confrdb
    namespace: confrdb
```
 
#### Deployment
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: confrdb
  namespace: confrdb
  labels:
    app: confrdb
spec:
  selector:
    matchLabels:
      app: confrdb
  template:
    metadata:
      labels:
        app: confrdb
    spec:
      securityContext:
        runAsNonRoot: true
      serviceAccountName: sa-confrdb
      containers:
      - name: confrdb
        image: docker.io/jnnkrdb/confrdb:v0.1.0
        imagePullPolicy: Always
        command:
        - /manager
        args:
        - --leader-elect
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
              - "ALL"
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 10m
            memory: 64Mi
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
      terminationGracePeriodSeconds: 10
```

### Example Deployments

In this section you can find some example deployments of the GlobalConfig and/or GlobalSecret resources.
  - [GlobalConfig](#globalconfig)
  - [GlobalSecret](#globalsecret)

#### GlobalConfig
```yaml
---
apiVersion: globals.jnnkrdb.de/v1alpha1
kind: GlobalConfig
metadata:
  name: gc-name
  namespace: default
spec:
  namespaces:
    avoidregex: 
      - default # matches namespace "default" -> namespace default will be avoided
      - prod. # matches namespaces like "production-financial", "prod-databases", "prod*" -> namespaces like "production-financial", "prod-databases" or "prod*" will be avoided
    matchregex: 
      - production-mssql # matches namespace "production-mssql", BUT since "prod."-regex is in the avoidregex-list, this namespace will not be matched
      - .dev # matches namespaces like "financials-dev", "databases-dev", "dev", etc. -> namespaces with the suffix "dev" will be matched
      - .internal. # matches namespaces like "test-internal-financials", "databases-internals", "internal", etc. -> namespaces, which contain the substring "internal" will be matched
  data: # the data section should be filled like the data-section of a normal configmap

    # kubernetes example of a configmap -> https://kubernetes.io/docs/concepts/configuration/configmap/
    # property-like keys; each key maps to a simple value
    player_initial_lives: "3"
    ui_properties_file_name: "user-interface.properties"

    # file-like keys
    game.properties: |
      enemy.types=aliens,monsters
      player.maximum-lives=5    
    user-interface.properties: |
      color.good=purple
      color.bad=yellow
      allow.textmode=true    
```

#### GlobalSecret
```yaml
---
apiVersion: globals.jnnkrdb.de/v1alpha1
kind: GlobalSecret
metadata:
  name: gs-name
spec:
  namespaces:
    avoidregex: []
    matchregex: 
      - "." # matches all namespaces
  type: kubernetes.io/dockerconfigjson # or other type, supported by kubernetes secrets -> https://kubernetes.io/docs/concepts/configuration/secret/
  data: # must be base64 encrypted by yourself, but like the globalconfig, this section is build like its underlying secret
    .dockerconfigjson: <base64 encrypted docker config json file>
```

## Configuration

The Operator package must be configured for each controller seperatly.
  - [Operator Arguments](#operator-arguments)

#### Operator Arguments

- `--leader-elect` (+Optional): determines whether or not to use leader election when starting the manager.

## RoadMap or Planned
- Validation for SecretTypes + Configuration
- High Availability Synchronization
- Prometheus Metrics
