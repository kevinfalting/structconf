package confhandler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestRSAHandlerMethods(t *testing.T) {
	handler := &confhandler.RSAHandler{}

	// Generate a new key pair
	privKeyPEM, pubKeyPEM, err := handler.NewPEMKeyPair(2048)
	if err != nil {
		t.Fatalf("failed to generate a new key pair: %v", err)
	}

	// Convert the private key back from PEM
	_, err = handler.PrivateKeyFromPEM(privKeyPEM)
	if err != nil {
		t.Fatalf("failed to convert the private key from PEM: %v", err)
	}

	// Convert the public key back from PEM
	_, err = handler.PublicKeyFromPEM(pubKeyPEM)
	if err != nil {
		t.Fatalf("failed to convert the public key from PEM: %v", err)
	}

	// Check if the keys are set in the handler
	if handler.PrivateKey == nil {
		t.Fatal("PrivateKey in the handler is not set")
	}

	if handler.PublicKey == nil {
		t.Fatal("PublicKey in the handler is not set")
	}

	// Create a test plaintext
	plaintext := "This is a test"

	// Encrypt the plaintext
	ciphertext, err := handler.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("failed to encrypt the plaintext: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := handler.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("failed to decrypt the ciphertext: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if string(decrypted) != plaintext {
		t.Fatalf("Decryption failed. Expected: %s, Got: %s", plaintext, string(decrypted))
	}

	type A struct {
		Password string
	}

	a := A{
		Password: fmt.Sprintf("secret://%s", string(ciphertext)),
	}
	fields, err := stronf.SettableFields(&a)
	if err != nil {
		t.Fatalf("failed to get SettableFields: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}

	// Handle the field
	decryptedValue, err := handler.Handle(context.Background(), fields[0], nil)
	if err != nil {
		t.Fatalf("failed to handle the field: %v", err)
	}

	// Check if the handled value matches the original plaintext
	if decryptedValue != plaintext {
		t.Fatalf("Handle method failed. Expected: %s, Got: %s", plaintext, decryptedValue)
	}

	plaintext = "Hello World!"

	ciphertext, err = handler.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("failed to encrypt the plaintext: %v", err)
	}

	decryptedValue, err = handler.Handle(context.Background(), fields[0], fmt.Sprintf("secret://%s", string(ciphertext)))
	if err != nil {
		t.Fatalf("failed to handle the field: %v", err)
	}

	if decryptedValue != plaintext {
		t.Fatalf("Handle method failed. Expected: %s, Got: %s", plaintext, decryptedValue)
	}

	var a1 A
	fields, err = stronf.SettableFields(&a1)
	if err != nil {
		t.Fatalf("failed to get SettableFields: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}

	_, err = handler.Handle(context.Background(), fields[0], nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
