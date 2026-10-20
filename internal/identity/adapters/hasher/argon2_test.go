package hasher

import "testing"

func TestHashVerify(t *testing.T) {
	h := NewArgon2()
	hash, err := h.Hash("s3cret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := h.Verify("s3cret-pass", hash); err != nil || !ok {
		t.Fatalf("expected match: ok=%v err=%v", ok, err)
	}
	if ok, _ := h.Verify("wrong", hash); ok {
		t.Fatal("expected mismatch")
	}
	hash2, _ := h.Hash("s3cret-pass")
	if hash == hash2 {
		t.Fatal("salt must make hashes differ")
	}
}