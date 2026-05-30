package keyring

import (
	"testing"

	"github.com/lestrrat-go/jwx/v2/jwe"
)

// ref: defradb-rs/keyring/tests/crypto_tests.rs -> test_go_decrypt_jwe
// go test -v -test.fullpath=true -run ^TestRsKeyring_CreateJWE$ github.com/sourcenetwork/defradb/keyring
func TestRsKeyring_CreateJWE(t *testing.T) {
	var plain = []byte("payload")
	var password = []byte("secret")
	cipher, err := jwe.Encrypt(plain, jwe.WithKey(keyEncryptionAlgorithm, password))
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	println("RsKeyring plain=payload password=secret cipher=", string(cipher))
}

// ref: defradb-rs/keyring/tests/crypto_tests.rs -> test_go_decrypt_jwe
// go test -v -test.fullpath=true -run ^TestRsKeyring_DecryptJWE$ github.com/sourcenetwork/defradb/keyring
func TestRsKeyring_DecryptJWE(t *testing.T) {
	var cipher = []byte("eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJlbmMiOiJBMjU2R0NNIiwicDJjIjoxMDAwMCwicDJzIjoiN2czeXpxdTlzSmpvZlMwX3NjcUpuQSJ9.UJD-eWbdpTuaQAdMO0HxzEmIlRzHwsvWSMzfyYTu5hXsYBIBAbmDnw.rPFof2XUfcL5a1Sa.rVb0Gla8Ng.fUHZKq1nGpum1RXDbEi5Rg")
	var password = []byte("secret")
	plain, err := jwe.Decrypt(cipher, jwe.WithKey(keyEncryptionAlgorithm, password))
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	println("RsKeyring plain=", string(plain), "plain(should be)=payload password=secret cipher=", string(cipher))

	// assert equals plain and "payload"
	if string(plain) != "payload" {
		t.Fatalf("expected plain to be 'payload', got '%s'", string(plain))
	}
}
