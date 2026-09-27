package selfupdate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
	if _, err := checksumFor("bad", strings.Repeat("g", 64)+"  bad\n"); err == nil {
		t.Fatal("expected non-hex checksum error")
	}
}

func TestLatestReleaseRequiresStableSemVerTag(t *testing.T) {
	for _, tag := range []string{"v1.2.3", "v0.0.0"} {
		t.Run(tag, func(t *testing.T) {
			got, err := latestRelease(releaseTestClient(tag))
			if err != nil || got.TagName != tag {
				t.Fatalf("latestRelease() = %#v, %v", got, err)
			}
		})
	}
	for _, tag := range []string{"", "main", "v1.2", "v01.2.3", "v1.2.3-rc.1", "v1.2.3/../../other"} {
		t.Run("reject-"+tag, func(t *testing.T) {
			if _, err := latestRelease(releaseTestClient(tag)); err == nil {
				t.Fatalf("latestRelease accepted invalid tag %q", tag)
			}
		})
	}
}

func releaseTestClient(tag string) *http.Client {
	body, _ := json.Marshal(githubRelease{TagName: tag})
	return &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader(body)),
			Request:    request,
		}, nil
	})}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }
