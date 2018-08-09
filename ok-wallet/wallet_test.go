package main

import (
	"testing"
)

func TestGetAddressByPrivateKey(t *testing.T) {
	sk := "97a4f67362116f9011448e881113213ec4cfe9a676605791a97b2838cf0f3486388ea2ae690da5116f8be18cdba570fc72df9fb3c09cd853a60e3737599dbd27"
	address := "d874365aaca9f09dcad33009daa76c62b4c036fdd1c3c9d3ac2297d5a29c9787927341218ccc"
	//pk := "388ea2ae690da5116f8be18cdba570fc72df9fb3c09cd853a60e3737599dbd27"

	addr, err := getAddressByPrivateKey(sk)

	if err != nil {
		t.Error("getAddressByPrivateKey return error :", err)
	}
	if address != addr {
		msg := "getAddressByPrivateKey return address which is not expected.\n" +
			"expected is: " + address + "\n" +
			"result is: " + addr + "\n"
		t.Error(msg)
	}
}
