package secrets

import (
	"testing"
	"time"
)

func TestNewRegistry(t *testing.T) {
	r1 := NewRegistry()
	r2 := NewRegistry()
	if r1 == nil || r2 == nil {
		t.Fatalf("Registry should be non nil")
	}
	if r1 != r2 {
		t.Fatalf("Registry should be a singleton: %v != %v", r1, r2)
	}
}

func TestRegistry(t *testing.T) {
	const (
		testfile    = "../skptesting/enckey"
		anotherfile = "../skptesting/static.eskip"
	)
	var (
		enc1 Encryption
		err  error
	)
	r1 := NewRegistry()
	defer r1.Close()

	enc1, err = r1.NewEncrypter(time.Second, testfile)
	if err != nil {
		t.Fatalf("Failed to create encrypter: %v", err)
	}

	enc2, err := r1.NewEncrypter(time.Second, testfile)
	if err != nil {
		t.Fatalf("Failed to create second encrypter: %v", err)
	}
	if enc1 != enc2 {
		t.Fatal("Failed to get the same encrypter")
	}

	enc3, err := r1.NewEncrypter(time.Second, anotherfile)
	if err != nil {
		t.Fatalf("Failed to create third encrypter: %v", err)
	}
	if enc1 == enc3 {
		t.Fatal("Failed to get different encrypter")
	}

	if _, err := r1.NewEncrypter(time.Second, "does-not-exist"); err == nil {
		t.Fatal("Create encrypter should fail if file does not exist")
	}
}
