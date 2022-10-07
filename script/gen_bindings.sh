#!/bin/bash

# Generate go contract bindings.
# ref: https://geth.ethereum.org/docs/dapp/native-bindings

set -eou pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"

echo ""
echo "TAIKO_MONO_DIR: ${TAIKO_MONO_DIR}"
echo "TAIKO_CLIENT_DIR: ${TAIKO_CLIENT_DIR}"
echo ""

cd ${TAIKO_CLIENT_DIR} &&
  make all &&
  cd -

cd ${TAIKO_MONO_DIR}/packages/protocol &&
  yarn clean &&
  yarn compile &&
  cd -

ABIGEN_BIN=$TAIKO_CLIENT_DIR/build/bin/abigen

echo ""
echo "Start generating go contract bindings..."
echo ""

cat ${TAIKO_MONO_DIR}/packages/protocol/artifacts/contracts/L1/TaikoL1.sol/TaikoL1.json |
	jq .abi |
	${ABIGEN_BIN} --abi - --type TaikoL1Client --pkg bindings --out $DIR/../bindings/gen_taiko_l1.go

cat ${TAIKO_MONO_DIR}/packages/protocol/artifacts/contracts/L2/V1TaikoL2.sol/V1TaikoL2.json |
	jq .abi |
	${ABIGEN_BIN} --abi - --type V1TaikoL2Client --pkg bindings --out $DIR/../bindings/gen_v1_taiko_l2.go

git -C ${TAIKO_MONO_DIR} log --format="%H" -n 1 >./bindings/.githead

echo "Go contract bindings generated!"
