#!/bin/bash

# govchain Volunteer Node Setup
# Allows volunteers to join the network as validators without tokens

set -e

echo "================================"
echo "govchain Volunteer Node Setup"
echo "================================"
echo ""

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <node-name> <genesis-file-url>"
    echo ""
    echo "Example:"
    echo "  $0 volunteer-node-1 https://raw.githubusercontent.com/org/govchain/main/genesis.json"
    echo ""
    exit 1
fi

NODE_NAME="$1"
GENESIS_URL="$2"

echo "📝 Node Name: $NODE_NAME"
echo "🌐 Genesis URL: $GENESIS_URL"
echo ""

# Check if blockchain is built
if [ ! -f "./build/govchaind" ]; then
    echo "❌ Error: Blockchain not built. Please run: ignite chain build"
    exit 1
fi

# Initialize node
echo "🔧 Initializing volunteer node..."
$(dirname "$0")/../build/govchaind "$NODE_NAME" --chain-id govchain

# Download genesis file
echo "📥 Downloading genesis file..."
curl -s "$GENESIS_URL" > "$HOME/.govchain/config/genesis.json"

# Create validator key
echo "🔑 Creating validator key..."
./build/govchaind keys add validator --keyring-backend test

# Get validator address
VALIDATOR_ADDR=$(./build/govchaind keys show validator -a --keyring-backend test)

echo ""
echo "✅ Volunteer node setup complete!"
echo "================================"
echo ""
echo "📝 Node Details:"
echo "  Name: $NODE_NAME"
echo "  Validator Address: $VALIDATOR_ADDR"
echo ""
echo "🚀 To start your volunteer node:"
echo "  ./build/govchaind start"
echo ""
echo "🌐 To become a validator:"
echo "  ./build/govchaind tx staking create-validator \\"
echo "    --amount=1000000stake \\"
echo "    --pubkey=\$(./build/govchaind tendermint show-validator) \\"
echo "    --moniker=\"$NODE_NAME\" \\"
echo "    --chain-id=govchain \\"
echo "    --from=validator \\"
echo "    --keyring-backend=test"
echo ""
