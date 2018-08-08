package main

import (
  "errors"
  "encoding/hex"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/types"
	"gitlab.com/NebulousLabs/Sia/node/api"
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

//reference: Wallet.registerTransaction
func startTransaction() (modules.TransactionBuilder, error) {
  t := types.Transaction{}
  var parents []types.Transaction = nil

	pBytes := encoding.Marshal(parents)
	var pCopy []types.Transaction
	err := encoding.Unmarshal(pBytes, &pCopy)
	if err != nil {
		panic(err)
	}
	tBytes := encoding.Marshal(t)
	var tCopy types.Transaction
	err = encoding.Unmarshal(tBytes, &tCopy)
	if err != nil {
		panic(err)
	}
	return &transactionBuilder{
		parents:     pCopy,
		transaction: tCopy,

		wallet: w,
	}
}

//reference: Wallet.SendSiacoins
//fee can get by /tpool/fee, max for general. and I will mul 750.
func createRawTransaction(amount_s string, dest_s string, fee string) (string, error) {
  amount, ok := api.scanAmount(amount_s)
  if !ok {
    return "", errors.New("could not read amount from '" + amount_s + "'")
  }
  dest, err := scanAddress(dest_s)
  if err != nil {
    return "", errors.New("could not read address from '" + dest_s + "'")
  }

	//_, tpoolFee := w.tpool.FeeEstimation()
	//tpoolFee = tpoolFee.Mul64(750) // Estimated transaction size in bytes
	output := types.SiacoinOutput{
		Value:      amount,
		UnlockHash: dest,
	}

	txnBuilder, err := startTransaction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			txnBuilder.Drop()
		}
	}()
	err = txnBuilder.FundSiacoins(amount.Add(tpoolFee))
	if err != nil {
		//w.log.Println("Attempt to send coins has failed - failed to fund transaction:", err)
		return nil, build.ExtendErr("unable to fund transaction", err)
	}
	txnBuilder.AddMinerFee(tpoolFee)
	txnBuilder.AddSiacoinOutput(output)
	//txnSet, err := txnBuilder.Sign(true)
}

func main() {}
