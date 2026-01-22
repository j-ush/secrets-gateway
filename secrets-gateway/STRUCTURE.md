/secrets-gateway
├── .github/workflows/      # CI for building images and running security scans
├── cmd/
│   └── gateway/            # Main entry point for the Go implementation 
├── internal/
│   ├── auth/               # IRSA and token caching logic [cite: 135, 180]
│   ├── proxy/              # Request forwarding and validation logic [cite: 184]
│   └── audit/              # Standardized logging for accountability [cite: 116]
├── deploy/
│   ├── helm/               # Production-ready Helm charts [cite: 375]
│   │   ├── templates/
│   │   │   ├── networkpolicy.yaml  # L5: Ingress restriction [cite: 167]
│   │   │   └── rbac.yaml           # L3: ESO modification prevention [cite: 158]
│   │   └── values.yaml
│   └── manifests/          # Raw YAML for ClusterSecretStores [cite: 140, 149]
├── examples/
│   ├── level-1-operators/  # ESO and cert-manager examples [cite: 49]
│   └── level-4-llm-agents/ # Simulation scripts for LLM context injection [cite: 110, 305]
├── scripts/
│   ├── connection-count.sh # Tools to measure 99%+ reduction [cite: 377]
│   └── latency-bench.sh    # Evaluation scripts for polling delay [cite: 377]
├── ARTIFACT_APPENDIX.md    # Required for USENIX Artifact Evaluation 
├── LICENSE                 # Apache License 2.0 [cite: 390]
└── README.md               # Setup, IRSA requirements, and security claims [cite: 381]
