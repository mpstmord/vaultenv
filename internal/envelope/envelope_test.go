package envelope_test

import (
	"bytes"
	"testing"

	"github.com/your-org/vaultenv/internal/envelope"
)

func key32() []byte {
	return bytes.Repeat([]byte{0x42}, 32)
}

func TestNew_Valid(t *testing.T) {
	c, err := envelope.New(key32())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Cipher")
	}
}

func TestNew_InvalidKeySize(t *testing.T) {
	for _, size := range []int{0, 16, 24, 31, 33, 64} {
		_, err := envelope.New(make([]byte, size))
		if err == nil {
			t.Errorf("expected error for key size %d", size)
		}
	}
}

func TestSeal_And_Open_RoundTrip(t *testing.T) {
	c, _ := envelope.New(key32())
	plaintext := []byte("super-secret-value")

	sealed, err := c.Seal(plaintext)
	if err != nil {
		t.Fatalf("Seal error: %v", err)
	}
	if bytes.Equal(sealed, plaintext) {
		t.Fatal("sealed output should differ from plaintext")
	}

	got, err := c.Open(sealed)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("got %q, want %q", got, plaintext)
	}
}

func TestSeal_ProducesUniqueOutputEachCall(t *testing.T) {
	c, _ := envelope.New(key32())
	plaintext := []byte("determinism-check")

	a, _ := c.Seal(plaintext)
	b, _ := c.Seal(plaintext)
	if bytes.Equal(a, b) {
		t.Fatal("two Seal calls should produce different ciphertext (random nonce)")
	}
}

func TestOpen_TamperedCiphertext(t *testing.T) {
	c, _ := envelope.New(key32())
	sealed, _ := c.Seal([]byte("value"))

	// flip a byte in the ciphertext portion
	sealed[len(sealed)-1] ^= 0xFF

	_, err := c.Open(sealed)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestOpen_TooShort(t *testing.T) {
	c, _ := envelope.New(key32())
	_, err := c.Open([]byte{0x01, 0x02})
	if err == nil {
		t.Fatal("expected error for too-short input")
	}
}

func TestOpen_WrongKey(t *testing.T) {
	c1, _ := envelope.New(bytes.Repeat([]byte{0x01}, 32))
	c2, _ := envelope.New(bytes.Repeat([]byte{0x02}, 32))

	sealed, _ := c1.Seal([]byte("secret"))
	_, err := c2.Open(sealed)
	if err == nil {
		t.Fatal("expected error when opening with wrong key")
	}
}
