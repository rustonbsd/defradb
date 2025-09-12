#!/bin/bash
set -euo pipefail

cd /app
export DEFRA_KEYRING_SECRET="$(head -c 32 /dev/urandom | base32)"
defradb keyring generate --rootdir ./.nodeB
defradb start --rootdir ./.nodeB --url localhost:9182 --p2paddr /ip4/0.0.0.0/tcp/9172 &
sleep 2

defradb client schema add --url localhost:9182 '
  type Article {
    content: String
    published: Boolean
  }
'

# get the libp2p addes json and format into ./defradb client p2p replicator set -c Article 'JSON_HERE'
# use bash var to store the json:

#>WARNING:(ast) sonic only supports go1.17~1.23, but your environment is not suitable
#>------ Request Results ------
#> {
#>   "ID": "12D3KooWSpwwW74TNFoPMrUitj87jYNgK3BQY8z6meYRXMy3yhpA",
#>   "Addresses": [
#>     "/ip4/127.0.0.1/tcp/9172",
#>     "/ip4/192.168.178.31/tcp/9172"
#>   ]
#> }
P2P_ADDR_JSON=$(defradb client p2p info --url localhost:9182 | jq -c .)
echo ""
echo "sudo docker compose exec nodeA defradb client p2p connect '$P2P_ADDR_JSON'"
echo "sudo docker compose exec nodeA defradb client p2p replicator set -c Article '$P2P_ADDR_JSON'"
wait