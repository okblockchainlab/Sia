package main

import (
  "errors"
  "encoding/hex"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/types"
)

func getAddressByPrivateKey(ssk string) (string, error) {
  var sk crypto.SecretKey
  b, err := hex.DecodeString(ssk)
  if err != nil {
    return "", err
  }

  if len(b) != len(sk) {
    return "", errors.New("invalid securet key!")
  }

  copy(sk[:], b)

  pk := sk.PublicKey()

  ucs := types.UnlockConditions{
			PublicKeys:         []types.SiaPublicKey{types.Ed25519PublicKey(pk)},
			SignaturesRequired: 1,
	}

  return ucs.UnlockHash().String(), nil
}

func main() {}
