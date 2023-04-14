# Dapr Middleware component using WASM

Use TinyGo to compile the `filter.go` program:

```
tinygo build -o filter.wasm -scheduler=none --no-debug -target=wasi filter.go
```



Create configMap from binary: 

```
kubectl create configmap wasm-filter --from-file=filter.wasm
```

```
kubectl delete cm wasm-filter && kubectl create configmap wasm-filter --from-file=filter.wasm && kubectl delete pod -l=app.kubernetes.io/name=frontend-app
```

You need to delete the write-app pod, so the `filter.wasm` file gets reloaded. 



