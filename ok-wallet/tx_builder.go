package main

import (
	"errors"
	"gitlab.com/NebulousLabs/Sia/encoding"
	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/types"
)

type okTransactionBuilder struct {
	Parents       []types.Transaction `json:"parents"`
	SiacoinInputs []int               `json:"siacoininputs"`
	Transaction   types.Transaction   `json:"transaction"`
}

//reference: Wallet.registerTransaction
func startTransaction() (*okTransactionBuilder, error) {
	t := types.Transaction{}
	var parents []types.Transaction = nil

	pBytes := encoding.Marshal(parents)
	var pCopy []types.Transaction
	err := encoding.Unmarshal(pBytes, &pCopy)
	if err != nil {
		return nil, err
	}
	tBytes := encoding.Marshal(t)
	var tCopy types.Transaction
	err = encoding.Unmarshal(tBytes, &tCopy)
	if err != nil {
		return nil, err
	}
	return &okTransactionBuilder{
		Parents:     pCopy,
		Transaction: tCopy,
	}, nil
}

//refund_ucs may nil
func (tb *okTransactionBuilder) FundSiacoins(amount types.Currency, spending []SpendingTransaction, from_ucs types.UnlockConditions, refund_ucs *types.UnlockConditions) error {
	//check value
	var fund types.Currency
	for _, sp := range spending {
		for si := range sp.SpendingOutputs {
			fund = fund.Add(sp.Tx.SiacoinOutputs[si].Value)
		}
	}
	if fund.Cmp(amount) < 0 {
		return modules.ErrLowBalance
	}

	parentTxn := types.Transaction{
		FileContracts:         []types.FileContract{},
		FileContractRevisions: []types.FileContractRevision{},
		StorageProofs:         []types.StorageProof{},
		SiafundInputs:         []types.SiafundInput{},
		SiafundOutputs:        []types.SiafundOutput{},
		MinerFees:             []types.Currency{},
		ArbitraryData:         [][]byte{},
		TransactionSignatures: []types.TransactionSignature{},
	}
	parentUnlockConditions := from_ucs

	for _, sp := range spending {
		for _, si := range sp.SpendingOutputs {
			sci := types.SiacoinInput{
				ParentID:         sp.Tx.SiacoinOutputID(uint64(si)),
				UnlockConditions: from_ucs,
			}
			parentTxn.SiacoinInputs = append(parentTxn.SiacoinInputs, sci)
		}
	}

	exactOutput := types.SiacoinOutput{
		Value:      amount,
		UnlockHash: parentUnlockConditions.UnlockHash(),
	}
	parentTxn.SiacoinOutputs = append(parentTxn.SiacoinOutputs, exactOutput)

	if !amount.Equals(fund) {
		if nil == refund_ucs {
			return errors.New("need refund UnlockConditions")
		}
		refundOutput := types.SiacoinOutput{
			Value:      fund.Sub(amount),
			UnlockHash: refund_ucs.UnlockHash(),
		}
		parentTxn.SiacoinOutputs = append(parentTxn.SiacoinOutputs, refundOutput)
	}

	//XXX: Sign all of the inputs to the parent transaction.
	//for _, sci := range parentTxn.SiacoinInputs {
	//addSignatures(&parentTxn, types.FullCoveredFields, sci.UnlockConditions, crypto.Hash(sci.ParentID), tb.wallet.keys[sci.UnlockConditions.UnlockHash()])
	//}

	// Add the exact output.
	newInput := types.SiacoinInput{
		ParentID:         parentTxn.SiacoinOutputID(0),
		UnlockConditions: parentUnlockConditions,
	}

	tb.Parents = append(tb.Parents, parentTxn)
	tb.SiacoinInputs = append(tb.SiacoinInputs, len(tb.Transaction.SiacoinInputs))
	tb.Transaction.SiacoinInputs = append(tb.Transaction.SiacoinInputs, newInput)
	return nil
}

//reference: wallet.AddMinerFee
func (tb *okTransactionBuilder) AddMinerFee(fee types.Currency) uint64 {
	tb.Transaction.MinerFees = append(tb.Transaction.MinerFees, fee)
	return uint64(len(tb.Transaction.MinerFees) - 1)
}

//reference: wallet.AddSiacoinOutput
func (tb *okTransactionBuilder) AddSiacoinOutput(output types.SiacoinOutput) uint64 {
	tb.Transaction.SiacoinOutputs = append(tb.Transaction.SiacoinOutputs, output)
	return uint64(len(tb.Transaction.SiacoinOutputs) - 1)
}
