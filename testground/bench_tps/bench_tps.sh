#!/usr/bin/env bash

set -euo pipefail

CONFIG=bench-tps-config.yml
LEDGER_DIR="test-ledger"

clean () {
  rm -rf "$LEDGER_DIR" "$CONFIG"
}

run () {
  # sudo ifconfig lo0 alias 127.0.0.1 up

  trap 'trap - SIGTERM && kill -- -$$' SIGINT SIGTERM EXIT

  tcpdump -i lo0 udp dst port 8003 -w tpu.pcap &

  solana-test-validator                \
    --ledger             "$LEDGER_DIR" \
    --bind-address       127.0.0.1     \
    --dynamic-port-range 8000-8200     \
    --rpc-port           8899          \
    --faucet-port        9900          \
    --faucet-sol         1000000000    \
    --reset                            \
    --quiet                            \
    &

  cat <<EOF > "$CONFIG"
json_rpc_url: http://127.0.0.1:8899
websocket_url: ws://127.0.0.1:8900
keypair_path: "$(pwd)/$LEDGER_DIR/faucet-keypair.json"
EOF

  until
    solana -C "$CONFIG" slot
  do sleep 0.1
  done

  solana-bench-tps                     \
    --config           "$CONFIG"       \
    --entrypoint        127.0.0.1:8000 \
    --faucet            127.0.0.1:9900 \
    --num-nodes         1              \
    --threads           4              \
    --tx_count          50000          \
    --duration          20             \
    --use-tpu-client                   \
    --tpu-disable-quic

  sleep 10

  exit 0
}

case "$1" in
  clean)
    clean
    ;;
  run)
    run
    ;;
  *)
    echo "Usage: $0 {run <duration>|clean}"
    exit 1
esac
