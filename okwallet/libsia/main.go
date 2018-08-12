package main

import (
	"gitlab.com/NebulousLabs/Sia/okwallet/okwallet"
)

func GetAddressByPrivateKey(ssk string) (address, ucs_s string, err error) {
	return okwallet.GetAddressByPrivateKey(ssk)
}

func CreateRawTransaction(amount_s string, fee_s string, from_ucs_s string, to_ucs_s string, refund_ucs_s string, spendingTxs_s string) (string, error) {
	return okwallet.CreateRawTransaction(amount_s, fee_s, from_ucs_s, to_ucs_s, refund_ucs_s, spendingTxs_s)
}

func SignRawTransaction(txBuilder_s string, secKeys_s string) (string, error) {
	return okwallet.SignRawTransaction(txBuilder_s, secKeys_s)
}

func main() {}
