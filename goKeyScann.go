package gokeypubscan

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func GenerateKey() {
	reader := rand.Reader
	// Generate a random 256-bit private key.
	private_key, err := ecdsa.GenerateKey(elliptic.P256(), reader)
	if err != nil {
		panic(err)
	}

	// Derive the public key from the private key.
	public_key := private_key.PublicKey

	// Encode the public key in uncompressed format.
	public_key_bytes := elliptic.Marshal(elliptic.P256(), public_key.X, public_key.Y)

	// Apply SHA-256 hash to the public key.
	sha256Hash := sha256.Sum256(public_key_bytes)

	// Apply RIPEMD-160 hash to the SHA-256 hash.
	ripemd160Hash := ripemd160.New()
	_, err = ripemd160Hash.Write(sha256Hash[:])
	if err != nil {
		panic(err)
	}
	ripemd160HashBytes := ripemd160Hash.Sum(nil)

	// Add version byte (0x00) in front of the RIPEMD-160 hash.
	extendedRipemd160HashBytes := append([]byte{0x00}, ripemd160HashBytes...)

	// Calculate the checksum by applying SHA-256 twice to the extended RIPEMD-160 hash.
	checksum := sha256.Sum256(extendedRipemd160HashBytes)
	checksum = sha256.Sum256(checksum[:])

	// Take the first 4 bytes of the second SHA-256 hash as the address checksum.
	var checksumBytes [4]byte
	copy(checksumBytes[:], checksum[:4])

	// Add the 4 checksum bytes to the extended RIPEMD-160 hash.
	finalBytes := append(extendedRipemd160HashBytes, checksumBytes[:]...)

	// Encode the final byte slice in Base58 to get the Bitcoin address.
	address := base58.Encode(finalBytes)

	// Print the private key, public key, and Bitcoin address to the console.
	fmt.Println("Private key:", private_key)
	fmt.Println("Public key:", public_key_bytes)
	fmt.Println("Bitcoin address:", address)
}
