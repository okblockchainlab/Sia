package main

// #include <jni.h>
import "C"

import (
	"gitlab.com/NebulousLabs/Sia/okwallet/okwallet"
	"regexp"
)

func setErrorResult(env *C.JNIEnv, errMsg string) C.jobjectArray {
	result := newStringObjectArray(env, 1)
	setObjectArrayStringElement(env, result, 0, errMsg)
	return result
}

func getAddressByPrivateKeyExecute(env *C.JNIEnv, cmd string, args []string) C.jobjectArray {
	if len(args) != 1 {
		return setErrorResult(env, "error: "+cmd+" wrong argument count")
	}

	addr, ucs, err := okwallet.GetAddressByPrivateKey(args[0])
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	result := newStringObjectArray(env, 2)
	setObjectArrayStringElement(env, result, 0, addr)
	setObjectArrayStringElement(env, result, 1, ucs)

	return result
}

func createRawTransactionExecute(env *C.JNIEnv, cmd string, args []string) C.jobjectArray {
	if len(args) != 6 {
		return setErrorResult(env, "error: "+cmd+" wrong argument count")
	}

	txBuilder, err := okwallet.CreateRawTransaction(args[0], args[1], args[2], args[3], args[4], args[5])
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	result := newStringObjectArray(env, 1)
	setObjectArrayStringElement(env, result, 0, txBuilder)

	return result
}

func SignRawTransactionExecute(env *C.JNIEnv, cmd string, args []string) C.jobjectArray {
	if len(args) != 2 {
		return setErrorResult(env, "error: "+cmd+" wrong argument count")
	}

	signedTx, err := okwallet.SignRawTransaction(args[0], args[1])
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	result := newStringObjectArray(env, 1)
	setObjectArrayStringElement(env, result, 0, signedTx)

	return result
}

//export Java_com_okcoin_vault_jni_sia_Siaj_execute
func Java_com_okcoin_vault_jni_sia_Siaj_execute(env *C.JNIEnv, _ C.jclass, _ C.jstring, jcommand C.jstring) C.jobjectArray {
	command, err := jstring2string(env, jcommand)
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	sepExp, err := regexp.Compile(`\s+`)
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	args := sepExp.Split(command, -1)
	if len(args) < 2 {
		return setErrorResult(env, "error: invalid command")
	}

	switch args[0] {
	case "getaddressbyprivatekey":
		return getAddressByPrivateKeyExecute(env, args[0], args[1:])
	case "createrawtransaction":
		return createRawTransactionExecute(env, args[0], args[1:])
	case "signrawtransaction":
		return SignRawTransactionExecute(env, args[0], args[1:])
	default:
		return setErrorResult(env, "error: unknown command: "+args[0])
	}
	return setErrorResult(env, "error: unknown command: "+args[0])
}

func main() {}
