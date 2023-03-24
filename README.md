# Dapr and Wazero :: Demo step-by-step tutorial - KubeCon EU 2023

On this step by step tutorial we will install Dapr into a Kubernetes Cluster, create two simple application to use Dapr Components and then customize how Dapr works by using the Wazero Runtime to extend Dapr behaviour using Webassembly. 

## Pre requisites and installation

We will be creating a local Kubernetes Cluster using KinD, installing Dapr and Redis using Helm. 
You will need to install KinD, `kubectl` and `helm` in your workstation. Then you can run the following commands: 

Create a local Kubernetes cluster with: 
```
kind create cluster
```

Let's create a Redis instance that our applications can use to store state or exchange messages: 

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update                            
helm install redis bitnami/redis --set image.tag=6.2 --set architecture=standalone
```

Finally, let's install Dapr into the Cluster: 

```
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.10.4 \
--namespace dapr-system \
--create-namespace \
--wait
```

## Deploying the applications and wiring things together

In this section, we will be deploying three applications that want to store and read data from a state store and publish and consume messages. 
To achieve this we will use the Dapr StateStore and PubSub components. So before deploying our applications let's configure these components to connect the Redis instance that we created before. 

The Dapr Statestore configuration looks like this: 
```
kind: Component
metadata:
  name: statestore
spec:
  type: state.redis
  version: v1
  metadata:
  - name: keyPrefix
    value: name
  - name: redisHost
    value: redis-master:6379
  - name: redisPassword
    secretKeyRef:
      name: redis
      key: redis-password
auth:
  secretStore: kubernetes
```

We can apply this resource to Kubernetes by running: 
```
kubectl apply -f resources/statestore.yaml
```

The PubSub Component looks like this: 
```
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: notifications-pubsub
spec:
  type: pubsub.redis
  version: v1
  metadata:
  - name: redisHost
    value: redis-master:6379
  - name: redisPassword
    secretKeyRef:
      name: redis
      key: redis-password
auth:
  secretStore: kubernetes
```

We can apply this resource to Kubernetes by running: 
```
kubectl apply -f resources/pubsub.yaml
```

Once we have the PubSub component configured, we can register Subscritions to define who and where notifications will be sent when new messages arrive to a certain topic. A Subscription resource look like this: 

```
apiVersion: dapr.io/v1alpha1
kind: Subscription
metadata:
  name: notifications-subscritpion
spec:
  topic: notifications 
  route: /notifications
  pubsubname: notifications-pubsub
```

Finally, let's deploy three applications that uses the Dapr StateStore and PubSub components. 

```
kubectl apply -f apps.yaml
```

