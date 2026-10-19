# KeyValue
**A distributed, in-memory key-value database written in Go (inspired by Redis).**

![Basic Structure](image.png)

## Features & Roadmap

**What Works:**
- [x] TCP-based client & peer-to-peer API
- [x] Concurrent request handling
- [x] Custom RESP2 protocol parser
- [x] Cluster state managed via consistent hashing (Hashring)

**Roadmap / To-Do:**
- [ ] Graceful node departure from the cluster
- [ ] Heartbeat / Keepalive checks between nodes
- [ ] Marking cluster as "down" when a node goes offline *(Prioritizing Consistency over Availability / CP model)*
- [ ] Virtual nodes on the Hashring (groundwork is implemented)
- [ ] Concurrency optimizations
- [ ] Dedicated CLI client for sending commands

---

## RESP2 Protocol Implementation

### Data Types
Based on the official [Redis Protocol Specification](https://redis.io/docs/latest/develop/reference/protocol-spec/):
- **RespString**: Corresponds to RESP Bulk String
- **RespInteger**: Corresponds to RESP Integer
- **RespArray**: Corresponds to RESP Array
- **RespError**: Corresponds to RESP Error
- **RespNull**: `RespString` of length `-1`

### Supported Commands

#### `GET [keys...]`
- Reads values of `<keys>`
- **Returns:** `RespArray` of corresponding values, or `RespError`

#### `SET <key> <value>`
- Sets `<key>` to `<value>`
- **Returns:** `"Success!"` (`RespString`), or `RespError`

#### `ADDSERVER <name> <address>`
- Connects to `<address>`, saves the connection in the cluster under `<name>`, and initiates a `PING` handshake.
- **Returns:** `"Success!"` (`RespString`), or `RespError`

#### `PING <name> <port>`
- Registers the incoming node in the local cluster state without echoing a `PING` back (acts as the receiver for `ADDSERVER`).

---

## Building & Running

### Option 1: Docker Compose (Recommended for Cluster Testing)
Spin up two independent nodes (`node1` and `node2`) on the same Docker network:
```bash
docker compose up --build