#!/bin/bash
# Evaluation scripts for polling delay measurements
# Benchmarks the latency impact of the secrets-gateway proxy

set -e

echo "=== Latency Benchmark ==="
echo "Measuring request latency through secrets-gateway"

GATEWAY_URL="${GATEWAY_URL:-http://secrets-gateway:8080}"
ITERATIONS="${ITERATIONS:-100}"
SECRET_PATH="${SECRET_PATH:-/v1/secrets/benchmark-test}"

echo ""
echo "Configuration:"
echo "  Gateway URL: $GATEWAY_URL"
echo "  Iterations: $ITERATIONS"
echo "  Secret Path: $SECRET_PATH"
echo ""

# Warm up
echo "Warming up..."
for i in {1..5}; do
    curl -s -o /dev/null "$GATEWAY_URL$SECRET_PATH" 2>/dev/null || true
done

# Run benchmark
echo "Running benchmark..."
TOTAL_TIME=0
MIN_TIME=999999
MAX_TIME=0

for i in $(seq 1 $ITERATIONS); do
    START=$(date +%s%N)
    curl -s -o /dev/null "$GATEWAY_URL$SECRET_PATH" 2>/dev/null || true
    END=$(date +%s%N)
    
    ELAPSED=$(( (END - START) / 1000000 ))  # Convert to milliseconds
    TOTAL_TIME=$((TOTAL_TIME + ELAPSED))
    
    if [ "$ELAPSED" -lt "$MIN_TIME" ]; then MIN_TIME=$ELAPSED; fi
    if [ "$ELAPSED" -gt "$MAX_TIME" ]; then MAX_TIME=$ELAPSED; fi
    
    # Progress indicator
    if [ $((i % 10)) -eq 0 ]; then
        echo "  Completed $i/$ITERATIONS iterations..."
    fi
done

AVG_TIME=$((TOTAL_TIME / ITERATIONS))

echo ""
echo "=== Results ==="
echo "  Total requests: $ITERATIONS"
echo "  Average latency: ${AVG_TIME}ms"
echo "  Min latency: ${MIN_TIME}ms"
echo "  Max latency: ${MAX_TIME}ms"
echo ""
echo "Target: <10ms average for cached secrets"
echo "=== Benchmark Complete ==="

