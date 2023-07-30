# GoKeyPubScan

Este es un generador de direcciones de Bitcoin en Golang, su implementación es sencilla

## Instalación

```bash
go get -u github.com/cristianino/go_key_scan
```

## Uso

```go
package main

import (
	"fmt"

	"github.com/cristianino/go_key_scan"
)

func main() {
	var keys Keys
	keys.GeneratePrivKey()
	keys.GeneratePublicKey()
	keys.GenerateAddress()

	// Print the private key, public key, and Bitcoin address to the console.
	fmt.Println("Private key:", keys.private)
	fmt.Println("Public key:", keys.public)
	fmt.Println("Bitcoin address:", keys.address)
}
```