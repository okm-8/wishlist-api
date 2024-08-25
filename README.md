# Wishlist API
Second Proof of Concept for Obelus One

## Installation
Build local image and run release in local k8s

```bash
make docker-build helm-install
```

## Public server
Forward public server from port 8080

```bash
make kube-forward-public 
```

## Private server
Forward private server from port 8081

```bash
make kube-forward-private
```

## Postgres
Forward postgres from port 8001

```bash
make kube-forward-postgres
```

## Redis
Forward postgres from port 8002

```bash
make kube-forward-redis
```

## Stop

```bash
make helm-delete
```
