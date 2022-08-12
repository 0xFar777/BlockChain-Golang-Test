package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
)

func Decodemain(address string) {
	res1, res2, res3 := base58.CheckDecode(address)
	fmt.Printf("%v--%v--%v\n", res1, res2, res3)
}

func Decodemain2(address string) {
	res := base58.Decode(address)
	fmt.Printf("%x\n", res)
}

func main() {
	Decodemain("1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird")
	Decodemain("1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird")
	Decodemain2("1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird")
	// Decodemain2("1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5smkNpq")
	Decodemain2("1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qmkNpq")
	Decodemain2("2Km9Nu5VJ4aDMtQ4Mc89WKuPmpn5qmkNpq")
}
