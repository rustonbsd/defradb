# Iroh Transport Integration Test (Experimental)

<!--
Integration Test: Cold Bootstrap & Replication Behind NAT (Docker)

This README documents an integration test that demonstrates DefraDB nodes achieving cold bootstrap and successful replication inside network isolated Docker containers, enabled solely by a minimal addition in net/host.go to incorporate the Iroh libp2p transport crate written for this proof of concept. For now "iroh transport" pretentds to be a tcp capable transport in libp2p. 

Scope:
- Show that a single, minimal transport-layer change is sufficient for replication between fresh (cold) nodes.
- Emphasize behavior in NATed/containerized conditions where implicit connectivity is non-trivial.
- Provide a reproducible scaffold others can extend into a formal test harness (test.sh to be added).

To Be Detailed (you will fill in next; I will then complete):
1. Precise summary of the net/host.go modification (what was added and why it is minimal).
2. Flow of replication-driven peer connectivity (command sequence, control plane vs data plane).
3. Container and network topology (number of nodes, networks, ports, isolation notes).
4. Step-by-step execution outline (startup, replication trigger, observation points).
5. Verification methodology (log signatures, peer IDs, data/state convergence checks).
6. Failure modes / limitations (what this test does not cover yet).
7. Cleanup and rerun considerations (ensuring true cold starts).

Assumptions:
- Go dependencies (including the Iroh transport) resolve without manual intervention.
- No persisted libp2p state between runs (ephemeral containers / volumes).
- No auxiliary discovery services (DHT, rendezvous, static bootstrap list) are preconfigured.

Artifacts to be added:
- test.sh (or equivalent) orchestrating build, container launch, replication invocation, and assertions.
- Optional docker-compose.yml (or inline docker CLI instructions).
- Sample log excerpts demonstrating first-contact connection establishment.
- Data validation snippet confirming replication success (e.g., querying replicated documents).

Next Action:
Provide the concrete diff / description of the host modification and the planned command sequence; I will integrate them into the appropriate sections and finalize prose.
-->
(Work in progress)

This directory will contain an integration test demonstrating cold bootstrap and replication across DefraDB nodes (inside Docker, behind NAT) enabled by a minimal addition of the Iroh libp2p transport in `net/host.go`. Current understanding: peer awareness/connection for the test scenario will originate from invoking the replication command (not from prior routing table state).

## Explanation (To Be Completed)

A detailed walkthrough will be added here covering:
- What minimal code change was introduced.
- How replication establishes peer connectivity.
- Expected container/network topology.
- Steps to run and observe successful replication behind NAT.
- Verification methodology (cold start, no prior peerstore).

(You can prompt me to expand this section when ready.)

## Status

Scaffolding only. Scripts, commands, and deep explanation pending.

## Disclaimer

Experimental; not part of the standard test suite yet.
