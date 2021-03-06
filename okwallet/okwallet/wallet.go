package okwallet

import (
	"encoding/json"
	"errors"
	"gitlab.com/NebulousLabs/Sia/types"
)

var (
	errBuilderAlreadySigned = errors.New("sign has already been called on this transaction builder, multiple calls can cause issues")
)

//return address, UnlockConditions, err
func GetAddressByPrivateKey(ssk string) (address, ucs_s string, err error) {
	address = ""
	ucs_s = ""

	sk, err := hexString2SecretKey(ssk)
	if err != nil {
		return
	}

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
func CreateRawTransaction(amount_s string, fee_s string, from_ucs_s string, to_ucs_s string, refund_ucs_s string, spendingTxs_s string) (string, error) {
	amount, ok := scanAmount(amount_s)
	if !ok {
		return "", errors.New("could not read amount from '" + amount_s + "'")
	}
	fee, ok := scanAmount(fee_s)
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

//reference: transactionBuilder.Sign
func SignRawTransaction(txBuilder_s string, secKeys_s string) (string, error) {
	var txBuilder okTransactionBuilder
	err := json.Unmarshal([]byte(txBuilder_s), &txBuilder)
	if err != nil {
		return "", err
	}

	secKeys, err := string2SecretKeys(secKeys_s)
	if err != nil {
		return "", err
	}

	txSet, err := txBuilder.Sign(secKeys)
	if err != nil {
		return "", err
	}

	txSet_b, err := json.Marshal(txSet)
	if err != nil {
		return "", err
	}

	return string(txSet_b), nil
}
