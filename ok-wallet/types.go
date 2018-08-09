package main

import (
	"gitlab.com/NebulousLabs/Sia/types"
)

type SpendingTransaction struct {
	Tx              types.Transaction `json:"transaction"`
	SpendingOutputs []int             `json:"outputs"`
}
