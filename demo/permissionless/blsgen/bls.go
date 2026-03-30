// blsgen derives BLS public key and proof-of-possession from an ECDSA nodekey (hex).
// Usage: blsgen <nodekey-hex>
// Output: two lines — "pub=0x..." and "pop=0x..."
package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: blsgen <nodekey-hex>")
		os.Exit(1)
	}
	raw, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode: %v\n", err)
		os.Exit(1)
	}
	priv, err := crypto.ToECDSA(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ecdsa: %v\n", err)
		os.Exit(1)
	}
	sk, err := bls.DeriveFromECDSA(priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bls derive: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("pub=0x%x\n", sk.PublicKey().Marshal())
	fmt.Printf("pop=0x%x\n", bls.PopProve(sk).Marshal())
}
