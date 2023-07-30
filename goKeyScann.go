package gokeypubscan

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Keys struct {
	address      string
	public       ecdsa.PublicKey
	publicValue  string
	publicBytes  []byte
	private      ecdsa.PrivateKey
	privateValue string
	hash256      [32]byte
	x            string
	y            string
	WalletBalance
	withLocalNode bool
}

type WalletBalance struct {
	Balance      float64 `json:"balance"`
	FinalBalance float64 `json:"final_balance"`
}

type NodeData struct {
	url     string
	port    string
	userRcp string
	passRcp string
	useSSL  bool
}

var accessNodeData NodeData

func (keys *Keys) GeneratePrivKey() {
	reader := rand.Reader
	// Generate a random 256-bit private key.
	private_key, err := ecdsa.GenerateKey(elliptic.P256(), reader)
	if err != nil {
		panic(err)
	}
	keys.private = *private_key
	keys.x = private_key.X.String()
	keys.y = private_key.Y.String()
	keys.privateValue = private_key.D.String()
}

func (keys *Keys) GeneratePublicKey() {
	// Derive the public key from the private key.
	public_key := keys.private.PublicKey
	keys.public = public_key

	// Encode the public key in uncompressed format.
	public_key_bytes := elliptic.Marshal(elliptic.P256(), public_key.X, public_key.Y)
	keys.publicBytes = public_key_bytes
	keys.publicValue = string(public_key_bytes)

	// Apply SHA-256 hash to the public key.
	sha256Hash := sha256.Sum256(public_key_bytes)

	keys.hash256 = sha256Hash

}

func (keys *Keys) GenerateAddress() {
	// Apply RIPEMD-160 hash to the SHA-256 hash.
	ripemd160Hash := ripemd160.New()
	_, err := ripemd160Hash.Write(keys.hash256[:])
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

	keys.address = address
}

func (keys *Keys) GetBalance() {

	if keys.withLocalNode {
		rpcUser := accessNodeData.userRcp
		rpcPass := accessNodeData.passRcp
		rpcHost := accessNodeData.url + accessNodeData.port // Dirección del nodo de Bitcoin
		useSSL := accessNodeData.useSSL

		// Dirección de la wallet que quieres consultar
		address := keys.address

		// Crea una configuración de conexión
		connCfg := &rpcclient.ConnConfig{
			Host:         rpcHost,
			User:         rpcUser,
			Pass:         rpcPass,
			DisableTLS:   useSSL,
			HTTPPostMode: true,
		}

		// Crea un cliente de RPC para conectar al nodo de Bitcoin
		client, err := rpcclient.New(connCfg, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Shutdown()

		// Obtiene el saldo de la dirección de la wallet
		balance, err := client.GetBalance(address)
		if err != nil {
			log.Fatal(err)
		}
		// Estructura para almacenar la respuesta JSON
		var walletBalance WalletBalance
		walletBalance.Balance = float64(balance)
		walletBalance.FinalBalance = balance.ToBTC()
		keys.WalletBalance = walletBalance
		return
	}

	result, err := http.Get("https://api.blockcypher.com/v1/btc/main/addrs/" + keys.address + "/balance")
	if err != nil {
		log.Print("Error consultado saldo.")
		return
	}
	defer result.Body.Close()
	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println("Error al leer la respuesta:", err)
		return
	}
	// Estructura para almacenar la respuesta JSON
	var walletBalance WalletBalance

	// Decodificar la respuesta JSON en la estructura
	err = json.Unmarshal(body, &walletBalance)
	if err != nil {
		log.Println("Error al decodificar la respuesta JSON:", err)
		return
	}

	keys.WalletBalance = walletBalance
}

func (keys *Keys) isLocalNode(nodeData NodeData) {
	keys.withLocalNode = true
	accessNodeData = nodeData
}
