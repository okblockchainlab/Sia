package main

import (
  "errors"
  "encoding/hex"
  "encoding/json"
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

		wallet: nil,
	}
}

//reference: Wallet.SendSiacoins
//fee can get by /tpool/fee, max for general. and I will mul 750.
/*siacoinoutputs:
[
  {
    transactions: {
    },
    use_outputs: [1, 2, 3]
  },
]
*/

func createRawTransaction(amount_s string, dest_s string, fee_s string, from_unlock_cond_s string, to_unlock_cond_s string, refund_unlock_cond_s string, siacoinoutputs string) (string, error) {
  amount, ok := api.scanAmount(amount_s)
  if !ok {
    return "", errors.New("could not read amount from '" + amount_s + "'")
  }
  dest, err := api.scanAddress(dest_s)
  if err != nil {
    return "", errors.New("could not read address from '" + dest_s + "'")
  }
  fee, err := api.scanAmount(fee_s)
  if err != nil {
    return "", errors.New("could not read fee from '" + fee_s + "'")
  }

  var scouts []SiacoinOutput
  err = json.Unmarshal(siacoinoutputs, &scouts)
  if err != nil {
    return "", err
  }

  var from_ucs types.UnlockConditions
  err = json.Unmarshal(from_unlock_cond_s, &from_ucs)
  if err != nil {
    return "", err
  }
  var to_ucs types.UnlockConditions
  err = json.Unmarshal(to_unlock_cond_s, &to_ucs)
  if err != nil {
    return "", err
  }

  //check value
  var fund types.Currency
  for _, sco := range scouts {
      fund.Add(sco.Value)
  }
  if fund.Cmp(amount.Add(fee)) < 0 {
      return modules.ErrLowBalance
  }

	parentTxn := types.Transaction{}
  for _, sco := range scouts {
		sci := types.SiacoinInput{
			ParentID:         t.SiacoinOutputID(xxxxxxi),
			UnlockConditions: from_ucs,
		}
		parentTxn.SiacoinInputs = append(parentTxn.SiacoinInputs, sci)
  }

	exactOutput := types.SiacoinOutput{
		Value:      amount,
		UnlockHash: to_ucs.UnlockHash(),
	}
	parentTxn.SiacoinOutputs = append(parentTxn.SiacoinOutputs, exactOutput)

  if !amount.Equals(fund) {
    var refund_ucs types.UnlockConditions
    err = json.Unmarshal(refun_unlock_cond_s, &refund_ucs)
    if err != nil {
      return "", err
    }
		refundOutput := types.SiacoinOutput{
			Value:      fund.Sub(amount),
			UnlockHash: refund_ucs.UnlockHash(),
		}
		parentTxn.SiacoinOutputs = append(parentTxn.SiacoinOutputs, refundOutput)
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
