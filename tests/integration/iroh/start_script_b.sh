#!/bin/bash
set -euo pipefail

cd /app
export DEFRA_KEYRING_SECRET="$(head -c 32 /dev/urandom | base32)"
defradb keyring generate --rootdir ./.nodeB
RUST_LOG=error,irohffi::ffi=debug defradb start --rootdir ./.nodeB --url localhost:9182 --p2paddr /ip4/0.0.0.0/tcp/9172 &
sleep 2

defradb client schema add --url localhost:9182 '
  type Article {
    content: String
    published: Boolean
  }
'
P2P_ADDR_JSON=$(defradb client p2p info --url localhost:9182 | jq -c .)
echo ""
echo "sudo docker compose exec nodeA defradb client p2p connect '$P2P_ADDR_JSON'"
echo "sudo docker compose exec nodeA defradb client p2p replicator set -c Article '$P2P_ADDR_JSON'"
wait