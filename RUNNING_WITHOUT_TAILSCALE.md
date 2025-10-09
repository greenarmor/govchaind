# Running a GovChain Node without Tailscale

This guide explains how to run a `govchaind` node using Docker without the Tailscale sidecar. This is useful for local development or for environments where you are managing network exposure manually.

## Prerequisites

- **Docker**: Ensure Docker is installed and running on your system.
- **GovChain Docker Image**: You must have the `govchaind` Docker image available locally. You can build it from the source using `make docker-build` or pull it from the GitHub Container Registry.

---

## Option 1: Using a Simplified Docker Compose File

This is the recommended method as it is easy to manage. Create a file named `docker-compose.standalone.yml` with the following content:

```yaml
version: '3.8'

services:
  govchaind:
    image: ghcr.io/bettergovph/govchaind:latest # Or your locally built image, e.g., govchaind:latest
    container_name: govchaind-node
    volumes:
      - govchaind-data:/home/nonroot/.govchain
    environment:
      - MONIKER="My Standalone Node"
    ports:
      - "26656:26656" # P2P port
      - "26657:26657" # RPC port
      - "1317:1317"    # REST API port
      - "9090:9090"    # gRPC port
    networks:
      - govchain_network

networks:
  govchain_network:
    driver: bridge

volumes:
  govchaind-data:
```

### How to Use

1.  **Save the File**: Save the content above as `docker-compose.standalone.yml` in the root of the project.

2.  **Start the Node**:
    ```bash
    docker-compose -f docker-compose.standalone.yml up -d
    ```

3.  **Monitor Logs**:
    ```bash
    docker logs -f govchaind-node
    ```

4.  **Stop the Node**:
    ```bash
    docker-compose -f docker-compose.standalone.yml down
    ```

---

## Option 2: Using `docker run`

If you prefer not to use Docker Compose, you can run the container directly with the `docker run` command.

### Command

```bash
docker run -d --name govchaind-node \
  -p 26656:26656 \
  -p 26657:26657 \
  -p 1317:1317 \
  -p 9090:9090 \
  -v govchaind-data:/home/nonroot/.govchain \
  -e MONIKER="My Standalone Node" \
  ghcr.io/bettergovph/govchaind:latest
```

### Explanation

- `-d`: Runs the container in detached mode.
- `--name govchaind-node`: Assigns a name to the container for easy reference.
- `-p ...`: Maps the necessary ports from your host machine to the container.
- `-v govchaind-data:/home/nonroot/.govchain`: Creates a named volume `govchaind-data` to persist blockchain data across container restarts.
- `-e MONIKER="..."`: Sets the moniker for your node.
- `ghcr.io/bettergovph/govchaind:latest`: The Docker image to run.

---

## How It Works: Automatic Public IP Detection

The `docker-entrypoint.sh` script inside the container is designed to be resilient and versatile.

1.  It first checks for a Tailscale sidecar. If a Tailscale IP is found, it will be used as the `external_address`.
2.  If no Tailscale IP is detected, the script automatically attempts to discover the host machine's public IP address using an external service.
3.  If a public IP is found, it configures `external_address` in `config.toml` with that IP (e.g., `external_address = "YOUR_PUBLIC_IP:26656"`). This allows your node to correctly advertise itself to the network, which is essential for VPS or public cloud deployments.
4.  If neither IP can be determined, the node will start without a pre-configured `external_address`.
