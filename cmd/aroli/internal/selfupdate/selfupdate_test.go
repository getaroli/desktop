package selfupdate

import (
	"strings"
	"testing"
)

func TestCLIAssetName(t *testing.T) {
	cases := []struct {
		goos, arch, want string
		ok               bool
	}{
		{"linux", "amd64", "aroli-linux-amd64", true},
		{"linux", "arm64", "aroli-linux-arm64", true},
		{"darwin", "arm64", "", false},
		{"linux", "386", "", false},
	}
	for _, test := range cases {
		got, err := cliAssetName(test.goos, test.arch)
		if (err == nil) != test.ok || got != test.want {
			t.Errorf("cliAssetName(%q, %q) = %q, %v", test.goos, test.arch, got, err)
		}
	}
}

func TestChecksumFor(t *testing.T) {
	checksum := strings.Repeat("a", 64)
	contents := checksum + "  aroli-linux-amd64\n" + strings.Repeat("b", 64) + " *aroli-linux-arm64\n"
	got, err := checksumFor("aroli-linux-arm64", contents)
	if err != nil || got != strings.Repeat("b", 64) {
		t.Fatalf("checksumFor returned %q, %v", got, err)
	}
	if _, err := checksumFor("missing", contents); err == nil {
		t.Fatal("expected missing checksum error")
	}
	if _, err := checksumFor("bad", "short  bad\n"); err == nil {
		t.Fatal("expected malformed checksum error")
	}
}
