# Dapr and Wazero :: Demo step-by-step tutorial - KubeCon EU 2023

On this step-by-step tutorial we will install Dapr into a Kubernetes Cluster, install four application  applications that use Dapr Components to interact with available infrastructure. Once we get things working we will use the Dapr Middleware Compponent that integrates with the Wazero Runtime to customize how the application behaves by extending the infrastructure.

This tutorial is divided into three parts: 
- [Prerequisites and Installation]()
- [Installing the Applications and wiring things together]()
- [Extending Infrastructure with Wazero and WebAssembly]()

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

Finally, let's deploy three applications that uses the Dapr StateStore and PubSub components. This are normal/regular Kubernetes applications, using Kubernetes `Deployments`. To make these apps Dapr-aware we just need to add some Dapr annotations, for example for the `Read App`:

```
annotations:
  dapr.io/app-id: read-app
  dapr.io/app-port: "8080"
  dapr.io/enabled: "true"
```

Let's deploy these apps with: 
```
kubectl apply -f resources/apps.yaml
```

Now you can access the Frontend application by using `kubectl port-forward`:

```
kubectl port-forward svc/frontend 8080:80
```

And then pointing your browser to [http://localhost:8080](http://localhost:8080)


## Extending Infrastructure with Wazero and WebAssembly

Now that we have our application up and running, let's use Wazero and the [Dapr WASM Middleware Component]() to extend our application infrastructure with a middleware filter.

You can find the [filter source code here](apps/middleware/).

This is a very simple filter that reads the body of the requests and do string replacements on the contents. 

To apply this filter to the `Write App` we need to first compile the filter source code written in Go using `tinyGo` (add link). 

This generates a `.wasm` file that we can run everywhere with any WASM runtime. For this tutorial we will be using the Wazero WebAssembly runtime that is already integrated with Dapr. 

To make this `.wasm` file available in our cluster we will be using Kubernetes `ConfigMaps` which then will be mounted as a volume for the Dapr sidecar to use. 

To configure this new filter we need to define two things: 
- A Dapr Middleware component
- A Dapr Configuration

Once we have the Middleware component and the configuration, we only need to update our `Write App` to use this configuration. 



