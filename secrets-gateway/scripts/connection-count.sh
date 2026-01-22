#!/bin/bash
# Tools to measure 99%+ reduction in direct cloud connections
# Evaluates the connection pooling efficiency of the secrets-gateway

set -e

echo "=== Connection Count Analysis ==="
echo "Measuring connections to AWS Secrets Manager"

# Count active connections from the gateway pod
GATEWAY_POD=$(kubectl get pods -l app=secrets-gateway -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -n "$GATEWAY_POD" ]; then
    echo ""
    echo "Gateway pod: $GATEWAY_POD"
    
    # Count established connections to port 443 (AWS endpoints)
    CONN_COUNT=$(kubectl exec "$GATEWAY_POD" -- netstat -an 2>/dev/null | grep -c ":443.*ESTABLISHED" || echo "0")
    echo "Active connections to cloud: $CONN_COUNT"
    
    # Count total pods that would otherwise need connections
    TOTAL_PODS=$(kubectl get pods --all-namespaces -o json | jq '.items | length')
    echo "Total pods in cluster: $TOTAL_PODS"
    
    # Calculate reduction
    if [ "$TOTAL_PODS" -gt 0 ]; then
        REDUCTION=$(echo "scale=2; (1 - $CONN_COUNT / $TOTAL_PODS) * 100" | bc)
        echo "Connection reduction: ${REDUCTION}%"
    fi
else
    echo "No gateway pod found. Run this script in a cluster with secrets-gateway deployed."
    echo ""
    echo "Expected metrics:"
    echo "  - Without gateway: 1 connection per pod needing secrets"
    echo "  - With gateway: 1-2 pooled connections total"
    echo "  - Reduction: >99% for clusters with 100+ pods"
fi

echo ""
echo "=== Analysis Complete ==="

