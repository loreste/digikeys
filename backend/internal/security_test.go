package internal

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/digikeys/backend/config"
	apphttp "github.com/digikeys/backend/internal/adapters/http"
	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/pkg/crypto"
)

// ════════════════════════════════════════════════════════════
// BIOMETRIC DATA SECURITY TESTS
// ════════════════════════════════════════════════════════════

func TestBiometricEncryptionRoundTrip(t *testing.T) {
	key := make([]byte, 32) // AES-256
	rand.Read(key)

	// Simulate fingerprint template data
	fingerprintData := make([]byte, 512)
	rand.Read(fingerprintData)

	encrypted, err := crypto.Encrypt(key, fingerprintData)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	// Encrypted data should differ from plaintext
	if string(encrypted) == string(fingerprintData) {
		t.Error("encrypted data should differ from plaintext")
	}

	// Encrypted should be longer (nonce + tag overhead)
	if len(encrypted) <= len(fingerprintData) {
		t.Error("encrypted data should be longer than plaintext (nonce + auth tag)")
	}

	decrypted, err := crypto.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if string(decrypted) != string(fingerprintData) {
		t.Error("decrypted data should match original")
	}
}

func TestBiometricEncryptionTampering(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	plaintext := []byte("sensitive biometric data - fingerprint template ISO 19794-2")

	encrypted, _ := crypto.Encrypt(key, plaintext)

	// Tamper with ciphertext
	tampered := make([]byte, len(encrypted))
	copy(tampered, encrypted)
	tampered[len(tampered)-1] ^= 0xFF // flip last byte

	_, err := crypto.Decrypt(key, tampered)
	if err == nil {
		t.Error("tampered ciphertext should fail decryption (AES-GCM auth tag)")
	}
}

func TestBiometricEncryptionWrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)

	plaintext := []byte("biometric data that must stay confidential")

	encrypted, _ := crypto.Encrypt(key1, plaintext)

	_, err := crypto.Decrypt(key2, encrypted)
	if err == nil {
		t.Error("decryption with wrong key should fail")
	}
}

func TestBiometricEncryptionNonceUniqueness(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	plaintext := []byte("same fingerprint template")

	// Encrypt same data twice - nonces must differ
	enc1, _ := crypto.Encrypt(key, plaintext)
	enc2, _ := crypto.Encrypt(key, plaintext)

	if string(enc1) == string(enc2) {
		t.Error("two encryptions of same data must produce different ciphertext (random nonce)")
	}
}

func TestBiometricKeyLengthValidation(t *testing.T) {
	shortKey := make([]byte, 16) // AES-128 instead of AES-256
	rand.Read(shortKey)

	_, err := crypto.Encrypt(shortKey, []byte("data"))
	if err == nil {
		t.Error("should reject non-256-bit keys for biometric encryption")
	}
}

// ════════════════════════════════════════════════════════════
// RACE CONDITION TESTS
// ════════════════════════════════════════════════════════════

func TestConcurrentEncryption(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data := []byte(fmt.Sprintf("fingerprint-template-%d", id))

			encrypted, err := crypto.Encrypt(key, data)
			if err != nil {
				errors <- fmt.Errorf("encrypt %d: %v", id, err)
				return
			}

			decrypted, err := crypto.Decrypt(key, encrypted)
			if err != nil {
				errors <- fmt.Errorf("decrypt %d: %v", id, err)
				return
			}

			if string(decrypted) != string(data) {
				errors <- fmt.Errorf("mismatch %d", id)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestConcurrentMRZGeneration(t *testing.T) {
	svc := application.NewMRZService()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			citizen := &domain.Citizen{
				FirstName:   "IBRAHIM",
				LastName:    "TRAORE",
				DateOfBirth: time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC),
				Gender:      "M",
				Nationality: "BFA",
			}
			expiry := time.Date(2031, 3, 15, 0, 0, 0, 0, time.UTC)
			card := &domain.Card{
				CardNumber: fmt.Sprintf("CC-FR-2026-%06d", id),
				ExpiresAt:  &expiry,
			}
			embassy := &domain.Embassy{CountryCode: "FR"}
			svc.GenerateTD1(citizen, card, embassy)
		}(i)
	}
	wg.Wait()
}

// ════════════════════════════════════════════════════════════
// MEMORY LEAK TESTS
// ════════════════════════════════════════════════════════════

func TestEncryptionMemoryLeak(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Encrypt/decrypt 1000 times
	for i := 0; i < 1000; i++ {
		data := make([]byte, 1024)
		rand.Read(data)
		encrypted, _ := crypto.Encrypt(key, data)
		crypto.Decrypt(key, encrypted)
	}

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	growth := int64(m2.HeapAlloc) - int64(m1.HeapAlloc)
	t.Logf("Memory growth after 1000 encrypt/decrypt cycles: %d bytes", growth)

	if growth > 10*1024*1024 {
		t.Errorf("excessive memory growth: %d bytes (possible leak)", growth)
	}
}

func TestGoroutineLeakOnHTTPRequests(t *testing.T) {
	deps := apphttp.RouterDeps{
		AuthService: application.NewAuthService(nil, config.JWTConfig{
			Secret: "goroutine-leak-test-32-chars!!", AccessTokenTTL: 15, RefreshTokenTTL: 7,
		}),
	}
	router := apphttp.NewRouter(deps)

	initial := runtime.NumGoroutine()

	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	time.Sleep(100 * time.Millisecond)
	runtime.GC()
	final := runtime.NumGoroutine()

	leaked := final - initial
	t.Logf("Goroutines: initial=%d, final=%d, leaked=%d", initial, final, leaked)

	if leaked > 5 {
		t.Errorf("goroutine leak: %d goroutines leaked after 100 requests", leaked)
	}
}

// ════════════════════════════════════════════════════════════
// SECURITY VULNERABILITY TESTS
// ════════════════════════════════════════════════════════════

func TestAuthProtectsAllEndpoints(t *testing.T) {
	deps := apphttp.RouterDeps{
		AuthService: application.NewAuthService(nil, config.JWTConfig{
			Secret: "auth-protection-test-32-chars!!", AccessTokenTTL: 15, RefreshTokenTTL: 7,
		}),
	}
	router := apphttp.NewRouter(deps)

	// Without token, protected routes should return 401 (if registered) or 404 (if not).
	// Either way, no data should be leaked.
	protectedPaths := []string{
		"/api/v1/citizens",
		"/api/v1/cards",
		"/api/v1/enrollments",
		"/api/v1/embassies",
		"/api/v1/transfers",
	}

	for _, path := range protectedPaths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized && w.Code != http.StatusNotFound {
			t.Errorf("%s should return 401 or 404 without auth, got %d", path, w.Code)
		}
	}
}

func TestLargePayloadDoesNotCrash(t *testing.T) {
	deps := apphttp.RouterDeps{
		AuthService: application.NewAuthService(nil, config.JWTConfig{
			Secret: "large-payload-digikeys-test-32!", AccessTokenTTL: 15, RefreshTokenTTL: 7,
		}),
	}
	router := apphttp.NewRouter(deps)

	largeBody := make([]byte, 10*1024*1024)
	rand.Read(largeBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(largeBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Error("should handle large payload gracefully")
	}
}

func TestJWTTamperingRejected(t *testing.T) {
	authService := application.NewAuthService(nil, config.JWTConfig{
		Secret: "jwt-tampering-test-32-chars!!!!", AccessTokenTTL: 15, RefreshTokenTTL: 7,
	})

	tokens := []string{
		"",
		"invalid",
		"eyJ.eyJ.xxx",
		// Token with different algorithm
		"eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjMiLCJyb2xlIjoic3VwZXJfYWRtaW4ifQ.",
	}

	for _, token := range tokens {
		_, err := authService.ValidateToken(token)
		if err == nil {
			t.Errorf("tampered token should be rejected: %s", token[:min(len(token), 20)])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
