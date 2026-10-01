package pseudoref

import (
	"strings"
	"testing"
)

func TestRef(t *testing.T) {
	key := []byte("test-key")
	a := Ref(key, "tenant-1")
	if len(a) != Length {
		t.Fatalf("len = %d, want %d", len(a), Length)
	}
	if a != Ref(key, "tenant-1") {
		t.Error("not stable for the same key and id")
	}
	if a == Ref(key, "tenant-2") {
		t.Error("different ids collided")
	}
	if a == Ref([]byte("other-key"), "tenant-1") {
		t.Error("different keys produced the same ref")
	}
	if strings.Contains("tenant-1", a) || strings.Contains(a, "tenant") {
		t.Error("ref leaks the id")
	}
	// Ids sharing a long prefix (e.g. a participant fingerprint) still differ.
	p := "u-joey::12202dbfcf5f9a41157bf4c3672cff7c6e5d60123317f39708fdde1f658115535ab8"
	if Ref(key, p) == Ref(key, strings.Replace(p, "joey", "jimmy", 1)) {
		t.Error("shared-prefix ids collided")
	}
	if Ref(nil, "tenant-1") != "" || Ref(key, "") != "" {
		t.Error("empty key or id must yield an empty ref, never the raw id")
	}
}
