# Tiny Test Helm Chart

A Helm chart for deploying the tiny-test application to any Kubernetes cluster.

## Features

- **Security-hardened**: Non-root user, dropped capabilities, seccomp profile, read-only root filesystem
- **Resource-efficient**: Minimal resource requests (10m CPU, 16Mi memory) with sensible limits
- **Production-ready**: Health checks, readiness probes, configurable replicas
- **Flexible**: Works on any Kubernetes cluster (EKS, AKS, GKE, Minikube, kind, k3s)
- **Small footprint**: ~2.13 MB container image

## Quick Start

```bash
# Install with default values
helm upgrade --install tiny-test ./tiny-test

# Install to specific namespace
helm upgrade --install tiny-test ./tiny-test \
  --namespace my-namespace \
  --create-namespace

# Install with custom values
helm upgrade --install tiny-test ./tiny-test \
  --set replicaCount=3 \
  --set service.type=LoadBalancer
```

## Configuration

See [values.yaml](values.yaml) for all available options.

### Common Configurations

#### Development/Testing (NodePort)

```bash
helm upgrade --install tiny-test ./tiny-test \
  --set service.type=NodePort
```

#### Production (LoadBalancer with replicas)

```bash
helm upgrade --install tiny-test ./tiny-test \
  --set replicaCount=3 \
  --set service.type=LoadBalancer \
  --set resources.requests.memory=32Mi \
  --set resources.limits.memory=64Mi
```

#### With Ingress

```bash
helm upgrade --install tiny-test ./tiny-test \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.hosts[0].host=tiny-test.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

#### Specific Image Version

```bash
helm upgrade --install tiny-test ./tiny-test \
  --set image.tag=v1.0.0
```

## Values

| Key | Description | Default |
|-----|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.registry` | Docker registry | `docker.io` |
| `image.repository` | Image repository | `bansikah/tiny-test` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `container.port` | Container port | `8080` |
| `resources.requests.cpu` | CPU request | `10m` |
| `resources.requests.memory` | Memory request | `16Mi` |
| `resources.limits.cpu` | CPU limit | `50m` |
| `resources.limits.memory` | Memory limit | `32Mi` |
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class | `""` |

See [values-examples.yaml](values-examples.yaml) for more example configurations.

## Security

The chart implements Kubernetes security best practices:

- **Pod Security**: Non-root user (UID 65532), dropped all capabilities
- **Seccomp Profile**: RuntimeDefault
- **Resource Limits**: Prevents resource exhaustion
- **Health Checks**: Liveness and readiness probes
- **Service Account**: Dedicated service account created

## Verification

After installation:

```bash
# Check all resources
helm list
kubectl get all -l app.kubernetes.io/name=tiny-test

# View logs
kubectl logs -l app=tiny-test

# Test health endpoint
kubectl port-forward svc/tiny-test-tiny-test 8080:80
curl http://localhost:8080/healthz
```

## Uninstallation

```bash
helm uninstall tiny-test

# If installed to specific namespace
helm uninstall tiny-test -n my-namespace
```

## Chart Testing

```bash
# Lint the chart
helm lint ./tiny-test

# Dry run to see generated manifests
helm install tiny-test ./tiny-test --dry-run --debug

# Template to see rendered YAML
helm template tiny-test ./tiny-test
```
