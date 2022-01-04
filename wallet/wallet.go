package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/devholic77/duckcoin/utils"
)

// const (
// 	signature   string = "8a4b76cf0007d0481ef4b066fbc515cbefc065a343e3d60f7c3af13d8aef6b032faddf2fc2f5e4afc29e28ec4edcbf02968f0a7638d7de308a57783b2dd218d4"
// 	hashMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
// 	privateKey  string = "30770201010420b18093557b09c041cb6806db3b541e24ea39d3baa528b0d79e6889250b9ab10ea00a06082a8648ce3d030107a14403420004b5231efa3617d805ca35027ae818353a80d3862368293956c5b5c6084626965141b3e49cbaf1d4da3bbbe4f13a1b428a88f999d67d70d0cc1f618b61497b5a67"
// )

var w *wallet

const walletFileName string = "duckcoin.wallet"

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

func persistKey(key *ecdsa.PrivateKey) {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)

	err = os.WriteFile(walletFileName, keyBytes, 0664)
	utils.HandleErr(err)
}

func restoreKey() *ecdsa.PrivateKey {
	bytes, err := os.ReadFile(walletFileName)
	utils.HandleErr(err)

	key, err := x509.ParseECPrivateKey(bytes)
	utils.HandleErr(err)

	return key
}

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	return !os.IsNotExist(err)
}

func Sign(payload string, w *wallet) string {
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadBytes)
	utils.HandleErr(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(signature string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	var bigA, bigB = big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

func Verify(signature string, payload string, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			// restore from wallet
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privateKey
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Start() {
	// privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	// fmt.Printf("%x\n\n", keyAsBytes)
	// utils.HandleErr(err)

	// message := "i love you"
	// hashedMessage := utils.Hash(message)
	// hashAsBytes, err := hex.DecodeString(hashedMessage)
	// fmt.Printf("%x\n\n", hashAsBytes)
	// utils.HandleErr(err)

	// r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	// utils.HandleErr(err)
	// signature := append(r.Bytes(), s.Bytes()...)

	// fmt.Printf("%x\n", signature)

	// ok := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)
	// fmt.Println(ok)

	// privByte, err := hex.DecodeString(privateKey)
	// utils.HandleErr(err)

	// private, err := x509.ParseECPrivateKey(privByte)
	// utils.HandleErr(err)

	// sigBytes, err := hex.DecodeString(signature)
	// utils.HandleErr(err)

	// rBytes := sigBytes[:len(sigBytes)/2]
	// sBytes := sigBytes[len(sigBytes)/2:]

	// var bigR, bigS = big.Int{}, big.Int{}
	// bigR.SetBytes(rBytes)
	// bigS.SetBytes(sBytes)

	// hashBytes, err := hex.DecodeString(hashMessage)
	// ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)
	// fmt.Println(ok)

}
