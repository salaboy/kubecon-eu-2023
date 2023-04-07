# Dapr Middleware component using WASM

Use TinyGo to compile the router.go program:

```
tinygo build -o router.wasm -scheduler=none --no-debug -target=wasi router.go
```


Create configMap from binary: 

```
kubectl delete cm wasm-router && kubectl create configmap wasm-router --from-file=router.wasm && kubectl delete pod -l=app.kubernetes.io/name=write-app
```

You need to delete the write-app pod, so the router.wasm file gets reloaded. 



