package okwallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/types"
	"math/big"
)

func hexString2SecretKey(s string) (*crypto.SecretKey, error) {
	var sk crypto.SecretKey
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(b) != len(sk) {
		return nil, errors.New("invalid securet key!")
	}

	copy(sk[:], b)
	return &sk, nil
}

func toStringArray(s string) ([]string, error) {
	var sArr []string
	err := json.Unmarshal([]byte(s), &sArr)
	if err != nil {
		return nil, err
	}

	return sArr, nil
}

func string2SecretKeys(s string) ([]crypto.SecretKey, error) {
	sks, err := toStringArray(s)
	if err != nil {
		return nil, err
	}

	var secKeys []crypto.SecretKey
	for _, s := range sks {
		sk, err := hexString2SecretKey(s)
		if err != nil {
			return nil, err
		}

		secKeys = append(secKeys, *sk)
	}

	return secKeys, nil
}

func scanAmount(amount string) (types.Currency, bool) {
	// use SetString manually to ensure that amount does not contain
	// multiple values, which would confuse fmt.Scan
	i, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return types.Currency{}, ok
	}
	return types.NewCurrency(i), true
}
