#!/bin/sh

BUILD_TIME=`date`
GIT_REVISION=`git rev-parse --short HEAD`
GIT_DIRTY=`git diff-index --quiet HEAD -- || echo "✗-"`

# ldflags is copied from Makefile
ldflags="-X gitlab.com/NebulousLabs/Sia/build.GitRevision=${GIT_DIRTY}${GIT_REVISION} -X \"gitlab.com/NebulousLabs/Sia/build.BuildTime=${BUILD_TIME}\""

# pkgs is copied from Makefile
pkgs="./build ./cmd/siac ./cmd/siad ./compatibility ./crypto ./encoding ./modules ./modules/consensus ./modules/explorer ./modules/gateway ./modules/host ./modules/host/contractmanager ./modules/renter ./modules/renter/contractor ./modules/renter/hostdb ./modules/renter/hostdb/hosttree ./modules/renter/proto ./modules/miner ./modules/wallet ./modules/transactionpool ./node ./node/api ./persist ./siatest ./siatest/consensus ./siatest/renter ./siatest/wallet ./node/api/server ./sync ./types"


go build -tags="netgo" -a -ldflags="-s -w ${ldflags}" ${pkgs}


go build -o libsia.dylib -buildmode=c-shared -tags="netgo" -a -ldflags="-s -w ${ldflags}" ./okwallet/libsia
