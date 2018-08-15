#!/bin/sh

if [ -z "$JAVA_HOME" ]; then
  printf "No JAVA_HOME detected! "
  printf "Setup JAVA_HOME before build: export JAVA_HOME=/path/to/java\\n"
  exit 1
fi

BUILD_TIME=`date`
GIT_REVISION=`git rev-parse --short HEAD`
GIT_DIRTY=`git diff-index --quiet HEAD -- || echo "✗-"`

EXT=so
TARGET_OS=`uname -s`
case "$TARGET_OS" in
  Darwin)
    EXT=dylib
    export CGO_CFLAGS="-I${JAVA_HOME}/include -I${JAVA_HOME}/include/darwin"
    ;;
  Linux)
    EXT=so
    export CGO_CFLAGS="-I${JAVA_HOME}/include -I${JAVA_HOME}/include/linux"
    ;;
  *)
  echo "Unknown platform!" >&2
  exit 1
esac


# ldflags is copied from Makefile
ldflags="-X gitlab.com/NebulousLabs/Sia/build.GitRevision=${GIT_DIRTY}${GIT_REVISION} -X \"gitlab.com/NebulousLabs/Sia/build.BuildTime=${BUILD_TIME}\""

# pkgs is copied from Makefile
pkgs="./build ./cmd/siac ./cmd/siad ./compatibility ./crypto ./encoding ./modules ./modules/consensus ./modules/explorer ./modules/gateway ./modules/host ./modules/host/contractmanager ./modules/renter ./modules/renter/contractor ./modules/renter/hostdb ./modules/renter/hostdb/hosttree ./modules/renter/proto ./modules/miner ./modules/wallet ./modules/transactionpool ./node ./node/api ./persist ./siatest ./siatest/consensus ./siatest/renter ./siatest/wallet ./node/api/server ./sync ./types"


go build -tags="netgo" -a -ldflags="-s -w ${ldflags}" ${pkgs}


go build -o libsia.${EXT} -buildmode=c-shared -tags="netgo" -a -ldflags="-s -w ${ldflags}" ./okwallet/libsia
[ $? -ne 0 ] && exit 1
nm -D libsia.${EXT} |grep "[ _]Java_com_okcoin"
