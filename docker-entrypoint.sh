#!/bin/sh
set -e

MONIKER=${MONIKER:-"GovChain Validator"}

# Ensure the config directory exists
mkdir -p "/home/nonroot/.govchain/config"
chown nonroot:nonroot "/home/nonroot/.govchain/config"

# Check if genesis.json exists. If not, initialize and then replace genesis.json
if [ ! -f "/home/nonroot/.govchain/config/genesis.json" ]; then
  # Initialize the node to create default config files (including a dummy genesis.json)
  # Use a temporary moniker and chain-id as they will be overwritten by the downloaded genesis
  gosu nonroot env 'HOME=/home/nonroot' govchaind init "$MONIKER" --chain-id "temp-chain" --home "/home/nonroot/.govchain"

  # Download the actual genesis.json
  gosu nonroot env 'HOME=/home/nonroot' curl -sL "https://raw.githubusercontent.com/bettergovph/govchain/refs/heads/main/genesis.json" -o "/home/nonroot/.govchain/config/genesis.json"

  # Set the minimum gas price in app.toml
  gosu nonroot env 'HOME=/home/nonroot' sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' "/home/nonroot/.govchain/config/app.toml"
fi

# Wait for Tailscale sidecar to write its IP to the shared file
echo "Waiting for Tailscale IP file from sidecar..."
TS_IP_FILE="/var/run/tailscale-ip/ts_ip"
TS_IP=""
for i in $(seq 1 30); do
  if [ -f "$TS_IP_FILE" ]; then
    TS_IP=$(cat "$TS_IP_FILE")
    if [ -n "$TS_IP" ]; then
      echo "Tailscale IP found: $TS_IP"
      break
    fi
  fi
  echo "Waiting for Tailscale IP file... ($i/30)"
  sleep 2
done

if [ -n "$TS_IP" ]; then
  echo "Tailscale IP: $TS_IP. Updating config.toml..."
  gosu nonroot sed -i "s/^external_address = \"\"/external_address = \"$TS_IP:26656\"/" "/home/nonroot/.govchain/config/config.toml"
  echo "config.toml updated with Tailscale IP."
else
  echo "Could not get Tailscale IP from sidecar. Proceeding without updating external_address."
fi

# Add persistent peer to config.toml
PERSISTENT_PEER="4d153c889d9f0f4b670d2f548994fcdde208240e@157.90.134.175:26656"
gosu nonroot sed -i "s/^persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEER\"/" "/home/nonroot/.govchain/config/config.toml"
echo "Persistent peer added to config.toml."

# Start govchaind
gosu nonroot govchaind "$@"
