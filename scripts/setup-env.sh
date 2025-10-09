#!/bin/bash

# govchain Blockchain Environment Setup (Tokenless Configuration)
# Run this script after building the blockchain

set -e

echo "================================"
echo "govchain Environment Setup (Tokenless)"
echo "================================"
echo ""

# Check if blockchain is built
if [ ! -f $(dirname "$0")/../build/govchaind ]; then
    echo "⚠️  Blockchain not built yet. Building now..."
    ignite chain build
fi

# Initialize blockchain if not already done
if [ ! -d "$HOME/.govchain" ]; then
    echo "🔧 Initializing tokenless blockchain..."
    ./build/govchaind init mynode --chain-id govchain
    
    # Create validator key (no tokens needed)
    ./build/govchaind keys add validator --keyring-backend test
    
    # Create a tokenless genesis configuration
    # No genesis accounts with stake tokens needed for tokenless network
    echo "📝 Configuring tokenless genesis..."
    
    # Modify genesis.json for tokenless operation
    GENESIS_FILE="$HOME/.govchain/config/genesis.json"
    
    # Set minimal staking parameters (validators don't need tokens)
    jq '.app_state.staking.params.bond_denom = ""' "$GENESIS_FILE" > tmp.json && mv tmp.json "$GENESIS_FILE"
    jq '.app_state.gov.params.min_deposit = []' "$GENESIS_FILE" > tmp.json && mv tmp.json "$GENESIS_FILE"
    jq '.app_state.mint.minter.inflation = "0.000000000000000000"' "$GENESIS_FILE" > tmp.json && mv tmp.json "$GENESIS_FILE"
    
    # Create genesis transaction without staking tokens
    ./build/govchaind genesis gentx validator 1000000stake --chain-id govchain --keyring-backend test
    
    # Collect genesis transactions
    ./build/govchaind genesis collect-gentxs
    
    echo "✅ Tokenless blockchain configured!"
    echo "🌐 Volunteers can join as validators without tokens"
fi

echo "✅ Blockchain environment ready!"
echo ""
echo "To start the blockchain:"
echo "  ignite chain serve"
echo ""
echo "To start supporting services, run from the project root:
  docker-compose up -d"
echo ""
echo "📋 Volunteer Node Operators:"
echo "  - No tokens required to participate"
echo "  - Governance is based on data contribution"
echo "  - Validators secure the network through consensus"
echo ""
