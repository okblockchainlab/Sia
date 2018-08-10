package main

import (
	"bytes"
	"errors"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/encoding"
	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/types"
)

type okTransactionBuilder struct {
	Signed        bool                `json:"signed"`
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
		Signed:      false,
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

func (tb *okTransactionBuilder) Sign(secKeys []crypto.SecretKey) ([]types.Transaction, error) {
	if tb.Signed {
		return nil, errBuilderAlreadySigned
	}

	// Create the coveredfields struct.
	coveredFields := types.CoveredFields{WholeTransaction: true}
	// TransactionSignatures don't get covered by the 'WholeTransaction' flag,
	// and must be covered manually.
	for i := range tb.Transaction.TransactionSignatures {
		coveredFields.TransactionSignatures = append(coveredFields.TransactionSignatures, uint64(i))
	}

	// Sign all of the inputs to the parent transaction.
	for i := range tb.Parents {
		for _, sci := range tb.Parents[i].SiacoinInputs {
			_, err := addSignatures(&tb.Parents[i], types.FullCoveredFields, sci.UnlockConditions, crypto.Hash(sci.ParentID), secKeys)
			if err != nil {
				return nil, err
			}
		}
	}

	// For each siacoin input in the transaction that we added, provide a
	// signature.
	for _, inputIndex := range tb.SiacoinInputs {
		input := tb.Transaction.SiacoinInputs[inputIndex]
		//newSigIndices := addSignatures(&tb.Transaction, coveredFields, input.UnlockConditions, crypto.Hash(input.ParentID), secKeys)
		_, err := addSignatures(&tb.Transaction, coveredFields, input.UnlockConditions, crypto.Hash(input.ParentID), secKeys)
		if err != nil {
			return nil, err
		}
		//tb.transactionSignatures = append(tb.transactionSignatures, newSigIndices...)
		tb.Signed = true // Signed is set to true after one successful signature to indicate that future signings can cause issues.
	}

	// Get the transaction set and delete the transaction from the registry.
	txnSet := append(tb.Parents, tb.Transaction)
	return txnSet, nil
}

// addSignatures will sign a transaction using a spendable key, with support
// for multisig spendable keys. Because of the restricted input, the function
// is compatible with both siacoin inputs and siafund inputs.
// reference: addSignatures in transactionbuilder.go
func addSignatures(txn *types.Transaction, cf types.CoveredFields, uc types.UnlockConditions, parentID crypto.Hash, secKeys []crypto.SecretKey) (newSigIndices []int, err error) {
	err = nil
	newSigIndices = nil
	// Try to find the matching secret key for each public key - some public
	// keys may not have a match. Some secret keys may be used multiple times,
	// which is why public keys are used as the outer loop.
	totalSignatures := uint64(0)
	for i, siaPubKey := range uc.PublicKeys {
		// Search for the matching secret key to the public key.
		for _, sk := range secKeys {
			pubKey := sk.PublicKey()
			if !bytes.Equal(siaPubKey.Key, pubKey[:]) {
				continue
			}

			// Found the right secret key, add a signature.
			sig := types.TransactionSignature{
				ParentID:       parentID,
				CoveredFields:  cf,
				PublicKeyIndex: uint64(i),
			}
			newSigIndices = append(newSigIndices, len(txn.TransactionSignatures))
			txn.TransactionSignatures = append(txn.TransactionSignatures, sig)
			sigIndex := len(txn.TransactionSignatures) - 1
			sigHash := txn.SigHash(sigIndex)
			encodedSig := crypto.SignHash(sigHash, sk)
			txn.TransactionSignatures[sigIndex].Signature = encodedSig[:]

			// Count that the signature has been added, and break out of the
			// secret key loop.
			totalSignatures++
			break
		}

		// If there are enough signatures to satisfy the unlock conditions,
		// break out of the outer loop.
		if totalSignatures == uc.SignaturesRequired {
			break
		}
	}

	if newSigIndices == nil {
		err = errors.New("didn't match a secret key for transaction input")
	}
	return newSigIndices, err
}
