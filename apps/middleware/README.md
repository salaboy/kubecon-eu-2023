# Dapr Middleware component using WASM

Use TinyGo to compile the router.go program:

```
tinygo build -o router.wasm -scheduler=none --no-debug -target=wasi router.go
```


Create configMap from binary: 

```
kubectl delete cm wasm-router && kubectl create configmap wasm-router --from-file=router.wasm
```
