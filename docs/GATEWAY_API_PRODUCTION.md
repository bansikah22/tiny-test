# Gateway API Production Setup Guide

This guide explains how to deploy the tiny-test application using Gateway API in production environments.

## Tested With

Tested with Gateway API v1beta1 controllers such as Envoy Gateway and Istio. Verify whether your controller supports v1beta1 or v1 before applying manifests, and adjust `apiVersion` accordingly.

## What is Gateway API?

Gateway API is the modern Kubernetes standard for traffic routing, replacing legacy Ingress. It provides:

- More expressive routing capabilities
- Better separation of concerns (infrastructure vs. application teams)
- Support for multiple protocols (HTTP, TLS, TCP, UDP)
- Vendor-neutral API across cloud providers

## Prerequisites

1. Kubernetes cluster (1.24+)
2. Gateway API CRDs installed
3. Gateway API controller/implementation installed

## Supported Gateway API Implementations

Choose one based on your platform:

### Cloud Providers

- **AWS**: AWS Gateway API Controller
- **Azure**: Azure Application Gateway for Containers
- **GCP**: GKE Gateway Controller (built-in)

### Open Source

- **Envoy Gateway**: General-purpose, high-performance
- **Istio**: Full service mesh with Gateway API support
- **Contour**: VMware-maintained Envoy-based controller
- **Traefik**: Simple and feature-rich

## Installation Steps

### Step 1: Install Gateway API CRDs (if not already installed)

```bash
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml
```

### Step 2: Install a Gateway API Controller

Example with Envoy Gateway:

```bash
kubectl apply -f https://github.com/envoyproxy/gateway/releases/latest/download/install.yaml
```

Wait for the controller to be ready before applying Gateway/HTTPRoute:

```bash
kubectl wait --timeout=5m -n envoy-gateway-system deployment/envoy-gateway --for=condition=Available
```

### Step 3: Verify GatewayClass

Check that your GatewayClass is available:

```bash
kubectl get gatewayclass
```

Example output:
```
NAME   CONTROLLER                      AGE
eg     gateway.envoyproxy.io/gateway   5m
```

### Step 4: Deploy the Application

Deploy tiny-test using Helm or manifests:

```bash
# With Helm
helm upgrade --install tiny-test ./helm/tiny-test

# Or with raw manifests
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### Step 5: Create Gateway and HTTPRoute

Update `k8s/gateway.yaml` with your GatewayClass name:

```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: tiny-test-gateway
spec:
  gatewayClassName: eg  # Change to your GatewayClass name
  listeners:
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes:
        namespaces:
          from: Same
```

Apply the Gateway and HTTPRoute:

```bash
kubectl apply -f k8s/gateway.yaml
kubectl apply -f k8s/httproute.yaml
```

### Step 6: Verify Gateway Status

Check that the Gateway is accepted and programmed:

```bash
kubectl get gateway tiny-test-gateway
```

Expected output:
```
NAME                CLASS   ADDRESS         PROGRAMMED   AGE
tiny-test-gateway   eg      10.96.xxx.xxx   True         30s
```

### Step 7: Get External Address

For cloud load balancers:

```bash
kubectl get gateway tiny-test-gateway -o jsonpath='{.status.addresses[0].value}'
```

### Step 8: Configure DNS

Point your domain to the Gateway's external IP:

```bash
# Example: Create A record
tiny-test.example.com  A  <GATEWAY-EXTERNAL-IP>
```

### Step 9: Test Access

```bash
curl http://tiny-test.example.com/healthz
```

Expected output:
```
ok
```

## Platform-Specific Examples

### GKE (Google Kubernetes Engine)

```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: tiny-test-gateway
spec:
  gatewayClassName: gke-l7-global-external-managed
  listeners:
    - name: http
      protocol: HTTP
      port: 80
```

### Azure AKS

```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: tiny-test-gateway
spec:
  gatewayClassName: azure-application-gateway
  listeners:
    - name: http
      protocol: HTTP
      port: 80
```

### AWS EKS

```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: tiny-test-gateway
spec:
  gatewayClassName: amazon-vpc-lattice
  listeners:
    - name: http
      protocol: HTTP
      port: 80
```

## Troubleshooting

### Gateway Not Programmed

Check Gateway status:
```bash
kubectl describe gateway tiny-test-gateway
```

Look for conditions and events.

### HTTPRoute Not Attached

Check HTTPRoute status:
```bash
kubectl get httproute tiny-test-route -o yaml
```

Look at `status.parents` for acceptance conditions.

### Service Not Found

Ensure the backend service exists:
```bash
kubectl get svc tiny-test
```

### No External IP

For cloud providers, ensure:
- Correct GatewayClass name
- Cloud controller manager is running
- Sufficient IAM permissions

## Local Testing (Minikube/kind)

For local development, Gateway API requires a controller installation. However, this can be complex for local testing. Instead, use:

**Option 1: Port-forward**
```bash
kubectl port-forward service/tiny-test 8080:80
```

**Option 2: NodePort**
```bash
kubectl patch service tiny-test -p '{"spec":{"type":"NodePort"}}'
minikube service tiny-test --url
```

**Option 3: Ingress (for local testing)**
```bash
minikube addons enable ingress
kubectl apply -f k8s/ingress-example.yaml  # Traditional Ingress for local dev
```

## Migration from NGINX Ingress

If migrating from NGINX Ingress:

1. Keep existing Ingress resources running
2. Install Gateway API controller
3. Create Gateway and HTTPRoute resources
4. Test Gateway API routing
5. Update DNS to point to Gateway
6. Verify traffic flow
7. Delete old Ingress resources

## Security Considerations

- Use HTTPS with TLS termination in production
- Configure rate limiting and authentication policies
- Apply network policies to restrict traffic
- Use cert-manager for automatic certificate management

## Next Steps

- Add HTTPS/TLS support with cert-manager
- Configure rate limiting policies
- Set up monitoring and observability
- Implement blue-green or canary deployments
