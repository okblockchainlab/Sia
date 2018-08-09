package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/node/api"
	"gitlab.com/NebulousLabs/Sia/types"
)

//return address, UnlockConditions, err
func getAddressByPrivateKey(ssk string) (address, ucs_s string, err error) {
	address = ""
	ucs_s = ""

	var sk crypto.SecretKey
	b, err := hex.DecodeString(ssk)
	if err != nil {
		return
	}

	if len(b) != len(sk) {
		err = errors.New("invalid securet key!")
		return
	}

	copy(sk[:], b)

	pk := sk.PublicKey()

	ucs := types.UnlockConditions{
		PublicKeys:         []types.SiaPublicKey{types.Ed25519PublicKey(pk)},
		SignaturesRequired: 1,
	}

	ucs_b, err := json.Marshal(ucs)
	if err != nil {
		return
	}

	ucs_s = string(ucs_b)
	address = ucs.UnlockHash().String()
	err = nil
	return
}

//reference: Wallet.SendSiacoins
//fee can get by /tpool/fee, max for general. and mul 750.
func createRawTransaction(amount_s string, fee_s string, from_ucs_s string, to_ucs_s string, refund_ucs_s string, spendingTxs_s string) (string, error) {
	amount, ok := api.ScanAmount(amount_s)
	if !ok {
		return "", errors.New("could not read amount from '" + amount_s + "'")
	}
	fee, ok := api.ScanAmount(fee_s)
	if !ok {
		return "", errors.New("could not read fee from '" + fee_s + "'")
	}

	var spendingTx []SpendingTransaction
	err := json.Unmarshal([]byte(spendingTxs_s), &spendingTx)
	if err != nil {
		return "", err
	}

	var from_ucs types.UnlockConditions
	err = json.Unmarshal([]byte(from_ucs_s), &from_ucs)
	if err != nil {
		return "", err
	}
	var to_ucs types.UnlockConditions
	err = json.Unmarshal([]byte(to_ucs_s), &to_ucs)
	if err != nil {
		return "", err
	}
	var refund_ucs *types.UnlockConditions = nil
	if len(refund_ucs_s) != 0 {
		refund_ucs = &types.UnlockConditions{}
		err = json.Unmarshal([]byte(refund_ucs_s), &refund_ucs)
		if err != nil {
			return "", err
		}
	}

	output := types.SiacoinOutput{
		Value:      amount,
		UnlockHash: to_ucs.UnlockHash(),
	}

	txnBuilder, err := startTransaction()
	if err != nil {
		return "", err
	}

	err = txnBuilder.FundSiacoins(amount.Add(fee), spendingTx, from_ucs, refund_ucs)
	if err != nil {
		return "", err
	}

	txnBuilder.AddMinerFee(fee)
	txnBuilder.AddSiacoinOutput(output)

	result, err := json.Marshal(txnBuilder)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func main() {}
