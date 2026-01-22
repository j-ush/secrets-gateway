# Secrets Gateway

A secure Kubernetes secrets management proxy that reduces direct cloud connections by 99%+ while providing defense-in-depth security layers.

## Overview

The secrets-gateway acts as a centralized proxy for accessing cloud secrets managers (AWS Secrets Manager, etc.) from Kubernetes workloads. It leverages IRSA (IAM Roles for Service Accounts) for secure authentication and provides multiple security layers.

## Security Architecture

### Defense Layers

1. **L1: IRSA Authentication** - Pod identity-based access to cloud secrets
2. **L2: Token Caching** - Efficient credential management with automatic refresh
3. **L3: RBAC Enforcement** - Prevention of ESO modification by unauthorized entities
4. **L4: Audit Logging** - Comprehensive logging for accountability
5. **L5: Network Policies** - Ingress restriction to authorized namespaces only

## Prerequisites

- Kubernetes 1.25+
- AWS EKS with IRSA enabled
- Helm 3.x
- Go 1.21+ (for development)

## Installation

### Using Helm

```bash
helm install secrets-gateway ./deploy/helm \
  --namespace secrets-gateway \
  --create-namespace \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=arn:aws:iam::ACCOUNT:role/secrets-gateway-role
```

### IRSA Setup

1. Create an IAM role with Secrets Manager read access
2. Configure the trust relationship for your EKS cluster's OIDC provider
3. Set the role ARN in the Helm values

## Usage

### For External Secrets Operator

See `examples/level-1-operators/` for integration examples.

### Security Testing

Run the LLM context injection simulation:

```bash
./examples/level-4-llm-agents/context-injection-sim.sh
```

## Evaluation

### Connection Count

```bash
./scripts/connection-count.sh
```

### Latency Benchmark

```bash
./scripts/latency-bench.sh
```

## Development

```bash
# Build
cd cmd/gateway
go build -o secrets-gateway

# Run locally
./secrets-gateway
```

## License

Apache License 2.0 - see [LICENSE](LICENSE)

