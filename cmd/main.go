package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	value := "1234567890abcdef"

	// Your payload (16 bytes)
	payload := []byte(value) // 16 bytes

	// Generate a new Ed25519 key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	log.Println("Public Key: ", hex.EncodeToString(pubKey))
	log.Println("Private Key: ", hex.EncodeToString(privKey))

	// Sign the payload
	signature := ed25519.Sign(privKey, payload)

	// Verify the signature (optional, just to show it's correct)
	valid := ed25519.Verify(pubKey, payload, signature)
	fmt.Println("Signature valid:", valid)

	// Print results
	fmt.Printf("Payload: %x\n", payload)
	fmt.Printf("Signature: %x\n", signature)
	fmt.Printf("Public Key: %x\n", pubKey)

	hexString := hex.EncodeToString(signature)
	log.Println("Length of signature:")

	// Generate and save the QR code as a PNG file
	err = qrcode.WriteFile(fmt.Sprintf("%s-%s", value, hexString), qrcode.Medium, 256, "signature.png")
	if err != nil {
		log.Fatal("Failed to create QR code:", err)
	}

	fmt.Println("QR code saved to signature.png")
}
