# Artifact Appendix

## Abstract

This artifact contains the implementation of the secrets-gateway, a secure proxy for Kubernetes secrets management that reduces direct cloud connections by 99%+ while providing defense-in-depth security layers.

## Artifact Check-List

- **Algorithm:** IRSA-based authentication with token caching
- **Program:** Go 1.21+
- **Compilation:** `go build ./...`
- **Run-time environment:** Kubernetes 1.25+, AWS EKS with IRSA enabled
- **Hardware:** Standard Kubernetes nodes (2 CPU, 4GB RAM minimum per node)
- **Metrics:** Connection count reduction, request latency, security policy enforcement
- **Output:** Audit logs, metrics endpoints
- **Experiments:** Automated via scripts in `/scripts` directory
- **How much disk space required?:** ~100MB for container image
- **How much time is needed to prepare workflow?:** ~30 minutes
- **How much time is needed to complete experiments?:** ~1 hour
- **Publicly available?:** Yes, Apache 2.0 License

## Description

### How to Access

Clone the repository and follow the README.md for setup instructions.

### Hardware Dependencies

- Kubernetes cluster (EKS recommended)
- AWS account with Secrets Manager access

### Software Dependencies

- Go 1.21+
- kubectl
- Helm 3.x
- AWS CLI configured with appropriate permissions

## Installation

```bash
# Build the gateway
cd cmd/gateway
go build -o secrets-gateway

# Deploy with Helm
helm install secrets-gateway ./deploy/helm \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=<your-role-arn>
```

## Experiment Workflow

1. Deploy the secrets-gateway to your cluster
2. Run connection count analysis: `./scripts/connection-count.sh`
3. Run latency benchmarks: `./scripts/latency-bench.sh`
4. Test security layers with example manifests in `/examples`

## Evaluation and Expected Results

- **Connection Reduction:** >99% reduction in direct cloud connections
- **Latency:** <10ms average for cached secrets
- **Security:** All 5 defense layers functional (IRSA, RBAC, audit, network policy, validation)

## Notes

For detailed security architecture and threat model, see the main paper.

