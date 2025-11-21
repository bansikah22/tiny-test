# Tiny Test App

A minimal Golang web application designed for testing Kubernetes deployments. The application provides a lightweight container image (~2.13 MB) with a simple web UI and health check endpoints, optimized with UPX compression.

## Quick Start

Deploy to your Kubernetes cluster in seconds:

```bash
kubectl apply -f k8s/
kubectl port-forward service/tiny-test 8080:80
```

Open `http://localhost:8080` in your browser. The image is automatically pulled from Docker Hub.

## Features

- Ultra-lightweight container image (~2.13 MB)
- Simple web UI with deployment status, pod name, version, uptime, and request statistics
- Health check endpoint at `/healthz`
- Version endpoint at `/version`
- Info endpoint at `/info` with comprehensive JSON response including metrics
- Metrics endpoint at `/metrics` with Prometheus-style format
- Built-in request tracking and statistics per endpoint
- Uptime tracking and display
- Kubernetes-ready with deployment and service manifests

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster access
- `kubectl` configured and connected to your cluster

### Deploy

The deployment uses the pre-built image from Docker Hub (`bansikah/tiny-test:latest`). No local build required.

```bash
kubectl apply -f k8s/
```

This will create:
- A Deployment with 1 replica
- A ClusterIP Service exposing port 80

### Access the Application

#### Port Forward (Recommended for Testing)

```bash
kubectl port-forward service/tiny-test 8080:80
```

Then open `http://localhost:8080` in your browser.

#### NodePort

To expose via NodePort, modify `k8s/service.yaml`:

```yaml
spec:
  type: NodePort
```

Then apply:

```bash
kubectl apply -f k8s/service.yaml
kubectl get service tiny-test
```

Access via any node IP on the assigned NodePort.

#### Ingress

Create an Ingress resource pointing to the `tiny-test` service:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tiny-test
spec:
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: tiny-test
            port:
              number: 80
```

Apply with:

```bash
kubectl apply -f ingress.yaml
```

### Verify Deployment

Check deployment status:

```bash
kubectl get deployment tiny-test
kubectl get pods -l app=tiny-test
kubectl get service tiny-test
```

View pod logs:

```bash
kubectl logs -l app=tiny-test
```

### Image Details

- **Docker Hub**: `bansikah/tiny-test:latest`
- **Image Size**: ~2.13 MB
- **Port**: 8080 (exposed as port 80 in service)

## Endpoints

- `GET /` - Web UI showing deployment information, uptime, and request statistics
- `GET /healthz` - Health check endpoint (returns "ok")
- `GET /version` - Returns JSON with version information
- `GET /info` - Returns JSON with pod information including uptime and request metrics
- `GET /metrics` - Returns Prometheus-style metrics (request counts, uptime)

## Environment Variables

- `APP_VERSION` - Application version (default: "1.0.0")
- `POD_NAME` - Pod name (default: "unknown", automatically set in Kubernetes)
- `PORT` - Server port (default: "8080")

## Development

### Local Development

1. Ensure Go 1.21+ is installed
2. Install dependencies:

```bash
go mod download
```

3. Run locally:

```bash
go run main.go
```

4. Access at `http://localhost:8080`

## Security

The application implements comprehensive security best practices:

- Kubernetes security contexts (non-root user, dropped capabilities, seccomp profile)
- Minimal attack surface (scratch base image, no shell, no utilities)
- Resource limits to prevent resource exhaustion
- Health checks for reliability

See [SECURITY.md](SECURITY.md) for detailed security documentation.

## Image Size Optimization

The image uses a multi-stage Docker build with aggressive optimizations:

1. Build stage uses `golang:alpine` to compile the static binary
2. Build flags: `-w -s` to strip debug symbols and symbol table
3. `-trimpath` to remove file system paths from binaries
4. `strip --strip-all` to remove additional symbols
5. UPX compression with `--best --lzma` to compress the binary by ~68%
6. Final stage uses `scratch` (empty base image) to minimize size
7. Only the compressed binary is included in the final image

**Current Image Size**: ~2.13 MB

**Optimization Results**:
- Original binary: 6.7 MB
- After strip: 6.4 MB  
- After UPX compression: 2.1 MB (31.66% of original)
- Final image: 2.13 MB

This represents a **68% reduction** from the original size.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
