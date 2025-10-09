# iroh Integration Test for DefraDB

This integration test demonstrates that the [iroh](https://github.com/n0-computer/iroh) transport layer is fully functional with DefraDB's peer-to-peer replication system. It creates two isolated Docker containers, each running a DefraDB node, and verifies that they can communicate and replicate data using the iroh transport layer.

## What This Test Does

The test creates a complete end-to-end demonstration of DefraDB's P2P replication over iroh:

1. **Builds two isolated containers** (`nodeA` and `nodeB`) with DefraDB (with the build in [go-libp2p-iroh-transport](https://github.com/rustonbsd/go-libp2p-iroh-transport/tree/v0.1.1) go<>rust ffi bindings)
2. **Starts DefraDB instances** on separate networks with different ports
3. **Creates identical schemas** on both nodes (an `Article` collection)
4. **Establishes P2P connection** between nodes using iroh multiaddresses (`/iroh/[32-byte-z32-encoded-ed25519-pub-key]`)
5. **Sets up replication** for the Article collection
6. **Verifies data sync** by creating documents on one node and querying them from the other

## Architecture

iroh uses a sophisticated networking stack built on **QUIC** with **relay servers** and **NAT hole-punching** to establish peer-to-peer connections. The system supports **QUIC Multipath**, allowing connections to use multiple network paths simultaneously for resilience and performance.

### How iroh Connections Work

1. **Initial Connection via Relay**: Nodes first connect through a relay server using QUIC over UDP
2. **NAT Traversal (Hole Punching)**: Nodes coordinate through the relay to establish direct connections
3. **Multipath QUIC**: Once connected, traffic can flow over multiple paths:
   - Direct connection (after successful hole punching)
   - Relay path (fallback when direct connection fails)
4. **Connection Healing**: If network changes (wifi outage, switching to cellular, etc.), connections automatically migrate to new paths

### Test Setup

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Container: nodeA                               │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  DefraDB Instance                                                     │  │
│  │  - HTTP API: localhost:9181                                           │  │
│  │  - P2P Port: 0.0.0.0:9171                                             │  │
│  │  - Iroh Addr: /iroh/[nodeA-32byte-ed25519-pubkey]                     │  │
│  │  - Network: defra-net-a                                               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  iroh Transport Layer                                                 │  │
│  │  - QUIC Protocol                                                      │  │
│  │  - NAT Hole Punching                                                  │  │
│  │  - Multipath Support                                                  │  │
│  │  - libiroh.so (Rust FFI)                                              │  │
│  │  - go-libp2p-iroh-transport                                           │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       │ Direct QUIC Connection
                                       │ (after hole punch)
                                       │ or via relay fallback
                                       ↕
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Container: nodeB                               │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  DefraDB Instance                                                     │  │
│  │  - HTTP API: localhost:9182                                           │  │
│  │  - P2P Port: 0.0.0.0:9172                                             │  │
│  │  - Iroh Addr: /iroh/[nodeB-32byte-ed25519-pubkey]                     │  │
│  │  - Network: defra-net-b                                               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  iroh Transport Layer                                                 │  │
│  │  - QUIC Protocol                                                      │  │
│  │  - NAT Hole Punching                                                  │  │
│  │  - Multipath Support                                                  │  │
│  │  - libiroh.so (Rust FFI)                                              │  │
│  │  - go-libp2p-iroh-transport                                           │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```
## Quick Start

### 1. Build the Containers

```bash
cd tests/integration/iroh
make build
```

This command:
- Builds the DefraDB binary with iroh support
- Creates two Docker images (`nodeA` and `nodeB`)
- Copies the required libraries (`libwasmer.so` and `libiroh.so`)

### 2. Start the Nodes

```bash
make run
```

This will start both containers. Watch the console output: `nodeB` will display connection commands that you need to run manually.

### 3. Connect the Nodes

After `make run`, you'll see output similar to this from `nodeB`:

```
nodeB       | WARNING:(ast) sonic only supports go1.17~1.23, but your environment is not suitable
nodeB       | Oct  9 10:29:33.203 INF http Request Method=GET Path=/api/v0/p2p/info Status=200 LengthBytes=190 ElapsedTime=126.031µs
nodeB       |
nodeB       | sudo docker compose exec nodeA defradb client p2p connect '{"ID":"12D3KooWFCtjB6nSjdFE1W7jhRZ4MpwkAJMtzvHdBMJUABt6CzYv","Addresses":["/ip4/127.0.0.1/tcp/9172","/ip4/172.23.0.2/tcp/9172","/iroh/kah4icli3ck7g4n5y7wnocxqrqnxmiut7ukrsht6ge4zuvewwytq"]}'
nodeB       | sudo docker compose exec nodeA defradb client p2p replicator set -c Article '{"ID":"12D3KooWFCtjB6nSjdFE1W7jhRZ4MpwkAJMtzvHdBMJUABt6CzYv","Addresses":["/ip4/127.0.0.1/tcp/9172","/ip4/172.23.0.2/tcp/9172","/iroh/kah4icli3ck7g4n5y7wnocxqrqnxmiut7ukrsht6ge4zuvewwytq"]}'
```

**Copy and run both commands** from the output (the first one is just a formality). Note the added `/iroh/...` multiaddress is availale and the only possible path to replicate between a and b since they don't share any other network. 

#### Connect Command
```bash
sudo docker compose exec nodeA defradb client p2p connect '<JSON_FROM_OUTPUT>'
```

#### Set Replicator Command
```bash
sudo docker compose exec nodeA defradb client p2p replicator set -c Article '<JSON_FROM_OUTPUT>'
```

After running these commands, you should see connection logs indicating successful iroh connection:

```
nodeA       | 2025-10-09T10:29:53.301709Z  INFO connect{me=86393509ea alpn="/libp2p/iroh/0.1.0" remote=500fc40968}:discovery{me=86393509ea node=500fc40968}:add_
```

The `alpn="/libp2p/iroh/0.1.0"` confirms the iroh protocol is active.

### 4. Test Replication

Now verify that data replicates between nodes using the iroh transport:

#### Add a document on Node A
```bash
make graphql-a-add
```

Expected output:
```json
{
  "data": {
    "create_Article": {
      "_docID": "bae-...",
      "content": "Hello world",
      "published": true
    }
  }
}
```

#### Query from Node A (should see the document)
```bash
make graphql-a-get
```

#### Query from Node B (should see the replicated document via iroh)
```bash
make graphql-b-get
```

Expected output on Node B:
```json
{
  "data": {
    "Article": [
      {
        "_docID": "bae-...",
        "content": "Hello world",
        "published": true
      }
    ]
  }
}
```

**If Node B returns the document created on Node A, iroh replication is working!**

## What This Proves

This integration test demonstrates:

1. **iroh Transport Integration**: DefraDB successfully uses the iroh transport layer (visible in multiaddresses with `/iroh/...` prefix and with debug logs the read write chatter)
2. **P2P Discovery**: Nodes can discover each other using iroh's peer discovery
3. **Data Replication**: Documents created on one node are successfully replicated to another node via iroh
4. **Schema Synchronization**: Both nodes maintain consistent schema definitions
5. **End-to-End Functionality**: The complete stack (DefraDB → go-libp2p-iroh-transport → libiroh.so → iroh protocol) works correctly

## Additional Testing: go-libp2p-iroh-transport Module

The underlying Go module that enables iroh transport in DefraDB is [`go-libp2p-iroh-transport`](https://github.com/rustonbsd/go-libp2p-iroh-transport). This module has its own e2e test.

To test the transport module independently:

```bash
# Clone the repository
git clone https://github.com/rustonbsd/go-libp2p-iroh-transport
cd go-libp2p-iroh-transport

go test .
```

## File Structure

```
tests/integration/iroh/
├── README.md             # This file
├── Makefile              # Build and test automation
├── docker-compose.yml    # Container orchestration
├── Dockerfile            # Multi-stage build for both nodes
├── start_script_a.sh     # Node A initialization script
└── start_script_b.sh     # Node B initialization script
```

## How It Works

### Build Process

1. **DefraDB Build** (`defrabuild` service):
   - Builds DefraDB with iroh support enabled
   - Copies `libiroh.so` from the go-libp2p-iroh-transport module
   - Copies `libwasmer.so` for WASM lens support
   - Creates the `defra:latest` image

2. **Node Images** (`nodeA` and `nodeB`):
   - Based on `golang:1.24` runtime
   - Copies DefraDB binary and required libraries
   - Installs `jq` for JSON processing
   - Sets up separate initialization scripts

### Runtime Process

**Node A** (`start_script_a.sh`):
1. Generates keyring with random secret
2. Starts DefraDB on port 9181 (HTTP) and 9171 (P2P, iroh unrelated)
3. Adds the Article schema
4. Waits for connections

**Node B** (`start_script_b.sh`):
1. Generates keyring with random secret
2. Starts DefraDB on port 9182 (HTTP) and 9172 (P2P, iroh unrelated)
3. Adds the Article schema
4. Retrieves its P2P info (including iroh address)
5. Outputs connection commands for Node A
6. Waits for connections

### Network Isolation

Each node runs on its own Docker network (`defra-net-a` and `defra-net-b`) to simulate real-world P2P scenarios where nodes are not on the same local network. They communicate exclusively through the iroh transport layer.

## For Maintainers

### Taking Over This Test

If you're maintaining this test:

1. **STack**:
   - DefraDB (this repository)
   - go-libp2p-iroh-transport (Go bindings)
    - libiroh.so (Rust FFI library)
   - iroh protocol (Rust implementation)

2. **Key files**:
   - `go.mod`: Version of `go-libp2p-iroh-transport`
   - `tools/defradb.containerfile`: Library paths and build process
   - `net/host.go`: iroh transport initialization

3. **Common issues**:
   - CGO linking errors: Check `CGO_LDFLAGS` in containerfile

### Testing Checklist

Before merging changes:

- [x] `make build` completes without errors
- [x] Both nodes start successfully
- [x] Connection commands include `/iroh/...` addresses
- [x] Nodes connect successfully (check logs for iroh ALPN)
- [x] Documents replicate from Node A to Node B
- [x] Documents replicate from Node B to Node A
- [x] `go test` passes in go-libp2p-iroh-transport repository

## References

- **go-libp2p-iroh-transport**: https://github.com/rustonbsd/go-libp2p-iroh-transport
