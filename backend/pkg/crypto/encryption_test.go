package crypto

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func validKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return key
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := validKey()
	plaintext := []byte("Carte consulaire DIGIKEYS - Données sensibles biométriques")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should differ from plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted text does not match original: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptEmptyData(t *testing.T) {
	key := validKey()
	plaintext := []byte{}

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt of empty data failed: %v", err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt of empty data failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Fatalf("expected empty decrypted data, got %d bytes", len(decrypted))
	}
}

func TestEncryptInvalidKeyLength(t *testing.T) {
	shortKey := make([]byte, 16) // AES-128, not AES-256
	_, err := Encrypt(shortKey, []byte("test"))
	if err == nil {
		t.Fatal("expected error for 16-byte key, got nil")
	}

	longKey := make([]byte, 64)
	_, err = Encrypt(longKey, []byte("test"))
	if err == nil {
		t.Fatal("expected error for 64-byte key, got nil")
	}
}

func TestDecryptInvalidKeyLength(t *testing.T) {
	_, err := Decrypt(make([]byte, 16), []byte("someciphertext"))
	if err == nil {
		t.Fatal("expected error for invalid key length in Decrypt")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := validKey()
	key2 := validKey()

	ciphertext, err := Encrypt(key1, []byte("secret data"))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(key2, ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	key := validKey()
	plaintext := []byte("tamper test data")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Flip a byte in the ciphertext (after the nonce)
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0xFF

	_, err = Decrypt(key, tampered)
	if err == nil {
		t.Fatal("expected error when decrypting tampered ciphertext")
	}
}

func TestDecryptTooShortCiphertext(t *testing.T) {
	key := validKey()
	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	key := validKey()
	plaintext := []byte("same data each time")

	ct1, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt 1 failed: %v", err)
	}
	ct2, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt 2 failed: %v", err)
	}

	if bytes.Equal(ct1, ct2) {
		t.Fatal("two encryptions of the same plaintext should produce different ciphertexts (random nonce)")
	}
}
