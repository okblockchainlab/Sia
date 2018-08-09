package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/node/api"
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

//reference: Wallet.SendSiacoins
//fee can get by /tpool/fee, max for general. and mul 750.
func createRawTransaction(amount_s string, fee_s string, from_ucs_s string, to_ucs_s string, refund_ucs_s string, spendingTxs_s string) (string, error) {
	amount, ok := api.scanAmount(amount_s)
	if !ok {
		return "", errors.New("could not read amount from '" + amount_s + "'")
	}
	fee, err := api.scanAmount(fee_s)
	if err != nil {
		return "", errors.New("could not read fee from '" + fee_s + "'")
	}

	var spendingTx []SpendingTransaction
	err = json.Unmarshal(spendingTxs_s, &spendingTx)
	if err != nil {
		return "", err
	}

	var from_ucs types.UnlockConditions
	err = json.Unmarshal(from_ucs_s, &from_ucs)
	if err != nil {
		return "", err
	}
	var to_ucs types.UnlockConditions
	err = json.Unmarshal(to_ucs_s, &to_ucs)
	if err != nil {
		return "", err
	}
	var refund_ucs types.UnlockConditions
	if len(refund_ucs_s) != 0 {
		err = json.Unmarshal(refund_ucs_s, &refund_ucs)
		if err != nil {
			return "", err
		}
	} else {
		refund_ucs = nil
	}

	output := types.SiacoinOutput{
		Value:      amount,
		UnlockHash: to_ucs.UnlockHash(),
	}

	txnBuilder, err := startTransaction()
	if err != nil {
		return nil, err
	}

	txnBuilder.FundSiacoins()

	txnBuilder.AddMinerFee(fee)
	txnBuilder.AddSiacoinOutput(output)

	result, err := json.Marshal(txnBuilder)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func main() {}
