#!/bin/bash
set -euo pipefail

cd /app
export DEFRA_KEYRING_SECRET="$(head -c 32 /dev/urandom | base32)"
defradb keyring generate --rootdir ./.nodeA
# start defradb in the background
defradb start --rootdir ./.nodeA --url localhost:9181 --p2paddr /ip4/0.0.0.0/tcp/9171 &
# wait for defradb to start
sleep 2

defradb client schema add --url localhost:9181 '
  type Article {
    content: String
    published: Boolean
  }
'
echo "Defradb nodeA started and schema added."
wait