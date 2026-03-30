#!/bin/bash
# join.sh — Permissionless validator join demo
# A new node generates a key, stakes, registers, and joins the active validator set.
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
RPC=http://localhost:8551

# Funded account: first derivation of test,junk mnemonic (large KAIA balance)
# Private key from: cast wallet derive-key --mnemonic "test junk" --hd-path "m/44'/60'/0'/0/0"
FUNDED_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
FUNDED_ADDR=$(cast wallet address --private-key "$FUNDED_KEY")

echo "=== Permissionless Validator Join Demo ==="
echo "Funded account: $FUNDED_ADDR"
echo ""

# ABv2 is at the fixed precompile address 0x400
ABV2=0x0000000000000000000000000000000000000400
# CnStakingFactory is resolved from Registry (0x401)
REGISTRY=0x0000000000000000000000000000000000000401
echo "[1/6] Resolving contract addresses..."
FACTORY=$(cast call "$REGISTRY" "getActiveAddr(string)(address)" "CnStakingFactory" --rpc-url "$RPC")
echo "  CnStakingFactory: $FACTORY"
echo "  ABv2:             $ABV2"

# 2. Generate CN5 key
echo "[2/6] Generating CN5 key..."
mkdir -p "$SCRIPT_DIR/cn5"
(cd "$SCRIPT_DIR/cn5" && kgen --file)
CN5_NODEID=$(python3 -c "import json; print(json.load(open('$SCRIPT_DIR/cn5/keys/node_info.json'))['Address'])")
CN5_NODEKEY="0x$(cat "$SCRIPT_DIR/cn5/keys/nodekey")"

if [ -z "$CN5_NODEID" ] || [ "$CN5_NODEID" = "0x0000000000000000000000000000000000000000" ]; then
  echo "ERROR: CN5_NODEID is empty or zero. Check kgen output at $SCRIPT_DIR/cn5/keys/"
  exit 1
fi
echo "  CN5 node ID: $CN5_NODEID"

# 3. Deploy CnStakingV4 + PublicDelegation for CN5
# ABI: deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs)
# INITIAL_LOCKUP = 1e9 wei must be sent with the call
echo "[3/6] Deploying CnStaking + PublicDelegation for CN5..."
DEPLOY_OUT=$(cast send "$FACTORY" \
  "deployCnStakingWithPD(address,(address,address,uint256,string))" \
  "$CN5_NODEID" "($FUNDED_ADDR,$FUNDED_ADDR,0,\"\")" \
  --private-key "$FUNDED_KEY" --rpc-url "$RPC" \
  --value 1000000000 --json)

CN5_STAKING=$(echo "$DEPLOY_OUT" | python3 -c "import sys,json; print(json.load(sys.stdin)['logs'][0]['address'])")
if [ -z "$CN5_STAKING" ] || [ "$CN5_STAKING" = "0x0000000000000000000000000000000000000000" ]; then
  echo "ERROR: CN5_STAKING is empty. deployCnStakingWithPD may have failed."
  exit 1
fi
CN5_PD=$(cast call "$CN5_STAKING" "publicDelegation()(address)" --rpc-url "$RPC")
echo "  CnStaking:        $CN5_STAKING"
echo "  PublicDelegation: $CN5_PD"

# 4. Stake MIN_STAKE (5,000,000 KAIA) via PublicDelegation
echo "[4/6] Staking 5,000,000 KAIA via PublicDelegation..."
MIN_STAKE=$(cast call "$ABV2" "MIN_STAKE()(uint256)" --rpc-url "$RPC" | cut -d' ' -f1)
MIN_STAKE=$(cast --from-wei "$MIN_STAKE")
echo "  MIN_STAKE: $MIN_STAKE KAIA"
cast send "$CN5_PD" "stake()" \
  --value "${MIN_STAKE}ether" \
  --private-key "$FUNDED_KEY" --rpc-url "$RPC" > /dev/null

# 5. Register CN5 in ABv2 (msg.sender becomes the node's manager)
echo "[5/6] Registering CN5 in ABv2..."
CN5_NODEKEY_RAW=$(cat "$SCRIPT_DIR/cn5/keys/nodekey")
BLS_INFO=$(cd "$SCRIPT_DIR/blsgen" && go run bls.go "$CN5_NODEKEY_RAW")
BLS_PUB=$(echo "$BLS_INFO" | grep '^pub=' | cut -d= -f2)
BLS_POP=$(echo "$BLS_INFO" | grep '^pop=' | cut -d= -f2)
cast send "$ABV2" \
  "createNode(address,address,address,address,(bytes,bytes),string)" \
  "$CN5_NODEID" "$CN5_STAKING" "$CN5_PD" "$FUNDED_ADDR" "($BLS_PUB,$BLS_POP)" "" \
  --private-key "$FUNDED_KEY" --rpc-url "$RPC" > /dev/null

# 6. Signal ready — MUST use CN5's nodekey (onlyNodeId: msg.sender == nodeId)
echo "[6/6] Signaling readyCandidate (using CN5 nodekey)..."
# Fund CN5_NODEID with gas first
cast send "$CN5_NODEID" --value 1ether \
  --private-key "$FUNDED_KEY" --rpc-url "$RPC" > /dev/null
# Call readyCandidate signed by CN5 nodekey
cast send "$ABV2" "readyCandidate(address)" "$CN5_NODEID" \
  --private-key "$CN5_NODEKEY" --rpc-url "$RPC" > /dev/null

echo ""
echo "CN5 ($CN5_NODEID) is CandReady."
echo "Waiting for next epoch (~100 blocks) to join validator set..."
echo ""

# Poll until CN5 appears in validator set (5-minute timeout = 60 × 5s)
ITERS=0
MAX_ITERS=60
while true; do
  VALS=$(cast rpc kaia_getNodeByState '["latest",["ValActive"]]' --rpc-url "$RPC" 2>/dev/null \
    | python3 -c "import sys,json; print(len(json.load(sys.stdin)))" 2>/dev/null || echo 0)
  BLOCK=$(cast block-number --rpc-url "$RPC" 2>/dev/null || echo 0)
  printf "\r[block %s] validators: %s  " "$BLOCK" "$VALS"
  if [ "$VALS" = "5" ]; then
    echo ""
    echo ""
    echo "SUCCESS: CN5 joined the validator set! Validators: $VALS"
    break
  fi
  ITERS=$((ITERS + 1))
  if [ "$ITERS" -ge "$MAX_ITERS" ]; then
    echo ""
    echo "TIMEOUT: CN5 did not appear after 5 minutes."
    echo "Check: is CN5 node process running? Verify BOOTNODES in cn5/conf/kcnd.conf if running a real node."
    exit 1
  fi
  sleep 5
done
