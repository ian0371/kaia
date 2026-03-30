#!/bin/bash
# setup.sh — Start a 4-node permissionless devnet
set -e

DEPLOY_DIR="$HOME/kaia/kaiaspray/local-deploy"

if [ ! -d "$DEPLOY_DIR" ]; then
  echo "ERROR: kaiaspray local-deploy not found at $DEPLOY_DIR"
  echo "Clone kaiaspray and run ./install.sh first."
  exit 1
fi

cd "$DEPLOY_DIR"

# Set homi options: mnemonic for reproducible keys, short vrank epoch for demo
if grep -q "HOMI_ADDITIONAL_OPTIONS" properties.sh; then
  sed -i.bak 's|HOMI_ADDITIONAL_OPTIONS=.*|HOMI_ADDITIONAL_OPTIONS="--mnemonic test,junk --vrank-epoch 30"|' properties.sh
else
  echo 'HOMI_ADDITIONAL_OPTIONS="--mnemonic test,junk --vrank-epoch 30"' >> properties.sh
fi

# Set 4 CNs
if grep -q "^NUMOFCN=" properties.sh; then
  sed -i.bak 's|^NUMOFCN=.*|NUMOFCN=4|' properties.sh
else
  echo 'NUMOFCN=4' >> properties.sh
fi

echo "[1/3] Setting up genesis and node configs..."
bash 0_kaia_setup.sh

echo "[2/3] Initializing nodes..."
bash 2_initialize_nodes.sh

echo "[3/3] Starting network..."
bash 3_ccstart.sh

echo ""
echo "Network is up!"
echo "  RPC: http://localhost:8551"
echo "  Validators: 4"
echo ""
echo "Next: open demo/permissionless/dashboard.html, then run join.sh"
