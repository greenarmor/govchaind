# Dockerized govchaind Node with Tailscale Support

This document provides instructions for building and running a Dockerized `govchaind` node with integrated Tailscale support, allowing for secure networking behind NAT.

## 1. Prerequisites

-   Docker installed and running on your system.
-   A Tailscale account and a valid authentication key (`TS_AUTHKEY`) if you plan to use automatic Tailscale login.

## 2. Building the Docker Image

First, navigate to the root directory of the `govchaind` project where the `Dockerfile` is located.

```bash
docker build -t govchaind:latest .
```

This command will build the Docker image and tag it as `govchaind:latest`.

## 3. Running the Docker Container

To run the `govchaind` node with Tailscale, you need to provide specific Docker flags and environment variables.

**Important:** The container needs `NET_ADMIN` capability and access to the `/dev/net/tun` device for Tailscale to function correctly.

```bash
docker run -d \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -p 26656:26656 \
  -p 26657:26657 \
  --name govchaind-node \
  -v govchaind-data:/home/nonroot/.govchain \
  -e MONIKER="my-node" \
  -e CHAIN_ID="govchain-testnet" \
  -e TS_AUTHKEY="tskey-auth-YOUR_ACTUAL_TAILSCALE_AUTH_KEY" \
  govchaind:latest
```

**Replace `tskey-auth-YOUR_ACTUAL_TAILSCALE_AUTH_KEY` with your actual Tailscale authentication key.**

### Environment Variables:

-   `MONIKER`: The moniker (name) for your `govchaind` node.
-   `CHAIN_ID`: The chain ID for your `govchaind` node.
-   `TS_AUTHKEY`: (Optional) Your Tailscale authentication key for automatic login. If not provided, you will need to log in manually.

### Docker Flags:

-   `-d`: Runs the container in detached mode (in the background).
-   `--cap-add=NET_ADMIN`: Grants the container network administration capabilities, required by Tailscale.
-   `--device=/dev/net/tun`: Provides access to the TUN device, essential for WireGuard (used by Tailscale).
-   `-p 26656:26656`: Maps the Cosmos P2P port from the container to the host.
-   `-p 26657:26657`: Maps the Cosmos RPC port from the container to the host.
-   `--name govchaind-node`: Assigns a name to your container for easy reference.
-   `-v govchaind-data:/home/nonroot/.govchain`: Mounts a Docker volume for persistent storage of `govchaind` data and configuration.

## 4. Tailscale Login (Manual)

If you did not provide a `TS_AUTHKEY` during `docker run`, you will need to log in manually:

1.  Get the container name (if you didn't use `--name`):
    ```bash
docker ps
    ```
2.  Execute the `tailscale up` command inside the container:
    ```bash
docker exec -it govchaind-node tailscale up
    ```
3.  Follow the instructions in your terminal to authenticate via your web browser.

## 5. Configuration

### Persistent Peers and Seeds

Once your `govchaind` node is connected to Tailscale, it will advertise its Tailscale IP address as its external P2P address in `config.toml`. You can find this IP by checking the container logs or running `docker exec govchaind-node tailscale ip -4`.

To configure persistent peers or seeds using Tailscale IPs, you would typically edit the `config.toml` file located in your mounted volume (`govchaind-data`).

1.  Stop the `govchaind` container:
    ```bash
docker stop govchaind-node
    ```
2.  Access the `config.toml` file in your `govchaind-data` volume. You can do this by running a temporary container:
    ```bash
docker run --rm -v govchaind-data:/home/nonroot/.govchain -it alpine cat /home/nonroot/.govchain/config/config.toml
    ```
3.  Edit the `config.toml` file (e.g., using `docker cp` to copy it out, edit, and copy back, or by mounting the volume directly to your host for editing).
4.  Restart the `govchaind` container:
    ```bash
docker start govchaind-node
    ```

### Security Best Practices

-   **Use short-lived Tailscale auth keys:** Generate ephemeral keys for your nodes to minimize the risk of compromise.
-   **Restrict access:** Use Tailscale ACLs (Access Control Lists) to limit which devices can connect to your validator node.
-   **Monitor logs:** Regularly check Docker container logs and Tailscale logs for any unusual activity.
-   **Keep software updated:** Regularly rebuild your Docker image to ensure `govchaind` and Tailscale are running the latest versions with security patches.

## 6. Troubleshooting

-   **Container exits immediately:** Check `docker logs govchaind-node` for error messages. Common issues include:
    -   `permission denied`: Ensure the mounted volume has correct permissions (handled by `chown` in Dockerfile).
    -   `validator set is empty`: This is a blockchain-level error, not a Docker error. Ensure your `genesis.json` is correctly configured with validators.
-   **Tailscale connection issues:**
    -   `invalid key`: Ensure your `TS_AUTHKEY` is correct and has not expired.
    -   `TUN device not found` or `NET_ADMIN` errors: Ensure you are running the container with `--cap-add=NET_ADMIN --device=/dev/net/tun`.
    -   Check Tailscale status inside the container: `docker exec -it govchaind-node tailscale status`.
-   **Peer discovery problems:** Verify that your `external_address` in `config.toml` is correctly set to the Tailscale IP. Ensure your Tailscale network allows P2P connections.

## 7. Removing the Container and Data

To stop and remove the container:

```bash
docker stop govchaind-node
docker rm govchaind-node
```

To remove the persistent data volume (this will delete all blockchain data!):

```bash
docker volume rm govchaind-data
```
