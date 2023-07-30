package gokeypubscan

import (
	"fmt"
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	var keys Keys
	keys.GeneratePrivKey()
	keys.GeneratePublicKey()
	keys.GenerateAddress()

	// Print the private key, public key, and Bitcoin address to the console.
	fmt.Println("Private key:", keys.private)
	fmt.Println("Public key:", keys.public)
	fmt.Println("Bitcoin address:", keys.address)

	var keys2 Keys
	keys2.GeneratePrivKey()
	keys2.GeneratePublicKey()
	keys2.GenerateAddress()

	// Print the private key, public key, and Bitcoin address to the console.
	fmt.Println("Private key:", keys2.private)
	fmt.Println("Public key:", keys2.public)
	fmt.Println("Bitcoin address:", keys2.address)
}
