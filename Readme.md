# KeyValue
### A distributed in-memory K/V database (basically a redis clone)
---

What works
 - API Via TCP Connections
 - Concurrent Request Handling
 - Custom Resp2 parser
 - Cluster Managed via Hashring

What doesn't work / What still needs to be implemented
 - Nodes leaving the cluster
 - Keepalive between nodes
 - Marking cluster as "down" when one node goes offline
  - I chose Partition Tolerance over Availability simply because it is easier to implement
 - Adding virtual nodes to the Hashring (Groundwork is laid)
 - Optimizations regarding concurrency
 - Tool to make sending commands easier

---

![Basic Structure](image.png)

# RESP Parser Implementation
## Data Types, See spec [here](https://redis.io/docs/latest/develop/reference/protocol-spec/)
   - **RespString**:  Corresponds to Resp Bulk String
   - **RespInteger**: Corresponds to Resp Integer
   - **RespArray**:   Corresponds to Resp Array
   - **RespError**:   Corresponds to Resp Error
   - **RespNull**:    RespString, but of length -1

## Request Types:
### GET \[keys\]
   - Reads values of \<keys\>
   - Returns a RespArray of corresponding values or the error
### SET \<key\> \<value\>
   - Sets \<key\> to \<value\>
   - Returns "Success!" RespString or the error
### ADDSERVER \<name\> \<address\>
   - Opens a connection to \<address\>, if successful saves connection in cluster under \<name\> and sends PING Request to \<address\>
   - Returns "Success!" RespString or the error
### PING \<name\> \<port\>
   - Same as ADDSERVER, but does not send PING

---

## Building
### There are multiple ways provided to build this project.

**Nix**<br>
This repository provides a nix flake with outputs
 - Build with `nix build`. The resulting binary will be in ./result/bin/keyvalue
 - Or run directly with `nix run`, Target a specific architecture with `nix run .#packages.<architecture>.keyvalue`

**Docker**<br>
This respository also provides a Dockerfile
1. Build the image `docker build . -t keyvalue:latest`
2. Run the image `docker run -p 8080:8080 keyvalue`

Or run `docker compose up --build` to spin up two independent nodes (node1, node2) that sit in the same docker network. <br>
Within that network node1's address is node1, listens on port 8080 and node2's address is node2, listens on port 8080 <br>
node1 is exposed to the host network on port 8080, node2 on port 8081 <br>
Connect them by sending an ADDSERVER Request 
- `printf "*3\r\n\$9\r\nADDSERVER\r\n\$5\r\nnode1\r\n\$10\r\nnode1:8080\r\n" | nc localhost 8081`
- (or `printf "*3\r\n\$9\r\nADDSERVER\r\n\$5\r\nnode2\r\n\$10\r\nnode2:8080\r\n" | nc localhost 8080`)

Lastly you can also just build the project normally via go's compiler <br>
`go build -o keyvalue ./cmd/server/`<br>
make it executable:<br>
`chmod +x ./keyvalue`<br>
run it, and pass a name of the node as the first and only argument:<br>
`./keyvalue <name>`
